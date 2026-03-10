package real

import "testing"

func TestTerminalPhaseFromObject_NoStatus(t *testing.T) {
	obj := map[string]any{
		"spec": map[string]any{
			"backupRef": "backup-123",
		},
	}

	phase, terminal := terminalPhaseFromObject(obj)

	if phase != "" {
		t.Fatalf("expected empty phase, got %q", phase)
	}
	if terminal {
		t.Fatal("expected terminal=false, got true")
	}
}

func TestTerminalPhaseFromObject_StatusWithoutPhase(t *testing.T) {
	obj := map[string]any{
		"status": map[string]any{},
	}

	phase, terminal := terminalPhaseFromObject(obj)

	if phase != "" {
		t.Fatalf("expected empty phase, got %q", phase)
	}
	if terminal {
		t.Fatal("expected terminal=false, got true")
	}
}

func TestTerminalPhaseFromObject_RunningIsNotTerminal(t *testing.T) {
	obj := map[string]any{
		"status": map[string]any{
			"phase": "RUNNING",
		},
	}

	phase, terminal := terminalPhaseFromObject(obj)

	if phase != "RUNNING" {
		t.Fatalf("expected phase RUNNING, got %q", phase)
	}
	if terminal {
		t.Fatal("expected terminal=false, got true")
	}
}

func TestTerminalPhaseFromObject_SucceededIsTerminal(t *testing.T) {
	obj := map[string]any{
		"status": map[string]any{
			"phase": "SUCCEEDED",
		},
	}

	phase, terminal := terminalPhaseFromObject(obj)

	if phase != "SUCCEEDED" {
		t.Fatalf("expected phase SUCCEEDED, got %q", phase)
	}
	if !terminal {
		t.Fatal("expected terminal=true, got false")
	}
}

func TestTerminalPhaseFromObject_FailedIsTerminal(t *testing.T) {
	obj := map[string]any{
		"status": map[string]any{
			"phase": "FAILED",
		},
	}

	phase, terminal := terminalPhaseFromObject(obj)

	if phase != "FAILED" {
		t.Fatalf("expected phase FAILED, got %q", phase)
	}
	if !terminal {
		t.Fatal("expected terminal=true, got false")
	}
}
