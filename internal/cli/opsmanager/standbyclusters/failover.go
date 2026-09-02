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
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/prompt"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

type FailoverOpts struct {
	opts
	cli.OutputOpts

	confirm bool
	watch   bool
	errW    io.Writer
}

var failoverTemplate = `Failover triggered for standby cluster '{{.ClusterName}}'. State: {{.State}}.
`

var failoverDoneTemplate = `Failover completed: standby cluster '{{.ClusterName}}' is now {{.State}}.
`

var watchPollInterval = 5 * time.Second

func (o *FailoverOpts) Run(ctx context.Context) error {
	if !o.confirm {
		p := prompt.NewConfirm(fmt.Sprintf(
			"Are you sure you want to trigger a failover of the standby cluster at %s? The standby cluster will be promoted to active.",
			o.store.FullPath(),
		))
		if err := survey.AskOne(p, &o.confirm); err != nil {
			return err
		}
		if !o.confirm {
			return errors.New("user-abort. Failover not triggered")
		}
	}

	written, err := o.triggerFailover(ctx)
	if err != nil {
		return err
	}

	if o.watch {
		final, err := o.watchUntilActive(ctx, written)
		if err != nil {
			return err
		}
		written = final
		o.Template = failoverDoneTemplate
	}
	return o.Print(written)
}

func (o *FailoverOpts) triggerFailover(ctx context.Context) (*standby.RemoteDRState, error) {
	return standby.TriggerFailover(ctx, o.store)
}

func (o *FailoverOpts) watchUntilActive(ctx context.Context, last *standby.RemoteDRState) (*standby.RemoteDRState, error) {
	_, _ = fmt.Fprintf(o.errW, "Watching %s (Ctrl+C to stop):\n", o.store.FullPath())
	printStateChange(o.errW, last)

	ticker := time.NewTicker(watchPollInterval)
	defer ticker.Stop()

	// heartbeat: one dot per poll while nothing changed
	dots := false

	for {
		select {
		case <-ctx.Done():
			_, _ = fmt.Fprintf(o.errW, "\nStopped watching; the failover is still in progress.\n")
			return last, nil
		case <-ticker.C:
		}

		state, _, err := o.store.FetchDRState(ctx)
		if err != nil {
			if ctx.Err() != nil {
				continue
			}
			return last, err
		}
		if state == nil {
			return last, fmt.Errorf("the DR state file at %s no longer exists", o.store.FullPath())
		}
		if state.State != last.State || state.Version != last.Version {
			if dots {
				_, _ = fmt.Fprintln(o.errW)
				dots = false
			}
			printStateChange(o.errW, state)
			last = state
		} else {
			_, _ = fmt.Fprint(o.errW, ".")
			dots = true
		}
		if state.State == standby.StateActive {
			if dots {
				_, _ = fmt.Fprintln(o.errW)
			}
			return last, nil
		}
	}
}

func printStateChange(w io.Writer, s *standby.RemoteDRState) {
	ts := s.LastModified
	if t, err := time.Parse(time.RFC3339, s.LastModified); err == nil {
		ts = t.Local().Format("15:04:05")
	}
	_, _ = fmt.Fprintf(w, "%s  %-20s (was %s, version %s)\n", ts, s.State, s.PreviousState, s.Version)
}

// mongocli ops-manager standby-clusters failover [--force] [--watch].
func FailoverBuilder() *cobra.Command {
	o := &FailoverOpts{}
	cmd := &cobra.Command{
		Use:   "failover",
		Short: "Trigger a failover of your standby cluster without Ops Manager.",
		Long: `Promotes the standby cluster to active by writing the PromoteStandby
state to the disaster recovery (DR) state file in S3.

The command fails unless the current state is 'Standby'. The write is guarded against
concurrent modifications.

Run 'mongocli ops-manager standby-clusters configure' first to set up access
to the DR state file.`,
		Example: `  # Trigger a failover of the standby cluster, with a confirmation prompt:
  mongocli ops-manager standby-clusters failover

  # Trigger a failover without the confirmation prompt and follow the
  # promotion progress until the cluster becomes active:
  mongocli ops-manager standby-clusters failover --force --watch`,
		Args: require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			o.errW = cmd.ErrOrStderr()
			return o.PreRunE(
				o.initStore(cmd.Context()),
				o.InitOutput(cmd.OutOrStdout(), failoverTemplate),
			)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return o.Run(cmd.Context())
		},
	}

	addSharedFlags(cmd, &o.opts)
	cmd.Flags().BoolVar(&o.confirm, flag.Force, false, usage.Force)
	cmd.Flags().BoolVarP(&o.watch, flag.Watch, flag.WatchShort, false, usage.StandbyWatch)
	cmd.Flags().StringVarP(&o.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, o.AutoCompleteOutputFlag())

	return cmd
}
