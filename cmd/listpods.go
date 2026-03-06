package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/gabrielribeirojb/crd-study/internal/kubeclient/real"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	listPodsNamespace string
)

var listPodsCmd = &cobra.Command{
	Use:   "listpods [podName]",
	Short: "List pods from the Kubernetes cluster (with autocomplete)",
	Args:  cobra.MaximumNArgs(1),

	// Autocomplete do argumento [podName]
	// O shell chama isso quando você aperta TAB.
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Se já tem 1 arg, não sugere mais nada.
		if len(args) >= 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// PEGA O NAMESPACE DO FLAG (respeita o que o usuário digitou)
		ns, _ := cmd.Flags().GetString("namespace")
		if ns == "" {
			ns = "default"
		}

		names, err := fetchPodNames(ns)
		if err != nil {
			// Em autocomplete, se der erro, só não sugere nada.
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Filtra por prefixo digitado (myp<TAB>)
		out := make([]string, 0, len(names))
		for _, n := range names {
			if strings.HasPrefix(n, toComplete) {
				out = append(out, n)
			}
		}

		sort.Strings(out)
		return out, cobra.ShellCompDirectiveNoFileComp
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientset, err := real.NewClientset(kubeconfigPath, kubeContext)
		if err != nil {
			return err
		}

		ns, _ := cmd.Flags().GetString("namespace")
		if ns == "" {
			ns = "default"
		}

		podList, err := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})

		if err != nil {
			return err
		}

		// Se veio [podName], filtra e mostra só ele
		var filter string
		if len(args) == 1 {
			filter = args[0]
		}

		for _, p := range podList.Items {
			if filter != "" && p.Name != filter {
				continue
			}
			fmt.Printf("%s\n", p.Name)
		}

		return nil
	},
}

// Função auxiliar usada pelo autocomplete
func fetchPodNames(ns string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientset, err := real.NewClientset(kubeconfigPath, kubeContext)
	if err != nil {
		return nil, err
	}

	podList, err := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(podList.Items))
	for _, p := range podList.Items {
		out = append(out, p.Name)
	}
	sort.Strings(out)
	return out, nil
}

func init() {
	// adiciona no root: `crd-study listpods`
	rootCmd.AddCommand(listPodsCmd)

	listPodsCmd.Flags().StringVarP(&listPodsNamespace, "namespace", "n", "default", "Namespace to query")
}
