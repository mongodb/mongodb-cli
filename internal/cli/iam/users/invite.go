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

package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/prerun"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
	"go.mongodb.org/ops-manager/opsmngr"
)

var inviteTemplate = "The user '{{.Username}}' has been invited.\nInvited users do not have access to the project until they accept the invitation.\n"

type InviteOpts struct {
	cli.OutputOpts
	cli.InputOpts
	username     string
	password     string
	country      string
	email        string
	mobile       string
	firstName    string
	lastName     string
	orgRoles     []string
	projectRoles []string
	store        store.UserCreator
}

func (opts *InviteOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *InviteOpts) newUserRequest() (*opsmngr.User, error) {
	userRoles, err := opts.createUserRole()
	if err != nil {
		return nil, err
	}

	user := &opsmngr.User{
		Username:     opts.username,
		Password:     opts.password,
		FirstName:    opts.firstName,
		LastName:     opts.lastName,
		EmailAddress: opts.email,
		MobileNumber: opts.mobile,
		Country:      opts.country,
		Roles:        userRoles,
	}

	return user, nil
}

func (opts *InviteOpts) Run() error {
	user, err := opts.newUserRequest()
	if err != nil {
		return err
	}

	r, err := opts.store.CreateUser(user)
	if err != nil {
		return err
	}

	return opts.Print(r)
}

const keyParts = 2

func (opts *InviteOpts) createUserRole() ([]*opsmngr.UserRole, error) {
	roles := make([]*opsmngr.UserRole, len(opts.orgRoles)+len(opts.projectRoles))

	i := 0
	for _, role := range opts.orgRoles {
		userRole, err := newUserOrgRole(role)
		if err != nil {
			return nil, err
		}
		roles[i] = userRole
		i++
	}

	for _, role := range opts.projectRoles {
		userRole, err := newUserProjectRole(role)
		if err != nil {
			return nil, err
		}
		roles[i] = userRole
		i++
	}

	return roles, nil
}

func (opts *InviteOpts) Prompt() error {
	if opts.password != "" {
		return nil
	}

	if !opts.IsTerminalInput() {
		_, err := fmt.Fscanln(opts.InReader, &opts.password)
		return err
	}

	prompt := &survey.Password{
		Message: "Password:",
	}
	return survey.AskOne(prompt, &opts.password)
}

func splitRole(role string) ([]string, error) {
	value := strings.Split(role, ":")
	if len(value) != keyParts {
		return nil, fmt.Errorf("unexpected role format: %s", role)
	}
	return value, nil
}

func newUserOrgRole(role string) (*opsmngr.UserRole, error) {
	value, err := splitRole(role)
	if err != nil {
		return nil, err
	}
	userRole := &opsmngr.UserRole{
		OrgID:    value[0],
		RoleName: strings.ToUpper(value[1]),
	}

	return userRole, nil
}

func newUserProjectRole(role string) (*opsmngr.UserRole, error) {
	value, err := splitRole(role)
	if err != nil {
		return nil, err
	}
	userRole := &opsmngr.UserRole{
		GroupID:  value[0],
		RoleName: strings.ToUpper(value[1]),
	}

	return userRole, nil
}

// mongocli iam users(s) invite --username username --password password --country country --email email
// --mobile mobile --firstName firstName --lastName lastName --team team1,team2 --orgRoles orgID:ROLE_NAME
// --projectRoles projectID:ROLE_NAME

func InviteBuilder() *cobra.Command {
	opts := &InviteOpts{}
	opts.Template = inviteTemplate
	cmd := &cobra.Command{
		Use:   "invite",
		Short: "Create a MongoDB user for your MongoDB application and invite the MongoDB user to your organizations and projects.",
		Long:  fmt.Sprintf(`A MongoDB user account grants access only to the the MongoDB application. To grant database access, create a database user with %s dbusers create.`, cli.ExampleEntryPoint()),
		Args:  require.NoArgs,
		Example: `  # Create the MongoDB user with the username user@example.com and invite them to the organization with the ID 5dd56c847a3e5a1f363d424d with ORG_OWNER access:
  mongocli iam users invite --email user@example.com --username user@example.com --orgRole 5dd56c847a3e5a1f363d424d:ORG_OWNER --firstName Example --lastName User --country US --output json
  
  # Create the MongoDB user with the username user@example.com and invite them to the project with the ID 5f71e5255afec75a3d0f96dc with GROUP_READ_ONLY access:
  mongocli iam users invite --email user@example.com --username user@example.com --projectRole 5f71e5255afec75a3d0f96dc:GROUP_READ_ONLY --firstName Example --lastName User --country US --output json`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if config.Service() != config.OpsManagerService {
				_ = cmd.MarkFlagRequired(flag.Country)
			}

			return prerun.ExecuteE(opts.InitOutput(cmd.OutOrStdout(), ""), opts.InitInput(cmd.InOrStdin()), opts.initStore(cmd.Context()))
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := opts.Prompt(); err != nil {
				return err
			}
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.username, flag.Username, "", usage.Username)
	cmd.Flags().StringVar(&opts.password, flag.Password, "", usage.Password)
	cmd.Flags().StringVar(&opts.country, flag.Country, "", usage.Country)
	cmd.Flags().StringVar(&opts.email, flag.Email, "", usage.Email)
	cmd.Flags().StringVar(&opts.mobile, flag.Mobile, "", usage.Mobile)
	cmd.Flags().StringVar(&opts.firstName, flag.FirstName, "", usage.FirstName)
	cmd.Flags().StringVar(&opts.lastName, flag.LastName, "", usage.LastName)
	cmd.Flags().StringSliceVar(&opts.orgRoles, flag.OrgRole, []string{}, usage.MCLIUserOrgRole)
	cmd.Flags().StringSliceVar(&opts.projectRoles, flag.ProjectRole, []string{}, usage.MCLIUserProjectRole)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	_ = cmd.MarkFlagRequired(flag.Username)
	_ = cmd.MarkFlagRequired(flag.Email)
	_ = cmd.MarkFlagRequired(flag.FirstName)
	_ = cmd.MarkFlagRequired(flag.LastName)

	return cmd
}
