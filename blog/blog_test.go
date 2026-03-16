package blog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateOutputDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		root    string
		outDir  string
		wantErr bool
		errMsg  string
	}{
		{"valid subdirectory", "/project", "/project/public", false, ""},
		{"empty output dir", "/project", "", true, "output directory cannot be empty"},
		{"same as root", "/project", "/project", true, "output directory cannot be the project root"},
		{"outside root", "/project", "/tmp/output", true, "output directory must be inside the project root"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			err := validateOutputDir(tt.root, tt.outDir)
			if tt.wantErr {
				r.Error(err)
				r.Contains(err.Error(), tt.errMsg)
			} else {
				r.NoError(err)
			}
		})
	}
}
