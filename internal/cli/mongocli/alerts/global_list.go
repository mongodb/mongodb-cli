// Copyright 2023 MongoDB Inc
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

package alerts

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

type GlobalListOpts struct {
	cli.ListOpts
	cli.OutputOpts
	status string
	store  store.GlobalAlertLister
}

func (opts *GlobalListOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *GlobalListOpts) Run() error {
	alertOpts := opts.newAlertsListOptions()
	r, err := opts.store.GlobalAlerts(alertOpts)
	if err != nil {
		return err
	}

	return opts.Print(r)
}

func (opts *GlobalListOpts) newAlertsListOptions() *opsmngr.AlertsListOptions {
	return &opsmngr.AlertsListOptions{
		Status:      opts.status,
		ListOptions: *opts.NewListOptions(),
	}
}

// mongocli om|cm alert(s) global list [--status status].
func GlobalListBuilder() *cobra.Command {
	opts := &GlobalListOpts{}
	opts.Template = listTemplate
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Returns all global alerts for the specified Ops Manager project.",
		Aliases: []string{"ls"},
		Args:    require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			opts.OutWriter = cmd.OutOrStdout()
			return opts.initStore(cmd.Context())()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().IntVar(&opts.PageNum, flag.Page, cli.DefaultPage, usage.Page)
	cmd.Flags().IntVar(&opts.ItemsPerPage, flag.Limit, cli.DefaultPageLimit, usage.Limit)
	cmd.Flags().StringVar(&opts.status, flag.Status, "", usage.Status)

	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	return cmd
}
