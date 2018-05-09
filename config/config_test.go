/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSeeleConfigInit(t *testing.T) {
	fmt.Printf("init config %v\n", AppConfig.String("run_mode"))
	LoadAppConfig("ini", "../conf/app.conf")
	fmt.Printf("init config load file\n%v\n", AppConfig.innerConfig)
	fmt.Printf("config run_mode %v\n", AppConfig.String("run_mode"))
	fmt.Printf("init config\n%v\n", AppConfig.innerConfig)

	runMode := AppConfig.String("run_mode")
	configs, err := AppConfig.innerConfig.GetSection(runMode)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("runMode %v configs is %v\n", runMode, configs)

}

func TestConfigPath(t *testing.T) {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath := filepath.Join(workPath, "conf", "app.conf")
	fmt.Printf("appConfigPath is %s", appConfigPath)
}
