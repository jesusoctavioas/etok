package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/apex/log"
	"github.com/leg100/stok/api/command"
	"github.com/leg100/stok/api/v1alpha1"
	"github.com/leg100/stok/pkg/k8s"
	"github.com/leg100/stok/util"
	"github.com/leg100/stok/util/slice"
	"github.com/spf13/cobra"
)

type runnerCmd struct {
	Path    string
	Tarball string

	Name      string
	Namespace string
	Kind      string
	Timeout   time.Duration
	Context   string
	NoWait    bool

	factory k8s.FactoryInterface
	args    []string
	cmd     *cobra.Command
}

func newRunnerCmd(f k8s.FactoryInterface) *cobra.Command {
	runner := &runnerCmd{}

	cmd := &cobra.Command{
		// TODO: what is the syntax for stating at least one command must be provided?
		Use:           "runner [command (args)]",
		Short:         "Run the stok runner",
		Long:          "The stok runner is intended to be run in on pod, started by the relevant stok command controller. When invoked, it extracts a tarball containing terraform configuration files. It then waits for the command's ClientReady condition to be true. And then it invokes the relevant command, typically a terraform command.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.doRunnerCmd(args); err != nil {
				return fmt.Errorf("runner: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&runner.Path, "path", ".", "Workspace config path")
	cmd.Flags().StringVar(&runner.Tarball, "tarball", "", "Extract specified tarball file to workspace path")

	cmd.Flags().BoolVar(&runner.NoWait, "no-wait", false, "Disable polling resource for client annotation")
	cmd.Flags().StringVar(&runner.Name, "name", "", "Name of command resource")
	cmd.Flags().StringVar(&runner.Namespace, "namespace", "default", "Namespace of command resource")
	cmd.Flags().StringVar(&runner.Kind, "kind", "", "Kind of command resource")
	cmd.Flags().StringVar(&runner.Context, "context", "", "Set kube context (defaults to kubeconfig current context)")
	cmd.Flags().DurationVar(&runner.Timeout, "timeout", 10*time.Second, "Timeout on waiting for client to confirm readiness")

	runner.factory = f
	runner.cmd = cmd
	return runner.cmd
}

func (r *runnerCmd) doRunnerCmd(args []string) error {
	if err := r.validate(); err != nil {
		return err
	}

	r.args = args

	if r.Tarball != "" {
		files, err := extractTarball(r.Tarball, r.Path)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{"files": files, "path": r.Path}).Debug("extracted tarball")
	}

	if !r.NoWait {
		rc, err := r.factory.NewClient(r.Context)
		if err != nil {
			return err
		}

		if err := r.handleSemaphore(rc, r.Kind, r.Name, r.Namespace, r.Timeout); err != nil {
			return err
		}
	}

	if err := r.run(os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}

func (r *runnerCmd) validate() error {
	if r.Kind == "" {
		return fmt.Errorf("missing flag: --kind <kind>")
	}

	if !slice.ContainsString(append(command.CommandKinds, "Workspace"), r.Kind) {
		return fmt.Errorf("invalid kind: %s", r.Kind)
	}

	return nil
}

// Extract Tarball with path 'src' to path 'dst'
func extractTarball(src, dst string) (int, error) {
	tarBytesBuffer := new(bytes.Buffer)
	tarBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return 0, err
	}

	_, err = tarBytesBuffer.Write(tarBytes)
	if err != nil {
		return 0, err
	}

	return util.Extract(tarBytesBuffer, dst)
}

func (r *runnerCmd) handleSemaphore(rc k8s.Client, kind, name, namespace string, timeout time.Duration) error {
	cache, err := rc.NewCache(r.Namespace)
	if err != nil {
		return err
	}

	informer, err := cache.GetInformerForKind(context.TODO(), v1alpha1.SchemeGroupVersion.WithKind(kind))
	if err != nil {
		return err
	}

	events := r.factory.GetEventsChannel()
	informer.AddEventHandler(k8s.EventHandlers(events))

	done := make(chan error)
	stop := make(chan struct{})
	go func() {
		if err := cache.Start(stop); err != nil {
			done <- err
		}
	}()

	// Wait until WaitAnnotation is unset, or return error if timeout or other error occurs
	timer := time.NewTimer(timeout)
	for {
		select {
		case e := <-events:
			switch obj := e.(type) {
			case command.Interface:
				if obj.GetName() != name || obj.GetNamespace() != namespace {
					// Ignore event for a different cmd
					break
				}

				if _, ok := obj.GetAnnotations()[v1alpha1.WaitAnnotationKey]; !ok {
					// No checkpoint annotation set, we're clear to go
					// Stop timer
					if !timer.Stop() {
						<-timer.C
					}
					// Stop cache
					close(stop)
					return nil
				}
			}
		case <-timer.C:
			return fmt.Errorf("timeout exceeded waiting for client hold to be released")
		case err = <-done:
			return err
		}
	}
}

// Run args, taking first arg as executable, and remainder as args to executable. Path sets the
// working directory of the executable; out and errout set stdout and stderr of executable.
func (r *runnerCmd) run(out, errout io.Writer) error {
	args := command.RunnerArgsForKind(r.Kind, r.args)

	log.WithFields(log.Fields{"command": args[0], "args": args[1:]}).Debug("running command")

	exe := exec.Command(args[0], args[1:]...)
	exe.Dir = r.Path
	exe.Stdin = os.Stdin
	exe.Stdout = out
	exe.Stderr = errout

	return exe.Run()
}
