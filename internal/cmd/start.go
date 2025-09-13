package cmd

import (
	"fmt"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the daemon process",
	Long:  `Spins up the daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.StartService()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Supadaemon is running!")
	},
}
