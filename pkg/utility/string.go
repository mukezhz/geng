package utility

import (
	"fmt"
	"path"
	"strings"
)

func PrependStringInPath(originalPath, prefix string) string {
	dir, fileName := path.Split(originalPath)

	if strings.HasPrefix(fileName, prefix) {
		return originalPath
	}

	newFileName := fmt.Sprintf("%s.%s", prefix, fileName)

	newPath := path.Join(dir, newFileName)
	return newPath
}
