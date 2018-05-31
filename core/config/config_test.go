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
	_, err := NewConfig("ini", "./testconfig/app.conf")
	if err != nil {
		assert.Errorf(t, err, "error occur %v", err)
	}
}

func TestNewConfigBad(t *testing.T) {
	_, err := NewConfig("ini", "./testconfig/app-bad.conf")
	if err == nil {
		assert.Errorf(t, err, "error occur %v", err)
	}
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
