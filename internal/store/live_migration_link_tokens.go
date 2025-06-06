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
)

//go:generate mockgen -destination=../mocks/mock_live_migration_link_tokens.go -package=mocks github.com/mongodb/mongodb-cli/mongocli/v2/internal/store LinkTokenDeleter

type LinkTokenDeleter interface {
	DeleteLinkToken(string) error
}

// DeleteLinkToken encapsulate the logic to manage different cloud providers.
func (s *Store) DeleteLinkToken(orgID string) error {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		_, err := s.client.LiveMigration.DeleteConnection(s.ctx, orgID)
		return err
	default:
		return fmt.Errorf("%w: %s", errUnsupportedService, s.service)
	}
}
