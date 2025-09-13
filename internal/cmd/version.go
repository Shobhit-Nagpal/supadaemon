package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Supadaemon",
	Long:  `All software has versions. This is Supadaemon's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Supadaemon v%s\n", VERSION)
	},
}
