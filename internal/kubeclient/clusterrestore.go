package kubeclient

import (
	"context"

	gdchv1 "github.com/gabrielribeirojb/crd-study/pkg/apis/gdch/v1"
)

type ClusterRestoreClient interface {
	Get(ctx context.Context, namespace, name string) (*gdchv1.ClusterRestore, error)
	List(ctx context.Context, namespace string) ([]gdchv1.ClusterRestore, error)
	Create(ctx context.Context, namespace string, obj *gdchv1.ClusterRestore) error
	UpdateStatusPhase(ctx context.Context, namespace, name, phase string) error
}
