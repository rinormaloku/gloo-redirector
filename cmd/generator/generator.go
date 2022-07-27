package generator

import (
	"bytes"
	"fmt"
	"github.com/rinormaloku/gloo-redirector/pkg/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func Generate(sourceFilePath, templateFilePath string, envGenerator func(data domain.InputData) (bytes.Buffer, error)) error {
	sourcePayload, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		log.WithError(err).Error("can't read source file")
		return nil
	}

	var templatePayload []byte
	if templateFilePath != "" {
		templatePayload, err = ioutil.ReadFile(templateFilePath)
		if err != nil {
			log.WithError(err).Error("can't read template file")
			return nil
		}
	}
	content, err := envGenerator(
		domain.InputData{
			CsvFile:      sourcePayload,
			TemplateFile: templatePayload,
		},
	)
	if err != nil {
		log.WithError(err).Error("can't generate redirections")
		return err
	}

	fmt.Println(content.String())
	return nil
}
