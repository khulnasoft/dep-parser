package npm

import (
	"os"
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

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		file     string // Test input file
		want     []types.Library
		wantDeps []types.Dependency
	}{
		{
			name:     "lock version v1",
			file:     "testdata/package-lock_v1.json",
			want:     npmV1Libs,
			wantDeps: npmDeps,
		},
		{
			name:     "lock version v2",
			file:     "testdata/package-lock_v2.json",
			want:     npmV2Libs,
			wantDeps: npmDeps,
		},
		{
			name:     "lock version v3",
			file:     "testdata/package-lock_v3.json",
			want:     npmV2Libs,
			wantDeps: npmDeps,
		},
		{
			name:     "lock version v3 with workspace",
			file:     "testdata/package-lock_v3_with_workspace.json",
			want:     npmV3WithWorkspaceLibs,
			wantDeps: npmV3WithWorkspaceDeps,
		},
		{
			name:     "lock version v3 with workspace and without direct deps field",
			file:     "testdata/package-lock_v3_without_root_deps_field.json",
			want:     npmV3WithoutRootDepsField,
			wantDeps: npmV3WithoutRootDepsFieldDeps,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := openFile(t, tt.file)
			defer f.Close()

			got, deps, err := NewParser().Parse(f)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
			if tt.wantDeps != nil {
				assert.Equal(t, tt.wantDeps, deps)
			}
		})
	}
}
