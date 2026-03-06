package cmd

import (
	"github.com/spf13/cobra"
)

// Fake API (continua existindo)
var apiBaseURL string

// Novas flags globais para Kubernetes real
var (
	kubeconfigPath string
	kubeContext    string
)

// rootCmd é o "comando base".
var rootCmd = &cobra.Command{
	Use:   "crd-study",
	Short: "CRD study CLI",
}

func init() {
	// Fake API
	rootCmd.PersistentFlags().StringVar(
		&apiBaseURL,
		"api",
		"http://localhost:8080",
		"Base URL of fake API server",
	)

	// NOVAS FLAGS (Kubernetes real)
	rootCmd.PersistentFlags().StringVar(
		&kubeconfigPath,
		"kubeconfig",
		"",
		"Path to kubeconfig file (default: ~/.kube/config)",
	)

	rootCmd.PersistentFlags().StringVar(
		&kubeContext,
		"context",
		"",
		"Kubeconfig context to use (default: current-context)",
	)
}
