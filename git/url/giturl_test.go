package url

import (
	"strings"
	"testing"
)

// go test -run TestSplitRepo -v ./src/util
func TestSplitRepo(t *testing.T) {
	testCases := []*struct {
		Name      string
		RepoURL   string
		Domain    string
		Path      string
		ExpectErr string
	}{
		{
			Name:      "param error",
			RepoURL:   "",
			ExpectErr: "empty repoURL",
		},
		{
			Name:      "incomplete https",
			RepoURL:   "https://gitlab:TOKEN@xxxx",
			ExpectErr: "unrecognized repoURL",
		},
		{
			Name:      "incomplete https",
			RepoURL:   "https://gitlab:TOKEN@xxxx.git",
			ExpectErr: "unrecognized repoURL",
		},
		{
			Name:      "incomplete https",
			RepoURL:   "https://gitlab:TOKEN@aaa.com",
			ExpectErr: "unrecognized repoURL",
		},
		{
			Name:      "incomplete https",
			RepoURL:   "https://aaa.com",
			ExpectErr: "unrecognized repoURL",
		},
		{
			Name:    "plain https",
			RepoURL: "https://aaa.com/a.git",
			Domain:  "aaa.com",
			Path:    "/a",
		},
		{
			Name:      "token https",
			RepoURL:   "https://my:TOKEN@aaa.com/a.git",
			Domain:    "aaa.com",
			Path:      "/a",
			ExpectErr: "",
		},
		{
			Name:      "incomplete ssh",
			RepoURL:   "gitlab@xxxx",
			ExpectErr: "",
		},
		{
			Name:      "incomplete ssh",
			RepoURL:   "gitlab@xxxx.git",
			ExpectErr: "",
		},
		{
			Name:    "ssh",
			RepoURL: "gitlab@aaa.com:a.git",
			Domain:  "aaa.com",
			Path:    "/a",
		},
		{
			Name:      "ssh with num port",
			RepoURL:   "gitlab@aaa.com:2234/a.git",
			Domain:    "aaa.com:2234",
			Path:      "/a",
			ExpectErr: "",
		},
		{
			Name:      "ssh with empty port, has slash",
			RepoURL:   "gitlab@aaa.com:/a.git",
			Domain:    "aaa.com",
			Path:      "/a",
			ExpectErr: "",
		},
		{
			Name:      "ssh with empty port",
			RepoURL:   "gitlab@aaa.com:/a.git",
			Domain:    "aaa.com",
			Path:      "/a",
			ExpectErr: "",
		},
		{
			Name:      "an actual https example",
			RepoURL:   "https://git.some/sp/some-service/be/pricing.git",
			Domain:    "git.some",
			Path:      "/sp/some-service/be/pricing",
			ExpectErr: "",
		},
		{
			Name:      "an actual https token example",
			RepoURL:   "https://gitlab-ci-token:SOM3243TOKEN@git.some/sp/some-service/be/pricing",
			Domain:    "git.some",
			Path:      "/sp/some-service/be/pricing",
			ExpectErr: "",
		},
		{
			Name:      "with token",
			RepoURL:   "https://gitlab:XYZ-@git.some/sp/some-service/be/pricing",
			Domain:    "git.some",
			Path:      "/sp/some-service/be/pricing",
			ExpectErr: "",
		},
		{
			Name:      "an actual ssh example",
			RepoURL:   "gitlab@git.some:sp/some-service/be/pricing",
			Domain:    "git.some",
			Path:      "/sp/some-service/be/pricing",
			ExpectErr: "",
		},
		{
			Name:      "h5-with-leading-space",
			RepoURL:   " https://git.some/sp/some-service/fe/merch_loan/h5",
			Domain:    "git.some",
			Path:      "/sp/some-service/fe/merch_loan/h5",
			ExpectErr: "unrecognized repoURL",
		},
	}

	for i, testCase := range testCases {
		t.Logf(">>>>>>>>>>\nRUN  case[%d]: %s -> %s", i, testCase.Name, testCase.RepoURL)

		domain, path, err := SplitRepoURL(testCase.RepoURL)
		if testCase.ExpectErr != "" {
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if !strings.Contains(errMsg, testCase.ExpectErr) {
				t.Fatalf("case[%d] call expect err:%v, actual:%v", i, testCase.ExpectErr, errMsg)
			}
			t.Logf("PASS case[%d]: %s", i, testCase.Name)
			continue
		}

		if domain != testCase.Domain {
			t.Fatalf("expect %s = %+v, actual:%+v", `domain`, testCase.Domain, domain)
		}
		if path != testCase.Path {
			t.Fatalf("expect %s = %+v, actual:%+v", `path`, testCase.Path, path)
		}

		t.Logf("PASS case[%d]: %s", i, testCase.Name)
	}
}
