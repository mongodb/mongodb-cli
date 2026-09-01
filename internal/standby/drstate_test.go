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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFailover_FromStandby(t *testing.T) {
	before := time.Now().UTC()
	s := &RemoteDRState{
		State:         StateStandby,
		ClusterName:   "rs1",
		Version:       "1",
		LastModified:  "2026-08-19T12:25:23Z",
		SchemaVersion: "1",
	}

	require.NoError(t, s.ApplyFailover())
	assert.Equal(t, StatePromoteStandby, s.State)
	assert.Equal(t, StateStandby, s.PreviousState)
	assert.Equal(t, "2", s.Version)
	assert.Equal(t, "rs1", s.ClusterName, "clusterName must be preserved")
	assert.Equal(t, "1", s.SchemaVersion, "schemaVersion must be preserved")
	parsed, err := time.Parse(time.RFC3339, s.LastModified)
	require.NoError(t, err, "lastModified must remain RFC3339")
	assert.False(t, parsed.Before(before.Truncate(time.Second)), "lastModified must advance to now")
}

func TestApplyFailover_RejectsOtherStates(t *testing.T) {
	for _, state := range []string{"Active", "PromoteStandby", "DemoteToStandby", "StandbyReadyToPromote", "RevertToActive", ""} {
		t.Run(state, func(t *testing.T) {
			s := &RemoteDRState{State: state, Version: "3"}
			require.Error(t, s.ApplyFailover())
			assert.Equal(t, state, s.State, "state must be untouched on error")
			assert.Equal(t, "3", s.Version, "version must be untouched on error")
		})
	}
}

func TestApplyFailover_ClearsPlanned(t *testing.T) {
	s := &RemoteDRState{
		State:           StateStandby,
		ClusterName:     "rs1",
		Version:         "1",
		SchemaVersion:   "1",
		Planned:         true,
		SyncDestination: "B",
	}

	require.NoError(t, s.ApplyFailover())
	assert.False(t, s.Planned)
	assert.Equal(t, "B", s.SyncDestination, "other planned fields must be preserved")
}

func TestIncrementVersion(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", "1"},
		{"1", "2"},
		{"9", "10"},
		{"not-a-number", "1"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, incrementVersion(tt.in), "input %q", tt.in)
	}
}

func TestRemoteDRState_MarshalRoundTrip(t *testing.T) {
	raw := `{"state":"PromoteStandby","previousState":"Standby","clusterName":"rs1","version":"2","lastModified":"2026-08-19T12:25:23Z","schemaVersion":"1","planned":true,"syncDestination":"B","plannedFailoverTo":"A"}`
	s, err := UnmarshalRemoteDRState([]byte(raw))
	require.NoError(t, err)

	out, err := s.Marshal()
	require.NoError(t, err)

	again, err := UnmarshalRemoteDRState(out)
	require.NoError(t, err)
	assert.Equal(t, *s, *again, "marshal must preserve every field the agent reads")
}
