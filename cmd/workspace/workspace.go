package workspace

import (
	cmdutil "github.com/leg100/stok/cmd/util"
	"github.com/spf13/cobra"
)

func WorkspaceCmd(opts *cmdutil.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Stok workspace management",
	}

	newCmd, _ := NewCmd(opts)
	cmd.AddCommand(newCmd)

	cmd.AddCommand(
		ListCmd(opts),
		DeleteCmd(opts),
		ShowCmd(opts),
		SelectCmd(opts),
		DeleteCmd(opts),
	)

	return cmd
}
