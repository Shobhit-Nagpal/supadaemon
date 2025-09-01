package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the daemon process",
	Long:  `Spins up the daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start daemon process")
	},
}
