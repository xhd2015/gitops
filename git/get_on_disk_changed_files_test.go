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
