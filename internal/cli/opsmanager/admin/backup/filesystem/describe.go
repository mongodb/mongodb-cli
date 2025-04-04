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

package filesystem

import (
	"context"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

var describeTemplate = `NAME	PATH	WT COMPRESSION	MMAPV1 COMPRESSION
{{.ID}}	{{.StorePath}}	{{.WTCompressionSetting}}	{{.MMAPV1CompressionSetting}}
`

type DescribeOpts struct {
	cli.OutputOpts
	store store.FileSystemsDescriber
	ID    string
}

func (opts *DescribeOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *DescribeOpts) Run() error {
	r, err := opts.store.DescribeFileSystem(opts.ID)
	if err != nil {
		return err
	}

	return opts.Print(r)
}

// mongocli ops-manager admin backup fileSystem(s) describe <name>.
func DescribeBuilder() *cobra.Command {
	opts := &DescribeOpts{}
	opts.Template = describeTemplate
	cmd := &cobra.Command{
		Use:     "describe <ID>",
		Aliases: []string{"get"},
		Short:   "Get a file system configuration.",
		Args:    require.ExactArgs(1),
		Annotations: map[string]string{
			"IDDesc": "Configuration identifier.",
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			opts.OutWriter = cmd.OutOrStdout()
			return opts.initStore(cmd.Context())()
		},
		RunE: func(_ *cobra.Command, args []string) error {
			opts.ID = args[0]
			return opts.Run()
		},
	}

	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	return cmd
}
