//go:build xiao_ble

package main

import (
	"device/nrf"
	"fmt"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"tinygo.org/x/bluetooth"
)

var (
	led1 = machine.LED_RED
	led2 = machine.LED_GREEN
	led3 = machine.LED_BLUE
)

// kb is kept as the concrete *SwitchableKeyboard: methods beyond UpDowner
// (SetBatteryLevel etc.) are not reachable through d.Keyboard.
var kb *keyboard.SwitchableKeyboard

// batADC reads the battery voltage divider on P0.31 (enabled via P0.14).
var batADC = machine.ADC{Pin: machine.P0_31}

func init() {
	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	led1.High()
	led2.High()
	led3.High()

	nrf.USBD.USBPULLUP.Set(0) // ホストにまだ見せない

	p014Disable() // 分圧回路オフ (P0.31 過電圧防止)
}

func callback(layer int) {
	//led1.Set(layer != 0)
	//led2.Set(layer != 0)
	led3.Set(layer == 0)
}

// setupKeyboard routes the keyboard and the joystick mouse through a
// USB/BLE switchable output.
func setupKeyboard(d *keyboard.Device) {
	kb = keyboard.NewBLEKeyboard()
	kb.MuteBLEOnUSB = true
	kb.OverrideCtrlH = true // override ctrl-h to BackSpace (USB/BLE 両方)
	d.Keyboard = kb
	d.Mouse = keyboard.NewBLEMouse(kb)
}

// enableUSB shows the device to the USB host. It runs after d.Init() so the
// host never sees a half-configured device (VID/PID and HID are set by then).
func enableUSB() {
	nrf.USBD.USBPULLUP.Set(1)
}

// allowIdle permits the main loop's idle (slow scan) mode only on battery:
// while USB power is present there is nothing to save, and development is
// nicer with the keyboard always at the fast cadence. Only called at the
// slow cadence, so the SoftDevice call stays off the fast path.
func allowIdle() bool {
	return !bluetooth.USBVBusPresent()
}

// tickBoard runs board-specific periodic work from the main loop (cnt
// advances every 500us): once a second, measure the battery voltage and
// report it through the BLE battery service. SetBatteryLevel only notifies
// the central when the percentage actually changes.
//
// The measurement runs on the main loop on purpose: machine.ADC.Get() blocks
// (~350us here) and the joystick already does blocking reads from the same
// loop, so all SAADC access stays serialized. Reading from a goroutine would
// race the joystick on the SAADC registers.
func tickBoard(cnt int) {
	switch cnt % 2000 { // 2000 ticks = 1s
	case 0:
		p014Enable() // divider on; settle until the read below
	case 20: // 10ms later
		// On nrf52, ADC.Configure writes the *global* SAADC registers
		// (CH[0].CONFIG etc.), so the battery settings have to be applied
		// before the read and reverted right after, or the joystick reads
		// would inherit them. Reference is the full-scale voltage in mV
		// (calculateVoltage must match). SampleTime 40 (the longest, for
		// source impedances up to 800kohm) suits the ~340kohm effective
		// impedance of the 1M+510k divider; Samples 8 reduces noise.
		batADC.Configure(machine.ADCConfig{
			Reference:  3000,
			SampleTime: 40,
			Samples:    8,
		})
		raw := batADC.Get()
		batADC.Configure(machine.ADCConfig{}) // back to the joystick settings
		p014Disable()

		vbat := calculateVoltage(raw)
		pct := keyboard.BatteryPercent(vbat)
		kb.SetBatteryLevel(pct)
		if bluetooth.USBVBusPresent() {
			fmt.Printf("ADC Raw: %04X, Battery Voltage: %.2f V, Percent: %d\n", raw, vbat, pct)
		}
	}
}

// ---- P0.14 (battery voltage divider enable) ----
//
// The XIAO BLE (Sense) enables the battery voltage divider (into P0.31) via
// P0.14. Driving P0.14 as Output+High (3.3V) can push P0.31 beyond the
// nRF52840's absolute maximum rating (3.6V) depending on the divider state,
// so P0.14 may only be driven as Output when Low (sink, enable); to disable,
// switch it to Input (high-Z) instead of driving it High.

func p014Enable() {
	machine.P0_14.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.P0_14.Low()
}

func p014Disable() {
	machine.P0_14.Configure(machine.PinConfig{Mode: machine.PinInput})
}

// calculateVoltage converts a raw ADC reading (16-bit left-aligned) of the
// battery voltage divider into the battery voltage in volts.
func calculateVoltage(raw uint16) float64 {
	// Full scale (65535) corresponds to the 3.0V ADC reference
	// (ADCConfig.Reference above must match).
	vIn := (float64(raw) / 65535.0) * 3.0

	// Undo the divider: V_in = V_bat * (510k / (1M + 510k)),
	// so V_bat = V_in * (1510 / 510).
	return vIn * (1510.0 / 510.0)
}
