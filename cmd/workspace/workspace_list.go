package workspace

import (
	"fmt"
	"os"

	"github.com/leg100/etok/cmd/flags"
	cmdutil "github.com/leg100/etok/cmd/util"
	"github.com/leg100/etok/pkg/env"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListCmd(opts *cmdutil.Options) *cobra.Command {
	var path, kubeContext string
	var namespace = "default"
	var workspace = "default"

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all workspaces",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			client, err := opts.Create(kubeContext)
			if err != nil {
				return err
			}

			etokenv, err := env.ReadEtokEnv(path)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				namespace = etokenv.Namespace()
				workspace = etokenv.Workspace()
			}

			// List across all namespaces
			workspaces, err := client.WorkspacesClient("").List(cmd.Context(), metav1.ListOptions{})
			if err != nil {
				return err
			}

			var prefix string
			for _, ws := range workspaces.Items {
				if ws.Namespace == namespace && ws.Name == workspace {
					prefix = "*"
				} else {
					prefix = ""
				}
				fmt.Fprintf(opts.Out, "%s\t%s/%s\n", prefix, ws.GetNamespace(), ws.GetName())
			}

			return nil
		},
	}

	flags.AddPathFlag(cmd, &path)
	flags.AddKubeContextFlag(cmd, &kubeContext)

	return cmd
}
