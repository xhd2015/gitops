package git

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

var repoDir string

func init() {
	repoDir = os.Getenv("TEST_REPO_DIR")
	if repoDir == "" {
		// panic(fmt.Errorf("requires TEST_REPO_DIR"))
	}
}

// go test -run TestRevListAll -v ./git/gitops
func TestRevListAll(t *testing.T) {
	commits, err := RevListAll(repoDir, "")
	if err != nil {
		t.Fatal(err)
	}
	commitsJSON, _ := json.Marshal(commits)
	t.Logf("%s", string(commitsJSON))
	ioutil.WriteFile("test-res", commitsJSON, 0777)
}
