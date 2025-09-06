package cmd

import (
	"fmt"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable the daemon process",
	Long:  `Enables the daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.EnableService()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
