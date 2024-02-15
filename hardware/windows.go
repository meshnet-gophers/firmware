//go:build windows && amd64

package hardware

import (
	"github.com/meshnet-gophers/firmware/sx126x"
	"machine"
)

type Pin uint8

func (p Pin) Configure(config machine.PinConfig) {}

func (p Pin) High() {

}

func (p Pin) Low() {

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

	return loraRadio, nil
}
