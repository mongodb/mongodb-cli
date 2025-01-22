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

package teams

import (
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

const describeTemplate = `ID	NAME
{{.ID}}	{{.Name}}
`

type DescribeOpts struct {
	cli.GlobalOpts
	cli.OutputOpts
	store store.TeamDescriber
	name  string
	id    string
}

func (opts *DescribeOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *DescribeOpts) Run() error {
	var r interface{}
	var err error

	if opts.name != "" {
		r, err = opts.store.TeamByName(opts.ConfigOrgID(), opts.name)
	}

	if opts.id != "" {
		r, err = opts.store.TeamByID(opts.ConfigOrgID(), opts.id)
	}

	if err != nil {
		return err
	}

	return opts.Print(r)
}

func (opts *DescribeOpts) validate() error {
	if opts.id == "" && opts.name == "" {
		return errors.New("must supply one of 'id' or 'username'")
	}

	if opts.id != "" && opts.name != "" {
		return errors.New("cannot supply both 'id' and 'username'")
	}

	return nil
}

// mongocli iam team(s) describe --id id --name name --orgId orgId.
func DescribeBuilder() *cobra.Command {
	opts := &DescribeOpts{}
	cmd := &cobra.Command{
		Use:         "describe",
		Aliases:     []string{"get"},
		Annotations: map[string]string{"output": describeTemplate},
		Example: `  # Return the JSON-formatted details for the the team with the ID 5e44445ef10fab20b49c0f31 in the organization with ID 5e2211c17a3e5a48f5497de3:
  mongocli iam teams describe --id 5e44445ef10fab20b49c0f31 --projectId 5e1234c17a3e5a48f5497de3 --output json
  
  # Return the JSON-formatted details for the the team with the name myTeam in the organization with ID 5e2211c17a3e5a48f5497de3:
  mongocli iam teams describe --name myTeam --projectId 5e1234c17a3e5a48f5497de3 --output json`,
		Short: "Return the details for the specified team for your organization.",
		Long: `You can return the details for a team using the team's ID or the team's name. You must specify either the id option or the name option.

` + fmt.Sprintf(usage.RequiredRole, "Organization Member"),
		Args: require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return opts.PreRunE(
				opts.ValidateOrgID,
				opts.initStore(cmd.Context()),
				opts.InitOutput(cmd.OutOrStdout(), describeTemplate),
				opts.validate,
			)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.name, flag.Name, "", usage.TeamName)
	cmd.Flags().StringVar(&opts.id, flag.ID, "", usage.TeamID)

	cmd.Flags().StringVar(&opts.OrgID, flag.OrgID, "", usage.OrgID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	return cmd
}
