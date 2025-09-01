package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon process",
	Long:  `Stops the daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stop daemon process")
	},
}
