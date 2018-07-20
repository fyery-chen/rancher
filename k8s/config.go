package k8s

import (
	"context"
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfig(ctx context.Context, k8sMode string, kubeConfig string) (bool, context.Context, *rest.Config, error) {
	var (
		cfg *rest.Config
		err error
	)

	switch k8sMode {
	case "external":
		cfg, err = getExternal(kubeConfig)
	default:
		return false, nil, nil, fmt.Errorf("invalid k8s-mode %s", k8sMode)
	}

	return false, ctx, cfg, err
}

func getExternal(kubeConfig string) (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", kubeConfig)
}
