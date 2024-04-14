package utility

import (
	"runtime"
	"strings"
)

func IgnoreWindowsPath(p string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(p, "\\", "/")
	}
	return p
}
