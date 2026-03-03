package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/gabrielribeirojb/crd-study/internal/kubeclient/httpclient"
	gdchv1 "github.com/gabrielribeirojb/crd-study/pkg/apis/gdch/v1"
)

var applyFile string
var verbose bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a ClusterRestore desired state (create if missing)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) ler YAML
		b, err := os.ReadFile(applyFile)
		if err != nil {
			return err
		}

		var desired gdchv1.ClusterRestore
		if err := yaml.Unmarshal(b, &desired); err != nil {
			return err
		}

		ns := desired.Namespace
		if ns == "" {
			ns = "default"
		}
		name := desired.Name
		if name == "" {
			return fmt.Errorf("metadata.name is required")
		}
		backupRef := desired.Spec.BackupRef
		if backupRef == "" {
			return fmt.Errorf("spec.backupRef is required")
		}

		// 2) criar HTTP client usando flag --api (apiBaseURL)
		c := httpclient.New(apiBaseURL, verbose)

		// 3) GET (existe?)
		_, status, err := c.GetClusterRestore(context.Background(), ns, name)
		if err != nil {
			return err
		}

		// 4) se não existe, cria
		if status == 404 {
			_, _, err := c.CreateClusterRestore(context.Background(), ns, name, backupRef)
			if err != nil {
				return err
			}
			fmt.Printf("APPLIED: CREATE %s/%s\n", ns, name)
			return nil
		}

		// 5) existe → noop
		fmt.Printf("APPLIED: NOOP %s/%s already exists\n", ns, name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "restore.yaml", "Path to ClusterRestore YAML")
	applyCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose HTTP logs")
}
