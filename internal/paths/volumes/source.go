package volumes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/frantjc/sequence/pkg/github/actions"
)

var regExp = regexp.MustCompile("[^a-zA-Z0-9_.-]")

func normalize(s string) string {
	return strings.TrimPrefix(regExp.ReplaceAllLiteralString(s, "-"), "-")
}

func GetActionSource(action *actions.Reference) string {
	return normalize(fmt.Sprintf("sqnc-actions-%s-%s-%s", action.Owner, action.Repository, action.Version))
}

func GetWorkspace(s string) string {
	return normalize(fmt.Sprintf("%s-workspace", s))
}

func GetRunnerTmp(s string) string {
	return normalize(fmt.Sprintf("%s-runner-tmp", s))
}

func GetRunnerToolCache(s string) string {
	return normalize(fmt.Sprintf("%s-runner-toolcache", s))
}

func GetGitHub(s string) string {
	return normalize(fmt.Sprintf("%s-github", s))
}
