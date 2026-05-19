package harness

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAgentsFile_Missing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("BOTCTL_HOME", tmp)
	got, err := loadAgentsFile()
	if err != nil {
		t.Fatalf("loadAgentsFile() with no file returned err: %v", err)
	}
	if got != "" {
		t.Errorf("loadAgentsFile() with no file = %q, want \"\"", got)
	}
}

func TestLoadAgentsFile_Present(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("BOTCTL_HOME", tmp)

	content := "  always prefer jq over python3\nnever use sleep in bash\n  "
	if err := os.WriteFile(filepath.Join(tmp, "AGENTS.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	got, err := loadAgentsFile()
	if err != nil {
		t.Fatalf("loadAgentsFile() returned err: %v", err)
	}
	want := "always prefer jq over python3\nnever use sleep in bash"
	if got != want {
		t.Errorf("loadAgentsFile() = %q, want %q", got, want)
	}
}

func TestLoadAgentsFile_Empty(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("BOTCTL_HOME", tmp)

	if err := os.WriteFile(filepath.Join(tmp, "AGENTS.md"), []byte("   \n\n  "), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	got, err := loadAgentsFile()
	if err != nil {
		t.Fatalf("loadAgentsFile() returned err: %v", err)
	}
	if got != "" {
		t.Errorf("loadAgentsFile() with whitespace-only file = %q, want \"\"", got)
	}
}

func TestBuildSystemPrompt_Basic(t *testing.T) {
	got := buildSystemPrompt("/ws", "Available skills: foo", 0, "")
	mustContain(t, got, "You are an autonomous agent managed by `botctl`.")
	mustContain(t, got, "Workspace directory: /ws")
	mustContain(t, got, "Available skills: foo")
	mustContain(t, got, "Your full instructions are in the user message below.")
	mustNotContain(t, got, "## Shared Instructions")
	mustNotContain(t, got, "maximum of")
}

func TestBuildSystemPrompt_WithMaxTurns(t *testing.T) {
	got := buildSystemPrompt("/ws", "", 25, "")
	mustContain(t, got, "maximum of 25 turns")
}

func TestBuildSystemPrompt_WithSharedInstructions(t *testing.T) {
	got := buildSystemPrompt("/ws", "", 0, "rule one\nrule two")
	mustContain(t, got, "## Shared Instructions\nrule one\nrule two")
	// Shared instructions block must come BEFORE the closing "user message below" line
	sharedIdx := strings.Index(got, "## Shared Instructions")
	closingIdx := strings.Index(got, "Your full instructions are in the user message below.")
	if sharedIdx < 0 || closingIdx < 0 || sharedIdx >= closingIdx {
		t.Errorf("Shared Instructions should appear before closing line. shared=%d closing=%d\ngot:\n%s", sharedIdx, closingIdx, got)
	}
}

func TestBuildSystemPrompt_WithBoth(t *testing.T) {
	got := buildSystemPrompt("/ws", "Available skills: foo", 10, "rule one")
	mustContain(t, got, "maximum of 10 turns")
	mustContain(t, got, "## Shared Instructions\nrule one")
}

func mustContain(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected substring %q in:\n%s", substr, s)
	}
}

func mustNotContain(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("did not expect substring %q in:\n%s", substr, s)
	}
}
