//go:build waveshare_rp2040_lora

package rp2040_lora

import (
	"github.com/meshnet-gophers/firmware/sx126x"
	"machine"
)

const LED = machine.LED

func init() {
	//machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func ConfigureLoRa() (*sx126x.Device, error) {

	err := machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
	})
	if err != nil {
		return nil, err
	}
	machine.LORA_BUSY.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.LORA_RESET.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LORA_ANT_SW.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.LORA_RESET.High() // not in reset

	loraRadio := sx126x.New(machine.SPI1)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)
	controller := sx126x.NewRadioControl(
		machine.LORA_CS, machine.LORA_BUSY, machine.LORA_DIO1,
		machine.LORA_ANT_SW, machine.NoPin, machine.NoPin)
	controller.SetRfSwitchMode(sx126x.RFSWITCH_RX)
	err = loraRadio.SetRadioController(controller)
	if err != nil {
		return nil, err
	}
	loraRadio.ExecSetCommand(sx126x.SX126X_CMD_SET_DIO2_AS_RF_SWITCH_CTRL, []uint8{sx126x.SX126X_DIO2_AS_RF_SWITCH})

	return loraRadio, nil
}
