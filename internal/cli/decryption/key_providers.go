// Copyright 2022 MongoDB Inc
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

package decryption

import (
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/spf13/cobra"
)

func KeyProvidersBuilder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keyProviders",
		Aliases: cli.GenerateAliases("keyProviders", "keys"),
		Short:   "Manage your key collections.",
	}

	cmd.AddCommand(KeyProvidersListBuilder())

	return cmd
}
