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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalClusterStatus_Fixture(t *testing.T) {
	data, err := os.ReadFile("testdata/cluster_status_v1.bson")
	require.NoError(t, err)

	status, err := UnmarshalClusterStatus(data)
	require.NoError(t, err)
	require.NotNil(t, status)

	assert.Equal(t, uint32(1), status.SchemaVersion)
	assert.Equal(t, Timestamp{T: 1793875780, I: 1}, status.ClusterStatus.EarliestSafeTimestamp)
	assert.Equal(t, Timestamp{T: 1793875770, I: 1}, status.ClusterStatus.MinCommittedOpTimestamp)
	// The older fixture omits sourceDurableTimestamp (the agent omits it while
	// ESTS is uninitialized), so it decodes as absent.
	assert.Nil(t, status.ClusterStatus.SourceDurableTimestamp)
	require.NotNil(t, status.ClusterStatus.EstimatedDataLossSeconds)
	assert.InEpsilon(t, 40.0, *status.ClusterStatus.EstimatedDataLossSeconds, 0.0001)

	require.Len(t, status.ShardStatus, 1)
	shard := status.ShardStatus[0]
	assert.Equal(t, "rs0", shard.ShardID)
	require.NotNil(t, shard.LagSeconds)
	assert.InEpsilon(t, 30.0, *shard.LagSeconds, 0.0001)
	assert.Equal(t, Timestamp{T: 1793875770, I: 1}, shard.CommittedOptime)
}

// The live fixture is a cluster_status_v1 object published by a real relay: it
// carries sourceDurableTimestamp, the last clusterStatus field, which the older
// fixture does not.
func TestUnmarshalClusterStatus_LiveFixture(t *testing.T) {
	data, err := os.ReadFile("testdata/cluster_status_v1_live.bson")
	require.NoError(t, err)

	status, err := UnmarshalClusterStatus(data)
	require.NoError(t, err)
	require.NotNil(t, status)

	assert.Equal(t, uint32(1), status.SchemaVersion)
	require.NotNil(t, status.ClusterStatus.SourceDurableTimestamp)
	assert.Equal(t, Timestamp{T: 1794321486, I: 0}, *status.ClusterStatus.SourceDurableTimestamp)
}

func TestUnmarshalClusterStatus_Malformed(t *testing.T) {
	data, err := os.ReadFile("testdata/cluster_status_v1.bson")
	require.NoError(t, err)

	truncated := data[:len(data)/2]
	status, err := UnmarshalClusterStatus(truncated)
	require.Error(t, err)
	assert.Nil(t, status)
}

func TestUnmarshalClusterStatus_Empty(t *testing.T) {
	status, err := UnmarshalClusterStatus(nil)
	require.ErrorIs(t, err, ErrEmptyClusterStatus)
	assert.Nil(t, status)
}

// setSchemaVersion overwrites the schemaVersion value byte in a copy of the
// fixture. schemaVersion is the first field, encoded as an int32, so locating
// the field name and stepping past its null terminator finds the value.
func setSchemaVersion(t *testing.T, data []byte, version byte) []byte {
	t.Helper()
	key := []byte("schemaVersion\x00")
	i := bytes.Index(data, key)
	require.NotEqual(t, -1, i, "fixture is missing the schemaVersion field")
	out := bytes.Clone(data)
	out[i+len(key)] = version
	return out
}

func TestUnmarshalClusterStatus_UnsupportedSchemaVersion(t *testing.T) {
	data, err := os.ReadFile("testdata/cluster_status_v1.bson")
	require.NoError(t, err)

	// 0 is not valid on the wire and 2 is a future version this build cannot read.
	for _, version := range []byte{0, 2} {
		status, err := UnmarshalClusterStatus(setSchemaVersion(t, data, version))
		require.Error(t, err)
		assert.Nil(t, status)
		assert.EqualError(t, err, fmt.Sprintf(
			"unsupported cluster status schemaVersion %d; this mongocli version reads schemaVersion %d — upgrade mongocli or use a matching relay build",
			version, ClusterStatusSchemaVersion,
		))
	}
}

// Timestamp's JSON shape is a round-trippable contract: UnmarshalJSON is the
// inverse of MarshalJSON, accepting {"seconds":…,"increment":…}.
func TestTimestamp_JSONRoundTrip(t *testing.T) {
	original := Timestamp{T: 1793875780, I: 7}

	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.JSONEq(t, `{"seconds":1793875780,"increment":7}`, string(data))

	var decoded Timestamp
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}

// The whole cluster status document round-trips through JSON: decoding the
// marshaled fixture must reproduce the same Timestamps (and the rest of the
// contract fields).
func TestClusterStatus_JSONRoundTrip(t *testing.T) {
	status := LoadTestClusterStatus(t)

	data, err := json.Marshal(status)
	require.NoError(t, err)

	var decoded ClusterStatus
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, status.SchemaVersion, decoded.SchemaVersion)
	assert.Equal(t, status.ClusterStatus.EarliestSafeTimestamp, decoded.ClusterStatus.EarliestSafeTimestamp)
	assert.Equal(t, status.ClusterStatus.MinCommittedOpTimestamp, decoded.ClusterStatus.MinCommittedOpTimestamp)
	require.NotNil(t, decoded.ClusterStatus.EstimatedDataLossSeconds)
	require.NotNil(t, status.ClusterStatus.EstimatedDataLossSeconds)
	assert.InEpsilon(t, *status.ClusterStatus.EstimatedDataLossSeconds, *decoded.ClusterStatus.EstimatedDataLossSeconds, 0.0001)
	require.Len(t, decoded.ShardStatus, 1)
	assert.Equal(t, status.ShardStatus[0].ShardID, decoded.ShardStatus[0].ShardID)
	assert.Equal(t, status.ShardStatus[0].CommittedOptime, decoded.ShardStatus[0].CommittedOptime)
	// The older fixture omits sourceDurableTimestamp: -o json must not invent it,
	// and it must decode back absent.
	assert.NotContains(t, string(data), "sourceDurableTimestamp")
	assert.Nil(t, decoded.ClusterStatus.SourceDurableTimestamp)
}

// A live document (which carries sourceDurableTimestamp) round-trips through
// JSON with the field preserved, in the contract shape.
func TestClusterStatus_LiveFixtureJSONRoundTrip(t *testing.T) {
	data, err := os.ReadFile("testdata/cluster_status_v1_live.bson")
	require.NoError(t, err)
	status, err := UnmarshalClusterStatus(data)
	require.NoError(t, err)

	out, err := json.Marshal(status)
	require.NoError(t, err)
	assert.Contains(t, string(out), `"sourceDurableTimestamp":{"seconds":1794321486,"increment":0}`)

	var decoded ClusterStatus
	require.NoError(t, json.Unmarshal(out, &decoded))
	require.NotNil(t, decoded.ClusterStatus.SourceDurableTimestamp)
	assert.Equal(t, *status.ClusterStatus.SourceDurableTimestamp, *decoded.ClusterStatus.SourceDurableTimestamp)
}
