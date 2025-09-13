package cmd

import (
	"fmt"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon process",
	Long:  `Stops the daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.StopService()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Supadaemon has stopped!")
	},
}
