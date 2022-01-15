package actions

import (
	"strings"
)

func Expand(s string, mapping func(string) string) string {
	return string(ExpandBytes([]byte(s), mapping))
}

func ExpandBytes(b []byte, mapping func(string) string) (p []byte) {
	i := 0
	for j := 0; j < len(b); j++ {
		if b[j] == '$' {
			if p == nil {
				p = make([]byte, 0, 2*len(b))
			}
			p = append(p, b[i:j]...)
			name, w := getGitHubName(b[j+1:])
			if name == "" && w > 0 {
				// encountered invalid syntax; eat the characters
			} else if name == "" {
				// valid syntax, but ${{ }} contained no name
				p = append(p, b[j])
			} else {
				p = append(p, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if p == nil {
		return b
	} else if i >= len(b) {
		i = len(b)
	}
	return append(p, b[i:]...)
}

func getGitHubName(b []byte) (s string, w int) {
	switch {
	case len(b) > 3:
		if b[0] == '{' && b[1] == '{' {
			i := 2
			for ; i+1 < len(b) && b[i] != '}'; i++ {
			}
			if b[i] == '}' && i+1 < len(b) && b[i+1] != '}' {
				return "", 0 // bad syntax
			}
			return strings.TrimSpace(string(b[2:i])), i + 2
		}
	}

	return "", 0
}
