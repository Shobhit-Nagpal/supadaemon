package cmd

import (
	"fmt"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable the daemon process",
	Long:  `Disables the daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.DisableService()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
