package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/xgo/support/cmd"
)

func TestCompareBranchesFastForward(t *testing.T) {
	dir, cleanup := mustGetTmpDir()
	defer cleanup()

	// Create feature branch, add a commit on it
	cmd.Dir(dir).Run("git", "checkout", "-b", "feature")
	os.WriteFile(filepath.Join(dir, "feat.txt"), []byte("feature"), 0644)
	cmd.Dir(dir).Run("git", "add", "feat.txt")
	commit(dir, "test", "test@test.com", "feature commit")

	// Switch back to master
	cmd.Dir(dir).Run("git", "checkout", "master")

	result, err := CompareBranches(dir, "master", "feature")
	if err != nil {
		t.Fatalf("CompareBranches: %v", err)
	}
	if result.Relation != BranchRelationAIsAncestorOfB {
		t.Fatalf("expected BranchRelationAIsAncestorOfB, got %v", result.Relation)
	}
	if result.CommitsAheadB != 1 {
		t.Fatalf("expected CommitsAheadB=1, got %d", result.CommitsAheadB)
	}
	if result.CommitsAheadA != 0 {
		t.Fatalf("expected CommitsAheadA=0, got %d", result.CommitsAheadA)
	}
}

func TestCompareBranchesDiverged(t *testing.T) {
	dir, cleanup := mustGetTmpDir()
	defer cleanup()

	// Create feature branch, add a commit
	cmd.Dir(dir).Run("git", "checkout", "-b", "feature")
	os.WriteFile(filepath.Join(dir, "feat.txt"), []byte("feature"), 0644)
	cmd.Dir(dir).Run("git", "add", "feat.txt")
	commit(dir, "test", "test@test.com", "feature commit")

	// Switch back to master, add a different commit
	cmd.Dir(dir).Run("git", "checkout", "master")
	os.WriteFile(filepath.Join(dir, "main.txt"), []byte("main"), 0644)
	cmd.Dir(dir).Run("git", "add", "main.txt")
	commit(dir, "test", "test@test.com", "main commit")

	result, err := CompareBranches(dir, "master", "feature")
	if err != nil {
		t.Fatalf("CompareBranches: %v", err)
	}
	if result.Relation != BranchRelationDiverged {
		t.Fatalf("expected BranchRelationDiverged, got %v", result.Relation)
	}
	if result.CommitsAheadA != 1 {
		t.Fatalf("expected CommitsAheadA=1, got %d", result.CommitsAheadA)
	}
	if result.CommitsAheadB != 1 {
		t.Fatalf("expected CommitsAheadB=1, got %d", result.CommitsAheadB)
	}
	if result.MergeBase == "" {
		t.Fatal("expected non-empty MergeBase")
	}
	if result.DiffFileCount == 0 {
		t.Fatal("expected non-zero DiffFileCount")
	}
}

func TestCompareBranchesSameCommit(t *testing.T) {
	dir, cleanup := mustGetTmpDir()
	defer cleanup()

	result, err := CompareBranches(dir, "HEAD", "HEAD")
	if err != nil {
		t.Fatalf("CompareBranches: %v", err)
	}
	if result.Relation != BranchRelationSame {
		t.Fatalf("expected BranchRelationSame, got %v", result.Relation)
	}
}

func TestCompareBranchesBIsAncestorOfA(t *testing.T) {
	dir, cleanup := mustGetTmpDir()
	defer cleanup()

	// master has a commit that feature does not have
	cmd.Dir(dir).Run("git", "checkout", "-b", "feature")
	cmd.Dir(dir).Run("git", "checkout", "master")
	os.WriteFile(filepath.Join(dir, "main.txt"), []byte("main"), 0644)
	cmd.Dir(dir).Run("git", "add", "main.txt")
	commit(dir, "test", "test@test.com", "main commit")

	result, err := CompareBranches(dir, "master", "feature")
	if err != nil {
		t.Fatalf("CompareBranches: %v", err)
	}
	// feature is behind master → feature (B) is ancestor of master (A)
	// So B is ancestor of A = BranchRelationBIsAncestorOfA
	if result.Relation != BranchRelationBIsAncestorOfA {
		t.Fatalf("expected BranchRelationBIsAncestorOfA, got %v", result.Relation)
	}
	if result.CommitsAheadA != 1 {
		t.Fatalf("expected CommitsAheadA=1, got %d", result.CommitsAheadA)
	}
	if result.CommitsAheadB != 0 {
		t.Fatalf("expected CommitsAheadB=0, got %d", result.CommitsAheadB)
	}
}
