package repos

import (
	"context"
	"fmt"
	"strings"
)

// CommitMeta is the tree and parents of a commit already present in a cache.
type CommitMeta struct {
	Tree    string
	Parents []string
}

// PlanReviewPair picks the two commits whose merge-base is the review base.
//
// Different trees: merge-base(source, base).
// Identical trees: peel "it" to first parent (the merge that contains the
// other, else a merge tip) and take merge-base(it^1, other).
func PlanReviewPair(sourceSHA string, source CommitMeta, baseSHA string, base CommitMeta) (left, right string) {
	if sourceSHA == "" {
		return baseSHA, baseSHA
	}
	if baseSHA == "" || sourceSHA == baseSHA {
		if len(source.Parents) > 0 {
			return source.Parents[0], sourceSHA
		}
		return sourceSHA, sourceSHA
	}
	if source.Tree != "" && base.Tree != "" && source.Tree == base.Tree {
		if containsSHA(base.Parents, sourceSHA) && len(base.Parents) > 0 {
			return base.Parents[0], sourceSHA
		}
		if containsSHA(source.Parents, baseSHA) && len(source.Parents) > 0 {
			return source.Parents[0], baseSHA
		}
		if len(base.Parents) >= 2 {
			return base.Parents[0], sourceSHA
		}
		if len(source.Parents) >= 2 {
			return source.Parents[0], baseSHA
		}
		if len(base.Parents) > 0 {
			return base.Parents[0], sourceSHA
		}
		if len(source.Parents) > 0 {
			return source.Parents[0], baseSHA
		}
	}
	return sourceSHA, baseSHA
}

func containsSHA(list []string, sha string) bool {
	for _, p := range list {
		if p == sha {
			return true
		}
	}
	return false
}

// ReadCommitMeta reads tree + parent SHAs from a commit object (parents do
// not need to be present in the object store).
func ReadCommitMeta(cacheDir, sha string) (CommitMeta, error) {
	var meta CommitMeta
	if cacheDir == "" || sha == "" {
		return meta, fmt.Errorf("read commit meta: empty cache or sha")
	}
	tree, err := gitOutput(cacheDir, "rev-parse", sha+"^{tree}")
	if err != nil {
		return meta, fmt.Errorf("tree %s: %w", sha, err)
	}
	meta.Tree = strings.TrimSpace(tree)
	body, err := gitOutput(cacheDir, "cat-file", "-p", sha)
	if err != nil {
		return meta, fmt.Errorf("cat-file %s: %w", sha, err)
	}
	for _, line := range strings.Split(body, "\n") {
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "parent ") {
			meta.Parents = append(meta.Parents, strings.TrimSpace(strings.TrimPrefix(line, "parent ")))
		}
	}
	return meta, nil
}

// MergeBase returns git merge-base of left and right when both are connected.
func MergeBase(cacheDir, left, right string) (string, bool) {
	if cacheDir == "" || left == "" || right == "" {
		return "", false
	}
	out, err := gitOutput(cacheDir, "merge-base", left, right)
	if err != nil {
		return "", false
	}
	mb := strings.TrimSpace(out)
	return mb, mb != ""
}

var reviewDeepenPhases = []TargetedFetchOptions{
	{Deepen: 1, HistoryPhase: "review_deepen_1"},
	{Deepen: 16, HistoryPhase: "review_deepen_16"},
	{Deepen: 32, HistoryPhase: "review_deepen_32"},
	{Deepen: 64, HistoryPhase: "review_deepen_64"},
	{Deepen: 256, HistoryPhase: "review_deepen_256"},
	{Unshallow: true, HistoryPhase: "review_unshallow"},
}

// EnsureReviewBase fetches the merge-base of PlanReviewPair(source, base) into
// the targeted cache (blob:none deepen, not a full clone).
func EnsureReviewBase(ctx context.Context, cloneURL, cacheDir, sourceSHA, baseSHA string, opts TargetedFetchOptions) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if sourceSHA == "" {
		return baseSHA, nil
	}
	if baseSHA == "" {
		return "", nil
	}

	srcMeta, err := ReadCommitMeta(cacheDir, sourceSHA)
	if err != nil {
		return "", err
	}
	baseMeta, err := ReadCommitMeta(cacheDir, baseSHA)
	if err != nil {
		return "", err
	}
	left, right := PlanReviewPair(sourceSHA, srcMeta, baseSHA, baseMeta)
	if opts.Progress != nil {
		fmt.Fprintf(opts.Progress, "review pair %s %s\n", left, right)
	}

	if err := fetchReviewTips(ctx, cloneURL, cacheDir, []string{left, right}, opts); err != nil {
		return "", err
	}
	if mb, ok := MergeBase(cacheDir, left, right); ok {
		return mb, nil
	}

	for _, phase := range reviewDeepenPhases {
		phase.Auth = opts.Auth
		phase.Env = opts.Env
		phase.ReposRoot = opts.ReposRoot
		phase.Progress = opts.Progress
		phase.TargetRefPrefix = opts.TargetRefPrefix
		if _, err := FetchTargetedRefs(ctx, cloneURL, []string{left, right}, phase); err != nil {
			return "", err
		}
		if mb, ok := MergeBase(cacheDir, left, right); ok {
			return mb, nil
		}
	}
	// Deepen could not connect. First parent (left) is still the useful review left side.
	if left != sourceSHA && left != baseSHA {
		return left, nil
	}
	return "", fmt.Errorf("no merge-base for %s and %s after deepen", left, right)
}

func fetchReviewTips(ctx context.Context, cloneURL, cacheDir string, shas []string, opts TargetedFetchOptions) error {
	missing := make([]string, 0, len(shas))
	seen := map[string]bool{}
	for _, sha := range shas {
		if sha == "" || seen[sha] {
			continue
		}
		seen[sha] = true
		if _, err := gitOutput(cacheDir, "rev-parse", "--verify", sha+"^{commit}"); err == nil {
			if opts.Progress != nil {
				fmt.Fprintf(opts.Progress, "have %s\n", sha)
			}
			continue
		}
		missing = append(missing, sha)
	}
	if len(missing) == 0 {
		return nil
	}
	fetchOpts := TargetedFetchOptions{
		Auth:            opts.Auth,
		Env:             opts.Env,
		ReposRoot:       opts.ReposRoot,
		Progress:        opts.Progress,
		TargetRefPrefix: opts.TargetRefPrefix,
		Depth:           1,
	}
	_, err := FetchTargetedRefs(ctx, cloneURL, missing, fetchOpts)
	return err
}
