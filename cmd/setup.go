package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up daemon configs for Supadaemon",
	Long:  `Writes a service file to boot up daemon on startup`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setup daemon service")
	},
}
