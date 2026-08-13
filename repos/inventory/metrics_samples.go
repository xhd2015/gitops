package inventory

// Metric base names (v1). Prefix applied by stats SDK / scrape env elsewhere.
const (
	MetricGitCachePresent       = "git_cache_present"
	MetricGitCacheRootBytes     = "git_cache_root_bytes"
	MetricGitCacheBytes         = "git_cache_bytes"
	MetricGitCacheRepoCount     = "git_cache_repo_count"
	MetricGitCacheWorktreeCount = "git_cache_worktree_count"
	MetricGitCacheScanSeconds   = "git_cache_scan_seconds"
	MetricGitCacheIncomplete    = "git_cache_incomplete"

	LabelLayout = "layout"

	// LayoutLabelTmpWorktrees is the layout label for tmp (not a path-segment constant).
	LayoutLabelTmpWorktrees = "tmp-worktrees"
)

// MetricSample is one pure gauge plan entry (no Prometheus types).
type MetricSample struct {
	Name   string            // short metric base, e.g. git_cache_bytes
	Labels map[string]string // low cardinality only
	Value  float64
}

// sealedLayoutLabels is the emit order for layout-labelled series (unified four + tmp).
var sealedLayoutLabels = []string{
	LayoutStaticCache,
	LayoutWorktrees,
	LayoutCache,
	LayoutTargetedCache,
	LayoutLabelTmpWorktrees,
}

// MetricSamples materializes the sealed v1 series plan from Snapshot.
// Always returns 19 samples in stable order. Never returns nil.
// snap.Err is ignored.
func MetricSamples(snap Snapshot) []MetricSample {
	samples := make([]MetricSample, 0, 19)

	present := 0.0
	if snap.Present {
		present = 1
	}
	samples = append(samples, MetricSample{
		Name:  MetricGitCachePresent,
		Value: present,
	})

	samples = append(samples, MetricSample{
		Name:  MetricGitCacheRootBytes,
		Value: float64(snap.TotalBytes),
	})

	for _, layout := range sealedLayoutLabels {
		st := layoutStats(snap, layout)
		samples = append(samples,
			MetricSample{
				Name:   MetricGitCacheBytes,
				Labels: map[string]string{LabelLayout: layout},
				Value:  float64(st.Bytes),
			},
			MetricSample{
				Name:   MetricGitCacheRepoCount,
				Labels: map[string]string{LabelLayout: layout},
				Value:  float64(st.RepoCount),
			},
			MetricSample{
				Name:   MetricGitCacheWorktreeCount,
				Labels: map[string]string{LabelLayout: layout},
				Value:  float64(st.WorktreeCount),
			},
		)
	}

	samples = append(samples, MetricSample{
		Name:  MetricGitCacheScanSeconds,
		Value: snap.ScanDuration.Seconds(),
	})

	incomplete := 0.0
	if snap.Incomplete {
		incomplete = 1
	}
	samples = append(samples, MetricSample{
		Name:  MetricGitCacheIncomplete,
		Value: incomplete,
	})

	return samples
}

func layoutStats(snap Snapshot, layout string) LayoutStats {
	if layout == LayoutLabelTmpWorktrees {
		return snap.TmpWorktrees
	}
	if snap.Layouts == nil {
		return LayoutStats{}
	}
	return snap.Layouts[layout]
}
