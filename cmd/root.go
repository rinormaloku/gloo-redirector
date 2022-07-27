package cmd

import (
	"github.com/rinormaloku/gloo-redirector/cmd/edge"
	"github.com/rinormaloku/gloo-redirector/cmd/mesh"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gloo-redirector",
	Short: "Gloo redirector generates redirection configuration from a .csv file",
	Long: `Gloo Redirector generates 3xx redirection configuration for either Gloo Edge and Gloo Mesh.
Examples:
  # Generate Gloo Mesh redirection configuration using a file as a source with the default template
  gloo-redirector mesh generate --source /tmp/redirections.csv

  # Generate Gloo Edge redirection configuration using a file as a source with the default template
  gloo-redirector edge generate --source /tmp/redirections.csv
`,
}

// ExecuteCmd adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteCmd() {
	rootCmd.AddCommand(mesh.GetMeshCommand())
	rootCmd.AddCommand(edge.GetEdgeCommand())
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("an error occurred")
		os.Exit(1)
	}
}
