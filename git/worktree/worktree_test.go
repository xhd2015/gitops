package worktree

import "testing"

// go test -run TestAddWorkTree -v ./git/gitops
func TestAddWorkTree(t *testing.T) {
	tmpDir, remove, err := AcquireTemp("TODO", "master")
	if err != nil {
		t.Fatal(err)
	}
	defer remove()
	t.Logf("tmpDir: %s", tmpDir)
}
