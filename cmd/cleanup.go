package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleans up daemon configs for Supadaemon",
	Long:  `Removes the service and config files of daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cleanup daemon service")
	},
}
