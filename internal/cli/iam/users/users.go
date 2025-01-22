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

package users

import (
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/spf13/cobra"
)

func Builder() *cobra.Command {
	description := "Create, list and manage your Cloud Manager or Ops Manager users."

	const use = "users"
	cmd := &cobra.Command{
		Use:     use,
		Short:   "Manage your Ops Manager or Cloud Manager users.",
		Long:    description,
		Aliases: cli.GenerateAliases(use),
	}

	cmd.AddCommand(
		InviteBuilder(),
		DescribeBuilder(),
		DeleteBuilder(),
	)

	return cmd
}
