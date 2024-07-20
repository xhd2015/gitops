package git

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

// go test -run TestListBranchRef -v ./git/gitops
func TestListBranchRef(t *testing.T) {
	commits, err := ListBranchRef(repoDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	commitsJSON, _ := json.Marshal(commits)
	t.Logf("%s", string(commitsJSON))
	ioutil.WriteFile("/tmp/test-res.json", commitsJSON, 0777)
}
