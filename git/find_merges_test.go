package git

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

// go test -run TestFindMergePoints -v ./git/gitops
func TestFindMergePoints(t *testing.T) {
	testFindMergePoints(t, os.Getenv("REF"), os.Getenv("BASE"))
}

func testFindMergePoints(t *testing.T, ref string, base string) {
	commits, err := FindMergePoints(repoDir, ref, base)
	if err != nil {
		t.Fatal(err)
	}
	commitsJSON, _ := json.Marshal(commits)
	t.Logf("%s", string(commitsJSON))
	ioutil.WriteFile("/tmp/test-res.json", commitsJSON, 0777)
}
