// Copyright 2020 MongoDB Inc
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

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/ops-manager/opsmngr"
)

func TestRSConfig_SettingsRoundtrip(t *testing.T) {
	t.Parallel()
	settings := &map[string]any{
		"chainingAllowed":            true,
		"heartbeatTimeoutSecs":       float64(10),
		"electionTimeoutMillis":      float64(10000),
		"catchUpTimeoutMillis":       float64(-1),
		"catchUpTakeoverDelayMillis": float64(30000),
		"getLastErrorDefaults":       map[string]any{"w": float64(1), "wtimeout": float64(0)},
	}
	in := &opsmngr.AutomationConfig{
		Processes: []*opsmngr.Process{
			{
				Name:                        "myRS_1",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				FeatureCompatibilityVersion: "4.2",
				Hostname:                    "host0",
				Args26: opsmngr.Args26{
					NET:     opsmngr.Net{Port: 27000},
					Storage: &opsmngr.Storage{DBPath: "/data"},
					SystemLog: opsmngr.SystemLog{
						Destination: "file",
						Path:        "/data/mongodb.log",
					},
					Replication: &opsmngr.Replication{ReplSetName: "myRS"},
				},
			},
		},
		ReplicaSets: []*opsmngr.ReplicaSet{
			{
				ID:              "myRS",
				ProtocolVersion: "1",
				Members: []opsmngr.Member{
					{
						ID:           0,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Host:         "myRS_1",
						Priority:     1,
						Votes:        1,
					},
				},
				Settings:                           settings,
				WriteConcernMajorityJournalDefault: "true",
			},
		},
	}

	// Describe direction: automation config → RSConfig
	rsConfig := newRSConfig(in, "myRS")
	assert.NotNil(t, rsConfig)
	assert.Equal(t, settings, rsConfig.Settings)
	assert.Equal(t, "true", rsConfig.WriteConcernMajorityJournalDefault)

	// Update direction: RSConfig → opsmngr.ReplicaSet
	rs, err := newReplicaSet(rsConfig)
	require.NoError(t, err)
	assert.Equal(t, settings, rs.Settings)
	assert.Equal(t, "true", rs.WriteConcernMajorityJournalDefault)
}

func TestRSConfig_protocolVer(t *testing.T) {
	testCases := map[string]struct {
		config                *RSConfig
		wantedProtocolVersion string
		wantErr               bool
	}{
		"empty fcv": {
			config:                &RSConfig{},
			wantedProtocolVersion: "",
			wantErr:               true,
		},
		"post 4.0": {
			config:                &RSConfig{FeatureCompatibilityVersion: "4.0"},
			wantedProtocolVersion: "1",
			wantErr:               false,
		},
		"pre 4.0": {
			config:                &RSConfig{FeatureCompatibilityVersion: "3.6"},
			wantedProtocolVersion: "0",
			wantErr:               false,
		},
		"empty at parent but with FC in process": {
			config: &RSConfig{
				Processes: []*ProcessConfig{
					{
						FeatureCompatibilityVersion: "4.0",
					},
				},
			},
			wantedProtocolVersion: "1",
			wantErr:               false,
		},
	}
	for name, tc := range testCases {
		m := tc.config
		expected := tc.wantedProtocolVersion
		wantErr := tc.wantErr
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ver, err := m.protocolVer()
			if (err != nil) != wantErr {
				t.Fatalf("protocolVer() unexpected error: %v\n", err)
			}
			if ver != expected {
				t.Errorf("protocolVer() expected: %s but got: %s", expected, ver)
			}
		})
	}
}
