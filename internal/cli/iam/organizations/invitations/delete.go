// Copyright 2021 MongoDB Inc
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

package invitations

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

type DeleteOpts struct {
	*cli.DeleteOpts
	cli.GlobalOpts
	store store.OrganizationInvitationDeleter
}

func (opts *DeleteOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *DeleteOpts) Run() error {
	return opts.Delete(opts.store.DeleteInvitation, opts.ConfigOrgID())
}

// mongocli iam organization(s) invitation(s) delete <invitationId> [--force] [--orgId orgId].
func DeleteBuilder() *cobra.Command {
	opts := &DeleteOpts{
		DeleteOpts: cli.NewDeleteOpts("Invitation '%s' deleted\n", "Invitation not deleted"),
	}
	cmd := &cobra.Command{
		Use:     "delete <invitationId>",
		Aliases: []string{"rm"},
		Short:   "Remove the specified pending invitation to your organization.",
		Long:    fmt.Sprintf(usage.RequiredRole, "Organization User Admin"),
		Args:    require.ExactArgs(1),
		Annotations: map[string]string{
			"invitationIdDesc": "Unique 24-digit string that identifies the invitation.",
			"output":           opts.SuccessMessage(),
		},
		Example: `  # Remove the pending invitation with the ID 5dd56c847a3e5a1f363d424d from the organization with the ID 5f71e5255afec75a3d0f96dc:
  mongocli iam organizations invitations delete 5dd56c847a3e5a1f363d424d --orgId 5f71e5255afec75a3d0f96dc`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.PreRunE(opts.ValidateOrgID, opts.initStore(cmd.Context())); err != nil {
				return err
			}
			opts.Entry = args[0]
			return opts.Prompt()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().BoolVar(&opts.Confirm, flag.Force, false, usage.Force)
	cmd.Flags().StringVar(&opts.OrgID, flag.OrgID, "", usage.OrgID)

	return cmd
}
