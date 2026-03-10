package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "desired.yaml")

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp yaml: %v", err)
	}

	return path
}

func TestLoadDesired_SimpleFormat(t *testing.T) {
	yaml := `
namespace: demo
name: r1
backupRef: backup-123
`
	path := writeTempYAML(t, yaml)

	got, err := LoadDesired(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Namespace != "demo" {
		t.Fatalf("expected namespace demo, got %q", got.Namespace)
	}
	if got.Name != "r1" {
		t.Fatalf("expected name r1, got %q", got.Name)
	}
	if got.BackupRef != "backup-123" {
		t.Fatalf("expected backupRef backup-123, got %q", got.BackupRef)
	}
}

func TestLoadDesired_KubernetesLikeFormat(t *testing.T) {
	yaml := `
apiVersion: gdch.mycompany.io/v1
kind: ClusterRestore
metadata:
  namespace: default
  name: r1
spec:
  backupRef: backup-123
`
	path := writeTempYAML(t, yaml)

	got, err := LoadDesired(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Namespace != "default" {
		t.Fatalf("expected namespace default, got %q", got.Namespace)
	}
	if got.Name != "r1" {
		t.Fatalf("expected name r1, got %q", got.Name)
	}
	if got.BackupRef != "backup-123" {
		t.Fatalf("expected backupRef backup-123, got %q", got.BackupRef)
	}
}

func TestLoadDesired_DefaultNamespaceWhenEmpty(t *testing.T) {
	yaml := `
name: r1
backupRef: backup-123
`
	path := writeTempYAML(t, yaml)

	got, err := LoadDesired(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Namespace != "default" {
		t.Fatalf("expected namespace default, got %q", got.Namespace)
	}
}

func TestLoadDesired_DefaultNamespaceWhenEmptyInK8SFormat(t *testing.T) {
	yaml := `
apiVersion: gdch.mycompany.io/v1
kind: ClusterRestore
metadata:
  name: r1
spec:
  backupRef: backup-123
`
	path := writeTempYAML(t, yaml)

	got, err := LoadDesired(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Namespace != "default" {
		t.Fatalf("expected namespace default, got %q", got.Namespace)
	}
}

func TestLoadDesired_ErrorWhenNameMissing(t *testing.T) {
	yaml := `
namespace: default
backupRef: backup-123
`
	path := writeTempYAML(t, yaml)

	_, err := LoadDesired(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("expected name error, got %v", err)
	}
}

func TestLoadDesired_ErrorWhenBackupRefMissing(t *testing.T) {
	yaml := `
namespace: default
name: r1
`
	path := writeTempYAML(t, yaml)

	_, err := LoadDesired(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "backupRef is required") {
		t.Fatalf("expected backupRef error, got %v", err)
	}
}

func TestLoadDesired_InvalidYAML(t *testing.T) {
	yaml := `
namespace: default
name: r1
backupRef: [invalid
`
	path := writeTempYAML(t, yaml)

	_, err := LoadDesired(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadDesired_FileNotFound(t *testing.T) {
	_, err := LoadDesired("file-that-does-not-exist.yaml")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
