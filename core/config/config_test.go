/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig("ini", "../../conf/app.conf")
	if err != nil {
		assert.Errorf(t, err, "error occur %v", err)
	}
	fmt.Printf("config is %v\n", config)
	appName := config.String("app_name")
	fmt.Printf("appName is: %v\n", appName)

}

func TestRegister(t *testing.T) {
	adapterAll := adapters
	fmt.Printf("adapterAll is: %+v\n", adapterAll)

}

func TestExpandValueEnv(t *testing.T) {
	yourShell := ExpandValueEnv("${SHELL}")
	fmt.Printf("gopath is: %v", yourShell)
	assert.NotEmptyf(t, yourShell, "your shell is: %v\n", yourShell)

}
