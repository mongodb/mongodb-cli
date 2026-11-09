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

//go:build unit

package standby

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterStatusKey(t *testing.T) {
	tests := []struct {
		name      string
		drFileKey string
		want      string
		wantErr   bool
	}{
		{
			name:      "full key yields the sibling status key",
			drFileKey: "mycluster93825721173-tony-clarke-B/dr_status_rs0.json",
			want:      "mycluster93825721173-tony-clarke-B/cluster_status_v1.bson",
		},
		{
			name:      "multiple path segments keep the whole prefix",
			drFileKey: "org/proj/cluster-A/dr_status_rs0.json",
			want:      "org/proj/cluster-A/cluster_status_v1.bson",
		},
		{
			name:      "bare filename has no prefix root",
			drFileKey: "dr_status_rs0.json",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterStatusKey(tt.drFileKey)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cluster prefix directory")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestFetchClusterStatus_DecodeFixture exercises the fetch decode path without a
// network client: it derives the status key from a DR file key and decodes the
// fixture the object would contain.
func TestFetchClusterStatus_DecodeFixture(t *testing.T) {
	key, err := ClusterStatusKey("mycluster93825721173-tony-clarke-B/dr_status_rs0.json")
	require.NoError(t, err)
	assert.Equal(t, "mycluster93825721173-tony-clarke-B/cluster_status_v1.bson", key)

	data, err := os.ReadFile("testdata/cluster_status_v1.bson")
	require.NoError(t, err)

	status, err := UnmarshalClusterStatus(data)
	require.NoError(t, err)
	require.NotNil(t, status)
	assert.Equal(t, uint32(1), status.SchemaVersion)
	require.Len(t, status.ShardStatus, 1)
	assert.Equal(t, "rs0", status.ShardStatus[0].ShardID)
}

// ---------------------------------------------------------------------------
// S3 state store tests
//
// These tests drive the real aws-sdk-go-v2 S3 client against an in-process
// httptest server that speaks just enough of the S3 GET protocol. Unlike the
// MockStateStore tests (which stub the StateStore interface) and the
// integration tests (which run the built binary in a separate process), these
// exercise (*S3StateStore).FetchClusterStatus / FetchDRState directly, so they
// contribute to in-process coverage.
// ---------------------------------------------------------------------------

const (
	testS3Bucket = "standby-unit-bucket"
	testS3Region = "us-east-1"
)

// fakeS3 is a minimal path-style S3 endpoint. Requests arrive as
// GET /<bucket>/<key> because the store is built with UsePathStyle = true.
type fakeS3 struct {
	mu      sync.Mutex
	objects map[string][]byte // key -> body served with 200
	fail500 map[string]bool   // key -> serve 500 instead (covers the wrapped-error path)
	fail403 map[string]bool   // key -> serve 403 AccessDenied (the ambiguous absent/denied case)
}

func (f *fakeS3) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/"+testS3Bucket+"/")

	f.mu.Lock()
	fail := f.fail500[key]
	forbidden := f.fail403[key]
	f.mu.Unlock()

	if r.Method != http.MethodGet {
		writeS3Error(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "method not allowed")
		return
	}
	if fail {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", "injected failure")
		return
	}
	if forbidden {
		writeS3Error(w, http.StatusForbidden, "AccessDenied", "Access Denied")
		return
	}
	body, ok := f.objects[key]
	if !ok {
		writeS3Error(w, http.StatusNotFound, "NoSuchKey", "The specified key does not exist.")
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("ETag", `"test-etag"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// writeS3Error emits the S3 XML error envelope the SDK deserializes. The HTTP
// status is what s3StatusCode classifies on; the <Code> keeps it realistic.
func writeS3Error(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(
		w,
		"<Error><Code>%s</Code><Message>%s</Message><RequestId>test-request-id</RequestId></Error>",
		code, message,
	)
}

func startFakeS3(t *testing.T, objects map[string][]byte) (*fakeS3, string) {
	t.Helper()
	f := &fakeS3{objects: objects, fail500: map[string]bool{}, fail403: map[string]bool{}}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return f, srv.URL
}

// newTestS3Store builds the store through the same production constructor the
// CLI uses, pointing the SDK at the httptest server via the custom endpoint +
// path-style knobs and static test credentials.
func newTestS3Store(t *testing.T, endpoint, key string) *S3StateStore {
	t.Helper()
	store, err := NewS3StateStore(context.Background(), Credentials{
		AuthMode:        AuthModeStaticCredentials,
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		BucketName:      testS3Bucket,
		Region:          testS3Region,
		Endpoint:        endpoint,
		UsePathStyle:    true,
	}, key)
	require.NoError(t, err)
	return store
}

func TestS3StateStore_FetchClusterStatus(t *testing.T) {
	fixture, err := os.ReadFile("testdata/cluster_status_v1.bson")
	require.NoError(t, err)

	drKey := "cluster-A/dr_status_rs0.json"
	statusKey, err := ClusterStatusKey(drKey)
	require.NoError(t, err)
	require.Equal(t, "cluster-A/cluster_status_v1.bson", statusKey)

	t.Run("found decodes the cluster status", func(t *testing.T) {
		_, endpoint := startFakeS3(t, map[string][]byte{statusKey: fixture})
		store := newTestS3Store(t, endpoint, drKey)

		status, err := store.FetchClusterStatus(context.Background())
		require.NoError(t, err)
		require.NotNil(t, status)
		assert.Equal(t, uint32(1), status.SchemaVersion)
		require.NotNil(t, status.ClusterStatus.EstimatedDataLossSeconds)
		assert.Equal(t, 40.0, *status.ClusterStatus.EstimatedDataLossSeconds)
		require.Len(t, status.ShardStatus, 1)
		assert.Equal(t, "rs0", status.ShardStatus[0].ShardID)
		require.NotNil(t, status.ShardStatus[0].LagSeconds)
		assert.Equal(t, 30.0, *status.ShardStatus[0].LagSeconds)
	})

	t.Run("missing object reports not found", func(t *testing.T) {
		_, endpoint := startFakeS3(t, map[string][]byte{})
		store := newTestS3Store(t, endpoint, drKey)

		status, err := store.FetchClusterStatus(context.Background())
		require.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "cluster status object not found")
		assert.Contains(t, err.Error(), "s3://"+testS3Bucket+"/"+statusKey)
		assert.NotContains(t, err.Error(), "error fetching cluster status")
	})

	t.Run("server error is wrapped, not misclassified as not found", func(t *testing.T) {
		f, endpoint := startFakeS3(t, map[string][]byte{statusKey: fixture})
		f.fail500[statusKey] = true
		store := newTestS3Store(t, endpoint, drKey)

		status, err := store.FetchClusterStatus(context.Background())
		require.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "error fetching cluster status from s3://"+testS3Bucket+"/"+statusKey)
		assert.NotContains(t, err.Error(), "cluster status object not found")
	})

	t.Run("forbidden object reports not found with the access caveat", func(t *testing.T) {
		f, endpoint := startFakeS3(t, map[string][]byte{statusKey: fixture})
		f.fail403[statusKey] = true
		store := newTestS3Store(t, endpoint, drKey)

		status, err := store.FetchClusterStatus(context.Background())
		require.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "not found")
		assert.Contains(t, err.Error(), "lack access")
		assert.Contains(t, err.Error(), "s3://"+testS3Bucket+"/"+statusKey)
		// The original SDK cause must be preserved, not discarded.
		assert.Contains(t, err.Error(), "AccessDenied")
		// 403 must take the not-found-plus-caveat path, not the generic wrapped error.
		assert.NotContains(t, err.Error(), "error fetching cluster status")
	})
}

func TestS3StateStore_FetchDRState(t *testing.T) {
	drKey := "cluster-A/dr_status_rs0.json"
	drJSON := []byte(
		`{"state":"Standby","clusterName":"cluster-A","version":"1",` +
			`"lastModified":"2026-08-19T12:25:23Z","schemaVersion":"1"}`,
	)

	t.Run("found decodes the state and etag", func(t *testing.T) {
		_, endpoint := startFakeS3(t, map[string][]byte{drKey: drJSON})
		store := newTestS3Store(t, endpoint, drKey)

		state, etag, err := store.FetchDRState(context.Background())
		require.NoError(t, err)
		require.NotNil(t, state)
		assert.Equal(t, StateStandby, state.State)
		assert.Equal(t, "cluster-A", state.ClusterName)
		assert.Equal(t, "1", state.Version)
		assert.Equal(t, `"test-etag"`, etag)
	})

	t.Run("missing object returns nil state and no error", func(t *testing.T) {
		_, endpoint := startFakeS3(t, map[string][]byte{})
		store := newTestS3Store(t, endpoint, drKey)

		state, etag, err := store.FetchDRState(context.Background())
		require.NoError(t, err)
		assert.Nil(t, state)
		assert.Empty(t, etag)
	})

	t.Run("server error is wrapped, not treated as absent", func(t *testing.T) {
		f, endpoint := startFakeS3(t, map[string][]byte{drKey: drJSON})
		f.fail500[drKey] = true
		store := newTestS3Store(t, endpoint, drKey)

		state, etag, err := store.FetchDRState(context.Background())
		require.Error(t, err)
		assert.Nil(t, state)
		assert.Empty(t, etag)
		assert.Contains(t, err.Error(), "error fetching DR state from s3://"+testS3Bucket+"/"+drKey)
	})

	// A 403 is ambiguous (absent vs denied) and, unlike 404, is surfaced as an
	// error carrying both possibilities so the failover flow does not claim the
	// file is gone when it may merely be denied.
	t.Run("forbidden object returns the combined absent/access error", func(t *testing.T) {
		f, endpoint := startFakeS3(t, map[string][]byte{drKey: drJSON})
		f.fail403[drKey] = true
		store := newTestS3Store(t, endpoint, drKey)

		state, etag, err := store.FetchDRState(context.Background())
		require.Error(t, err)
		assert.Nil(t, state)
		assert.Empty(t, etag)
		assert.Contains(t, err.Error(), "not found")
		assert.Contains(t, err.Error(), "lack access")
		assert.Contains(t, err.Error(), "s3://"+testS3Bucket+"/"+drKey)
		// The original SDK cause must be preserved, not discarded.
		assert.Contains(t, err.Error(), "AccessDenied")
		// 403 must take the combined path, not the generic wrapped error.
		assert.NotContains(t, err.Error(), "error fetching DR state")
	})
}

// classifyGetObjectError must keep the caller's forbidden prose, wrap the
// original cause, and re-apply the auth hint so a real credential failure keeps
// both its cause and the aws-sso-login guidance.
func TestClassifyGetObjectError_ForbiddenKeepsCauseAndHint(t *testing.T) {
	cause := errors.New("refresh cached SSO token failed")
	err := &smithyhttp.ResponseError{
		Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusForbidden}},
		Err:      cause,
	}

	handled, out := classifyGetObjectError(err, "not-found-message", "forbidden-message")
	require.True(t, handled)
	require.Error(t, out)
	assert.Contains(t, out.Error(), "forbidden-message")
	assert.Contains(t, out.Error(), "refresh cached SSO token failed")
	assert.Contains(t, out.Error(), "aws sso login")
	assert.ErrorIs(t, out, cause, "the original cause must remain unwrappable")
}

// A 404 with an empty notFound message means "absent is not an error": the DR
// state file's (nil, "", nil) contract.
func TestClassifyGetObjectError_NotFoundEmptyMessageIsNil(t *testing.T) {
	err := &smithyhttp.ResponseError{
		Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
		Err:      errors.New("NoSuchKey"),
	}

	handled, out := classifyGetObjectError(err, "", "forbidden-message")
	require.True(t, handled)
	assert.NoError(t, out)
}

// Anything that is neither 404 nor 403 is left for the caller's own wrapping.
func TestClassifyGetObjectError_OtherIsUnhandled(t *testing.T) {
	err := &smithyhttp.ResponseError{
		Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusInternalServerError}},
		Err:      errors.New("InternalError"),
	}

	handled, out := classifyGetObjectError(err, "not-found-message", "forbidden-message")
	assert.False(t, handled)
	assert.NoError(t, out)
}
