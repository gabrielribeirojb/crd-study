package real

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientset cria um client-go "typed" (clientset) igual o kubectl usa.
// - kubeconfigPath: caminho pro arquivo kubeconfig (ou vazio pra default ~/.kube/config)
// - kubeContext: nome do contexto (ou vazio pra current-context)
func NewClientset(kubeconfigPath, kubeContext string) (*kubernetes.Clientset, error) {
	cfg, err := buildConfig(kubeconfigPath, kubeContext)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

// buildConfig centraliza como montamos a *rest.Config a partir do kubeconfig.
func buildConfig(kubeconfigPath, kubeContext string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		home, _ := os.UserHomeDir()
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	overrides := &clientcmd.ConfigOverrides{}
	if kubeContext != "" {
		overrides.CurrentContext = kubeContext
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
}
