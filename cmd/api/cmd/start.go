/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package cmd

import (
	"log"
	"sync"

	"github.com/seeleteam/monitor-api/config"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/seeleteam/monitor-api/server"
)

var (
	configFile *string
	g          errgroup.Group
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the monitor-api server",
	Long: `usage example:
	   monitor-api.exe start -c cmd\app.conf
		start the monitor-api.`,

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		// config init
		config.Init(*configFile)

		// init server, if modify the config should write above this line
		server.Start(&g)

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	configFile = startCmd.Flags().StringP("config", "c", "", "api config file (required)")
	startCmd.MarkFlagRequired("config")
}
