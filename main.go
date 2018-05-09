/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package main

import (
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/seeleteam/monitor-api/server"
)

var (
	g errgroup.Group
)

func main() {
	// init server, if modify the config should write above this line
	server.Start(&g)

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}
