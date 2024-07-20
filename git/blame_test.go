package git

import "testing"

func TestParseBlame(t *testing.T) {
	line := "^c9e7754 (<some@some.com> 1662996314 +0800 2) const root = createRoot(rootElement)"

	blameInfo, err := parsePlainBlame(line)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("blameInfo: %+v", blameInfo)
	if blameInfo.Line != 2 {
		t.Fatalf("expect %s = %+v, actual:%+v", `blameInfo.Line`, 2, blameInfo.Line)
	}
	if blameInfo.CommitHash != "c9e7754" {
		t.Fatalf("expect %s = %+v, actual:%+v", `blameInfo.CommitHash`, "c9e7754", blameInfo.CommitHash)
	}

	if blameInfo.AuthorEmail != "some@some.com" {
		t.Fatalf("expect %s = %+v, actual:%+v", `blameInfo.AuthorEmail`, "some@some.com", blameInfo.AuthorEmail)
	}
	if blameInfo.Timestamp != 1662996314 {
		t.Fatalf("expect %s = %+v, actual:%+v", `blameInfo.Timestamp`, 1662996314, blameInfo.Timestamp)
	}
	if !blameInfo.Boundary {
		t.Fatalf("expect %s = %+v, actual:%+v", `blameInfo.Boundary`, true, blameInfo.Boundary)
	}

}
