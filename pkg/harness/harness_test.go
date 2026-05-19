package harness

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAgentsFile_Missing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("BOTCTL_HOME", tmp)
	if got := loadAgentsFile(); got != "" {
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

	want := "always prefer jq over python3\nnever use sleep in bash"
	if got := loadAgentsFile(); got != want {
		t.Errorf("loadAgentsFile() = %q, want %q", got, want)
	}
}

func TestLoadAgentsFile_Empty(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("BOTCTL_HOME", tmp)

	if err := os.WriteFile(filepath.Join(tmp, "AGENTS.md"), []byte("   \n\n  "), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	if got := loadAgentsFile(); got != "" {
		t.Errorf("loadAgentsFile() with whitespace-only file = %q, want \"\"", got)
	}
}
