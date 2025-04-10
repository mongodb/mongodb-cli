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
	"fmt"

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
	status string
	store  store.AlertLister
}

func (opts *ListOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

var listTemplate = `ID	TYPE	STATUS{{range valueOrEmptySlice .Results}}
{{.ID}}	{{.EventTypeName}}	{{.Status}}{{end}}
`

func (opts *ListOpts) Run() error {
	listOpts := opts.newAlertsListOptions()
	r, err := opts.store.Alerts(opts.ConfigProjectID(), listOpts)

	if err != nil {
		return err
	}

	return opts.Print(r)
}

func (opts *ListOpts) newAlertsListOptions() *opsmngr.AlertsListOptions {
	o := &opsmngr.AlertsListOptions{
		Status:      opts.status,
		ListOptions: *opts.NewListOptions(),
	}

	return o
}

// mongocli atlas alerts list [--status status] [--projectId projectId] [--page N] [--limit N].
func ListBuilder() *cobra.Command {
	opts := new(ListOpts)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Return all alerts for your project.",
		Long:  fmt.Sprintf(usage.RequiredRole, "Project Read Only"),
		Example: fmt.Sprintf(`  # Return a JSON-formatted list of all alerts for the project with the ID 5e2211c17a3e5a48f5497de3:
  %s alerts list --projectId 5e2211c17a3e5a48f5497de3 --output json`, cli.ExampleEntryPoint()),
		Aliases: []string{"ls"},
		Annotations: map[string]string{
			"output": listTemplate,
		},
		Args: require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return opts.PreRunE(
				opts.ValidateProjectID,
				opts.initStore(cmd.Context()),
				opts.InitOutput(cmd.OutOrStdout(), listTemplate),
			)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}
	cmd.Flags().IntVar(&opts.PageNum, flag.Page, cli.DefaultPage, usage.Page)
	cmd.Flags().IntVar(&opts.ItemsPerPage, flag.Limit, cli.DefaultPageLimit, usage.Limit)
	cmd.Flags().StringVar(&opts.status, flag.Status, "", usage.Status)
	cmd.Flags().StringVar(&opts.ProjectID, flag.ProjectID, "", usage.ProjectID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())
	return cmd
}
