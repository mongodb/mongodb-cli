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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/mocks"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lagTestOpts(store standby.StateStore, out *bytes.Buffer) *LagOpts {
	return &LagOpts{
		opts:       opts{store: store},
		OutputOpts: cli.OutputOpts{OutWriter: out},
	}
}

func clusterStatusWithLag(seconds *float64) *standby.ClusterStatus {
	status := &standby.ClusterStatus{}
	status.ClusterStatus.EstimatedDataLossSeconds = seconds
	// Mirror the shared cluster_status_v1 fixture's timestamps exactly so the
	// derived standby lag (estimated data loss + ESTS - MCTS) is exercised with
	// the same values: 1793875780 - 1793875770 = 10s.
	status.ClusterStatus.EarliestSafeTimestamp = standby.Timestamp{T: 1793875780, I: 1}
	status.ClusterStatus.MinCommittedOpTimestamp = standby.Timestamp{T: 1793875770, I: 1}
	return status
}

func clusterStatusWithShards(seconds *float64, shards ...standby.ShardStatus) *standby.ClusterStatus {
	status := clusterStatusWithLag(seconds)
	status.ShardStatus = shards
	return status
}

func float64Ptr(v float64) *float64 { return &v }

// setGlobalOutput sets the profile-level 'output' config setting (what
// 'mongocli config set output <value>' writes and config.Output() reads),
// restoring the previous value when the test finishes.
func setGlobalOutput(t *testing.T, value string) {
	t.Helper()
	prev := config.Output()
	config.SetOutput(value)
	t.Cleanup(func() { config.SetOutput(prev) })
}

func TestLag_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(float64Ptr(40)), nil).Times(1)

	var out bytes.Buffer
	require.NoError(t, lagTestOpts(mockStore, &out).Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40s\nstandby lag: 50s\n", out.String())
}

func TestLag_RunFractional(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(float64Ptr(40.5)), nil).Times(1)

	var out bytes.Buffer
	require.NoError(t, lagTestOpts(mockStore, &out).Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40.5s\nstandby lag: 50.5s\n", out.String())
}

// Scenario C (uninitialized), without --shards: two lines, the reason in the
// bracket of line 1.
func TestLag_RunUninitialized(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(nil), nil).Times(1)

	var out bytes.Buffer
	require.NoError(t, lagTestOpts(mockStore, &out).Run(context.Background()))
	assert.Equal(t, "estimated data loss: NA (standby has not yet initialized replication)\nstandby lag: NA\n", out.String())
}

// The standby-lag total line is guarded: when the ESTS/MCTS pair cannot yield a
// total it prints "not available yet" instead of a wrapped number, so the line's
// shape stays predictable.
func TestLag_RunStandbyLagNotAvailable(t *testing.T) {
	tests := []struct {
		name string
		ests standby.Timestamp
		mcts standby.Timestamp
	}{
		{
			name: "MCTS zero (partial document)",
			ests: standby.Timestamp{T: 1793875780, I: 1},
			mcts: standby.Timestamp{T: 0, I: 0},
		},
		{
			name: "ESTS zero",
			ests: standby.Timestamp{T: 0, I: 0},
			mcts: standby.Timestamp{T: 1793875770, I: 1},
		},
		{
			name: "MCTS ahead of ESTS",
			ests: standby.Timestamp{T: 1793875770, I: 0},
			mcts: standby.Timestamp{T: 1793875780, I: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStore := mocks.NewMockStateStore(ctrl)

			status := clusterStatusWithLag(float64Ptr(40))
			status.ClusterStatus.EarliestSafeTimestamp = tt.ests
			status.ClusterStatus.MinCommittedOpTimestamp = tt.mcts
			mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(status, nil).Times(1)

			var out bytes.Buffer
			require.NoError(t, lagTestOpts(mockStore, &out).Run(context.Background()))
			assert.Equal(t, "estimated data loss: 40s\nstandby lag: not available yet\n", out.String())
		})
	}
}

func TestLag_RunFetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(nil, errors.New("boom")).Times(1)

	err := lagTestOpts(mockStore, &bytes.Buffer{}).Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

// Scenario A with --shards: the healthy three-line output. Shard values are the
// published lagSeconds field (relayed, never recomputed). The shards are
// deliberately out of shard-ID order: the CLI relays the order the object was
// decoded in rather than sorting, since the upstream projection already sorts.
func TestLag_RunShardsPreservesDecodedOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	status := clusterStatusWithShards(float64Ptr(40),
		standby.ShardStatus{ShardID: "shard1", LagSeconds: float64Ptr(12.5)},
		standby.ShardStatus{ShardID: "rs0", LagSeconds: float64Ptr(40)},
	)
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(status, nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.shards = true
	require.NoError(t, o.Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40s\nshard1: 12.5s of oplog pending apply\nrs0: 40s of oplog pending apply\nstandby lag: 50s\n", out.String())
}

// Scenario B with --shards: line 1 is the number, the shard lines still relay
// the published lagSeconds (a published value shows, an omitted one reads NA),
// and the total reads "not available yet".
func TestLag_RunShardsScenarioB(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	status := clusterStatusWithShards(float64Ptr(40),
		standby.ShardStatus{ShardID: "rs0", LagSeconds: float64Ptr(12.5)}, // published value
		standby.ShardStatus{ShardID: "shard1"},                            // lagSeconds omitted -> NA
	)
	status.ClusterStatus.MinCommittedOpTimestamp = standby.Timestamp{T: 0, I: 0} // MCTS unset -> scenario B
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(status, nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.shards = true
	require.NoError(t, o.Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40s\nrs0: 12.5s of oplog pending apply\nshard1: NA\nstandby lag: not available yet\n", out.String())
}

// Scenario C with --shards: three lines; the shard line names the one component
// encoded in the configured DR state file key and reads NA.
func TestLag_RunShardsScenarioC(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(nil), nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.shards = true
	o.key = "cluster-A/dr_status_rs0.json"
	require.NoError(t, o.Run(context.Background()))
	assert.Equal(t, "estimated data loss: NA (standby has not yet initialized replication)\nrs0: NA\nstandby lag: NA\n", out.String())
}

func TestComponentIDFromKey(t *testing.T) {
	tests := []struct{ key, want string }{
		{"cluster-A/dr_status_rs0.json", "rs0"},
		{"org/proj/cluster-B/dr_status_shard1.json", "shard1"},
		{"dr_status_rs0.json", "rs0"},
		{"cluster-A/other.json", "other"},
		{"cluster-A/dr_status_.json", "unknown"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, componentIDFromKey(tt.key), "key %q", tt.key)
	}
}

func TestLag_RunShardsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	status := clusterStatusWithLag(float64Ptr(40))
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(status, nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.shards = true
	require.NoError(t, o.Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40s\nstandby lag: 50s\n", out.String())
}

func TestLag_RunShardsFlagUnset(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	// Shard data present, but --shards not set: only the two cluster lines print.
	status := clusterStatusWithShards(float64Ptr(40),
		standby.ShardStatus{ShardID: "rs0", LagSeconds: float64Ptr(1)},
		standby.ShardStatus{ShardID: "shard1", LagSeconds: float64Ptr(2)},
	)
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(status, nil).Times(1)

	var out bytes.Buffer
	require.NoError(t, lagTestOpts(mockStore, &out).Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40s\nstandby lag: 50s\n", out.String())
}

func TestLag_RunJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(standby.LoadTestClusterStatus(t), nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.Output = "json"
	require.NoError(t, o.Run(context.Background()))

	var got struct {
		SchemaVersion  uint32 `json:"schemaVersion"`
		GenerationTime string `json:"generationTime"`
		ClusterStatus  struct {
			EstimatedDataLossSeconds *float64 `json:"estimatedDataLossSeconds"`
			ObservationTime          string   `json:"observationTime"`
			EarliestSafeTimestamp    struct {
				Seconds   uint32 `json:"seconds"`
				Increment uint32 `json:"increment"`
			} `json:"earliestSafeTimestamp"`
			MinCommittedOpTimestamp struct {
				Seconds   uint32 `json:"seconds"`
				Increment uint32 `json:"increment"`
			} `json:"minCommittedOpTimestamp"`
		} `json:"clusterStatus"`
		ShardStatus []struct {
			ShardID    string   `json:"shardId"`
			LagSeconds *float64 `json:"lagSeconds"`
		} `json:"shardStatus"`
	}
	require.NoError(t, json.Unmarshal(out.Bytes(), &got), "output was: %s", out.String())

	assert.Equal(t, uint32(1), got.SchemaVersion)
	require.NotNil(t, got.ClusterStatus.EstimatedDataLossSeconds)
	assert.Equal(t, 40.0, *got.ClusterStatus.EstimatedDataLossSeconds)
	assert.Equal(t, uint32(1793875780), got.ClusterStatus.EarliestSafeTimestamp.Seconds)
	assert.Equal(t, uint32(1), got.ClusterStatus.EarliestSafeTimestamp.Increment)
	assert.Equal(t, uint32(1793875770), got.ClusterStatus.MinCommittedOpTimestamp.Seconds)
	require.Len(t, got.ShardStatus, 1)
	assert.Equal(t, "rs0", got.ShardStatus[0].ShardID)
	require.NotNil(t, got.ShardStatus[0].LagSeconds)
	assert.Equal(t, 30.0, *got.ShardStatus[0].LagSeconds)

	// The standby lag total is a plaintext-only line; the JSON document is the
	// wire contract and must not carry it.
	assert.NotContains(t, out.String(), "standby lag")
}

// -o json already carries shardStatus, so --shards must not change the output.
func TestLag_RunJSONShardsFlagNoEffect(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(standby.LoadTestClusterStatus(t), nil).Times(2)

	var withShards, withoutShards bytes.Buffer
	o := lagTestOpts(mockStore, &withShards)
	o.Output = "json"
	o.shards = true
	require.NoError(t, o.Run(context.Background()))

	o2 := lagTestOpts(mockStore, &withoutShards)
	o2.Output = "json"
	require.NoError(t, o2.Run(context.Background()))

	assert.Equal(t, withoutShards.String(), withShards.String())
}

// -o json prints the full document even when the cluster lag is uninitialized.
func TestLag_RunJSONUninitialized(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(nil), nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.Output = "json"
	require.NoError(t, o.Run(context.Background()))

	var got map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &got))
	clusterStatus, ok := got["clusterStatus"].(map[string]any)
	require.True(t, ok)
	assert.NotContains(t, clusterStatus, "estimatedDataLossSeconds")
}

func TestLagBuilder(t *testing.T) {
	cmd := LagBuilder()
	require.Equal(t, "lag", cmd.Use)
	test.CmdValidator(t, cmd, 0, []string{
		flag.S3BucketName, flag.S3Key, flag.S3BucketEndpoint,
		flag.AWSRegion, flag.AWSAuthMode, flag.AWSAccessKey, flag.AWSSecretKey, flag.AWSRoleARN, flag.AWSProfile,
		flag.Shards, flag.Output,
	})
}

// lag only supports -o json (and the plaintext default); other -o formats must
// fail explicitly rather than silently falling back to plaintext.
func TestLag_RunUnsupportedOutputFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)
	// Validation happens before any store access.
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Times(0)

	for _, format := range []string{"json-path", "go-template", "go-template-file"} {
		var out bytes.Buffer
		o := lagTestOpts(mockStore, &out)
		o.Output = format
		err := o.Run(context.Background())
		require.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("unsupported output format %q for this command; supported values: json", format))
		assert.Empty(t, out.String())
	}
}

func TestLagBuilder_OutputCompletion(t *testing.T) {
	values, directive := completeLagOutput(nil, nil, "")
	assert.Equal(t, []string{"json"}, values)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

// When -o is unset, lag falls back to the global 'output' config setting, the
// same way the other OutputOpts commands do.
func TestLag_RunGlobalOutputConfigJSON(t *testing.T) {
	setGlobalOutput(t, "json")

	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(standby.LoadTestClusterStatus(t), nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	require.NoError(t, o.Run(context.Background()))

	var got map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &got), "output was: %s", out.String())
	_, ok := got["clusterStatus"]
	assert.True(t, ok, "expected a JSON cluster status document, got: %s", out.String())
}

// The global 'output' config plus --shards: JSON wins (shardStatus is already
// part of the document), so the flag must be a no-op, not add plaintext lines.
func TestLag_RunGlobalOutputConfigJSONWithShards(t *testing.T) {
	setGlobalOutput(t, "json")

	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(standby.LoadTestClusterStatus(t), nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.shards = true
	require.NoError(t, o.Run(context.Background()))

	var got map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &got), "output was: %s", out.String())
	_, ok := got["clusterStatus"]
	assert.True(t, ok, "expected a JSON cluster status document, got: %s", out.String())
	assert.NotContains(t, out.String(), "estimated data loss", "--shards must not mix plaintext lines into JSON output")
}

// A global output setting that this command cannot render must produce the same
// clear unsupported-format error as an explicit -o, not silent plaintext.
func TestLag_RunGlobalOutputConfigUnsupported(t *testing.T) {
	for _, format := range []string{"json-path", "go-template", "go-template-file"} {
		t.Run(format, func(t *testing.T) {
			setGlobalOutput(t, format)

			ctrl := gomock.NewController(t)
			mockStore := mocks.NewMockStateStore(ctrl)
			// Validation happens before any store access.
			mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Times(0)

			var out bytes.Buffer
			err := lagTestOpts(mockStore, &out).Run(context.Background())
			require.Error(t, err)
			assert.EqualError(t, err, fmt.Sprintf("unsupported output format %q for this command; supported values: json", format))
			assert.Empty(t, out.String())
		})
	}
}

// With neither -o nor a global output setting, behavior stays byte-identical:
// the labeled plaintext output.
func TestLag_RunNoGlobalOutputConfigPlaintext(t *testing.T) {
	setGlobalOutput(t, "")

	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(float64Ptr(40)), nil).Times(1)

	var out bytes.Buffer
	require.NoError(t, lagTestOpts(mockStore, &out).Run(context.Background()))
	assert.Equal(t, "estimated data loss: 40s\nstandby lag: 50s\n", out.String())
}

// An explicit -o wins over the global output config.
func TestLag_RunFlagOutputWinsOverGlobalConfig(t *testing.T) {
	// The global setting alone would be rejected as unsupported.
	setGlobalOutput(t, "json-path")

	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)
	mockStore.EXPECT().FetchClusterStatus(gomock.Any()).Return(clusterStatusWithLag(float64Ptr(40)), nil).Times(1)

	var out bytes.Buffer
	o := lagTestOpts(mockStore, &out)
	o.Output = "json"
	require.NoError(t, o.Run(context.Background()))

	var got map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &got), "output was: %s", out.String())
	_, ok := got["clusterStatus"]
	assert.True(t, ok, "expected a JSON cluster status document, got: %s", out.String())
}
