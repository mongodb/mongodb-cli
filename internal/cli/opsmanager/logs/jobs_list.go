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

package logs

import (
	"context"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
	"go.mongodb.org/ops-manager/opsmngr"
)

type JobsListOpts struct {
	cli.GlobalOpts
	cli.OutputOpts
	verbose bool
	store   store.LogJobLister
}

func (opts *JobsListOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

var listTemplate = `ID	CREATED AT	EXPIRES AT	STATUS	URL	REDACTED{{range valueOrEmptySlice .Results}}
{{.ID}}	{{.CreationDate}}	{{.ExpirationDate}}	{{.Status}}	{{.DownloadURL}}	{{.Redacted}}{{end}}
`

func (opts *JobsListOpts) Run() error {
	r, err := opts.store.LogCollectionJobs(opts.ConfigProjectID(), opts.newLogListOptions())
	if err != nil {
		return err
	}
	return opts.Print(r)
}

func (opts *JobsListOpts) newLogListOptions() *opsmngr.LogListOptions {
	return &opsmngr.LogListOptions{Verbose: opts.verbose}
}

// mongocli om logs jobs list --verbose verbose [--projectId projectId].
func JobsListOptsBuilder() *cobra.Command {
	opts := &JobsListOpts{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List log collection jobs for your project.",
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

	cmd.Flags().BoolVar(&opts.verbose, flag.Verbose, false, usage.Verbose)

	cmd.Flags().StringVar(&opts.ProjectID, flag.ProjectID, "", usage.ProjectID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	return cmd
}
