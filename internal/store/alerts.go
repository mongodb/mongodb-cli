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

package store

import (
	"fmt"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"go.mongodb.org/ops-manager/opsmngr"
)

//go:generate mockgen -destination=../mocks/mock_alerts.go -package=mocks github.com/mongodb/mongodb-cli/mongocli/v2/internal/store AlertDescriber,AlertLister,AlertAcknowledger

type AlertDescriber interface {
	Alert(string, string) (*opsmngr.Alert, error)
}

type AlertLister interface {
	Alerts(string, *opsmngr.AlertsListOptions) (*opsmngr.AlertsResponse, error)
}

type AlertAcknowledger interface {
	AcknowledgeAlert(string, string, *opsmngr.AcknowledgeRequest) (*opsmngr.Alert, error)
}

// Alert encapsulate the logic to manage different cloud providers.
func (s *Store) Alert(projectID, alertID string) (*opsmngr.Alert, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.Alerts.Get(s.ctx, projectID, alertID)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// Alerts encapsulate the logic to manage different cloud providers.
func (s *Store) Alerts(projectID string, opts *opsmngr.AlertsListOptions) (*opsmngr.AlertsResponse, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.Alerts.List(s.ctx, projectID, opts)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// AcknowledgeAlert encapsulate the logic to manage different cloud providers.
func (s *Store) AcknowledgeAlert(projectID, alertID string, body *opsmngr.AcknowledgeRequest) (*opsmngr.Alert, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.Alerts.Acknowledge(s.ctx, projectID, alertID, body)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}
