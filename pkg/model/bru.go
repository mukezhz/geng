package model

import "fmt"

type BruModel struct {
	Route       string
	Method      string
	Body        string
	Handler     string
	ModuleName  string
	Name        string
	Description string
}

func (b *BruModel) GetModulePath(f string) string {
	return "domain/" + b.ModuleName + "/" + fmt.Sprintf("%s.%s", b.ModuleName, f)
}
