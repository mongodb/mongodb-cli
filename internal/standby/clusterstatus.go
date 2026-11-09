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
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrEmptyClusterStatus is returned when there is no cluster status data to decode.
var ErrEmptyClusterStatus = errors.New("cluster status data is empty")

// ClusterStatusSchemaVersion is the only schemaVersion this mongocli release
// knows how to decode. It is part of the cluster_status_v1 wire contract and
// must be kept in sync with the agent that writes the document.
const ClusterStatusSchemaVersion = 1

// ClusterStatus mirrors the cluster_status_v1 document written to blob storage
// by the agent. It is the single source of truth for the contract: the BSON
// field names drive decoding and the JSON tags drive '-o json' rendering, so a
// contract change is made in exactly one place. The BSON field names are an
// external wire contract and must not be changed. SchemaVersion must equal
// ClusterStatusSchemaVersion for the document to be decodable.
type ClusterStatus struct {
	SchemaVersion  uint32             `bson:"schemaVersion" json:"schemaVersion"`
	GenerationTime primitive.DateTime `bson:"generationTime" json:"generationTime"`
	ClusterStatus  struct {
		EstimatedDataLossSeconds *float64           `bson:"estimatedDataLossSeconds,omitempty" json:"estimatedDataLossSeconds,omitempty"`
		ObservationTime          primitive.DateTime `bson:"observationTime" json:"observationTime"`
		EarliestSafeTimestamp    Timestamp          `bson:"earliestSafeTimestamp" json:"earliestSafeTimestamp"`
		MinCommittedOpTimestamp  Timestamp          `bson:"minCommittedOpTimestamp" json:"minCommittedOpTimestamp"`
	} `bson:"clusterStatus" json:"clusterStatus"`
	ShardStatus []ShardStatus `bson:"shardStatus" json:"shardStatus"`
}

// ShardStatus describes the replication lag of a single shard within a standby
// cluster.
type ShardStatus struct {
	ShardID         string              `bson:"shardId" json:"shardId"`
	ObservationTime *primitive.DateTime `bson:"observationTime,omitempty" json:"observationTime,omitempty"`
	CommittedOptime Timestamp           `bson:"committedOptime" json:"committedOptime"`
	LagSeconds      *float64            `bson:"lagSeconds,omitempty" json:"lagSeconds,omitempty"`
}

// Timestamp is a BSON timestamp that also carries the cluster_status_v1 JSON
// wire shape. It decodes exactly like primitive.Timestamp, but
// primitive.Timestamp has no MarshalJSON and would otherwise render as
// {"T":…,"I":…}; Timestamp renders the contract shape
// {"seconds":…,"increment":…} instead.
type Timestamp primitive.Timestamp

// MarshalJSON renders the timestamp in the cluster_status_v1 contract shape.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Seconds   uint32 `json:"seconds"`
		Increment uint32 `json:"increment"`
	}{Seconds: t.T, Increment: t.I})
}

// UnmarshalJSON decodes the cluster_status_v1 contract shape
// {"seconds":…,"increment":…}, mirroring MarshalJSON so the Timestamp round-trips
// through JSON. It is the inverse of MarshalJSON and accepts exactly that object.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var v struct {
		Seconds   uint32 `json:"seconds"`
		Increment uint32 `json:"increment"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = Timestamp{T: v.Seconds, I: v.Increment}
	return nil
}

// UnmarshalBSONValue decodes a BSON timestamp exactly as primitive.Timestamp
// would, keeping the wire-decode behavior unchanged.
func (t *Timestamp) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	var ts primitive.Timestamp
	if err := bson.UnmarshalValue(bt, data, &ts); err != nil {
		return err
	}
	*t = Timestamp(ts)
	return nil
}

// UnmarshalClusterStatus decodes a cluster_status_v1 BSON document. It rejects
// any document whose schemaVersion is not ClusterStatusSchemaVersion, since a
// newer wire contract may not be readable by this build.
func UnmarshalClusterStatus(data []byte) (*ClusterStatus, error) {
	if len(data) == 0 {
		return nil, ErrEmptyClusterStatus
	}

	var s ClusterStatus
	if err := bson.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error unmarshaling cluster status: %w", err)
	}
	if s.SchemaVersion != ClusterStatusSchemaVersion {
		return nil, fmt.Errorf(
			"unsupported cluster status schemaVersion %d; this mongocli version reads schemaVersion %d — upgrade mongocli or use a matching relay build",
			s.SchemaVersion, ClusterStatusSchemaVersion,
		)
	}
	return &s, nil
}
