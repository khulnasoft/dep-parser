package uv

import (
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"

	dio "github.com/khulnasoft/dep-parser/pkg/io"
	"github.com/khulnasoft/dep-parser/pkg/log"
	"github.com/khulnasoft/dep-parser/pkg/types"
	"github.com/khulnasoft/dep-parser/pkg/utils"
)

type Lockfile struct {
	Packages []struct {
		Name        string                 `toml:"name"`
		Version     string                 `toml:"version"`
		Source      Source                 `toml:"source"`
		Dependencies map[string]interface{} `toml:"dependencies"`
		Dev         bool                   `toml:"dev"`
	} `toml:"package"`
}

type Source struct {
	Editable string `toml:"editable"`
	Virtual  string `toml:"virtual"`
}

type Parser struct{}

func NewParser() types.Parser {
	return &Parser{}
}

func (p *Parser) Parse(r dio.ReadSeekerAt) ([]types.Library, []types.Dependency, error) {
	var lockfile Lockfile
	if _, err := toml.NewDecoder(r).Decode(&lockfile); err != nil {
		return nil, nil, xerrors.Errorf("failed to decode uv lock file: %w", err)
	}

	libVersions := parseVersions(lockfile)

	var libs []types.Library
	var deps []types.Dependency
	for _, pkg := range lockfile.Packages {
		if pkg.Dev {
			continue
		}

		pkgID := utils.PackageID(pkg.Name, pkg.Version)
		libs = append(libs, types.Library{
			ID:      pkgID,
			Name:    pkg.Name,
			Version: pkg.Version,
		})

		dependsOn := parseDependencies(pkg.Dependencies, libVersions)
		if len(dependsOn) != 0 {
			deps = append(deps, types.Dependency{
				ID:        pkgID,
				DependsOn: dependsOn,
			})
		}
	}
	return libs, deps, nil
}

func parseVersions(lockfile Lockfile) map[string][]string {
	libVersions := map[string][]string{}
	for _, pkg := range lockfile.Packages {
		if pkg.Dev {
			continue
		}
		libVersions[pkg.Name] = append(libVersions[pkg.Name], pkg.Version)
	}
	return libVersions
}

func parseDependencies(deps map[string]any, libVersions map[string][]string) []string {
	var dependsOn []string
	for name, versRange := range deps {
		if dep, err := parseDependency(name, versRange, libVersions); err != nil {
			log.Logger.Debugf("failed to parse uv dependency: %s", err)
		} else if dep != "" {
			dependsOn = append(dependsOn, dep)
		}
	}
	sort.Strings(dependsOn)
	return dependsOn
}

func parseDependency(name string, versRange any, libVersions map[string][]string) (string, error) {
	name = normalizePkgName(name)
	vers, ok := libVersions[name]
	if !ok {
		return "", xerrors.Errorf("no version found for %q", name)
	}

	for _, ver := range vers {
		var vRange string
		switch r := versRange.(type) {
		case string:
			vRange = r
		case map[string]interface{}:
			if v, exists := r["version"]; exists {
				vRange = v.(string)
			}
		}

		if matched, err := matchVersion(ver, vRange); err != nil {
			return "", xerrors.Errorf("failed to match version for %s: %w", name, err)
		} else if matched {
			return utils.PackageID(name, ver), nil
		}
	}
	return "", xerrors.Errorf("no matched version found for %q", name)
}

func normalizePkgName(name string) string {
	return strings.ToLower(name)
}
