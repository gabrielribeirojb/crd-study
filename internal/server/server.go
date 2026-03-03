package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

// store = nosso "banco" em memória, protegido por mutex (mu)
type store struct {
	mu   sync.Mutex
	data map[string]ClusterRestore // chave: "namespace/name"
}

func newStore() *store {
	return &store{data: make(map[string]ClusterRestore)}
}

func key(ns, name string) string {
	return ns + "/" + name
}

// Run sobe o servidor HTTP e bloqueia o processo (fica rodando)
func Run(addr string) error {
	s := newStore()
	mux := http.NewServeMux()

	// healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// GET /v1/namespaces/{ns}/clusterrestores/{name}
	// POST /v1/namespaces/{ns}/clusterrestores
	mux.HandleFunc("/v1/namespaces/", func(w http.ResponseWriter, r *http.Request) {
		// Exemplo de path:
		// /v1/namespaces/demo/clusterrestores/r1
		// /v1/namespaces/demo/clusterrestores
		path := strings.TrimPrefix(r.URL.Path, "/v1/namespaces/")
		parts := strings.Split(path, "/")

		// parts[0] = namespace
		if len(parts) < 2 {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		ns := parts[0]
		resource := parts[1]

		if resource != "clusterrestores" {
			http.Error(w, "unknown resource", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// precisa ter /{name}
			if len(parts) != 3 {
				http.Error(w, "missing name", http.StatusBadRequest)
				return
			}
			name := parts[2]
			handleGetClusterRestore(w, s, ns, name)

		case http.MethodPost:
			// POST não tem {name} no path
			if len(parts) != 2 {
				http.Error(w, "invalid path for POST", http.StatusBadRequest)
				return
			}
			handleCreateClusterRestore(w, r, s, ns)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("listening on", addr)
	return http.ListenAndServe(addr, mux)
}

func handleGetClusterRestore(w http.ResponseWriter, s *store, ns, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cr, ok := s.data[key(ns, name)]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cr)
}

func handleCreateClusterRestore(w http.ResponseWriter, r *http.Request, s *store, ns string) {
	// body esperado:
	// {"name":"r1","backupRef":"bkp-123"}
	type createReq struct {
		Name      string `json:"name"`
		BackupRef string `json:"backupRef"`
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.BackupRef == "" {
		http.Error(w, "name and backupRef are required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	k := key(ns, req.Name)
	if _, exists := s.data[k]; exists {
		http.Error(w, "already exists", http.StatusConflict)
		return
	}

	cr := ClusterRestore{
		Namespace: ns,
		Name:      req.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	cr.Spec.BackupRef = req.BackupRef
	cr.Status.Phase = "PENDING"

	s.data[k] = cr

	go func(key string) {
		time.Sleep(2 * time.Second)

		s.mu.Lock()
		cr := s.data[key]
		cr.Status.Phase = "RUNNING"
		s.data[key] = cr
		s.mu.Unlock()

		time.Sleep(2 * time.Second)

		s.mu.Lock()
		cr = s.data[key]
		cr.Status.Phase = "SUCCEEDED"
		s.data[key] = cr
		s.mu.Unlock()
	}(k)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cr)
}
