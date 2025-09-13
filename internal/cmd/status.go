package cmd

import (
	"fmt"

	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status of the daemon process",
	Long:  `Gets the status of daemon process for supadaemon`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.GetServiceStatus()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
