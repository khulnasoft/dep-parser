package lockfile

import (
	"io"
	"sort"
	"strings"

	"github.com/liamg/jfather"
	"golang.org/x/exp/maps"
	"golang.org/x/xerrors"

	dio "github.com/khulnasoft/dep-parser/pkg/io"
	"github.com/khulnasoft/dep-parser/pkg/log"
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

func NewParser() *Parser {
	return &Parser{}
}

func (Parser) Parse(r dio.ReadSeekerAt) ([]types.Package, []types.Dependency, error) {
	var lockfile sbtLockfile
	input, err := io.ReadAll(r)

	if err != nil {
		return nil, nil, xerrors.Errorf("failed to read sbt lockfile: %w", err)
	}
	if err := jfather.Unmarshal(input, &lockfile); err != nil {
		return nil, nil, xerrors.Errorf("JSON decoding failed: %w", err)
	}

	var libraries types.Packages

	for _, dep := range lockfile.Dependencies {
		if slices.ContainsFunc(dep.Configurations, isIncludedConfig) {
			name := dep.Organization + ":" + dep.Name
			libraries = append(libraries, types.Package{
				ID:      dependency.ID(types.Sbt, name, dep.Version),
				Name:    name,
				Version: dep.Version,
				Locations: []types.Location{
					{
						StartLine: dep.StartLine,
						EndLine:   dep.EndLine,
					},
				},
			})
		}
	}

	sort.Sort(libraries)
	return libraries, nil, nil
}

// UnmarshalJSONWithMetadata needed to detect start and end lines of deps
func (t *sbtLockfileDependency) UnmarshalJSONWithMetadata(node jfather.Node) error {
	if err := node.Decode(&t); err != nil {
		return err
	}
	// Decode func will overwrite line numbers if we save them first
	t.StartLine = node.Range().Start.Line
	t.EndLine = node.Range().End.Line
	return nil
}

func isIncludedConfig(config string) bool {
	return config == "compile" || config == "runtime"
}
