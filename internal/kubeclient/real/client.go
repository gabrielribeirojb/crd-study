package real

import (
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func buildRESTConfig(kubeconfigPath string, kubeContext string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{
		ExplicitPath: kubeconfigPath,
	}

	overrides := &clientcmd.ConfigOverrides{}
	if kubeContext != "" {
		overrides.CurrentContext = kubeContext
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		overrides,
	).ClientConfig()
}

func NewDynamicClient(kubeconfigPath string, kubeContext string) (dynamic.Interface, error) {
	cfg, err := buildRESTConfig(kubeconfigPath, kubeContext)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(cfg)
}
