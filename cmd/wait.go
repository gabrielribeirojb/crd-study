package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/gabrielribeirojb/crd-study/internal/config"
	"github.com/gabrielribeirojb/crd-study/internal/kubeclient/real"
)

var (
	waitFile    string
	waitTimeout time.Duration
	waitVerbose bool
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait until ClusterRestore reaches a terminal phase (WATCH, no polling)",
	RunE: func(cmd *cobra.Command, args []string) error {
		desired, err := config.LoadDesired(waitFile)
		if err != nil {
			return err
		}

		if desired.Namespace == "" || desired.Name == "" {
			return fmt.Errorf("desired.yaml precisa ter namespace e name (namespace=%q name=%q)", desired.Namespace, desired.Name)
		}

		ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
		defer cancel()

		fmt.Printf("Waiting (watch) for %s/%s (timeout=%s)\n", desired.Namespace, desired.Name, waitTimeout)

		phase, err := real.WaitClusterRestoreTerminalPhase(
			ctx,
			kubeconfigPath,
			kubeContext,
			desired.Namespace,
			desired.Name,
			waitVerbose,
		)
		if err != nil {
			return err
		}

		if phase == "SUCCEEDED" {
			fmt.Println("Done: restore succeeded")
			return nil
		}

		return fmt.Errorf("restore terminou com phase=%s", phase)
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
	waitCmd.Flags().StringVarP(&waitFile, "file", "f", "desired.yaml", "Path to desired state file")
	waitCmd.Flags().DurationVar(&waitTimeout, "timeout", 30*time.Second, "Max time to wait")
	waitCmd.Flags().BoolVarP(&waitVerbose, "verbose", "v", false, "Verbose logs")
}
