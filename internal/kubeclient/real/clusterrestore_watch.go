package real

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

func WaitClusterRestoreTerminalPhase(
	ctx context.Context,
	kubeconfigPath string,
	kubeContext string,
	namespace string,
	name string,
	verbose bool,
) (string, error) {
	c, err := NewDynamicClient(kubeconfigPath, kubeContext)
	if err != nil {
		return "", err
	}

	gvr := schema.GroupVersionResource{
		Group:    "gdch.mycompany.io",
		Version:  "v1",
		Resource: "clusterrestores",
	}

	opts := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	}

	u, err := c.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		phase, terminal := terminalPhaseFromObject(u.Object)
		if verbose {
			fmt.Printf("Initial GET phase=%q\n", phase)
		}
		if terminal {
			return phase, nil
		}
	}

	w, err := c.Resource(gvr).Namespace(namespace).Watch(ctx, opts)
	if err != nil {
		return "", err
	}
	defer w.Stop()

	if verbose {
		fmt.Printf("Watching %s/%s...\n", namespace, name)
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()

		case evt, ok := <-w.ResultChan():
			if !ok {
				return "", fmt.Errorf("watch channel closed unexpectedly")
			}

			// só processa eventos relevantes
			if evt.Type != watch.Added && evt.Type != watch.Modified {
				continue
			}

			u, ok := evt.Object.(*unstructured.Unstructured)
			if !ok || u == nil {
				continue
			}

			phase, terminal := terminalPhaseFromObject(u.Object)

			if verbose {
				fmt.Printf("Event=%s phase=%q\n", evt.Type, phase)
			}

			if terminal {
				return phase, nil
			}
		}
	}
}

func terminalPhaseFromObject(obj map[string]any) (phase string, terminal bool) {
	status, ok := obj["status"].(map[string]any)
	if !ok {
		return "", false
	}

	phase, _ = status["phase"].(string)

	switch phase {
	case "SUCCEEDED", "FAILED":
		return phase, true
	default:
		return phase, false
	}
}
