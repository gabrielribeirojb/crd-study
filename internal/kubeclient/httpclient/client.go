package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
	verbose bool
}

func New(baseURL string, verbose bool) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		verbose: verbose,
	}
}

// Retorno do GET/POST
type ClusterRestore struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`

	Spec struct {
		BackupRef string `json:"backupRef"`
	} `json:"spec"`

	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`

	CreatedAt string `json:"createdAt"`
}

func (c *Client) GetClusterRestore(ctx context.Context, ns, name string) (ClusterRestore, int, error) {
	url := fmt.Sprintf("%s/v1/namespaces/%s/clusterrestores/%s", c.baseURL, ns, name)
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ClusterRestore{}, 0, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return ClusterRestore{}, 0, err
	}
	defer resp.Body.Close()

	if c.verbose {
		fmt.Printf("HTTP GET %s -> %d (%s)\n", url, resp.StatusCode, time.Since(start))
	}

	// 404 é um "não existe", não é erro fatal
	if resp.StatusCode == http.StatusNotFound {
		return ClusterRestore{}, resp.StatusCode, nil
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return ClusterRestore{}, resp.StatusCode, fmt.Errorf("GET failed: %d: %s", resp.StatusCode, string(b))
	}

	var out ClusterRestore
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ClusterRestore{}, resp.StatusCode, err
	}
	return out, resp.StatusCode, nil
}

func (c *Client) CreateClusterRestore(ctx context.Context, ns, name, backupRef string) (ClusterRestore, int, error) {
	url := fmt.Sprintf("%s/v1/namespaces/%s/clusterrestores", c.baseURL, ns)
	start := time.Now()

	body := map[string]string{
		"name":      name,
		"backupRef": backupRef,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return ClusterRestore{}, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return ClusterRestore{}, 0, err
	}
	defer resp.Body.Close()

	if c.verbose {
		fmt.Printf("HTTP POST %s -> %d (%s)\n", url, resp.StatusCode, time.Since(start))
	}

	if resp.StatusCode != http.StatusCreated {
		rb, _ := io.ReadAll(resp.Body)
		return ClusterRestore{}, resp.StatusCode, fmt.Errorf("POST failed: %d: %s", resp.StatusCode, string(rb))
	}

	var out ClusterRestore
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ClusterRestore{}, resp.StatusCode, err
	}
	return out, resp.StatusCode, nil
}
