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

package security

import (
	"github.com/mongodb/mongocli/internal/cli/atlas/security/customercerts"
	"github.com/mongodb/mongocli/internal/cli/atlas/security/ldap"
	"github.com/spf13/cobra"
)

const security = "Manage security configuration for your project."

func Builder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: security,
	}

	cmd.AddCommand(
		customercerts.Builder(),
		ldap.Builder(),
	)

	return cmd
}
