package cmd

import (
	"fmt"
	"os"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up daemon configs for Supadaemon",
	Long:  `Writes a service file to boot up daemon on startup`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.SetupAll()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Service setup successful!")
	},
}
