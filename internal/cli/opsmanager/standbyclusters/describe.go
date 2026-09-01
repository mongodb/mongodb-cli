// Copyright 2026 MongoDB Inc
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

package standbyclusters

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

type DescribeOpts struct {
	opts
	cli.OutputOpts
}

var describeTemplate = `STATE	PREVIOUS STATE	CLUSTER	VERSION	LAST MODIFIED
{{.State}}	{{.PreviousState}}	{{.ClusterName}}	{{.Version}}	{{.LastModified}}
`

func (o *DescribeOpts) fetchState(ctx context.Context) (*standby.RemoteDRState, error) {
	state, _, err := o.store.FetchDRState(ctx)
	return state, err
}

func (o *DescribeOpts) Run(ctx context.Context) error {
	state, err := o.fetchState(ctx)
	if err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("no DR state file found at %s; check your standby-clusters configuration", o.store.FullPath())
	}
	return o.Print(state)
}

// mongocli ops-manager standby-clusters describe.
func DescribeBuilder() *cobra.Command {
	o := &DescribeOpts{}
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Return the current disaster recovery state of your standby cluster.",
		Long: `Reads the DR state file directly from Amazon S3 and prints it. Use
'-o json' for the full document.`,
		Example: `  # Return the current DR state of the standby cluster:
  mongocli ops-manager standby-clusters describe

  # Return the full DR state document as JSON:
  mongocli ops-manager standby-clusters describe -o json`,
		Args: require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return o.PreRunE(
				o.initStore(cmd.Context()),
				o.InitOutput(cmd.OutOrStdout(), describeTemplate),
			)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return o.Run(cmd.Context())
		},
	}

	addSharedFlags(cmd, &o.opts)
	cmd.Flags().StringVarP(&o.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, o.AutoCompleteOutputFlag())

	return cmd
}
