package swift

import (
	"github.com/khulnasoft/dep-parser/pkg/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Helper function to open a file and return its handle
func openFile(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.Open(path)
	assert.NoError(t, err)
	return f
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      []types.Library
	}{
		{
			name:      "happy path v1",
			inputFile: "testdata/happy-v1-Package.resolved",
			want: []types.Library{
				{
					ID:        "github.com/Quick/Nimble@9.2.1",
					Name:      "github.com/Quick/Nimble",
					Version:   "9.2.1",
					Locations: []types.Location{{StartLine: 4, EndLine: 12}},
				},
				{
					ID:        "github.com/ReactiveCocoa/ReactiveSwift@7.1.1",
					Name:      "github.com/ReactiveCocoa/ReactiveSwift",
					Version:   "7.1.1",
					Locations: []types.Location{{StartLine: 13, EndLine: 21}},
				},
			},
		},
		{
			name:      "happy path v2",
			inputFile: "testdata/happy-v2-Package.resolved",
			want: []types.Library{
				{
					ID:        "github.com/Quick/Nimble@9.2.1",
					Name:      "github.com/Quick/Nimble",
					Version:   "9.2.1",
					Locations: []types.Location{{StartLine: 21, EndLine: 29}},
				},
				{
					ID:        "github.com/Quick/Quick@7.2.0",
					Name:      "github.com/Quick/Quick",
					Version:   "7.2.0",
					Locations: []types.Location{{StartLine: 30, EndLine: 38}},
				},
				{
					ID:        "github.com/ReactiveCocoa/ReactiveSwift@7.1.1",
					Name:      "github.com/ReactiveCocoa/ReactiveSwift",
					Version:   "7.1.1",
					Locations: []types.Location{{StartLine: 39, EndLine: 47}},
				},
				{
					ID:        "github.com/mattgallagher/CwlCatchException@2.1.2",
					Name:      "github.com/mattgallagher/CwlCatchException",
					Version:   "2.1.2",
					Locations: []types.Location{{StartLine: 3, EndLine: 11}},
				},
				{
					ID:        "github.com/mattgallagher/CwlPreconditionTesting@2.1.2",
					Name:      "github.com/mattgallagher/CwlPreconditionTesting",
					Version:   "2.1.2",
					Locations: []types.Location{{StartLine: 12, EndLine: 20}},
				},
			},
		},
		{
			name:      "empty",
			inputFile: "testdata/empty-Package.resolved",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			f := openFile(t, tt.inputFile)
			defer f.Close()

			libs, _, err := parser.Parse(f)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, libs)
		})
	}
}
