package repos

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanReviewPairDifferentTrees(t *testing.T) {
	t.Parallel()
	left, right := PlanReviewPair(
		"src", CommitMeta{Tree: "t-src", Parents: []string{"p"}},
		"base", CommitMeta{Tree: "t-base", Parents: []string{"q"}},
	)
	if left != "src" || right != "base" {
		t.Fatalf("got %s %s, want src base", left, right)
	}
}

func TestPlanReviewPairIdenticalBaseMergesSource(t *testing.T) {
	t.Parallel()
	// 4e0e9354 merges 2f5d8ab4; same tree → review left is base^1.
	left, right := PlanReviewPair(
		"2f5d8ab4src", CommitMeta{Tree: "same", Parents: []string{"srcp1", "srcp2"}},
		"4e0e9354base", CommitMeta{Tree: "same", Parents: []string{"oldmaster", "2f5d8ab4src"}},
	)
	if left != "oldmaster" || right != "2f5d8ab4src" {
		t.Fatalf("got %s %s, want oldmaster 2f5d8ab4src", left, right)
	}
}

func TestPlanReviewPairIdenticalSourceMergesBase(t *testing.T) {
	t.Parallel()
	left, right := PlanReviewPair(
		"src", CommitMeta{Tree: "same", Parents: []string{"old", "base"}},
		"base", CommitMeta{Tree: "same", Parents: []string{"bp"}},
	)
	if left != "old" || right != "base" {
		t.Fatalf("got %s %s, want old base", left, right)
	}
}

func TestPlanReviewPairIdenticalPeelMergeTip(t *testing.T) {
	t.Parallel()
	left, right := PlanReviewPair(
		"src", CommitMeta{Tree: "same", Parents: []string{"s1"}},
		"base", CommitMeta{Tree: "same", Parents: []string{"b1", "b2"}},
	)
	if left != "b1" || right != "src" {
		t.Fatalf("got %s %s, want b1 src", left, right)
	}
}

func TestPlanReviewPairSameSHAPeelsFirstParent(t *testing.T) {
	t.Parallel()
	left, right := PlanReviewPair(
		"aa", CommitMeta{Tree: "t", Parents: []string{"p1", "p2"}},
		"aa", CommitMeta{Tree: "t", Parents: []string{"p1", "p2"}},
	)
	if left != "p1" || right != "aa" {
		t.Fatalf("got %s %s, want p1 aa", left, right)
	}
}

func TestReadCommitMetaAndMergeBase(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q", "-b", "master")
	run("config", "user.email", "t@t")
	run("config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", "a.txt")
	run("commit", "-q", "-m", "A")
	run("checkout", "-q", "-b", "feature")
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", "b.txt")
	run("commit", "-q", "-m", "B")
	feat := gitOut(t, dir, "rev-parse", "HEAD")
	run("checkout", "-q", "master")
	base := gitOut(t, dir, "rev-parse", "HEAD")
	run("merge", "-q", "--no-ff", "-m", "M", "feature")
	merge := gitOut(t, dir, "rev-parse", "HEAD")

	featMeta, err := ReadCommitMeta(dir, feat)
	if err != nil {
		t.Fatal(err)
	}
	mergeMeta, err := ReadCommitMeta(dir, merge)
	if err != nil {
		t.Fatal(err)
	}
	if len(mergeMeta.Parents) != 2 {
		t.Fatalf("merge parents=%v", mergeMeta.Parents)
	}
	if mergeMeta.Parents[1] != feat && mergeMeta.Parents[0] != feat {
		t.Fatalf("merge parents %v should include feature %s", mergeMeta.Parents, feat)
	}
	mb, ok := MergeBase(dir, feat, base)
	if !ok || mb != base {
		t.Fatalf("merge-base(feat, base)=%q ok=%v want %s", mb, ok, base)
	}
	// identical trees after merge: peel merge^1 vs feature
	if mergeMeta.Tree == featMeta.Tree {
		left, right := PlanReviewPair(feat, featMeta, merge, mergeMeta)
		if left != mergeMeta.Parents[0] {
			t.Fatalf("identical peel left=%s want %s", left, mergeMeta.Parents[0])
		}
		got, ok := MergeBase(dir, left, right)
		if !ok || got != base {
			t.Fatalf("merge-base(peel)=%q ok=%v want %s", got, ok, base)
		}
	}
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}
