package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/xgo/support/cmd"
)

func TestInspectWorktree_NotARepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitops-inspect-not-repo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	inspect, err := InspectWorktree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if inspect.IsRepo {
		t.Fatalf("IsRepo = true, want false")
	}
	if inspect.Branch != "" || inspect.CommitShort != "" || inspect.CommitMessage != "" {
		t.Fatalf("expected empty branch/commit, got %+v", inspect)
	}
	if !inspect.IsClean || inspect.Uncommitted != 0 {
		t.Fatalf("expected clean zero counts, got %+v", inspect)
	}
}

func TestInspectWorktree_CleanRepo(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	inspect, err := InspectWorktree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !inspect.IsRepo {
		t.Fatal("IsRepo = false")
	}
	if inspect.Branch != "master" {
		t.Fatalf("Branch = %q, want master", inspect.Branch)
	}
	if len(inspect.CommitShort) != 7 {
		t.Fatalf("CommitShort = %q, want 7-char hash", inspect.CommitShort)
	}
	if inspect.CommitMessage != "init" {
		t.Fatalf("CommitMessage = %q, want init", inspect.CommitMessage)
	}
	if !inspect.IsClean || inspect.Uncommitted != 0 {
		t.Fatalf("expected clean worktree, got %+v", inspect)
	}
}

func TestInspectWorktree_DirtyAllTypes(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "tracked.txt", "baseline")
	writeFile(dir, "to-delete.txt", "baseline")
	writeFile(dir, "to-rename.txt", "baseline")
	gitAddA(dir)
	gitCommit(dir, "baseline")

	writeFile(dir, "untracked.txt", "new")
	writeFile(dir, "tracked.txt", "modified")
	if err := os.Remove(filepath.Join(dir, "to-delete.txt")); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Dir(dir).Run("git", "mv", "to-rename.txt", "renamed.txt"); err != nil {
		t.Fatal(err)
	}

	inspect, err := InspectWorktree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if inspect.IsClean {
		t.Fatal("expected dirty worktree")
	}
	if inspect.Added != 1 {
		t.Fatalf("Added = %d, want 1", inspect.Added)
	}
	if inspect.Changed != 1 {
		t.Fatalf("Changed = %d, want 1", inspect.Changed)
	}
	if inspect.Renamed != 1 {
		t.Fatalf("Renamed = %d, want 1", inspect.Renamed)
	}
	if inspect.Deleted != 1 {
		t.Fatalf("Deleted = %d, want 1", inspect.Deleted)
	}
	if inspect.Uncommitted != 4 {
		t.Fatalf("Uncommitted = %d, want 4", inspect.Uncommitted)
	}
}

func TestInspectWorktree_DetachedHead(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	if err := cmd.Dir(dir).Run("git", "checkout", "--detach", "HEAD"); err != nil {
		t.Fatal(err)
	}

	inspect, err := InspectWorktree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if inspect.Branch != "(detached)" {
		t.Fatalf("Branch = %q, want (detached)", inspect.Branch)
	}
	if len(inspect.CommitShort) != 7 {
		t.Fatalf("CommitShort = %q, want 7-char hash", inspect.CommitShort)
	}
	if inspect.CommitMessage != "init" {
		t.Fatalf("CommitMessage = %q, want init", inspect.CommitMessage)
	}
	if !inspect.IsClean {
		t.Fatalf("expected clean detached worktree, got %+v", inspect)
	}
}

func TestInspectWorktree_UntrackedDirExpandsFileCount(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "solo.txt", "one")
	if err := os.MkdirAll(filepath.Join(dir, "nested", "pkg"), 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(dir, "nested/pkg/a.go", "a")
	writeFile(dir, "nested/pkg/b.go", "b")
	writeFile(dir, "nested/pkg/c.go", "c")

	inspect, err := InspectWorktree(dir)
	if err != nil {
		t.Fatal(err)
	}
	if inspect.Added != 4 {
		t.Fatalf("Added = %d, want 4 (1 file + 3 under nested/)", inspect.Added)
	}
	if inspect.Uncommitted != 2 {
		t.Fatalf("Uncommitted = %d, want 2 porcelain lines", inspect.Uncommitted)
	}
}

func TestClassifyPorcelainLine(t *testing.T) {
	cases := []struct {
		xy   string
		want porcelainChangeType
	}{
		{"??", porcelainAdded},
		{" M", porcelainChanged},
		{"M ", porcelainChanged},
		{"MM", porcelainChanged},
		{" D", porcelainDeleted},
		{"D ", porcelainDeleted},
		{"R ", porcelainRenamed},
		{" R", porcelainRenamed},
		{"A ", porcelainAdded},
		{" A", porcelainAdded},
		{"AD", porcelainDeleted},
		{"  ", porcelainNone},
	}
	for _, tc := range cases {
		got := classifyPorcelainLine(tc.xy)
		if got != tc.want {
			t.Fatalf("classifyPorcelainLine(%q) = %v, want %v", tc.xy, got, tc.want)
		}
	}
}