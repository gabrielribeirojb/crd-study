package controller

import (
	"context"

	"github.com/gabrielribeirojb/crd-study/internal/kubeclient/httpclient"
	"github.com/gabrielribeirojb/crd-study/internal/state"
)

type HTTPReader struct {
	client *httpclient.Client
}

func NewHTTPReader(c *httpclient.Client) *HTTPReader {
	return &HTTPReader{client: c}
}

func (r *HTTPReader) Get(ctx context.Context, namespace, name string) (state.CurrentState, error) {
	cr, status, err := r.client.GetClusterRestore(ctx, namespace, name)
	if err != nil {
		return state.CurrentState{}, err
	}

	// 404 = não existe (estado atual)
	if status == 404 {
		return state.CurrentState{Exists: false, Phase: ""}, nil
	}

	// sucesso
	return state.CurrentState{
		Exists: true,
		Phase:  cr.Status.Phase,
	}, nil
}
