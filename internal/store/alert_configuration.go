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
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/pointer"
	"go.mongodb.org/ops-manager/opsmngr"
)

//go:generate mockgen -destination=../mocks/mock_alert_configuration.go -package=mocks github.com/mongodb/mongodb-cli/mongocli/v2/internal/store AlertConfigurationLister,AlertConfigurationCreator,AlertConfigurationDeleter,AlertConfigurationUpdater,MatcherFieldsLister,AlertConfigurationEnabler,AlertConfigurationDisabler

type AlertConfigurationLister interface {
	AlertConfigurations(string, *opsmngr.ListOptions) ([]opsmngr.AlertConfiguration, error)
}

type AlertConfigurationCreator interface {
	CreateAlertConfiguration(*opsmngr.AlertConfiguration) (*opsmngr.AlertConfiguration, error)
}

type AlertConfigurationDeleter interface {
	DeleteAlertConfiguration(string, string) error
}

type AlertConfigurationUpdater interface {
	UpdateAlertConfiguration(*opsmngr.AlertConfiguration) (*opsmngr.AlertConfiguration, error)
}

type MatcherFieldsLister interface {
	MatcherFields() ([]string, error)
}

type AlertConfigurationEnabler interface {
	EnableAlertConfiguration(string, string) (*opsmngr.AlertConfiguration, error)
}

type AlertConfigurationDisabler interface {
	DisableAlertConfiguration(string, string) (*opsmngr.AlertConfiguration, error)
}

// AlertConfigurations encapsulate the logic to manage different cloud providers.
func (s *Store) AlertConfigurations(projectID string, opts *opsmngr.ListOptions) ([]opsmngr.AlertConfiguration, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.AlertConfigurations.List(s.ctx, projectID, opts)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// CreateAlertConfiguration encapsulate the logic to manage different cloud providers.
func (s *Store) CreateAlertConfiguration(alertConfig *opsmngr.AlertConfiguration) (*opsmngr.AlertConfiguration, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.AlertConfigurations.Create(s.ctx, alertConfig.GroupID, alertConfig)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// DeleteAlertConfiguration encapsulate the logic to manage different cloud providers.
func (s *Store) DeleteAlertConfiguration(projectID, id string) error {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		_, err := s.client.AlertConfigurations.Delete(s.ctx, projectID, id)
		return err
	default:
		return fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// MatcherFields encapsulate the logic to manage different cloud providers.
func (s *Store) MatcherFields() ([]string, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.AlertConfigurations.ListMatcherFields(s.ctx)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

func (s *Store) UpdateAlertConfiguration(alertConfig *opsmngr.AlertConfiguration) (*opsmngr.AlertConfiguration, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.AlertConfigurations.Update(s.ctx, alertConfig.GroupID, alertConfig.ID, alertConfig)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// EnableAlertConfiguration encapsulate the logic to manage different cloud providers.
func (s *Store) EnableAlertConfiguration(projectID, id string) (*opsmngr.AlertConfiguration, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.AlertConfigurations.EnableAnAlertConfig(s.ctx, projectID, id, pointer.Get(true))
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}

// DisableAlertConfiguration encapsulate the logic to manage different cloud providers.
func (s *Store) DisableAlertConfiguration(projectID, id string) (*opsmngr.AlertConfiguration, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.AlertConfigurations.EnableAnAlertConfig(s.ctx, projectID, id, pointer.Get(false))
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}
