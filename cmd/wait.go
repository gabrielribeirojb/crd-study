package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/gabrielribeirojb/crd-study/internal/config"
	"github.com/gabrielribeirojb/crd-study/internal/controller"
	"github.com/gabrielribeirojb/crd-study/internal/kubeclient/httpclient"
)

var (
	waitFile     string
	waitTimeout  time.Duration
	waitInterval time.Duration
	waitVerbose  bool
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait until ClusterRestore reaches a terminal phase",
	RunE: func(cmd *cobra.Command, args []string) error {
		desired, err := config.LoadDesired(waitFile) // deve retornar state.DesiredSpec
		if err != nil {
			return err
		}

		client := httpclient.New(apiBaseURL, waitVerbose)
		reader := controller.NewHTTPReader(client)

		r := controller.NewReconciler(reader)

		ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
		defer cancel()

		fmt.Printf("Waiting for %s/%s (timeout=%s interval=%s)\n",
			desired.Namespace, desired.Name, waitTimeout, waitInterval)

		return r.Wait(ctx, desired, waitInterval)
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
	waitCmd.Flags().StringVarP(&waitFile, "file", "f", "desired.yaml", "Path to desired state file")
	waitCmd.Flags().DurationVar(&waitTimeout, "timeout", 30*time.Second, "Max time to wait")
	waitCmd.Flags().DurationVar(&waitInterval, "interval", 1*time.Second, "Polling interval")
	waitCmd.Flags().BoolVarP(&waitVerbose, "verbose", "v", false, "Verbose HTTP logs")
}
