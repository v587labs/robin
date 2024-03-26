package rcobra

import (
	"github.com/spf13/cobra"
	"slices"
)

type runFn func(cmd *cobra.Command, args []string)

func (f runFn) RunE(cmd *cobra.Command, args []string) error {
	f(cmd, args)
	return nil
}

type runEs []func(cmd *cobra.Command, args []string) error

func (es runEs) runE(cmd *cobra.Command, args []string) error {
	for _, e := range es {
		if e == nil {
			continue
		}
		err := e(cmd, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es runEs) run(cmd *cobra.Command, args []string) {
	_ = es.runE(cmd, args)
}

func AddCommands(root *cobra.Command, commands ...*cobra.Command) {
	// parent per method
	preRunES := runEs{}
	postRunES := runEs{}

	for p := root; p != nil; p = p.Parent() {
		if p.PersistentPreRunE != nil {
			preRunES = append(preRunES, p.PersistentPreRunE)
		} else if p.PersistentPreRun != nil {
			preRunES = append(preRunES, runFn(p.PersistentPreRun).RunE)
		}
		if p.PersistentPostRunE != nil {
			postRunES = append(postRunES, p.PersistentPostRunE)
		} else if p.PersistentPostRun != nil {
			postRunES = append(postRunES, runFn(p.PersistentPostRun).RunE)
		}
	}

	slices.Reverse(preRunES)
	slices.Reverse(postRunES)

	for i, _ := range commands {
		if commands[i].PersistentPreRunE != nil {
			commands[i].PersistentPreRunE = append(preRunES, commands[i].PersistentPreRunE).runE
		} else if commands[i].PersistentPreRun != nil {
			commands[i].PersistentPreRun = append(preRunES, runFn(commands[i].PersistentPreRun).RunE).run
		}
		if commands[i].PersistentPostRunE != nil {
			commands[i].PersistentPostRunE = append(postRunES, commands[i].PersistentPostRunE).runE
		} else if commands[i].PersistentPostRun != nil {
			commands[i].PersistentPostRun = append(postRunES, runFn(commands[i].PersistentPostRun).RunE).run
		}
	}
	root.AddCommand(commands...)
}
