module github.com/sago35/tinygo-keyboard

go 1.26.0

require (
	github.com/itchio/lzma v0.0.0-20190703113020-d3e24e3e3d49
	github.com/tinygo-org/pio v0.3.0
	golang.org/x/exp v0.0.0-20261014143641-f0251eb4aa3e
	tinygo.org/x/bluetooth v0.15.0
	tinygo.org/x/drivers v0.37.0
	tinygo.org/x/tinydraw v0.4.0
	tinygo.org/x/tinyfont v0.7.0
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/saltosystems/winrt-go v0.0.0-20260317170058-9c2fec580d96 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/soypat/cyw43439 v0.1.0 // indirect
	github.com/soypat/lneto v0.4.0 // indirect
	github.com/soypat/seqs v0.0.0-20250124201400-0d65bc7c1710 // indirect
	github.com/tinygo-org/cbgo v0.0.4 // indirect
	golang.org/x/sys v0.11.0 // indirect
)

replace tinygo.org/x/bluetooth => github.com/sago35/bluetooth v0.15.1-0.20260822083226-597fa3c4b688
