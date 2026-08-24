package repos

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildURLFetchArgsDeepenUsesBlobNone(t *testing.T) {
	args := buildURLFetchArgs("/tmp/unused", "https://example/repo.git", URLFetchOptions{
		Deepen: 16,
		Refs:   []string{"abc"},
	}, true)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--filter=blob:none") {
		t.Fatalf("deepen args missing blob:none filter: %v", args)
	}
	if !strings.Contains(joined, "--deepen=16") {
		t.Fatalf("deepen args missing --deepen=16: %v", args)
	}
}

func TestBuildURLFetchArgsDepthOmitsBlobNone(t *testing.T) {
	args := buildURLFetchArgs("/tmp/unused", "https://example/repo.git", URLFetchOptions{
		Depth: 1,
		Refs:  []string{"abc"},
	}, false)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--filter=") {
		t.Fatalf("depth fetch must not filter (worktree needs blobs): %v", args)
	}
	if !strings.Contains(joined, "--depth=1") {
		t.Fatalf("depth args missing --depth=1: %v", args)
	}
}

func TestFetchFromURLDeepenFilterConnectsWithoutPullingAllBlobs(t *testing.T) {
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

	res, err := EnsureStaticBareCommits(context.Background(), cacheDir, []string{featureHead, masterHead}, EnsureStaticBareCommitsOptions{
		CloneURL:        cloneURL,
		Primary:         featureHead,
		EnsureConnected: true,
	})
	if err != nil {
		t.Fatalf("ensure connected: %v", err)
	}
	if res.Phase == "" {
		t.Fatal("empty phase")
	}

	ok, err := TargetedHistoryConnected(context.Background(), cacheDir, featureHead, masterHead)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("not connected after ensure phase=%s", res.Phase)
	}

	// Tip trees from the unfiltered depth-1 seed must still materialize for worktrees.
	wt := filepath.Join(t.TempDir(), "wt")
	if err := AddDetachedWorktree(cacheDir, wt, featureHead); err != nil {
		t.Fatalf("worktree after filtered deepen: %v", err)
	}
	if _, err := os.Stat(filepath.Join(wt, "feature.txt")); err != nil {
		t.Fatalf("expected feature.txt in worktree: %v", err)
	}
}
