// Copyright 2023 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

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

type EventListOpts struct {
	cli.ListOpts
	EventType []string
	MinDate   string
	MaxDate   string
}

func (opts *EventListOpts) newEventListOptions() *opsmngr.EventListOptions {
	return &opsmngr.EventListOptions{
		ListOptions: opsmngr.ListOptions{
			PageNum:      opts.PageNum,
			ItemsPerPage: opts.ItemsPerPage,
		},
		EventType: opts.EventType,
		MinDate:   opts.MinDate,
		MaxDate:   opts.MaxDate,
	}
}

type ListOpts struct {
	EventListOpts
	cli.OutputOpts
	orgID     string
	projectID string
	store     store.EventLister
}

func (opts *ListOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

var listTemplate = `ID	TYPE	CREATED{{range valueOrEmptySlice .Results}}
{{.ID}}	{{.EventTypeName}}	{{.Created}}{{end}}
`

func (opts *ListOpts) Run() error {
	listOpts := opts.newEventListOptions()

	var r *opsmngr.EventResponse
	var err error

	if opts.orgID != "" {
		r, err = opts.store.OrganizationEvents(opts.orgID, listOpts)
	} else if opts.projectID != "" {
		r, err = opts.store.ProjectEvents(opts.projectID, listOpts)
	}
	if err != nil {
		return err
	}

	return opts.Print(r)
}

// ListBuilder
//
//	mongocli atlas event(s) list
//
// [--projectId projectId]
// [--orgId orgId]
// [--page N]
// [--limit N]
// [--minDate minDate]
// [--maxDate maxDate].
func ListBuilder() *cobra.Command {
	opts := &ListOpts{}
	opts.Template = listTemplate
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Return all events for an organization or project.",
		Long:  fmt.Sprintf(usage.RequiredRole, "Project Read Only"),
		Deprecated: `  
  To return project events prefer
  mongocli atlas|ops-manager|cloud-manager events projects list [--projectId <projectId>]

  To return organization events prefer
  mongocli atlas|ops-manager|cloud-manager events organizations list [--orgId <orgId>]
`,
		Aliases: []string{"ls"},
		Args:    require.NoArgs,
		Annotations: map[string]string{
			"output": listTemplate,
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if opts.orgID != "" && opts.projectID != "" {
				return fmt.Errorf("both --%s and --%s set", flag.ProjectID, flag.OrgID)
			}
			if opts.orgID == "" && opts.projectID == "" {
				return fmt.Errorf("--%s or --%s must be set", flag.ProjectID, flag.OrgID)
			}
			opts.OutWriter = cmd.OutOrStdout()
			return opts.initStore(cmd.Context())()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().IntVar(&opts.PageNum, flag.Page, cli.DefaultPage, usage.Page)
	cmd.Flags().IntVar(&opts.ItemsPerPage, flag.Limit, cli.DefaultPageLimit, usage.Limit)

	cmd.Flags().StringSliceVar(&opts.EventType, flag.TypeFlag, nil, usage.Event)
	cmd.Flags().StringVar(&opts.MaxDate, flag.MaxDate, "", usage.MaxDate)
	cmd.Flags().StringVar(&opts.MinDate, flag.MinDate, "", usage.MinDate)

	cmd.Flags().StringVar(&opts.projectID, flag.ProjectID, "", usage.ProjectID)
	cmd.Flags().StringVar(&opts.orgID, flag.OrgID, "", usage.OrgID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, opts.AutoCompleteOutputFlag())

	return cmd
}
