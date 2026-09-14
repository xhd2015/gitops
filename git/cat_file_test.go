package git

import (
	"testing"
)

func TestCatFile(t *testing.T) {
	type args struct {
		ref  string
		file string
	}
	tests := []struct {
		name        string
		args        args
		wantOk      bool
		wantContent string
		wantErr     bool
	}{
		{
			name: "file ok",
			args: args{
				ref:  "master",
				file: "README.md",
			},
			wantOk:      true,
			wantContent: "test",
		},
		{
			name: "file missing",
			args: args{
				ref:  "master",
				file: "NOT_EXISTS.md",
			},
			wantOk: false,
		},
		{
			name: "working file ok",
			args: args{
				ref:  COMMIT_WORKING,
				file: "README.md",
			},
			wantOk:      true,
			wantContent: "test",
		},
		{
			name: "working path escape",
			args: args{
				ref:  COMMIT_WORKING,
				file: "../secret.txt",
			},
			wantErr: true,
		},
		{
			name: "abs path rejected",
			args: args{
				ref:  COMMIT_WORKING,
				file: "/etc/passwd",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, cleanUp := mustGetTmpDir()
			defer cleanUp()

			gotOk, gotContent, err := CatFile(tmpDir, tt.args.ref, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("CatFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("CatFile() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotContent != tt.wantContent {
				t.Errorf("CatFile() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}
