//go:build amd64

package mock

import (
	"fmt"
	"github.com/meshnet-gophers/firmware/sx126x"
	"machine"
)

type Pin uint8

func (p Pin) Configure(config machine.PinConfig) {}

func (p Pin) High() {
	fmt.Println("pin", p, "high")
}

func (p Pin) Low() {
	fmt.Println("pin", p, "low")
}

var LED = NoPin
var NoPin Pin = 0

type MockSPI struct{}

func (m MockSPI) Tx(w, r []byte) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (m MockSPI) Transfer(b byte) (byte, error) {
	//TODO implement me
	//panic("implement me")
	return 0, nil
}

func ConfigureLoRa() (*sx126x.Device, error) {

	loraRadio := sx126x.New(new(MockSPI))
	controller := new(MockRadioController)
	//controller.SetRfSwitchMode(sx126x.RFSWITCH_RX)
	err := loraRadio.SetRadioController(controller)
	if err != nil {
		return nil, err
	}
	//	loraRadio.ExecSetCommand(sx126x.SX126X_CMD_SET_DIO2_AS_RF_SWITCH_CTRL, []uint8{sx126x.SX126X_DIO2_AS_RF_SWITCH})

	return loraRadio, nil
}

type MockRadioController struct{}

func (m MockRadioController) Init() error {
	return nil
}

func (m MockRadioController) SetRfSwitchMode(mode int) error {
	return nil
}

func (m MockRadioController) SetNss(state bool) error {
	return nil
}

func (m MockRadioController) WaitWhileBusy() error {
	return nil
}

func (m MockRadioController) SetupInterrupts(handler func()) error {
	return nil
}
