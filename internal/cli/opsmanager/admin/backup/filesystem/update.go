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
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/opsmanager/admin/backup"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
	"go.mongodb.org/ops-manager/opsmngr"
)

var updateTemplate = "File System configuration '{{.ID}}' updated.\n"

type UpdateOpts struct {
	cli.OutputOpts
	backup.AdminOpts
	store                    store.FileSystemsUpdater
	mmapv1CompressionSetting string
	storePath                string
	wtCompressionSetting     string
}

func (opts *UpdateOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *UpdateOpts) Run() error {
	r, err := opts.store.UpdateFileSystems(opts.newFileSystemConfiguration())
	if err != nil {
		return err
	}
	return opts.Print(r)
}

func (opts *UpdateOpts) newFileSystemConfiguration() *opsmngr.FileSystemStoreConfiguration {
	return &opsmngr.FileSystemStoreConfiguration{
		BackupStore:              *opts.NewBackupStore(),
		MMAPV1CompressionSetting: opts.mmapv1CompressionSetting,
		StorePath:                opts.storePath,
		WTCompressionSetting:     opts.wtCompressionSetting,
	}
}

// mongocli ops-manager admin backup fileSystem(s) update <name> [--assignment][--encryptedCredentials]
// [--label label][--loadFactor loadFactor][--id ID][--storePath storePath][--mmapv1CompressionSetting mmapv1CompressionSetting]
// [--wtCompressionSetting wtCompressionSetting].
func UpdateBuilder() *cobra.Command {
	opts := &UpdateOpts{}
	opts.Template = updateTemplate
	cmd := &cobra.Command{
		Use:   "update <ID>",
		Short: "Update a file system configuration.",
		Args:  require.ExactArgs(1),
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

	cmd.Flags().BoolVar(&opts.Assignment, flag.Assignment, false, usage.FileSystemAssignment)
	cmd.Flags().BoolVar(&opts.EncryptedCredentials, flag.EncryptedCredentials, false, usage.EncryptedCredentials)
	cmd.Flags().StringSliceVar(&opts.Label, flag.Label, []string{}, usage.Label+usage.UpdateWarning)
	cmd.Flags().Int64Var(&opts.LoadFactor, flag.LoadFactor, 0, usage.LoadFactor)
	cmd.Flags().StringVar(&opts.mmapv1CompressionSetting, flag.MMAPV1CompressionSetting, "", usage.MMAPV1CompressionSetting)
	cmd.Flags().StringVar(&opts.wtCompressionSetting, flag.WTCompressionSetting, "", usage.WTCompressionSetting)
	cmd.Flags().StringVar(&opts.storePath, flag.StorePath, "", usage.StorePath)

	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	_ = cmd.MarkFlagRequired(flag.StorePath)
	_ = cmd.MarkFlagRequired(flag.MMAPV1CompressionSetting)
	_ = cmd.MarkFlagRequired(flag.WTCompressionSetting)

	return cmd
}
