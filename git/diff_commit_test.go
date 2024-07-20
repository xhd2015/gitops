package git

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

// go test -run TestDiffCommit -v ./git/gitops
func TestDiffCommit(t *testing.T) {
	details, err := DiffCommit(repoDir, "HEAD~8", "HEAD~9", nil)
	if err != nil {
		t.Fatal(err)
	}
	commitsJSON, _ := json.Marshal(details)
	t.Logf("%s", string(commitsJSON))
	ioutil.WriteFile("/tmp/test-res.json", commitsJSON, 0777)
}
