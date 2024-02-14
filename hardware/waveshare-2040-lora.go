//go:build waveshare_rp2040_lora

package hardware

import "machine"

const LED = machine.LED

func init() {
	//machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
}
