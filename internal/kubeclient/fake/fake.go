package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/gabrielribeirojb/crd-study/internal/kubeclient"
	gdchv1 "github.com/gabrielribeirojb/crd-study/pkg/apis/gdch/v1"
)

type Client struct {
	mu sync.Mutex
	db map[string]gdchv1.ClusterRestore
}

func New() kubeclient.ClusterRestoreClient {
	return &Client{
		db: make(map[string]gdchv1.ClusterRestore),
	}
}

func key(ns, name string) string {
	return ns + "/" + name
}

func (c *Client) Get(ctx context.Context, namespace, name string) (*gdchv1.ClusterRestore, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	obj, ok := c.db[key(namespace, name)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	copy := obj // cópia rasa ok aqui pq ainda não temos maps/slices próprios no spec/status
	return &copy, nil
}

func (c *Client) List(ctx context.Context, namespace string) ([]gdchv1.ClusterRestore, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]gdchv1.ClusterRestore, 0)
	prefix := namespace + "/"
	for k, v := range c.db {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out, nil
}

func (c *Client) Create(ctx context.Context, namespace string, obj *gdchv1.ClusterRestore) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	k := key(namespace, obj.Name)
	if _, exists := c.db[k]; exists {
		return fmt.Errorf("already exists")
	}

	clone := *obj
	clone.Namespace = namespace
	c.db[k] = clone
	return nil
}

func (c *Client) UpdateStatusPhase(ctx context.Context, namespace, name, phase string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	k := key(namespace, name)
	obj, ok := c.db[k]
	if !ok {
		return fmt.Errorf("not found")
	}

	obj.Status.Phase = phase
	c.db[k] = obj
	return nil
}
