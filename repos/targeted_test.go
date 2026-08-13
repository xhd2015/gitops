package repos

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTargetedCacheFetchAndLazyFile(t *testing.T) {
	remote, work := makeTargetedTestRemote(t)
	t.Setenv("HOME", t.TempDir())
	cloneURL := "file://" + filepath.ToSlash(remote)
	head := gitTestOutput(t, work, "rev-parse", "HEAD")

	result, err := FetchTargetedRefs(context.Background(), cloneURL, []string{head}, TargetedFetchOptions{Depth: 1})
	if err != nil {
		t.Fatalf("fetch exact commit: %v", err)
	}
	if !result.Fetched {
		t.Fatal("cold exact fetch should use the remote")
	}
	if result.Commits[head] != head {
		t.Fatalf("resolved commit = %q, want %q", result.Commits[head], head)
	}
	legacyDir, err := CacheDir(cloneURL)
	if err != nil {
		t.Fatal(err)
	}
	if result.CacheDir == legacyDir || !strings.Contains(result.CacheDir, string(filepath.Separator)+"targeted-cache"+string(filepath.Separator)) {
		t.Fatalf("targeted cache not isolated: %s", result.CacheDir)
	}

	ok, content, err := ReadTargetedFile(context.Background(), result.CacheDir, nil, nil, head, "keep-newlines.txt")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !ok || content != "line one\n\n" {
		t.Fatalf("content = %q, ok=%v", content, ok)
	}
	ok, _, err = ReadTargetedFile(context.Background(), result.CacheDir, nil, nil, head, "missing.txt")
	if err != nil || ok {
		t.Fatalf("missing file: ok=%v err=%v", ok, err)
	}

	warm, err := FetchTargetedRefs(context.Background(), cloneURL, []string{head}, TargetedFetchOptions{Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	if warm.Fetched {
		t.Fatal("immutable commit should be a cache hit")
	}
}

func TestTargetedBranchRefreshAndConcurrentHit(t *testing.T) {
	remote, work := makeTargetedTestRemote(t)
	t.Setenv("HOME", t.TempDir())
	cloneURL := "file://" + filepath.ToSlash(remote)

	first, err := FetchTargetedRefs(context.Background(), cloneURL, []string{"refs/heads/master"}, TargetedFetchOptions{Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	oldHead := first.Commits["refs/heads/master"]

	if err := os.WriteFile(filepath.Join(work, "branch.txt"), []byte("moved\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitTestRun(t, work, "add", "branch.txt")
	gitTestRun(t, work, "commit", "-m", "move branch")
	gitTestRun(t, work, "push", "origin", "master")
	newHead := gitTestOutput(t, work, "rev-parse", "HEAD")
	if newHead == oldHead {
		t.Fatal("test branch did not move")
	}

	refreshed, err := FetchTargetedRefs(context.Background(), cloneURL, []string{"refs/heads/master"}, TargetedFetchOptions{Depth: 1, MaxAge: -1})
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.Commits["refs/heads/master"] != newHead {
		t.Fatalf("refreshed head = %s, want %s", refreshed.Commits["refs/heads/master"], newHead)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := FetchTargetedRefs(context.Background(), cloneURL, []string{"refs/heads/master"}, TargetedFetchOptions{Depth: 1})
			if err == nil && got.Commits["refs/heads/master"] != newHead {
				err = fmt.Errorf("got head %s", got.Commits["refs/heads/master"])
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTargetedFetchReportsUnsupportedFilterWithoutBroadeningRefs(t *testing.T) {
	remote, _ := makeTargetedTestRemote(t)
	t.Setenv("HOME", t.TempDir())
	gitTestRun(t, remote, "config", "uploadpack.allowFilter", "false")
	cloneURL := "file://" + filepath.ToSlash(remote)

	result, err := FetchTargetedRefs(context.Background(), cloneURL, []string{"refs/heads/master"}, TargetedFetchOptions{Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.FilterAccepted {
		t.Fatal("server without upload-pack filter support reported filter accepted")
	}
	if _, err := gitOutput(result.CacheDir, "show-ref", "--verify", DefaultTargetRefPrefix+"/"+targetedRefHash("refs/heads/master")); err != nil {
		t.Fatalf("requested branch was not fetched: %v", err)
	}
	if _, err := gitOutput(result.CacheDir, "show-ref", "--verify", "refs/remotes/origin/other"); err == nil {
		t.Fatal("narrow fetch unexpectedly created unrelated remote refs")
	}
}

func TestRepoLockHonorsContext(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache", "repo")
	held, err := acquireRepoLock(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	defer releaseRepoLock(held)

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err = acquireRepoLockContext(ctx, cacheDir)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("lock error = %v, want deadline exceeded", err)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("context-aware lock returned too slowly: %s", time.Since(start))
	}
}

func TestRepoLockWaitDurationReflectsContention(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache", "repo")
	held, err := acquireRepoLock(cacheDir)
	if err != nil {
		t.Fatal(err)
	}

	const hold = 80 * time.Millisecond
	go func() {
		time.Sleep(hold)
		_ = releaseRepoLock(held)
	}()

	_, wait, err := acquireRepoLockContextTimed(context.Background(), cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	// Allow timer/scheduling slack; wait should be on the order of hold, not near-zero.
	if wait < hold/2 {
		t.Fatalf("lock wait %v too small (hold %v)", wait, hold)
	}
	if wait > 2*time.Second {
		t.Fatalf("lock wait %v unexpectedly large", wait)
	}
}

func TestTargetedRefsProgressivelyDeepen(t *testing.T) {
	remote, work := makeTargetedTestRemote(t)
	t.Setenv("HOME", t.TempDir())
	cloneURL := "file://" + filepath.ToSlash(remote)
	base := gitTestOutput(t, work, "rev-parse", "HEAD")

	gitTestRun(t, work, "checkout", "-b", "feature")
	if err := os.WriteFile(filepath.Join(work, "feature.txt"), []byte("feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitTestRun(t, work, "add", "feature.txt")
	gitTestRun(t, work, "commit", "-m", "feature")
	gitTestRun(t, work, "push", "origin", "feature")

	gitTestRun(t, work, "checkout", "master")
	for i := 0; i < 12; i++ {
		if err := os.WriteFile(filepath.Join(work, "branch.txt"), []byte(fmt.Sprintf("%d\n", i)), 0o644); err != nil {
			t.Fatal(err)
		}
		gitTestRun(t, work, "add", "branch.txt")
		gitTestRun(t, work, "commit", "-m", fmt.Sprintf("master %d", i))
	}
	gitTestRun(t, work, "push", "origin", "master")

	refs := []string{"refs/heads/feature", "refs/heads/master"}
	result, err := FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	if mergeBaseOrEmpty(result.CacheDir, result.Commits[refs[0]], result.Commits[refs[1]]) != "" {
		t.Fatal("depth-1 cache unexpectedly has merge base")
	}
	result, err = FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{Deepen: 4})
	if err != nil {
		t.Fatal(err)
	}
	if mergeBaseOrEmpty(result.CacheDir, result.Commits[refs[0]], result.Commits[refs[1]]) != "" {
		t.Fatal("deepen=4 unexpectedly reached merge base")
	}
	result, err = FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{Deepen: 16})
	if err != nil {
		t.Fatal(err)
	}
	if got := mergeBaseOrEmpty(result.CacheDir, result.Commits[refs[0]], result.Commits[refs[1]]); got != base {
		t.Fatalf("merge base = %s, want %s", got, base)
	}

	refreshed, err := FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{Depth: 1, MaxAge: -1})
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.Fetched {
		t.Fatal("unchanged remote refs should refresh freshness without a depth-1 fetch")
	}
	if got := mergeBaseOrEmpty(refreshed.CacheDir, refreshed.Commits[refs[0]], refreshed.Commits[refs[1]]); got != base {
		t.Fatalf("merge base after unchanged ref refresh = %s, want %s", got, base)
	}

	stale, err := FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{
		Depth:  1,
		MaxAge: -1,
		Env: map[string]string{
			"GIT_CONFIG_COUNT":   "1",
			"GIT_CONFIG_KEY_0":   "url.http://127.0.0.1:1/.insteadOf",
			"GIT_CONFIG_VALUE_0": cloneURL,
		},
	})
	if err != nil {
		t.Fatalf("serve cached refs during transient remote failure: %v", err)
	}
	if stale.Fetched {
		t.Fatal("transient remote failure should not fetch")
	}
	if got := mergeBaseOrEmpty(stale.CacheDir, stale.Commits[refs[0]], stale.Commits[refs[1]]); got != base {
		t.Fatalf("stale merge base = %s, want %s", got, base)
	}
}

func TestConcurrentDeepenRechecksHistoryAfterLock(t *testing.T) {
	remote, work := makeTargetedTestRemote(t)
	t.Setenv("HOME", t.TempDir())
	cloneURL := "file://" + filepath.ToSlash(remote)
	base := gitTestOutput(t, work, "rev-parse", "HEAD")

	gitTestRun(t, work, "checkout", "-b", "feature")
	if err := os.WriteFile(filepath.Join(work, "feature.txt"), []byte("feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitTestRun(t, work, "add", "feature.txt")
	gitTestRun(t, work, "commit", "-m", "feature")
	gitTestRun(t, work, "push", "origin", "feature")

	gitTestRun(t, work, "checkout", "master")
	for i := 0; i < 12; i++ {
		if err := os.WriteFile(filepath.Join(work, "branch.txt"), []byte(fmt.Sprintf("%d\n", i)), 0o644); err != nil {
			t.Fatal(err)
		}
		gitTestRun(t, work, "add", "branch.txt")
		gitTestRun(t, work, "commit", "-m", fmt.Sprintf("master %d", i))
	}
	gitTestRun(t, work, "push", "origin", "master")

	refs := []string{"refs/heads/feature", "refs/heads/master"}
	initial, err := FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	if mergeBaseOrEmpty(initial.CacheDir, initial.Commits[refs[0]], initial.Commits[refs[1]]) != "" {
		t.Fatal("depth-1 cache unexpectedly has merge base")
	}

	const callers = 50
	start := make(chan struct{})
	results := make(chan *TargetedFetchResult, callers)
	errs := make(chan error, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			result, err := FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{
				Deepen:       4,
				HistoryPhase: "deepen_4",
			})
			results <- result
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	close(errs)

	fetches := 0
	for result := range results {
		if result != nil && result.Fetched {
			fetches++
		}
	}
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if fetches != 1 {
		t.Fatalf("remote deepen fetches = %d, want 1", fetches)
	}
	if got := mergeBaseOrEmpty(initial.CacheDir, initial.Commits[refs[0]], initial.Commits[refs[1]]); got != "" {
		t.Fatalf("deepen=4 merge base = %s, want empty", got)
	}

	next, err := FetchTargetedRefs(context.Background(), cloneURL, refs, TargetedFetchOptions{
		Deepen:       16,
		HistoryPhase: "deepen_16",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !next.Fetched {
		t.Fatal("new history phase should fetch once")
	}
	if got := mergeBaseOrEmpty(initial.CacheDir, initial.Commits[refs[0]], initial.Commits[refs[1]]); got != base {
		t.Fatalf("merge base = %s, want %s", got, base)
	}
}

func makeTargetedTestRemote(t *testing.T) (remote, work string) {
	t.Helper()
	root := t.TempDir()
	work = filepath.Join(root, "work")
	remote = filepath.Join(root, "remote.git")
	gitTestRun(t, "", "init", "--initial-branch=master", work)
	gitTestRun(t, work, "config", "user.name", "Gitops Test")
	gitTestRun(t, work, "config", "user.email", "gitops@example.com")
	if err := os.WriteFile(filepath.Join(work, "keep-newlines.txt"), []byte("line one\n\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(work, "large.bin"), []byte(strings.Repeat("blob-data-", 128*1024)), 0o644); err != nil {
		t.Fatal(err)
	}
	gitTestRun(t, work, "add", ".")
	gitTestRun(t, work, "commit", "-m", "initial")
	gitTestRun(t, "", "clone", "--bare", work, remote)
	gitTestRun(t, remote, "config", "uploadpack.allowFilter", "true")
	gitTestRun(t, remote, "config", "uploadpack.allowAnySHA1InWant", "true")
	gitTestRun(t, work, "remote", "add", "origin", remote)
	return remote, work
}

func gitTestRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitTestOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func mergeBaseOrEmpty(dir, left, right string) string {
	cmd := exec.Command("git", "merge-base", left, right)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func targetedRefHash(ref string) string {
	targets, err := normalizeTargetedRefs([]string{ref}, "")
	if err != nil || len(targets) != 1 {
		return ""
	}
	return filepath.Base(targets[0].local)
}
