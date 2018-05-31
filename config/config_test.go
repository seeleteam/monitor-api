/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath(t *testing.T) {
	workPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	appConfigPath := filepath.Join(workPath, "conf", "app.conf")
	fmt.Printf("appConfigPath is %s\n", appConfigPath)
}
