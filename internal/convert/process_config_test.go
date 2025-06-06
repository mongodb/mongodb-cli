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

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/pointer"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/ops-manager/opsmngr"
)

func Test_newReplicaSetProcessConfig(t *testing.T) {
	var slaveDelay float64
	fipsMode := true
	omp := &opsmngr.Process{
		Args26: opsmngr.Args26{
			AuditLog: &opsmngr.AuditLog{
				Destination: "file",
				Path:        "/data/audit.log",
				Format:      "JSON",
				Filter:      "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
			},
			NET: opsmngr.Net{
				Port: 27017,
				TLS: &opsmngr.TLS{
					Mode:     "disabled",
					FIPSMode: &fipsMode,
				},
			},
			Replication: &opsmngr.Replication{
				ReplSetName: "myReplicaSet",
				OplogSizeMB: pointer.Get(10),
			},
			Storage: &opsmngr.Storage{
				DBPath:         "/data/db",
				DirectoryPerDB: pointer.Get(true),
				WiredTiger: &map[string]any{
					"collectionConfig": map[string]any{},
					"engineConfig": map[string]any{
						"cacheSizeGB": 1,
					},
					"indexConfig": map[string]any{},
				},
			},
			SystemLog: opsmngr.SystemLog{
				Destination: "file",
				Path:        "/data/log/mongodb.log",
			},
			SetParameter: &map[string]any{
				"param": "param",
			},
		},
		AuthSchemaVersion:           5,
		Disabled:                    false,
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogRotate: &opsmngr.LogRotate{
			SizeThresholdMB:  1000,
			TimeThresholdHrs: 24,
		},
		ManualMode:  false,
		Name:        "myReplicaSet_1",
		ProcessType: "mongod",
		Version:     "4.4.1-ent",
	}
	omm := &opsmngr.Member{
		ID:                 0,
		ArbiterOnly:        false,
		BuildIndexes:       true,
		Hidden:             false,
		Host:               "myReplicaSet_1",
		Priority:           1,
		SlaveDelay:         &slaveDelay,
		SecondaryDelaySecs: &slaveDelay,
		Votes:              1,
	}

	expected := &ProcessConfig{
		AuditLogPath:                "/data/audit.log",
		AuditLogDestination:         "file",
		AuditLogFormat:              "JSON",
		AuditLogFilter:              "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
		BuildIndexes:                pointer.Get(true),
		DBPath:                      "/data/db",
		DirectoryPerDB:              pointer.Get(true),
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogDestination:              "file",
		LogPath:                     "/data/log/mongodb.log",
		Name:                        "myReplicaSet_1",
		OplogSizeMB:                 pointer.Get(10),
		Port:                        27017,
		Priority:                    pointer.Get[float64](1),
		ProcessType:                 "mongod",
		SlaveDelay:                  pointer.Get[float64](0),
		SecondaryDelaySecs:          pointer.Get[float64](0),
		Version:                     "4.4.1-ent",
		Votes:                       pointer.Get[float64](1),
		ArbiterOnly:                 pointer.Get(false),
		Disabled:                    false,
		Hidden:                      pointer.Get(false),
		TLS: &TLS{
			Mode:     "disabled",
			FIPSMode: &fipsMode,
		},
		SetParameter: &map[string]any{
			"param": "param",
		},
		WiredTiger: &map[string]any{
			"collectionConfig": map[string]any{},
			"engineConfig": map[string]any{
				"cacheSizeGB": 1,
			},
			"indexConfig": map[string]any{},
		},
	}
	result := newReplicaSetProcessConfig(omm, omp)

	assert.Equal(t, expected, result)
}

func Test_newConfigRSProcess(t *testing.T) {
	p := &ProcessConfig{
		AuditLogPath:                "/data/audit.log",
		AuditLogDestination:         "file",
		AuditLogFormat:              "JSON",
		AuditLogFilter:              "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
		BuildIndexes:                pointer.Get(true),
		DBPath:                      "/data/db",
		DirectoryPerDB:              pointer.Get(true),
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogDestination:              "file",
		LogPath:                     "/data/log/mongodb.log",
		Name:                        "myReplicaSet_1",
		Port:                        27017,
		Priority:                    pointer.Get[float64](1),
		ProcessType:                 "mongod",
		SlaveDelay:                  pointer.Get[float64](0),
		SecondaryDelaySecs:          pointer.Get[float64](0),
		Version:                     "4.4.1-ent",
		Votes:                       pointer.Get[float64](1),
		ArbiterOnly:                 pointer.Get(false),
		Disabled:                    false,
		Hidden:                      pointer.Get(false),
		TLS:                         &TLS{Mode: "disabled"},
		WiredTiger: &map[string]any{
			"collectionConfig": map[string]any{},
			"engineConfig": map[string]any{
				"cacheSizeGB": 1,
			},
			"indexConfig": map[string]any{},
		},
		SetParameter: &map[string]any{
			"param": "param",
		},
	}

	want := &opsmngr.Process{
		Args26: opsmngr.Args26{
			AuditLog: &opsmngr.AuditLog{
				Destination: "file",
				Path:        "/data/audit.log",
				Format:      "JSON",
				Filter:      "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
			},
			NET: opsmngr.Net{
				Port: 27017,
				TLS:  &opsmngr.TLS{Mode: "disabled"},
			},
			Replication: &opsmngr.Replication{
				ReplSetName: "myReplicaSet",
			},
			Storage: &opsmngr.Storage{
				DBPath:         "/data/db",
				DirectoryPerDB: pointer.Get(true),
				WiredTiger: &map[string]any{
					"collectionConfig": map[string]any{},
					"engineConfig": map[string]any{
						"cacheSizeGB": 1,
					},
					"indexConfig": map[string]any{},
				},
			},
			SystemLog: opsmngr.SystemLog{
				Destination: "file",
				Path:        "/data/log/mongodb.log",
			},
			Sharding: &opsmngr.Sharding{ClusterRole: "configsvr"},
			SetParameter: &map[string]any{
				"param": "param",
			},
		},
		AuthSchemaVersion:           5,
		Disabled:                    false,
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogRotate: &opsmngr.LogRotate{
			SizeThresholdMB:  1000,
			TimeThresholdHrs: 24,
		},
		ManualMode:  false,
		Name:        "myReplicaSet_1",
		ProcessType: "mongod",
		Version:     "4.4.1-ent",
	}
	got := newConfigRSProcess(p, "myReplicaSet")
	assert.Equal(t, want, got)
}

func Test_newConfigRSProcess_audit(t *testing.T) {
	p := &ProcessConfig{
		AuditLogDestination:         "sysfile",
		BuildIndexes:                pointer.Get(true),
		DBPath:                      "/data/db",
		DirectoryPerDB:              pointer.Get(true),
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogDestination:              "file",
		LogPath:                     "/data/log/mongodb.log",
		Name:                        "myReplicaSet_1",
		Port:                        27017,
		Priority:                    pointer.Get[float64](1),
		ProcessType:                 "mongod",
		SlaveDelay:                  pointer.Get[float64](0),
		Version:                     "4.4.1-ent",
		Votes:                       pointer.Get[float64](1),
		ArbiterOnly:                 pointer.Get(false),
		Disabled:                    false,
		Hidden:                      pointer.Get(false),
		TLS:                         &TLS{Mode: "disabled"},
		WiredTiger: &map[string]any{
			"collectionConfig": map[string]any{},
			"engineConfig": map[string]any{
				"cacheSizeGB": 1,
			},
			"indexConfig": map[string]any{},
		},
		SetParameter: &map[string]any{
			"param": "param",
		},
	}

	want := &opsmngr.Process{
		Args26: opsmngr.Args26{
			AuditLog: &opsmngr.AuditLog{
				Destination: "sysfile",
			},
			NET: opsmngr.Net{
				Port: 27017,
				TLS:  &opsmngr.TLS{Mode: "disabled"},
			},
			Replication: &opsmngr.Replication{
				ReplSetName: "myReplicaSet",
			},
			Storage: &opsmngr.Storage{
				DBPath:         "/data/db",
				DirectoryPerDB: pointer.Get(true),
				WiredTiger: &map[string]any{
					"collectionConfig": map[string]any{},
					"engineConfig": map[string]any{
						"cacheSizeGB": 1,
					},
					"indexConfig": map[string]any{},
				},
			},
			SystemLog: opsmngr.SystemLog{
				Destination: "file",
				Path:        "/data/log/mongodb.log",
			},
			Sharding: &opsmngr.Sharding{ClusterRole: "configsvr"},
			SetParameter: &map[string]any{
				"param": "param",
			},
		},
		AuthSchemaVersion:           5,
		Disabled:                    false,
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogRotate: &opsmngr.LogRotate{
			SizeThresholdMB:  1000,
			TimeThresholdHrs: 24,
		},
		ManualMode:  false,
		Name:        "myReplicaSet_1",
		ProcessType: "mongod",
		Version:     "4.4.1-ent",
	}
	got := newConfigRSProcess(p, "myReplicaSet")
	assert.Equal(t, want, got)
}

func Test_newReplicaSetProcess(t *testing.T) {
	p := &ProcessConfig{
		AuditLogPath:                "/data/audit.log",
		AuditLogDestination:         "file",
		AuditLogFormat:              "JSON",
		AuditLogFilter:              "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
		BuildIndexes:                pointer.Get(true),
		DBPath:                      "/data/db",
		DirectoryPerDB:              pointer.Get(true),
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogDestination:              "file",
		LogPath:                     "/data/log/mongodb.log",
		Name:                        "myReplicaSet_1",
		Port:                        27017,
		Priority:                    pointer.Get[float64](1),
		ProcessType:                 "mongod",
		SlaveDelay:                  pointer.Get[float64](0),
		Version:                     "4.4.1-ent",
		Votes:                       pointer.Get[float64](1),
		ArbiterOnly:                 pointer.Get(false),
		Disabled:                    false,
		Hidden:                      pointer.Get(false),
		TLS:                         &TLS{Mode: "disabled"},
		WiredTiger: &map[string]any{
			"collectionConfig": map[string]any{},
			"engineConfig": map[string]any{
				"cacheSizeGB": 1,
			},
			"indexConfig": map[string]any{},
		},
		SetParameter: &map[string]any{
			"enableLocalhostAuthBypass": "false",
			"auditAuthorizationSuccess": "true",
		},
	}

	want := &opsmngr.Process{
		Args26: opsmngr.Args26{
			AuditLog: &opsmngr.AuditLog{
				Destination: "file",
				Path:        "/data/audit.log",
				Format:      "JSON",
				Filter:      "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
			},
			NET: opsmngr.Net{
				Port: 27017,
				TLS:  &opsmngr.TLS{Mode: "disabled"},
			},
			Replication: &opsmngr.Replication{
				ReplSetName: "myReplicaSet",
			},
			Storage: &opsmngr.Storage{
				DBPath:         "/data/db",
				DirectoryPerDB: pointer.Get(true),
				WiredTiger: &map[string]any{
					"collectionConfig": map[string]any{},
					"engineConfig": map[string]any{
						"cacheSizeGB": 1,
					},
					"indexConfig": map[string]any{},
				},
			},
			SystemLog: opsmngr.SystemLog{
				Destination: "file",
				Path:        "/data/log/mongodb.log",
			},
			SetParameter: &map[string]any{
				"enableLocalhostAuthBypass": "false",
				"auditAuthorizationSuccess": "true",
			},
		},
		AuthSchemaVersion:           5,
		Disabled:                    false,
		FeatureCompatibilityVersion: "4.4",
		Hostname:                    "n1.omansible.int",
		LogRotate: &opsmngr.LogRotate{
			SizeThresholdMB:  1000,
			TimeThresholdHrs: 24,
		},
		ManualMode:  false,
		Name:        "myReplicaSet_1",
		ProcessType: "mongod",
		Version:     "4.4.1-ent",
	}
	got := newReplicaSetProcess(p, "myReplicaSet")
	assert.Equal(t, want, got)
}

func Test_newMongosProcessConfig(t *testing.T) {
	p := &opsmngr.Process{
		Args26: opsmngr.Args26{
			NET: opsmngr.Net{
				Port: 27017,
			},
			SystemLog: opsmngr.SystemLog{
				Destination: "file",
				Path:        "/data/mongos/mongodb.log",
			},
			SetParameter: &map[string]any{
				"param": "param",
			},
			AuditLog: &opsmngr.AuditLog{
				Path:        "/data/mongos/audit.log",
				Destination: "file",
				Format:      "JSON",
				Filter:      "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
			},
		},
		AuthSchemaVersion:           5,
		Cluster:                     "myCluster",
		Disabled:                    false,
		FeatureCompatibilityVersion: "3.6",
		Hostname:                    "n1.omansible.int",
		LogRotate: &opsmngr.LogRotate{
			SizeThresholdMB:  1000,
			TimeThresholdHrs: 24,
		},
		ManualMode:  false,
		Name:        "myCluster_mongos_4",
		ProcessType: "mongos",
		Version:     "3.6.21-ent",
	}

	want := &ProcessConfig{
		AuditLogPath:                "/data/mongos/audit.log",
		AuditLogDestination:         "file",
		AuditLogFormat:              "JSON",
		AuditLogFilter:              "{ atype: { $in: [ \"createCollection\", \"dropCollection\" ] } }",
		FeatureCompatibilityVersion: "3.6",
		Hostname:                    "n1.omansible.int",
		LogDestination:              "file",
		LogPath:                     "/data/mongos/mongodb.log",
		Name:                        "myCluster_mongos_4",
		Port:                        27017,
		ProcessType:                 "mongos",
		Version:                     "3.6.21-ent",
		Disabled:                    false,
		SetParameter: &map[string]any{
			"param": "param",
		},
	}
	got := newMongosProcessConfig(p)
	assert.Equal(t, want, got)
}
