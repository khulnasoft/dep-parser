package jar

import (
	"fmt"
	"github.com/khulnasoft/dep-parser/pkg/types"
	"golang.org/x/xerrors"
)

var ArtifactNotFoundErr = xerrors.New("no artifact found")

type Properties struct {
	GroupID    string
	ArtifactID string
	Version    string
	FilePath   string // path to file containing these props
}

func (p Properties) Library() types.Library {
	return types.Library{
		Name:     fmt.Sprintf("%s:%s", p.GroupID, p.ArtifactID),
		Version:  p.Version,
		FilePath: p.FilePath,
	}
}

func (p Properties) Valid() bool {
	return p.GroupID != "" && p.ArtifactID != "" && p.Version != ""
}

func (p Properties) String() string {
	return fmt.Sprintf("%s:%s:%s", p.GroupID, p.ArtifactID, p.Version)
}
