package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/gabrielribeirojb/crd-study/internal/kubeclient/real"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var namespace string

var allNamespaces bool

var output string

var listClusterRestoresCmd = &cobra.Command{
	Use:   "clusterrestores",
	Short: "List ClusterRestore resources from the Kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := real.NewDynamicClient(kubeconfigPath, kubeContext)
		if err != nil {
			return err
		}

		// Esse GVR identifica o recurso na API:
		// group=gdch.mycompany.io, version=v1, resource=clusterrestores
		gvr := schema.GroupVersionResource{
			Group:    "gdch.mycompany.io",
			Version:  "v1",
			Resource: "clusterrestores",
		}

		ctx := context.Background()

		nsQuery := namespace
		if allNamespaces {
			nsQuery = v1.NamespaceAll
		}

		list, err := c.Resource(gvr).Namespace(nsQuery).List(ctx, v1.ListOptions{})
		if err != nil {
			return err
		}

		switch strings.ToLower(output) {
		case "json":
			b, err := json.MarshalIndent(list, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil

		case "yaml", "yml":
			// Faz JSON -> YAML pra garantir que "items" saia certinho
			jb, err := json.Marshal(list)
			if err != nil {
				return err
			}
			yb, err := yaml.JSONToYAML(jb)
			if err != nil {
				return err
			}
			fmt.Print(string(yb))
			return nil
		}

		if strings.ToLower(output) == "wide" {
			fmt.Printf("%-12s %-20s %-6s %-25s %-12s\n", "NAMESPACE", "NAME", "AGE", "CREATED", "BACKUPREF")
		} else {
			fmt.Printf("%-12s %-20s %-6s %-12s\n", "NAMESPACE", "NAME", "AGE", "BACKUPREF")
		}

		for _, item := range list.Items {
			ns := item.GetNamespace()
			if ns == "" {
				ns = namespace // fallback quando não vem (raro)
			}
			name := item.GetName()

			// AGE
			createdT := item.GetCreationTimestamp().Time
			created := createdT.Local().Format("2006-01-02 15:04:05 -0700")
			age := humanDuration(time.Since(createdT))

			// spec.backupRef
			spec, _ := item.Object["spec"].(map[string]any)
			backupRef, _ := spec["backupRef"].(string)

			if strings.ToLower(output) == "wide" {
				fmt.Printf("%-12s %-20s %-6s %-25s %-12s\n", ns, name, age, created, backupRef)
			} else {
				fmt.Printf("%-12s %-20s %-6s %-12s\n", ns, name, age, backupRef)
			}
		}

		return nil
	},
}

func humanDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

func init() {
	listClusterRestoresCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to query")
	listClusterRestoresCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list across all namespaces")
	listClusterRestoresCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: wide|yaml|json")
	listCmd.AddCommand(listClusterRestoresCmd)
}
