package domain

import (
	"text/template"
)

type InputData struct {
	CsvFile         []byte
	TemplateFile    []byte
	DefaultTemplate string
}

func (id *InputData) ParseTemplate() (*template.Template, error) {
	rtTemplate := template.New("resource")
	if len(id.TemplateFile) > 0 {
		return rtTemplate.Parse(string(id.TemplateFile))
	} else {
		return rtTemplate.Parse(id.DefaultTemplate)
	}
}
