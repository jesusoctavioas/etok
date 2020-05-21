// Code generated by go generate; DO NOT EDIT.
package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/leg100/stok/pkg/apis/stok/v1alpha1"
	"github.com/spf13/cobra"
)

var cmdImport = &cobra.Command{
	Use:   "import [global flags] -- [import args]",
	Short: "Run terraform import",
	Run: func(cmd *cobra.Command, args []string) {
		app := newApp("import", DoubleDashArgsHandler(os.Args))
		if err := app.run(&v1alpha1.Imp{}); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cmdImport.DisableFlagsInUseLine = true
	rootCmd.AddCommand(cmdImport)
}
