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
	"math"
	"strconv"
	"strings"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli/require"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

type LagOpts struct {
	opts
	cli.OutputOpts

	shards bool
}

// jsonOutputFormat is the only -o value this command supports. An empty -o
// keeps the default plaintext behavior.
//
// Kept as a local constant rather than reusing cli.IsJSONOutput(): that helper is
// just an equality test against the (unexported) "json" format, but this command
// needs a three-way split of the ConfigOutput()-resolved value — the plaintext
// default (""), json, and every other format (rejected as unsupported below) —
// which a boolean cannot express. The cli package's json format constant is
// unexported, so it is not referenceable either.
const jsonOutputFormat = "json"

func (o *LagOpts) Run(ctx context.Context) error {
	// Resolve the effective output format the same way the other OutputOpts
	// commands do: an explicit -o wins, otherwise fall back to the global
	// 'output' config setting. With neither set this stays "" (plaintext).
	output := o.ConfigOutput()
	if output != "" && output != jsonOutputFormat {
		return fmt.Errorf("unsupported output format %q for this command; supported values: json", output)
	}

	status, err := o.store.FetchClusterStatus(ctx)
	if err != nil {
		return err
	}

	if output == jsonOutputFormat {
		// -o json prints the full decoded Cluster Status Object directly, so
		// --shards has no additional effect: shardStatus is already part of the
		// document. The struct's json tags (and Timestamp.MarshalJSON) carry the
		// wire contract for both decode and display.
		return o.Print(status)
	}

	lag := status.ClusterStatus.EstimatedDataLossSeconds

	// The plaintext output always shows the same labeled lines, one set per
	// scenario; only the values change (a number, short words, or NA), so the
	// shape is stable for scripts.
	//
	// Line 1 — estimated data loss. NA with the reason when the standby is
	// uninitialized (the field is omitted until replication starts).
	if lag == nil {
		if _, err := fmt.Fprintln(o.OutWriter, "estimated data loss: NA (standby has not yet initialized replication)"); err != nil {
			return err
		}
	} else if _, err := fmt.Fprintf(o.OutWriter, "estimated data loss: %ss\n", formatLagSeconds(*lag)); err != nil {
		return err
	}

	// Line 2 (--shards only) — per-shard pending apply.
	if o.shards {
		if err := o.printShards(status, lag); err != nil {
			return err
		}
	}

	// Line 3 — standby lag (the cluster-wide total).
	return o.printStandbyLag(status, lag)
}

// standbyLagSeconds derives the cluster-wide standby lag: the estimated data
// loss plus the pending apply (ESTS - MCTS). The timestamp delta is computed in
// int64 so a not-yet-initialized or out-of-order pair never wraps a uint32. ok is
// false when the total cannot be computed — either timestamp is unset (0) or
// MCTS is ahead of ESTS. The caller must ensure EstimatedDataLossSeconds is
// non-nil.
func standbyLagSeconds(status *standby.ClusterStatus) (float64, bool) {
	ests := status.ClusterStatus.EarliestSafeTimestamp.T
	mcts := status.ClusterStatus.MinCommittedOpTimestamp.T
	if ests == 0 || mcts == 0 || ests < mcts {
		return 0, false
	}
	delta := int64(ests) - int64(mcts)
	return *status.ClusterStatus.EstimatedDataLossSeconds + float64(delta), true
}

// printShards prints the per-shard pending-apply lines, one per shard, in the
// order the shards were decoded (the upstream projection already orders them by
// shard ID). The value is the published lagSeconds field: the contract
// precomputes each shard's lag (the same way /status does) and the CLI relays
// it, never recomputes — the same single-source-of-truth rule as the rest of the
// struct. This matters during a planned failover, where the agent measures each
// shard against its terminal boundary, so a naive ESTS-committedOptime
// subtraction would be wrong. A shard whose lagSeconds is omitted (not
// standby-initialized) prints "<shardId>: NA". When the standby itself is
// uninitialized (EDL nil) there is no shard data, so a single "<componentId>: NA"
// line names the one component encoded in the configured DR state file key.
func (o *LagOpts) printShards(status *standby.ClusterStatus, edl *float64) error {
	if edl == nil {
		_, err := fmt.Fprintf(o.OutWriter, "%s: NA\n", o.componentID())
		return err
	}
	for _, s := range status.ShardStatus {
		if s.LagSeconds == nil {
			if _, err := fmt.Fprintf(o.OutWriter, "%s: NA\n", s.ShardID); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(o.OutWriter, "%s: %ss of oplog pending apply\n", s.ShardID, formatLagSeconds(*s.LagSeconds)); err != nil {
			return err
		}
	}
	return nil
}

// printStandbyLag prints the cluster-wide standby lag line: EDL + (ESTS - MCTS)
// when computable, "not available yet" when the boundaries are unset or out of
// order, and NA when the standby is uninitialized (EDL nil).
func (o *LagOpts) printStandbyLag(status *standby.ClusterStatus, edl *float64) error {
	if edl == nil {
		_, err := fmt.Fprintln(o.OutWriter, "standby lag: NA")
		return err
	}
	if total, ok := standbyLagSeconds(status); ok {
		_, err := fmt.Fprintf(o.OutWriter, "standby lag: %ss\n", formatLagSeconds(total))
		return err
	}
	_, err := fmt.Fprintln(o.OutWriter, "standby lag: not available yet")
	return err
}

// componentID returns the component id named by the configured DR state file
// key. The CLI can only know the one component the operator configured.
func (o *LagOpts) componentID() string {
	return componentIDFromKey(o.stateFileKey())
}

// componentIDFromKey extracts the component id encoded in a DR state file key:
// "<prefix>/dr_status_<componentId>.json" -> "<componentId>". If the filename
// does not follow that convention it falls back to the filename without its
// extension, and to "unknown" when there is nothing usable.
func componentIDFromKey(key string) string {
	name := key
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimSuffix(name, ".json")
	name = strings.TrimPrefix(name, "dr_status_")
	if name == "" {
		return "unknown"
	}
	return name
}

// formatLagSeconds renders the estimated data loss as whole seconds when the
// value is exact ("40") and with a single decimal otherwise ("40.5").
func formatLagSeconds(seconds float64) string {
	if seconds == math.Trunc(seconds) {
		return strconv.FormatFloat(seconds, 'f', 0, 64)
	}
	return strconv.FormatFloat(seconds, 'f', 1, 64)
}

// mongocli ops-manager standby-clusters lag.
func LagBuilder() *cobra.Command {
	o := &LagOpts{}
	cmd := &cobra.Command{
		Use:   "lag",
		Short: "Report the replication lag of your standby cluster without Ops Manager.",
		Long: `Reads the Cluster Status Object (cluster_status_v1.bson) directly from the
object store and reports the standby replication lag. It doesn't use the Ops
Manager API, so it remains available when Ops Manager is unreachable.

The plaintext output always shows the same labeled lines, whatever the object
supports; only the values change (a number, short words, or NA):

  estimated data loss: the operations the active cluster committed since the
  standby's last safe point but that have not yet been shipped. They would be
  lost in an unplanned failover. It is real elapsed time measured since the
  safe point. It reads NA until the standby has initialized replication.

  With '--shards', each shard line reports the oplog the shard has received but
  not yet applied, as published by the agent. This is the drain work that
  promotion must complete before the shard can take over. A shard whose value
  the agent has not published yet reads NA.

  standby lag = estimated data loss + pending apply — how far behind the
  standby is: the newest op it has applied was originally committed that many
  seconds ago on the source. It reads "not available yet" when the boundaries
  are missing or out of order, and NA when the standby is uninitialized.

Use '-o json' for the full cluster status document. With '-o json', '--shards'
has no additional effect because shardStatus is already part of the document.

Run 'mongocli ops-manager standby-clusters configure' first to set up access
to the DR state file of your standby cluster.`,
		Args: require.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return o.PreRunE(
				o.initStore(cmd.Context()),
				o.InitOutput(cmd.OutOrStdout(), ""),
			)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return o.Run(cmd.Context())
		},
	}

	addSharedFlags(cmd, &o.opts)
	cmd.Flags().BoolVar(&o.shards, flag.Shards, false, usage.StandbyShards)
	cmd.Flags().StringVarP(&o.Output, flag.Output, flag.OutputShort, "",
		"Output format. Valid values are json — or leave unset for the labeled summary.")
	_ = cmd.RegisterFlagCompletionFunc(flag.Output, completeLagOutput)

	return cmd
}

// completeLagOutput advertises only the -o values this command accepts. The
// shared cli.AutoCompleteOutputFlag() offers every template format, which lag
// rejects at runtime, so completion is restricted here to 'json' (plus the
// implicit plaintext default of an unset -o).
func completeLagOutput(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return []string{jsonOutputFormat}, cobra.ShellCompDirectiveNoFileComp
}
