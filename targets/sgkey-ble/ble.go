//go:build tinygo && nrf52840

package main

import (
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var reportIn bluetooth.Characteristic

var reportMap = []byte{
	0x05, 0x01, // USAGE_PAGE (Generic Desktop)
	0x09, 0x06, // USAGE (Keyboard)
	0xa1, 0x01, // COLLECTION (Application)
	0x85, 0x01, //   REPORT_ID (1)
	0x05, 0x07, //   USAGE_PAGE (Keyboard)
	0x19, 0x01, //   USAGE_MINIMUM
	0x29, 0x7f, //   USAGE_MAXIMUM
	0x15, 0x00, //   LOGICAL_MINIMUM (0)
	0x25, 0x01, //   LOGICAL_MAXIMUM (1)
	0x75, 0x01, //   REPORT_SIZE (1)
	0x95, 0x08, //   REPORT_COUNT (8)
	0x81, 0x02, //   INPUT (Data,Var,Abs)
	0x95, 0x01, //   REPORT_COUNT (1)
	0x75, 0x08, //   REPORT_SIZE (8)
	0x81, 0x01, //   INPUT (Cnst,Ary,Abs)
	0x95, 0x06, //   REPORT_COUNT (6)
	0x75, 0x08, //   REPORT_SIZE (8)
	0x15, 0x00, //   LOGICAL_MINIMUM (0)
	0x25, 0x65, //   LOGICAL_MAXIMUM (101)
	0x05, 0x07, //   USAGE_PAGE (Keyboard)
	0x19, 0x00, //   USAGE_MINIMUM (Reserved (no event indicated))
	0x29, 0x65, //   USAGE_MAXIMUM (Keyboard Application)
	0x81, 0x00, //   INPUT (Data,Ary,Abs)
	0xc0, // END_COLLECTION
}
