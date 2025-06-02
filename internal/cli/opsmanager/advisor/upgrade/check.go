package upgrade

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/fatih/color"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/store"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"go.mongodb.org/ops-manager/opsmngr"
)

type CheckOpts struct {
	cli.OutputOpts
	targetVersion string
	store         store.AdvisorUpgrade
}

func (opts *CheckOpts) initStore(ctx context.Context) func() error {
	return func() error {
		var err error
		opts.store, err = store.New(store.AuthenticatedPreset(config.Default()), store.WithContext(ctx))
		return err
	}
}

func (opts *CheckOpts) Run() error {
	steps, err := opts.store.AdvisorUpgradeCheck(opts.targetVersion)
	if err != nil {
		return err
	}

	printFullWidthWarning([]string{
		"IMPORTANT!",
		"You must execute the upgrade IN ORDER.",
		"",
		"• Work step-by-step: finish Step 1 before starting Step 2,",
		"  finish Step 2 before starting Step 3, and so on.",
		"• Inside each step follow the sequence shown in the table",
		"  (Operating-System  →  MongoDB  →  Ops Manager  →  Automation Agent).",
		"• Do NOT skip, reorder, or run steps in parallel.",
	})

	for i, s := range *steps {
		renderCheckOutput(i+1, &s)
	}

	return nil
}

// mongocli opsmanager advisor upgrade check --targetVersion version.
func CheckBuilder() *cobra.Command {
	opts := &CheckOpts{}
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Checks if your system meets the requirements to upgrade to a certain Ops Manager version.",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			opts.OutWriter = cmd.OutOrStdout()
			return opts.initStore(cmd.Context())()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.targetVersion, flag.TargetVersion, "", usage.TargetVersion)

	return cmd
}

func renderCheckOutput(n int, step *opsmngr.UpgradeCheckStep) {

	println("\n", color.New(color.FgYellow, color.Bold).Sprintf("Step %d", n))
	section := func(s string) string { return color.New(color.FgCyan, color.Bold).Sprint(s) }
	current := func(s string) string { return color.New(color.FgHiRed).Sprint(s) }
	target := func(s string) string { return color.New(color.FgHiGreen).Sprint(s) }
	multi := func(h []string) string { return strings.Join(h, "\n") }

	tbl := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{Formatting: tw.CellFormatting{
				MergeMode: tw.MergeBoth, Alignment: tw.AlignLeft,
			}},
		}),
	)
	arrow := "→"

	printOpsManager := func(blocks opsmngr.OpsManagerStep) {
		tbl.Append([]string{section("Ops Manager"), section("Ops Manager")})
		tbl.Append([]string{"Version:", fmt.Sprintf("%s %s %s", current(step.OpsManager.CurrentVersion), arrow, target(step.OpsManager.TargetVersion))})
		tbl.Append([]string{color.New(color.FgMagenta, color.Bold).Sprint("Hosts Need to be upgraded:"), multi(step.OpsManager.Hosts)})
	}

	printOS := func(blocks []opsmngr.OperatingSystemStep) {
		if len(blocks) == 0 {
			return
		}
		for i, b := range blocks {
			if i == 0 {
				tbl.Append([]string{section("OS"), section("OS")})
			}
			if b.BaseVersion != "" {
				tbl.Append([]string{"Base Version:", b.BaseVersion})
			}
			tbl.Append([]string{"Version:", fmt.Sprintf("%s %s %s", current(b.CurrentVersion), arrow, target(b.TargetVersion))})
			tbl.Append([]string{color.New(color.FgMagenta, color.Bold).Sprint("Hosts Need to be upgraded:"), multi(b.Hosts)})
		}
	}
	printMongo := func(blocks []opsmngr.MongoDBStep) {
		if len(blocks) == 0 {
			return
		}
		for i, b := range blocks {
			if i == 0 {
				tbl.Append([]string{section("MongoDB"), section("MongoDB")})
			}
			tbl.Append([]string{"Version:", fmt.Sprintf("%s %s %s", current(b.CurrentVersion), arrow, target(b.TargetVersion))})
			tbl.Append([]string{color.New(color.FgMagenta, color.Bold).Sprint("Hosts Need to be upgraded:"), multi(b.Hosts)})
		}
	}
	printAgent := func(blocks []opsmngr.AgentStep) {
		if len(blocks) == 0 {
			return
		}
		for i, b := range blocks {
			if i == 0 {
				tbl.Append([]string{section("Automation Agent"), section("Automation Agent")})
			}
			tbl.Append([]string{"Version:", fmt.Sprintf("%s %s %s", current(b.CurrentVersion), arrow, target(b.TargetVersion))})
			tbl.Append([]string{color.New(color.FgMagenta, color.Bold).Sprint("Hosts Need to be upgraded:"), multi(b.Hosts)})
		}
	}

	printOS(step.OperatingSystem)
	printMongo(step.MongoDB)
	printOpsManager(step.OpsManager)
	printAgent(step.Agent)

	tbl.Render()
}

func printFullWidthWarning(lines []string) {
	body := strings.Join(lines, "\n")

	// 3. configure and print the box
	b := box.New(box.Config{
		Px:           1,
		Py:           1,
		Type:         "Bold",
		Color:        "Yellow",
		TitlePos:     "Inside",
		TitleColor:   "Yellow",
		ContentAlign: "Left",
	})
	b.Print(" Notice ", body)
	fmt.Println() // blank line after the box
}
