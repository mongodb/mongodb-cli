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

//go:build integration

package standbyclusters_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

// These tests run the standby-clusters commands of the real mongocli binary
// against throwaway MinIO containers. Skipped when no docker daemon is available.

const (
	minioUser     = "mcli-it"
	minioPassword = "mcli-it-secret"
	minioBucket   = "mcli-standby-it"
	minioImage    = "minio/minio:RELEASE.2025-09-07T16-13-09Z"
	awsTestRegion = "us-east-1"
	testCluster   = "it-cluster"
)

var cliBinary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "mongocli-it-")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "create temp dir for mongocli binary: %v\n", err)
		os.Exit(1)
	}

	cliBinary = filepath.Join(dir, "mongocli")
	build := exec.CommandContext(context.Background(), "go", "build",
		"-o", cliBinary, "github.com/mongodb/mongodb-cli/mongocli/v2/cmd/mongocli")
	build.Stdout = os.Stderr
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "build mongocli binary: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

// Runs `mongocli ops-manager <args...>` with an isolated config home.
func runCLI(t *testing.T, home string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	args = append([]string{"ops-manager"}, args...)
	cmd := exec.CommandContext(context.Background(), cliBinary, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + home,
		"XDG_CONFIG_HOME=" + filepath.Join(home, ".config"),
	}
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func startMinIO(t *testing.T) (endpoint string, client *s3.Client) {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctr, err := minio.Run(context.Background(), minioImage,
		minio.WithUsername(minioUser),
		minio.WithPassword(minioPassword),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(ctr) })

	endpoint, err = ctr.Endpoint(context.Background(), "http")
	require.NoError(t, err)

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(awsTestRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(minioUser, minioPassword, "")),
	)
	require.NoError(t, err)
	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{Bucket: aws.String(minioBucket)})
	require.NoError(t, err)
	return endpoint, client
}

func seedState(t *testing.T, client *s3.Client, key string, state standby.RemoteDRState) {
	t.Helper()
	body, err := json.Marshal(state)
	require.NoError(t, err)
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(minioBucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	require.NoError(t, err)
}

func fetchState(t *testing.T, client *s3.Client, key string) standby.RemoteDRState {
	t.Helper()
	out, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(minioBucket),
		Key:    aws.String(key),
	})
	require.NoError(t, err)
	defer out.Body.Close()
	var state standby.RemoteDRState
	require.NoError(t, json.NewDecoder(out.Body).Decode(&state))
	return state
}

func standbyDoc() standby.RemoteDRState {
	return standby.RemoteDRState{
		State:         standby.StateStandby,
		ClusterName:   testCluster,
		Version:       "1",
		LastModified:  time.Now().UTC().Format(time.RFC3339),
		SchemaVersion: "1",
	}
}

func endpointFlags(endpoint string) []string {
	return []string{
		"--s3BucketName", minioBucket,
		"--awsRegion", awsTestRegion,
		"--awsAuthMode", "staticCredentials",
		"--awsAccessKey", minioUser,
		"--awsSecretKey", minioPassword,
		"--s3BucketEndpoint", endpoint,
	}
}

func TestConfigureDescribeFailover(t *testing.T) {
	endpoint, client := startMinIO(t)
	home := t.TempDir()
	key := "it/configure-describe-failover/dr_status.json"
	seedState(t, client, key, standbyDoc())

	// run configure with non-interactive flags.
	args := append([]string{"standby-clusters", "configure"}, append(endpointFlags(endpoint), "--s3Key", key)...)
	stdout, stderr, err := runCLI(t, home, args...)
	require.NoError(t, err, "stderr: %s", stderr)
	assert.Contains(t, stdout, "Verified access")

	// describe with no flags.
	stdout, stderr, err = runCLI(t, home, "standby-clusters", "describe")
	require.NoError(t, err, "stderr: %s", stderr)
	assert.Contains(t, stdout, standby.StateStandby)
	assert.Contains(t, stdout, testCluster)

	// describe -o json returns the full document.
	stdout, stderr, err = runCLI(t, home, "standby-clusters", "describe", "-o", "json")
	require.NoError(t, err, "stderr: %s", stderr)
	var doc standby.RemoteDRState
	require.NoError(t, json.Unmarshal([]byte(stdout), &doc))
	assert.Equal(t, standby.StateStandby, doc.State)
	assert.Equal(t, "1", doc.Version)

	// failover --force promotes the standby cluster in the file.
	stdout, stderr, err = runCLI(t, home, "standby-clusters", "failover", "--force")
	require.NoError(t, err, "stderr: %s", stderr)
	assert.Contains(t, stdout, "Failover triggered")
	written := fetchState(t, client, key)
	assert.Equal(t, standby.StatePromoteStandby, written.State)
	assert.Equal(t, standby.StateStandby, written.PreviousState)
	assert.Equal(t, "2", written.Version)

	// describe after the failover reports the promoted state.
	stdout, stderr, err = runCLI(t, home, "standby-clusters", "describe", "-o", "json")
	require.NoError(t, err, "stderr: %s", stderr)
	doc = standby.RemoteDRState{}
	require.NoError(t, json.Unmarshal([]byte(stdout), &doc))
	assert.Equal(t, standby.StatePromoteStandby, doc.State)
	assert.Equal(t, standby.StateStandby, doc.PreviousState)
	assert.Equal(t, "2", doc.Version)

	// a second failover is refused.
	_, stderr, err = runCLI(t, home, "standby-clusters", "failover", "--force")
	require.Error(t, err)
	assert.Contains(t, stderr, "a failover can only be triggered from")
}

func TestDescribeMissingKey(t *testing.T) {
	endpoint, _ := startMinIO(t)
	args := append([]string{"standby-clusters", "describe"}, endpointFlags(endpoint)...)
	_, stderr, err := runCLI(t, t.TempDir(), args...)
	require.Error(t, err)
	assert.Contains(t, stderr, "no DR state file configured")
}

func TestDescribeBadCredentials(t *testing.T) {
	endpoint, client := startMinIO(t)
	home := t.TempDir()
	key := "it/bad-credentials/dr_status.json"
	seedState(t, client, key, standbyDoc())

	flags := append([]string{"standby-clusters", "describe"},
		append(endpointFlags(endpoint), "--s3Key", key, "--awsSecretKey", "wrong")...)
	_, _, err := runCLI(t, home, flags...)
	require.Error(t, err)
}
