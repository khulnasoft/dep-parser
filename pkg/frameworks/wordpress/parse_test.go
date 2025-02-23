package wordpress

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/dep-parser/pkg/types"
)

// Helper function to open a file and return its handle
func openFile(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.Open(path)
	require.NoError(t, err)
	return f
}

func TestParseWordPress(t *testing.T) {
	tests := []struct {
		file    string // Test input file
		want    types.Library
		wantErr string
	}{
		{
			file: "testdata/version.php",
			want: types.Library{
				Name:    "wordpress",
				Version: "4.9.4-alpha",
			},
		},
		{
			file:    "testdata/versionFail.php",
			wantErr: "version.php could not be parsed",
		},
	}

	for _, tt := range tests {
		t.Run(path.Base(tt.file), func(t *testing.T) {
			f := openFile(t, tt.file)
			defer f.Close()

			got, err := Parse(f)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
