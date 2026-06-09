package git

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/xhd2015/xgo/support/cmd"
)

func assertFilesEqual(t *testing.T, got, want []string) {
	t.Helper()
	sort.Strings(got)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %d files %v, want %d files %v", len(got), got, len(want), want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("file[%d]: got %q, want %q\n got: %v\nwant: %v", i, got[i], want[i], got, want)
		}
	}
}

func writeFile(dir, name, content string) {
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		panic(err)
	}
}

func gitAdd(dir string, files ...string) {
	for _, f := range files {
		if err := cmd.Dir(dir).Run("git", "add", f); err != nil {
			panic(err)
		}
	}
}

func gitAddA(dir string) {
	if err := cmd.Dir(dir).Run("git", "add", "-A"); err != nil {
		panic(err)
	}
}

func gitCommit(dir, msg string) {
	if err := cmd.Dir(dir).Run("git", "-c", "user.email=test@test.com", "-c", "user.name=test", "commit", "-m", msg); err != nil {
		panic(err)
	}
}

func TestGetOnDiskChangedFiles_NoChanges(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_ModifiedUnstaged(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	writeFile(dir, "README.md", "modified content")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"README.md"})
}

func TestGetOnDiskChangedFiles_ModifiedStaged(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	writeFile(dir, "README.md", "modified content")
	gitAdd(dir, "README.md")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"README.md"})
}

func TestGetOnDiskChangedFiles_ModifiedBoth(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	writeFile(dir, "README.md", "modified v1")
	gitAdd(dir, "README.md")
	writeFile(dir, "README.md", "modified v2")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"README.md"})
}

func TestGetOnDiskChangedFiles_NewUntracked(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	writeFile(dir, "new.go", "package main")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"new.go"})
}

func TestGetOnDiskChangedFiles_NewStaged(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	writeFile(dir, "staged.go", "package main")
	gitAdd(dir, "staged.go")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"staged.go"})
}

func TestGetOnDiskChangedFiles_DeletedUnstaged(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	if err := os.Remove(filepath.Join(dir, "README.md")); err != nil {
		t.Fatal(err)
	}
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_DeletedStaged(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	if err := cmd.Dir(dir).Run("git", "rm", "-q", "README.md"); err != nil {
		t.Fatal(err)
	}
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_StagedRenameNoChange(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	if err := cmd.Dir(dir).Run("git", "mv", "README.md", "RENAMED.md"); err != nil {
		t.Fatal(err)
	}
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"RENAMED.md"})
}

func TestGetOnDiskChangedFiles_StagedRenameWithEdit(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	if err := cmd.Dir(dir).Run("git", "mv", "README.md", "RENAMED.md"); err != nil {
		t.Fatal(err)
	}
	writeFile(dir, "RENAMED.md", "modified after rename")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"RENAMED.md"})
}

func TestGetOnDiskChangedFiles_UntrackedRenameViaMv(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()
	if err := os.Rename(filepath.Join(dir, "README.md"), filepath.Join(dir, "RENAMED2.md")); err != nil {
		t.Fatal(err)
	}
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"RENAMED2.md"})
}

func TestGetOnDiskChangedFiles_Mixed(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "mod.go", "original")
	writeFile(dir, "del.go", "original")
	writeFile(dir, "ren.go", "original")
	gitAddA(dir)
	gitCommit(dir, "add baseline files")

	writeFile(dir, "mod.go", "modified")
	writeFile(dir, "added.go", "new staged")
	gitAdd(dir, "added.go")
	writeFile(dir, "untracked.go", "untracked")
	if err := os.Remove(filepath.Join(dir, "del.go")); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Dir(dir).Run("git", "mv", "ren.go", "ren_new.go"); err != nil {
		t.Fatal(err)
	}

	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"mod.go", "added.go", "untracked.go", "ren_new.go"})
}

func TestGetOnDiskChangedFiles_SubdirectoryFile(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	subDir := filepath.Join(dir, "sub", "pkg")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(dir, "sub/pkg/foo.go", "original")
	gitAddA(dir)
	gitCommit(dir, "add subdir file")

	writeFile(dir, "sub/pkg/foo.go", "modified")
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"sub/pkg/foo.go"})
}

func TestGetOnDiskChangedFiles_MultipleModified(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	for _, f := range []string{"a.go", "b.go", "c.go"} {
		writeFile(dir, f, "original")
	}
	gitAddA(dir)
	gitCommit(dir, "add files")

	for _, f := range []string{"a.go", "b.go", "c.go"} {
		writeFile(dir, f, "modified")
	}
	files, err := GetOnDiskChangedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"a.go", "b.go", "c.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_ModifiedSinceCommit(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "clean.go", "original")
	gitAddA(dir)
	gitCommit(dir, "second commit")

	writeFile(dir, "mod.go", "package main\nimport \"flag\"\nfunc main() { flag.Parse() }\n")
	writeFile(dir, "clean.go", "modified")

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"mod.go", "clean.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_NoChanges(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "clean.go", "original")
	gitAddA(dir)
	gitCommit(dir, "second commit")

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_CompareWith_OlderCommit(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "a.go", "package a\n")
	gitAddA(dir)
	gitCommit(dir, "add a.go")

	writeFile(dir, "b.go", "package b\n")
	gitAddA(dir)
	gitCommit(dir, "add b.go")

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD~1"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"b.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_InvalidRef(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	_, err := GetOnDiskChangedFiles(dir, CompareWith("nonexistent123"))
	if err == nil {
		t.Fatal("expected error for invalid ref")
	}
}

func TestGetOnDiskChangedFiles_CompareWith_ExcludesDeleted(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "a.go", "package main\n")
	gitAddA(dir)
	gitCommit(dir, "add a.go")

	if err := os.Remove(filepath.Join(dir, "a.go")); err != nil {
		t.Fatal(err)
	}

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_CompareWith_StagedChanges(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "mod.go", "original")
	gitAddA(dir)
	gitCommit(dir, "commit mod.go")

	writeFile(dir, "mod.go", "modified")
	gitAdd(dir, "mod.go")

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"mod.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_UntrackedFile(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)

	writeFile(dir, "new.go", "package main\n")

	files, err = GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"new.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_RenamedFile(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "old.go", "package main\n")
	gitAddA(dir)
	gitCommit(dir, "commit old.go")

	if err := cmd.Dir(dir).Run("git", "mv", "old.go", "new.go"); err != nil {
		t.Fatal(err)
	}

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"new.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_MixedChanges(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "mod.go", "original")
	writeFile(dir, "del.go", "original")
	gitAddA(dir)
	gitCommit(dir, "add files")

	writeFile(dir, "mod.go", "modified")
	if err := os.Remove(filepath.Join(dir, "del.go")); err != nil {
		t.Fatal(err)
	}
	writeFile(dir, "untracked.go", "package main\n")

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"mod.go", "untracked.go"})
}

func TestGetOnDiskChangedFiles_CompareWith_NoDiffBetweenCleanCommits(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "a.go", "package main\n")
	gitAddA(dir)
	gitCommit(dir, "add a.go")

	if err := cmd.Dir(dir).Run("git", "-c", "user.email=test@test.com", "-c", "user.name=test", "commit", "--allow-empty", "-m", "empty commit"); err != nil {
		t.Fatal(err)
	}

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD~1"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_CompareWith_BranchName(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	files, err := GetOnDiskChangedFiles(dir, CompareWith("master"))
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, nil)
}

func TestGetOnDiskChangedFiles_ResolvePathsToFiles_UntrackedDir(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	subDir := filepath.Join(dir, "view")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(dir, "view/a.go", "package view\n")
	writeFile(dir, "view/b.go", "package view\n")

	files, err := GetOnDiskChangedFiles(dir, ResolvePathsToFiles())
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"view/a.go", "view/b.go"})
}

func TestGetOnDiskChangedFiles_ResolvePathsToFiles_OnlyFiles(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "new.go", "package main\n")

	files, err := GetOnDiskChangedFiles(dir, ResolvePathsToFiles())
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"new.go"})
}

func TestGetOnDiskChangedFiles_ResolvePathsToFiles_WithCompareWith(t *testing.T) {
	dir, clean := mustGetTmpDir()
	defer clean()

	writeFile(dir, "clean.go", "original")
	gitAddA(dir)
	gitCommit(dir, "second commit")

	subDir := filepath.Join(dir, "view")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(dir, "view/a.go", "package view\n")

	files, err := GetOnDiskChangedFiles(dir, CompareWith("HEAD"), ResolvePathsToFiles())
	if err != nil {
		t.Fatal(err)
	}
	assertFilesEqual(t, files, []string{"view/a.go"})
}
