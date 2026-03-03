package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gabrielribeirojb/crd-study/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a fake Kubernetes API server for ClusterRestore",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fake API listening on http://localhost:8080")
		return server.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
