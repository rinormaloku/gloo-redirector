package mesh

import (
	"fmt"
	"github.com/rinormaloku/gloo-redirector/cmd/generator"
	rd "github.com/rinormaloku/gloo-redirector/pkg/redirections"
	"github.com/rinormaloku/gloo-redirector/pkg/templates"
	"github.com/spf13/cobra"
)

var sourceFilePath string
var templateFilePath string

// generateCmd represents the generate command
var meshCmd = &cobra.Command{
	Use:   "mesh",
	Short: "Gloo Mesh commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called in mesh")
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Gloo Mesh configuration to redirect traffic",
	Long: `Generate Gloo Mesh configuration to redirect traffic according to the input csv file

Examples:
  # Generate Gloo Mesh redirection configuration using a file as a source with the default template
  gloo-redirector mesh generate --source /tmp/redirections.csv

  # Generate Gloo Mesh redirection configuration using a file as a source with a custom template
  gloo-redirector mesh generate --source /tmp/redirections.csv --template /tmp/template.yaml
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generator.Generate(sourceFilePath, templateFilePath, rd.GenerateMeshRedirections)
	},
}

var printTemplateCmd = &cobra.Command{
	Use:   "print-template",
	Short: "Prints the default template to generate Gloo Mesh configuration",
	Long: `Prints the default template to generate Gloo Mesh configuration

Examples:
  # Print the default template
  gloo-redirector mesh print-template
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(templates.GlooMeshResourceTemplate)
		return nil
	},
}

func GetMeshCommand() *cobra.Command {
	generateCmd.Flags().StringVarP(&sourceFilePath, "source", "s", "", "The .csv file with the redirections")
	generateCmd.Flags().StringVarP(&templateFilePath, "template", "t", "",
		"The template file that contains customized VirtualGateway and RouteTable to match your environment")
	meshCmd.AddCommand(generateCmd)
	meshCmd.AddCommand(printTemplateCmd)
	return meshCmd
}
