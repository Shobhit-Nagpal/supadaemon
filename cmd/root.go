package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "supadaemon",
	Short: "Supabase project monitoring daemon",
	Long: `supadaemon is a background service that continuously monitors your Supabase projects.
    
    It runs as a daemon process, periodically checking the health and status of your 
    Supabase instances at configurable intervals (default: every 5 hours).
    
    Perfect for:
    • Monitoring project uptime and availability  
    • Tracking database connection health
    • Getting alerts on service disruptions
    • Automated health checks for production workloads
    
    Run 'supadaemon start' to begin monitoring, or 'supadaemon --help' for all options.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Daemon logic here
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
