package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PathFromReader(r io.Reader) (string, error) {
	var (
		lines   []string
		path    = ""
		scanner = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	for _, line := range lines {
		if !shouldIgnore(line) {
			cleaned := filepath.Clean(line)
			if !strings.Contains(path, cleaned) {
				path = fmt.Sprintf("%s:%s", path, cleaned)
			}
		}
	}

	return path, nil
}

func PathFromFile(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}

	return PathFromReader(f)
}

func shouldIgnore(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
