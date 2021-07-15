package ipc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindSocketPath(chain string, kind string, baseDir string) string {
	name := fmt.Sprintf("%s-%s", chain, kind)
	found := ""

	filepath.Walk(baseDir, func(fullPath string, info os.FileInfo, err error) error {
		if found == "" && strings.Contains(info.Name(), name) {
			found = fullPath
		}
		return nil
	})

	return found
}

func FindDesicionsSocketPath(chain string, baseDir string) string {
	return FindSocketPath(chain, TypeDesicions, baseDir)
}
