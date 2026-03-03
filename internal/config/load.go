package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/gabrielribeirojb/crd-study/internal/state"
)

// loadDesiredFile é o "formato do YAML" (espelho do arquivo)
type loadDesiredSimple struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	BackupRef string `yaml:"backupRef"`
}

type loadDesiredK8S struct {
	Metadata struct {
		Namespace string `yaml:"namespace"`
		Name      string `yaml:"name"`
	} `yaml:"metadata"`

	Spec struct {
		BackupRef string `yaml:"backupRef"`
	} `yaml:"spec"`
}

func LoadDesired(path string) (state.DesiredSpec, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return state.DesiredSpec{}, err
	}

	// 1) Tenta ler como formato simples
	var s loadDesiredSimple
	if err := yaml.Unmarshal(b, &s); err != nil {
		return state.DesiredSpec{}, err
	}

	desired := state.DesiredSpec{
		Namespace: s.Namespace,
		Name:      s.Name,
		BackupRef: s.BackupRef,
	}

	if desired.Namespace == "" && desired.Name == "" && desired.BackupRef == "" {
		var k loadDesiredK8S
		if err := yaml.Unmarshal(b, &k); err != nil {
			return state.DesiredSpec{}, err
		}

		desired.Namespace = k.Metadata.Namespace
		desired.Name = k.Metadata.Name
		desired.BackupRef = k.Spec.BackupRef
	}

	// 3) Defaults + validação
	if desired.Namespace == "" {
		desired.Namespace = "default"
	}
	if desired.Name == "" {
		return state.DesiredSpec{}, fmt.Errorf("name is required (use either top-level 'name' or 'metadata.name')")
	}
	if desired.BackupRef == "" {
		return state.DesiredSpec{}, fmt.Errorf("backupRef is required (use either top-level 'backupRef' or 'spec.backupRef')")
	}

	return desired, nil
}
