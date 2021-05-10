package templates

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
)

type FileDetails struct {
	// Target go file to create
	TargetFile string
	// Package name (default = last directory in path)
	Package string
}

func CreateTemplate(templateContent string, fileDetails FileDetails, details interface{}) (err error) {
	var t *template.Template
	funcMap := createFunctions()
	if t, err = template.New("").Funcs(funcMap).Parse(templateContent); err == nil {
		var outputFile *os.File
		if outputFile, err = os.Create(fileDetails.TargetFile); err == nil {
			defer outputFile.Close()
			err = t.Execute(outputFile, createData(fileDetails, details))
		}
	}
	return err
}

func createData(fileDetails FileDetails, details interface{}) (data map[string]interface{}) {
	data = make(map[string]interface{})
	mapstructure.Decode(details, &data)
	// Default package
	pkg := fileDetails.Package
	if pkg == "" {
		pkg = filepath.Base(filepath.Dir(fileDetails.TargetFile))
	}
	data["Package"] = pkg
	return data
}

func createFunctions() template.FuncMap {
	return template.FuncMap{
		"ToLowerCamel": strcase.ToLowerCamel,
	}
}
