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

type AdvisorUpgrade interface {
	AdvisorUpgradeCheck(version string) (*[]opsmngr.UpgradeCheckStep, error)
}

// AdvisorUpgradeCheck ...
func (s *Store) AdvisorUpgradeCheck(version string) (*[]opsmngr.UpgradeCheckStep, error) {
	switch s.service {
	case config.OpsManagerService:
		result, _, err := s.client.Advisor.CheckUpgrade(s.ctx, version)
		return result, err
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}
