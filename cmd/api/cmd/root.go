/**
*  @file
*  @copyright defined in dashboard-api/LICENSE
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "monitor-api",
	Short: "api command for starting the monitor-api",
	Long:  `use "api help [<command>]" for detailed usage`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
