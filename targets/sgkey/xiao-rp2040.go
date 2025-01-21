//go:build xiao_rp2040

package main

import (
	"machine"
)

func init() {
	i2c = machine.I2C1
}
