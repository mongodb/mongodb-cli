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

package standbyclusters

import (
	"testing"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials_AuthModeInference(t *testing.T) {
	tests := []struct {
		name     string
		opts     opts
		expected standby.AuthMode
	}{
		{
			name:     "explicit auth mode wins",
			opts:     opts{authMode: string(standby.AuthModeCredentialsChain), roleARN: "arn:aws:iam::123:role/r"},
			expected: standby.AuthModeCredentialsChain,
		},
		{
			name:     "role ARN infers assumeRole",
			opts:     opts{roleARN: "arn:aws:iam::123:role/r"},
			expected: standby.AuthModeAssumeRole,
		},
		{
			name:     "access key infers staticCredentials",
			opts:     opts{accessKey: "AKIAEXAMPLE"},
			expected: standby.AuthModeStaticCredentials,
		},
		{
			name:     "nothing set falls back to the credentials chain",
			opts:     opts{},
			expected: standby.AuthModeCredentialsChain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.opts.credentials().AuthMode)
		})
	}
}

func TestCredentials_RegionAndPathStyle(t *testing.T) {
	t.Run("defaults to us-east-1 and virtual-hosted style", func(t *testing.T) {
		c := (&opts{}).credentials()
		assert.Equal(t, defaultAWSRegion, c.Region)
		assert.False(t, c.UsePathStyle)
	})
	t.Run("custom endpoint implies path style", func(t *testing.T) {
		c := (&opts{region: "eu-west-1", endpoint: "http://localhost:9000"}).credentials()
		assert.Equal(t, "eu-west-1", c.Region)
		assert.True(t, c.UsePathStyle)
	})
}

func TestStateFileKey(t *testing.T) {
	assert.Equal(t, "cluster/dr_status_cluster.json", (&opts{key: "cluster/dr_status_cluster.json"}).stateFileKey())
	assert.Empty(t, (&opts{}).stateFileKey())
}

func TestAskModeCredentials_ClearsCrossModeFields(t *testing.T) {
	t.Run("static clears a stale role ARN", func(t *testing.T) {
		o := &ConfigureOpts{opts: opts{
			authMode:  string(standby.AuthModeStaticCredentials),
			accessKey: "AKIAEXAMPLE",
			secretKey: "secret",
		}}
		c := standby.Credentials{
			AuthMode:        standby.AuthModeStaticCredentials,
			BucketName:      "bucket",
			Region:          defaultAWSRegion,
			AccessKeyID:     "AKIAEXAMPLE",
			SecretAccessKey: "secret",
			RoleArn:         "arn:aws:iam::123456789012:role/stale",
		}
		require.NoError(t, o.askModeCredentials(&c))
		assert.Empty(t, c.RoleArn)
		require.NoError(t, c.Validate())
	})
	t.Run("assumeRole clears stale static credentials", func(t *testing.T) {
		o := &ConfigureOpts{opts: opts{
			roleARN:    "arn:aws:iam::123456789012:role/dr",
			awsProfile: "profile",
		}}
		c := standby.Credentials{
			AuthMode:        standby.AuthModeAssumeRole,
			BucketName:      "bucket",
			Region:          defaultAWSRegion,
			RoleArn:         "arn:aws:iam::123456789012:role/dr",
			AccessKeyID:     "AKIASTALE",
			SecretAccessKey: "stale",
		}
		require.NoError(t, o.askModeCredentials(&c))
		assert.Empty(t, c.AccessKeyID)
		assert.Empty(t, c.SecretAccessKey)
		require.NoError(t, c.Validate())
	})
	t.Run("credentialsChain clears all mode-specific fields", func(t *testing.T) {
		o := &ConfigureOpts{opts: opts{awsProfile: "profile"}}
		c := standby.Credentials{
			AuthMode:        standby.AuthModeCredentialsChain,
			BucketName:      "bucket",
			Region:          defaultAWSRegion,
			RoleArn:         "arn:aws:iam::123456789012:role/stale",
			AccessKeyID:     "AKIASTALE",
			SecretAccessKey: "stale",
		}
		require.NoError(t, o.askModeCredentials(&c))
		assert.Empty(t, c.RoleArn)
		assert.Empty(t, c.AccessKeyID)
		assert.Empty(t, c.SecretAccessKey)
		require.NoError(t, c.Validate())
	})
}

func TestConfigureBuilder(t *testing.T) {
	test.CmdValidator(
		t,
		ConfigureBuilder(),
		0,
		[]string{
			flag.S3BucketName, flag.S3Key, flag.S3BucketEndpoint,
			flag.AWSRegion, flag.AWSAuthMode, flag.AWSAccessKey, flag.AWSSecretKey, flag.AWSRoleARN, flag.AWSProfile,
		},
	)
}
