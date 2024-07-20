package git

import "testing"

// go test -run TestListFile -v ./git/gitops
func TestListFile(t *testing.T) {
	files, err := ListFile("../..", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("files: %+v", files)
}
