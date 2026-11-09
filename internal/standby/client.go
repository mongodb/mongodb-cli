// Copyright 2026 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package standby

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// ErrPreconditionFailed mirrors the agent's blobstore error: returned when an
// If-Match conditional write loses a race against a concurrent writer.
var ErrPreconditionFailed = errors.New("precondition failed")

// ErrObjectDoesNotExist is returned when the DR state object is not found.
var ErrObjectDoesNotExist = errors.New("object does not exist")

type StateStore interface {
	// FetchDRState retrieves the DR status file. Returns the parsed state and
	// the ETag for use in conditional writes, or (nil, "", nil) if absent (404).
	// A 403 is ambiguous (the object may be absent or access may be denied) and
	// is returned as an error carrying both possibilities.
	FetchDRState(ctx context.Context) (*RemoteDRState, string, error)
	// FetchClusterStatus retrieves the cluster status object published next to
	// the DR status file. Returns an error if the object is absent.
	FetchClusterStatus(ctx context.Context) (*ClusterStatus, error)
	// ConditionalPutDRState writes the DR status file guarded by If-Match on
	// the expected ETag. Returns ErrPreconditionFailed on a lost race.
	ConditionalPutDRState(ctx context.Context, newState *RemoteDRState, expectedETag string) error
	// FullPath returns the human-readable s3:// path for messages.
	FullPath() string
}

type S3StateStore struct {
	client *s3.Client
	bucket string
	key    string
}

func NewS3StateStore(ctx context.Context, creds Credentials, key string) (*S3StateStore, error) {
	raw, err := newRawS3Client(ctx, creds)
	if err != nil {
		return nil, err
	}
	return &S3StateStore{client: raw, bucket: creds.BucketName, key: key}, nil
}

// FetchDRState retrieves the DR status file from blob storage.
// Returns the parsed state and the ETag for use in conditional writes.
// Returns (nil, "", nil) if the object does not exist (404). A 403 is treated
// as an error (see the case below) because it cannot be distinguished from a
// missing object without s3:ListBucket.
func (s *S3StateStore) FetchDRState(ctx context.Context) (*RemoteDRState, string, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
	})
	if err != nil {
		// A 403 is deliberately surfaced as an error rather than the 404
		// (nil, "", nil) contract: (nil, "", nil) means "definitely absent", and
		// collapsing 403 into it would make every caller claim the file is gone —
		// e.g. the failover watcher would report "the DR state file ... no longer
		// exists" for a file that may still exist but be denied.
		// classifyGetObjectError keeps the original cause and the auth hint.
		// TriggerFailover aborts either way, so the failover flow is no less safe.
		forbidden := fmt.Sprintf(
			"DR state file not found at %s (or the credentials lack access to it): check your standby-clusters configuration; also verify the credentials' S3 access",
			s.FullPath(),
		)
		if handled, out := classifyGetObjectError(err, "", forbidden); handled {
			return nil, "", out
		}
		return nil, "", fmt.Errorf("error fetching DR state from %s: %w", s.FullPath(), withAuthHint(err))
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading DR state body from %s: %w", s.FullPath(), err)
	}

	state, err := UnmarshalRemoteDRState(body)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing DR state from %s: %w", s.FullPath(), err)
	}

	etag := ""
	if result.ETag != nil {
		etag = *result.ETag
	}
	return state, etag, nil
}

// ConditionalPutDRState writes the DR status file to blob storage, guarded by
// If-Match on the expected ETag. Returns ErrPreconditionFailed if the object
// was modified since the last read, allowing the caller to retry.
func (s *S3StateStore) ConditionalPutDRState(ctx context.Context, newState *RemoteDRState, expectedETag string) error {
	data, err := newState.Marshal()
	if err != nil {
		return fmt.Errorf("error marshaling DR state for %s: %w", s.FullPath(), err)
	}

	input := &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &s.key,
		IfMatch:     aws.String(expectedETag),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
	}
	if _, err := s.client.PutObject(ctx, input); err != nil {
		statusCode := s3StatusCode(err)
		if statusCode == http.StatusPreconditionFailed || statusCode == http.StatusConflict {
			return fmt.Errorf(
				"PutObjectIfMatch precondition failed; key=%s, bucket=%s: %w",
				s.key, s.bucket, ErrPreconditionFailed,
			)
		}
		return fmt.Errorf("PutObjectIfMatch failed; key=%s, bucket=%s: %w", s.key, s.bucket, err)
	}
	return nil
}

func (s *S3StateStore) FullPath() string {
	return fmt.Sprintf("s3://%s/%s", s.bucket, s.key)
}

// TriggerFailover reads the current remote state, applies the Standby ->
// PromoteStandby failover, and writes it back.
// On a lost race (ErrPreconditionFailed) it re-reads and retries once.
func TriggerFailover(ctx context.Context, store StateStore) (*RemoteDRState, error) {
	for attempt := 0; attempt < 2; attempt++ {
		current, etag, err := store.FetchDRState(ctx)
		if err != nil {
			return nil, err
		}
		if current == nil {
			return nil, fmt.Errorf("no DR state file found at %s; check your standby-clusters configuration", store.FullPath())
		}

		if err := current.ApplyFailover(); err != nil {
			return nil, err
		}

		if err := store.ConditionalPutDRState(ctx, current, etag); err != nil {
			if errors.Is(err, ErrPreconditionFailed) {
				continue
			}
			return nil, err
		}
		return current, nil
	}
	return nil, fmt.Errorf("the DR state file at %s changed while writing; re-run the command", store.FullPath())
}

func withAuthHint(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "refresh cached SSO token failed") || strings.Contains(msg, "InvalidGrantException"):
		return fmt.Errorf("%w\n\nyour AWS SSO session has expired or is missing, run 'aws sso login' (optionally with --profile) and try again", err)
	case strings.Contains(msg, "failed to refresh cached credentials"), strings.Contains(msg, "no EC2 IMDS role found"), strings.Contains(msg, "NoCredentialProviders"):
		return fmt.Errorf("%w\n\nno usable AWS credentials found, configure them via 'mongocli ops-manager standby-clusters configure', AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY, or an AWS profile ('aws sso login' for SSO)", err)
	default:
		return err
	}
}

func s3StatusCode(err error) int {
	var respErr *smithyhttp.ResponseError
	if !errors.As(err, &respErr) {
		return -1
	}
	return respErr.HTTPStatusCode()
}

// classifyGetObjectError centralizes the S3 GetObject error handling shared by
// FetchDRState and FetchClusterStatus. S3 answers a GetObject for an absent key
// with 404 NoSuchKey when the caller may list the bucket, but with 403
// AccessDenied when it may not — so a 403 is ambiguous: the object may be absent
// or merely denied.
//
// It returns handled=true for the 404 and 403 cases:
//   - 404: notFound is returned verbatim. An empty notFound means "absent is not
//     an error" and yields a nil error (the DR state file's (nil, "", nil)
//     contract).
//   - 403: the caller's forbidden prose is kept, with the original cause wrapped
//     (%w) and the auth hint re-applied, so a real credential failure keeps both
//     its cause and the SSO/credentials guidance.
//
// handled=false means the error is neither 404 nor 403, and the caller applies
// its own generic wrapping.
func classifyGetObjectError(err error, notFound, forbidden string) (handled bool, out error) {
	switch s3StatusCode(err) {
	case http.StatusNotFound:
		if notFound == "" {
			return true, nil
		}
		return true, errors.New(notFound)
	case http.StatusForbidden:
		return true, fmt.Errorf("%s: %w", forbidden, withAuthHint(err))
	default:
		return false, nil
	}
}
