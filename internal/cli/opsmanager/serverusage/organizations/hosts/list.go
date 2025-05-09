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

package hosts

import (
	"context"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
	"go.mongodb.org/ops-manager/opsmngr"
)

type ListOpts struct {
	cli.GlobalOpts
	cli.OutputOpts
	cli.ListOpts
	store     store.OrganizationHostAssignmentLister
	startDate string
	endDate   string
}

func (opts *ListOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

var listTemplate = `HOSTNAME	SERVER TYPE	MEM SIZE MB	CHARGEABLE{{- range valueOrEmptySlice .Results}}	
{{.Hostname}}	{{.ServerType.Name}}	{{.MemSizeMB}}	{{.IsChargeable}}{{end}}
`

func (opts *ListOpts) Run() error {
	r, err := opts.store.OrganizationHostAssignments(opts.ConfigOrgID(), opts.newServerTypeOptions())
	if err != nil {
		return err
	}

	return opts.Print(r)
}

func (opts *ListOpts) newServerTypeOptions() *opsmngr.ServerTypeOptions {
	return &opsmngr.ServerTypeOptions{
		ListOptions: *opts.NewListOptions(),
		StartDate:   opts.startDate,
		EndDate:     opts.endDate,
	}
}

// mongocli om serverUsage project(s) hosts list [--orgId orgId].
func ListBuilder() *cobra.Command {
	opts := &ListOpts{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all host assignments for one organization.",
		Args:    require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return opts.PreRunE(
				opts.ValidateOrgID,
				opts.initStore(cmd.Context()),
				opts.InitOutput(cmd.OutOrStdout(), listTemplate),
			)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.startDate, flag.StartDate, "", usage.ServerUsageStartDate)
	cmd.Flags().StringVar(&opts.endDate, flag.EndDate, "", usage.ServerUsageEndDate)

	cmd.Flags().IntVar(&opts.PageNum, flag.Page, cli.DefaultPage, usage.Page)
	cmd.Flags().IntVar(&opts.ItemsPerPage, flag.Limit, cli.DefaultPageLimit, usage.Limit)

	cmd.Flags().StringVar(&opts.OrgID, flag.OrgID, "", usage.OrgID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	_ = cmd.MarkFlagRequired(flag.StartDate)
	_ = cmd.MarkFlagRequired(flag.EndDate)

	return cmd
}
