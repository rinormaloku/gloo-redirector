package cmd

import (
	"fmt"
	"github.com/rinormaloku/gloo-redirector/domain"
	"github.com/rinormaloku/gloo-redirector/pkg/redirections"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var sourceFilePath string
var templateFilePath string

func init() {
	generateCmd.Flags().StringVarP(&sourceFilePath, "source", "s", "", "The .csv file with the redirections")
	generateCmd.Flags().StringVarP(&templateFilePath, "template", "t", "",
		"The template file that contains customized VirtualGateway and RouteTable to match your environment")
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate yaml VirtualService from .csv file",
	Long:  "Take a .csv file in input and generate a VirtualService",
	RunE: func(cmd *cobra.Command, args []string) error {

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
		content, err := redirections.Generate(
			domain.InputData{
				File:     sourcePayload,
				Template: templatePayload,
			},
		)
		if err != nil {
			log.WithError(err).Error("can't generate redirections")
			return err
		}

		fmt.Println(content.String())
		return nil
	},
}
