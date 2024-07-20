package git

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

// go test -run TestListCommitRelativeToBase -v ./git/gitops
func TestListCommitRelativeToBase(t *testing.T) {
	// exists, merged, commits, err := ListCommitRelativeToBase(repoDir, "origin/release-v1.11.0", "origin/master2") // exists = false
	exists, merged, commits, err := ListCommitRelativeToBase(repoDir, "origin/release-v1.11.0", "origin/master") // merged = true, 1 commits
	// exists, merged, commits, err := ListCommitRelativeToBase(repoDir, "origin/dev-1.15.0", "origin/master") // merged = false,have 2 commits
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("exists=%v, merged=%v", exists, merged)
	commitsJSON, _ := json.Marshal(commits)
	t.Logf("%s", string(commitsJSON))
	ioutil.WriteFile("/tmp/test-res.json", commitsJSON, 0777)
}

// go test -run TestResolveDiffCommit -v ./git/gitops
func TestResolveDiffCommit(t *testing.T) {
	// exists, merged, headCommit, baseCommit, err := ResolveDiffCommit(repoDir, "origin/release-v1.11.0", "origin/master2") // exists = false
	// exists, merged, headCommit, baseCommit, err := ResolveDiffCommit(repoDir, "origin/release-v1.11.0", "origin/master") // merged = true, headCommit,baseCommit
	exists, merged, headCommit, baseCommit, err := ResolveDiffCommit(repoDir, "origin/dev-1.15.0", "origin/master") // merged = false,have 2 commits
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("exists=%v, merged=%v,headCommit=%v,baseCommit=%v", exists, merged, headCommit, baseCommit)
}
