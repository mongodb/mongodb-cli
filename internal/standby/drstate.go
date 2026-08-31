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
	"fmt"
	"strconv"
	"time"
)

const (
	StateStandby        = "Standby"
	StatePromoteStandby = "PromoteStandby"
)

type RemoteDRState struct {
	State         string `json:"state"`
	PreviousState string `json:"previousState,omitempty"`
	ClusterName   string `json:"clusterName"`
	Version       string `json:"version"`
	LastModified  string `json:"lastModified"`
	SchemaVersion string `json:"schemaVersion"`

	Planned               bool   `json:"planned,omitempty"`
	SyncDestination       string `json:"syncDestination,omitempty"`
	PlannedFailoverTo     string `json:"plannedFailoverTo,omitempty"`
	CancelPlannedFailover string `json:"cancelPlannedFailover,omitempty"`
}

func UnmarshalRemoteDRState(data []byte) (*RemoteDRState, error) {
	var s RemoteDRState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error unmarshaling Remote DR state: %w", err)
	}
	return &s, nil
}

func (s *RemoteDRState) Marshal() ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("error marshaling Remote DR state: %w", err)
	}
	return data, nil
}

func (s *RemoteDRState) ApplyFailover() error {
	if s.State != StateStandby {
		return fmt.Errorf("current state is %q; a failover can only be triggered from %q", s.State, StateStandby)
	}

	s.PreviousState = s.State
	s.State = StatePromoteStandby
	s.Version = incrementVersion(s.Version)
	s.LastModified = time.Now().UTC().Format(time.RFC3339)
	s.Planned = false
	return nil
}

func incrementVersion(v string) string {
	n, err := strconv.Atoi(v)
	if err != nil {
		return "1"
	}
	return strconv.Itoa(n + 1)
}
