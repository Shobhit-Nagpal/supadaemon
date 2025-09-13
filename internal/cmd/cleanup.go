package cmd

import (
	"fmt"
	"os"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleans up daemon configs for Supadaemon",
	Long:  `Removes the service and config files of daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.CleanupAll()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Service cleanup successful!")
	},
}
