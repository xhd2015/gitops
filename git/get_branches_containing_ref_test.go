package git

import (
	"reflect"
	"testing"
)

func TestGetBranchesContainingRef(t *testing.T) {
	tests := []struct {
		name         string
		dir          string
		ref          string
		skip         bool
		wantBranches []string
		wantErr      bool
	}{
		{
			name: "merged branch",
			// dir:          "TODO-make a git dir containing such merged branches",
			// skip:         true,
			dir:          "TODO",
			ref:          "v2.50.2.t12",
			wantBranches: []string{"release-v2.16.0"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
				return
			}
			gotBranches, err := GetBranchesContainingRef(tt.dir, tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBranchesContainingRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBranches, tt.wantBranches) {
				t.Errorf("GetBranchesContainingRef() = %v, want %v", gotBranches, tt.wantBranches)
			}
		})
	}
}
