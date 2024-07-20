package git

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/xhd2015/gitops/model"
)

// go test -run TestGrepLine -v ./git/gitops
func TestGrepLine(t *testing.T) {
	lines, err := GrepLines(repoDir, "origin/master", &model.GrepLineOptions{
		IgnoreCase: true,
		Patterns:   []string{"interestRate"},
		Files:      []string{"src/*.go"},
	})
	if err != nil {
		t.Fatal(err)
	}
	commitsJSON, _ := json.Marshal(lines)
	t.Logf("%s", string(commitsJSON))
	ioutil.WriteFile("/tmp/test-res.json", commitsJSON, 0777)
}
