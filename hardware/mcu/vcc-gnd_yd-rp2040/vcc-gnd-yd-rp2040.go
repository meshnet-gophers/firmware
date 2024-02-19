//go:build vcc_gnd_yd_rp2040

package vcc_gnd_yd_rp2040

import (
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/ws2812"
)

var (
	RED   = color.RGBA{R: 0xff, G: 0x00, B: 0x00}
	GREEN = color.RGBA{R: 0x00, G: 0xff, B: 0x00}
	BLUE  = color.RGBA{R: 0x00, G: 0x00, B: 0xff}
)

func init() {
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	go func() {
		ws := ws2812.New(machine.LED)
		for {
			err := ws.WriteColors([]color.RGBA{RED})
			if err != nil {
				println("error setting color:", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
