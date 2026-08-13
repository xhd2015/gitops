package repos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// Concurrent static depth-1 fetch + connectivity deepen on one static bare must
// not hit git shallow.lock (the race that failed report 149277-class jobs).
func TestEnsureStaticBareCommitsConcurrentStaticAndPartialNoShallowLock(t *testing.T) {
	remote, featureHead, masterHead := makeStaticCommitsTestRemote(t)
	reposRoot := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	cloneURL := "file://" + filepath.ToSlash(remote)

	cacheDir, err := EnsureStaticBareCache(context.Background(), cloneURL, BareCacheOptions{
		Depth:     1,
		ReposRoot: reposRoot,
	})
	if err != nil {
		t.Fatalf("ensure static bare: %v", err)
	}

	const callers = 24
	start := make(chan struct{})
	errs := make(chan error, callers)
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			<-start
			var runErr error
			if i%2 == 0 {
				// Static path: FetchCommits only (depth-1 missing SHAs).
				runErr = FetchCommits(context.Background(), cacheDir, []string{featureHead}, FetchCommitsOptions{
					CloneURL: cloneURL,
				})
			} else {
				// Partial path: same bare, connectivity deepen under one flock.
				_, runErr = EnsureStaticBareCommits(context.Background(), cacheDir, []string{featureHead, masterHead}, EnsureStaticBareCommitsOptions{
					CloneURL:        cloneURL,
					Primary:         featureHead,
					EnsureConnected: true,
				})
			}
			if runErr != nil {
				errs <- runErr
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "shallow.lock") {
			t.Fatalf("shallow.lock race: %v", err)
		}
		t.Fatalf("concurrent ensure/fetch: %v", err)
	}

	ok, err := TargetedHistoryConnected(context.Background(), cacheDir, featureHead, masterHead)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("histories not connected after concurrent ensure")
	}
}

func TestEnsureStaticBareCommitsRecheckSkipsAfterPeerConnected(t *testing.T) {
	remote, featureHead, masterHead := makeStaticCommitsTestRemote(t)
	reposRoot := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	cloneURL := "file://" + filepath.ToSlash(remote)

	cacheDir, err := EnsureStaticBareCache(context.Background(), cloneURL, BareCacheOptions{
		Depth:     1,
		ReposRoot: reposRoot,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Cold: missing SHAs + deepen.
	first, err := EnsureStaticBareCommits(context.Background(), cacheDir, []string{featureHead, masterHead}, EnsureStaticBareCommitsOptions{
		CloneURL:        cloneURL,
		Primary:         featureHead,
		EnsureConnected: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if first.Phase == "fetch_commits" {
		// May be deepen_* or connected; fetch_commits alone is wrong if others needed history.
		ok, _ := TargetedHistoryConnected(context.Background(), cacheDir, featureHead, masterHead)
		if !ok {
			t.Fatalf("first phase=%s but not connected", first.Phase)
		}
	}

	// Warm: should recheck under lock and return connected without error.
	second, err := EnsureStaticBareCommits(context.Background(), cacheDir, []string{featureHead, masterHead}, EnsureStaticBareCommitsOptions{
		CloneURL:        cloneURL,
		Primary:         featureHead,
		EnsureConnected: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if second.Phase != "connected" {
		t.Fatalf("warm phase=%s want connected", second.Phase)
	}
}

func makeStaticCommitsTestRemote(t *testing.T) (remote, featureHead, masterHead string) {
	t.Helper()
	root := t.TempDir()
	work := filepath.Join(root, "work")
	remote = filepath.Join(root, "remote.git")
	gitTestRun(t, "", "init", "--initial-branch=master", work)
	gitTestRun(t, work, "config", "user.name", "Gitops Test")
	gitTestRun(t, work, "config", "user.email", "gitops@example.com")
	if err := os.WriteFile(filepath.Join(work, "base.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitTestRun(t, work, "add", ".")
	gitTestRun(t, work, "commit", "-m", "base")

	gitTestRun(t, work, "checkout", "-b", "feature")
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(filepath.Join(work, "feature.txt"), []byte(fmt.Sprintf("%d\n", i)), 0o644); err != nil {
			t.Fatal(err)
		}
		gitTestRun(t, work, "add", "feature.txt")
		gitTestRun(t, work, "commit", "-m", fmt.Sprintf("feature %d", i))
	}
	featureHead = gitTestOutput(t, work, "rev-parse", "HEAD")

	gitTestRun(t, work, "checkout", "master")
	for i := 0; i < 12; i++ {
		if err := os.WriteFile(filepath.Join(work, "master.txt"), []byte(fmt.Sprintf("%d\n", i)), 0o644); err != nil {
			t.Fatal(err)
		}
		gitTestRun(t, work, "add", "master.txt")
		gitTestRun(t, work, "commit", "-m", fmt.Sprintf("master %d", i))
	}
	masterHead = gitTestOutput(t, work, "rev-parse", "HEAD")

	gitTestRun(t, "", "clone", "--bare", work, remote)
	gitTestRun(t, remote, "config", "uploadpack.allowAnySHA1InWant", "true")
	return remote, featureHead, masterHead
}
