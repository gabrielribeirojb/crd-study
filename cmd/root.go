package cmd

import "github.com/spf13/cobra"

// apiBaseURL é global ao pacote cmd
// Todos os comandos (serve, apply, etc) podem usar
var apiBaseURL string

// rootCmd é o "comando base".
var rootCmd = &cobra.Command{
	Use:   "crd-study",
	Short: "CRD study CLI",
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&apiBaseURL,
		"api",
		"http://localhost:8080",
		"Base URL of fake API server",
	)
}
