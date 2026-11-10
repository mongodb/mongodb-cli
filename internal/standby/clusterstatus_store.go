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
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// clusterStatusObjectName is the fixed filename of the cluster status object
// published next to the DR state file in the same cluster prefix directory.
const clusterStatusObjectName = "cluster_status_v1.bson"

// ClusterStatusKey derives the blob-storage key of the cluster status object
// from the configured DR state file key. The status object lives at the root of
// the same cluster prefix directory as the DR file, so the DR filename is
// replaced and the directory prefix is preserved.
//
// Returns an error if drFileKey has no "/": a bare filename means the key
// carries no cluster prefix directory, so the status object location would be
// meaningless.
func ClusterStatusKey(drFileKey string) (string, error) {
	idx := strings.LastIndex(drFileKey, "/")
	if idx < 0 {
		return "", fmt.Errorf(
			"invalid DR state file key %q: the configured key must include the cluster prefix directory, for example \"<cluster-prefix>/dr_status_rs0.json\"",
			drFileKey,
		)
	}
	return drFileKey[:idx+1] + clusterStatusObjectName, nil
}

// FetchClusterStatus retrieves and decodes the cluster_status_v1 object that the
// agent publishes alongside the DR state file, in the same bucket.
func (s *S3StateStore) FetchClusterStatus(ctx context.Context) (*ClusterStatus, error) {
	key, err := ClusterStatusKey(s.key)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("s3://%s/%s", s.bucket, key)

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		// A 404 (absent) and a 403 (S3's ambiguous AccessDenied: the object may be
		// missing or merely denied when the caller lacks s3:ListBucket) share the
		// same not-found diagnosis; classifyGetObjectError centralizes that split
		// with FetchDRState. The 403 variant adds the access caveat and keeps the
		// original cause, so a real credential failure stays actionable.
		notFound := fmt.Sprintf(
			"cluster status object not found at %s; the standby pair must be switchover-capable (linked) and running a build that publishes cluster status",
			path,
		)
		forbidden := fmt.Sprintf(
			"cluster status object not found at %s (or the credentials lack access to it): the standby pair must be switchover-capable (linked) and running a build that publishes cluster status; also verify the credentials' S3 access",
			path,
		)
		if handled, out := classifyGetObjectError(err, notFound, forbidden); handled {
			return nil, out
		}
		return nil, fmt.Errorf("error fetching cluster status from %s: %w", path, withAuthHint(err))
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading cluster status body from %s: %w", path, err)
	}

	status, err := UnmarshalClusterStatus(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing cluster status from %s: %w", path, err)
	}
	return status, nil
}
