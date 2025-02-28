package lockfile

import (
	"io"
	"sort"

	"github.com/liamg/jfather"
	"golang.org/x/exp/maps"
	"golang.org/x/xerrors"

	dio "github.com/khulnasoft/dep-parser/pkg/io"
	"github.com/khulnasoft/dep-parser/pkg/types"
	"github.com/khulnasoft/dep-parser/pkg/utils"
)

// lockfile format defined at: https://stringbean.github.io/sbt-dependency-lock/file-formats/version-1.html
type sbtLockfile struct {
	Version      int                     `json:"lockVersion"`
	Dependencies []sbtLockfileDependency `json:"dependencies"`
}

type sbtLockfileDependency struct {
	Organization   string   `json:"org"`
	Name           string   `json:"name"`
	Version        string   `json:"version"`
	Configurations []string `json:"configurations"`
	StartLine      int
	EndLine        int
}

type Parser struct{}

func NewParser() types.Parser {
	return &Parser{}
}

func (p *Parser) Parse(r dio.ReadSeekerAt) ([]types.Library, []types.Dependency, error) {
	var lockfile sbtLockfile
	input, err := io.ReadAll(r)

	if err != nil {
		return nil, nil, xerrors.Errorf("failed to read sbt lockfile: %w", err)
	}
	if err := jfather.Unmarshal(input, &lockfile); err != nil {
		return nil, nil, xerrors.Errorf("JSON decoding failed: %w", err)
	}

	libs := map[string]types.Library{}

	for _, dep := range lockfile.Dependencies {
		if containsConfig(dep.Configurations, "compile") || containsConfig(dep.Configurations, "runtime") {
			name := dep.Organization + ":" + dep.Name
			lib := types.Library{
				ID:      utils.PackageID(name, dep.Version),
				Name:    name,
				Version: dep.Version,
				Locations: []types.Location{
					{
						StartLine: dep.StartLine,
						EndLine:   dep.EndLine,
					},
				},
			}
			libs[lib.Name] = lib
		}
	}

	// Convert map to slice
	libSlice := maps.Values(libs)
	sort.Sort(types.Libraries(libSlice))

	return libSlice, nil, nil
}

// UnmarshalJSONWithMetadata needed to detect start and end lines of deps
func (t *sbtLockfileDependency) UnmarshalJSONWithMetadata(node jfather.Node) error {
	if err := node.Decode(&t); err != nil {
		return err
	}
	t.StartLine = node.Range().Start.Line
	t.EndLine = node.Range().End.Line
	return nil
}

// Alternative to slices.ContainsFunc for Go < 1.21
func containsConfig(configs []string, match string) bool {
	for _, c := range configs {
		if c == match {
			return true
		}
	}
	return false
}
