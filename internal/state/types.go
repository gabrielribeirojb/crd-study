package state

// DesiredSpec = "estado desejado" (vem do YAML)
type DesiredSpec struct {
	Namespace string
	Name      string
	BackupRef string
}

// CurrentState = "estado atual observado" (vem da API)
type CurrentState struct {
	Exists bool
	Phase  string
}
