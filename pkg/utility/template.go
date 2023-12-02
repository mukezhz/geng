package utility

import (
	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/model"
	"os"
	"text/template"
)

func GenerateFromTemplate(templateFile, outputFile string, data model.ModuleData) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			color.Redln("Error closing:", err)
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}
