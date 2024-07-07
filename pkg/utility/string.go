package utility

import (
	"fmt"
	"path"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func ToPasalCase(s string) string {
	s = strings.TrimSpace(s)
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return ""
	}
	titleCaser := cases.Title(language.Und)
	for i := 0; i < len(parts); i++ {
		parts[i] = titleCaser.String(parts[i])
	}

	return strings.Join(parts, "")
}

func ToCamelCase(s string) string {
	pascalCase := ToPasalCase(s)
	return fmt.Sprintf("%s%s", strings.ToLower(string(pascalCase[0])), pascalCase[1:])
}

func SanitizeEndpoint(endpoint string) string {
	splitted := strings.Split(endpoint, "/")
	splitted = splitted[1:]

	return "/" + strings.Join(splitted, "/")
}
