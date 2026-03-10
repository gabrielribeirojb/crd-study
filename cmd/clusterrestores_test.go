package cmd

import (
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func makeClusterRestore(name, namespace, backupRef string, created time.Time) unstructured.Unstructured {
	obj := unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "gdch.mycompany.io/v1",
			"kind":       "ClusterRestore",
			"metadata": map[string]any{
				"name":              name,
				"namespace":         namespace,
				"creationTimestamp": created.Format(time.RFC3339),
			},
			"spec": map[string]any{
				"backupRef": backupRef,
			},
		},
	}

	obj.SetName(name)
	obj.SetNamespace(namespace)
	obj.SetCreationTimestamp(metav1.NewTime(created))

	return obj
}

func TestRenderClusterRestoreTable_Default(t *testing.T) {
	now := time.Now().Add(-2 * time.Hour)

	list := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			makeClusterRestore("r1", "default", "backup-123", now),
		},
	}

	out, err := renderClusterRestoreTable(list, "default", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(out, "NAMESPACE") {
		t.Fatalf("expected header NAMESPACE in output, got:\n%s", out)
	}
	if !strings.Contains(out, "NAME") {
		t.Fatalf("expected header NAME in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BACKUPREF") {
		t.Fatalf("expected header BACKUPREF in output, got:\n%s", out)
	}
	if !strings.Contains(out, "default") {
		t.Fatalf("expected namespace default in output, got:\n%s", out)
	}
	if !strings.Contains(out, "r1") {
		t.Fatalf("expected resource name r1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "backup-123") {
		t.Fatalf("expected backupRef backup-123 in output, got:\n%s", out)
	}
}

func TestRenderClusterRestoreTable_Wide(t *testing.T) {
	now := time.Now().Add(-2 * time.Hour)

	list := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			makeClusterRestore("r1", "default", "backup-123", now),
		},
	}

	out, err := renderClusterRestoreTable(list, "default", "wide")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(out, "CREATED") {
		t.Fatalf("expected header CREATED in output, got:\n%s", out)
	}
	if !strings.Contains(out, "backup-123") {
		t.Fatalf("expected backupRef backup-123 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "r1") {
		t.Fatalf("expected resource name r1 in output, got:\n%s", out)
	}
}

func TestRenderClusterRestoreJSON(t *testing.T) {
	now := time.Now().Add(-2 * time.Hour)

	list := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			makeClusterRestore("r1", "default", "backup-123", now),
		},
	}

	out, err := renderClusterRestoreJSON(list)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(out, `"items"`) {
		t.Fatalf("expected items in json output, got:\n%s", out)
	}
	if !strings.Contains(out, `"name": "r1"`) {
		t.Fatalf("expected name r1 in json output, got:\n%s", out)
	}
	if !strings.Contains(out, `"backupRef": "backup-123"`) {
		t.Fatalf("expected backupRef backup-123 in json output, got:\n%s", out)
	}
}

func TestRenderClusterRestoreYAML(t *testing.T) {
	now := time.Now().Add(-2 * time.Hour)

	list := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			makeClusterRestore("r1", "default", "backup-123", now),
		},
	}

	out, err := renderClusterRestoreYAML(list)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(out, "items:") {
		t.Fatalf("expected items in yaml output, got:\n%s", out)
	}
	if !strings.Contains(out, "name: r1") {
		t.Fatalf("expected name r1 in yaml output, got:\n%s", out)
	}
	if !strings.Contains(out, "backupRef: backup-123") {
		t.Fatalf("expected backupRef backup-123 in yaml output, got:\n%s", out)
	}
}
