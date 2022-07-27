package edge

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
var edgeCmd = &cobra.Command{
	Use:   "edge",
	Short: "Gloo Edge commands",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Gloo Edge configuration to redirect traffic",
	Long: `Generate Gloo Edge configuration to redirect traffic according to the input csv file

Examples:
  # Generate Gloo Edge redirection configuration using a file as a source with the default template
  gloo-redirector edge generate --source /tmp/redirections.csv

  # Generate Gloo Edge redirection configuration using a file as a source with a custom template
  gloo-redirector edge generate --source /tmp/redirections.csv --template /tmp/template.yaml
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generator.Generate(sourceFilePath, templateFilePath, rd.GenerateEdgeRedirections)
	},
}

var printTemplateCmd = &cobra.Command{
	Use:   "print-template",
	Short: "Prints the default template to generate Gloo Mesh configuration",
	Long: `Prints the default template to generate Gloo Mesh configuration

Examples:
  # Print the default template
  gloo-redirector edge print-template
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(templates.GlooEdgeResourceTemplate)
		return nil
	},
}

func GetEdgeCommand() *cobra.Command {
	generateCmd.Flags().StringVarP(&sourceFilePath, "source", "s", "", "The .csv file with the redirections")
	generateCmd.Flags().StringVarP(&templateFilePath, "template", "t", "",
		"The template file that contains customized VirtualGateway and RouteTable to match your environment")
	edgeCmd.AddCommand(generateCmd)
	edgeCmd.AddCommand(printTemplateCmd)
	return edgeCmd
}
