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

package settings

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

type DisableOpts struct {
	cli.GlobalOpts
	cli.OutputOpts
	alertID string
	store   store.AlertConfigurationDisabler
}

func (opts *DisableOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

var disableTemplate = "Alert configuration '{{.ID}}' disabled\n"

func (opts *DisableOpts) Run() error {
	r, err := opts.store.DisableAlertConfiguration(opts.ConfigProjectID(), opts.alertID)
	if err != nil {
		return err
	}

	return opts.Print(r)
}

// mongocli atlas alerts disable <ID> --projectId projectId.
func DisableBuilder() *cobra.Command {
	opts := new(DisableOpts)
	cmd := &cobra.Command{
		Use:   "disable <alertConfigId>",
		Short: "Disables one alert configuration for the specified project.",
		Args:  require.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return opts.PreRunE(
				opts.ValidateProjectID,
				opts.initStore(cmd.Context()),
				opts.InitOutput(cmd.OutOrStdout(), disableTemplate),
			)
		},
		Example: fmt.Sprintf(`  # Disable the alert configuration with the ID 5d1113b25a115342acc2d1aa in the project with the ID 5e2211c17a3e5a48f5497de3:
  %s alerts settings disable 5d1113b25a115342acc2d1aa --projectId 5e2211c17a3e5a48f5497de3`, cli.ExampleEntryPoint()),
		Annotations: map[string]string{
			"alertConfigIdDesc": "ID of the alert configuration you want to disable.",
		},
		RunE: func(_ *cobra.Command, args []string) error {
			opts.alertID = args[0]
			return opts.Run()
		},
	}
	cmd.OutOrStdout()
	cmd.Flags().StringVar(&opts.ProjectID, flag.ProjectID, "", usage.ProjectID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	return cmd
}
