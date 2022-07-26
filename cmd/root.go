package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gloo-redirector",
	Short: "Gloo redirector generates redirection configuration from a .csv file",
	Long:  `gloo-redirector generates redirection configuration for 3xx and 4xx redirections. For either Gloo Edge and Gloo Mesh.`,
}

// ExecuteCmd adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteCmd() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("an error occurred")
		os.Exit(1)
	}
}
