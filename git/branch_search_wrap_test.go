package git

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestWrapGitBranchContainsErr_includesStderrOn129(t *testing.T) {
	ee := &exec.ExitError{ProcessState: nil}
	// craft via sh
	cmd := exec.Command("sh", "-c", "echo 'error: malformed object name refs/heads/x' >&2; exit 129")
	_, err := cmd.Output()
	if err == nil {
		t.Fatal("expected error")
	}
	wrapped := wrapGitBranchContainsErr("refs/heads/x", err, "error: malformed object name refs/heads/x")
	if !strings.Contains(wrapped.Error(), "exit status 129") {
		t.Fatalf("want exit status 129 in %v", wrapped)
	}
	if !strings.Contains(wrapped.Error(), "malformed object name") {
		t.Fatalf("want stderr in %v", wrapped)
	}
	_ = ee
	_ = errors.New("")
}
