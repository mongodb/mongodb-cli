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

package accesslists

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
	"go.mongodb.org/ops-manager/opsmngr"
)

const createTemplate = "Created new access list entry(s).\n"

type CreateOpts struct {
	cli.GlobalOpts
	cli.OutputOpts
	apyKey string
	ips    []string
	cidrs  []string
	store  store.OrganizationAPIKeyAccessListCreator
}

func (opts *CreateOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *CreateOpts) newAccessListAPIKeysReq() ([]*opsmngr.AccessListAPIKeysReq, error) {
	req := make([]*opsmngr.AccessListAPIKeysReq, 0, len(opts.ips)+len(opts.cidrs))
	if len(opts.ips) == 0 && len(opts.cidrs) == 0 {
		return nil, fmt.Errorf("either --%s, --%s must be set", flag.IP, flag.CIDR)
	}
	for _, v := range opts.ips {
		entry := &opsmngr.AccessListAPIKeysReq{
			IPAddress: v,
		}
		req = append(req, entry)
	}

	for _, v := range opts.cidrs {
		entry := &opsmngr.AccessListAPIKeysReq{
			CidrBlock: v,
		}
		req = append(req, entry)
	}

	return req, nil
}

func (opts *CreateOpts) Run() error {
	req, err := opts.newAccessListAPIKeysReq()
	if err != nil {
		return err
	}

	var result *opsmngr.AccessListAPIKeys

	result, err = opts.store.CreateOrganizationAPIKeyAccessList(opts.ConfigOrgID(), opts.apyKey, req)
	if err != nil {
		return err
	}
	return opts.Print(result)
}

// mongocli iam organizations|orgs apiKey(s)|apikeys accessList create [--apiKey keyId] [--orgId orgId] [--ip ip] [--cidr cidr].
func CreateBuilder() *cobra.Command {
	opts := new(CreateOpts)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an IP access list entry for your API Key.",
		Long: fmt.Sprintf(`To view possible values for the apiKey option, run mongocli iam organizations apiKeys list.

%s`, fmt.Sprintf(usage.RequiredRole, "Read Write")),
		Example: `  # Create access list entries for two IP addresses for the API key with the ID 5f24084d8dbffa3ad3f21234 in the organization with the ID 5a1b39eec902201990f12345:
  mongocli iam organizations apiKeys accessLists create --apiKey 5f24084d8dbffa3ad3f21234 --cidr 192.0.2.0/24,198.51.100.0/24 --orgId 5a1b39eec902201990f12345 --output json`,
		Annotations: map[string]string{
			"output": createTemplate,
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return opts.PreRunE(
				opts.ValidateOrgID,
				opts.initStore(cmd.Context()),
				opts.InitOutput(cmd.OutOrStdout(), createTemplate),
			)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.apyKey, flag.APIKey, "", usage.APIKey)
	cmd.Flags().StringSliceVar(&opts.cidrs, flag.CIDR, []string{}, usage.AccessListCIDREntry)
	cmd.Flags().StringSliceVar(&opts.ips, flag.IP, []string{}, usage.APIAccessListIPEntry)

	cmd.Flags().StringVar(&opts.OrgID, flag.OrgID, "", usage.OrgID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	_ = cmd.MarkFlagRequired(flag.APIKey)

	return cmd
}
