//go:build waveshare_rp2040_lora

package main

import (
	"github.com/meshnet-gophers/firmware/hardware/mcu/waveshare/rp2040-lora"
	"github.com/meshnet-gophers/firmware/node"
	"log/slog"
	"machine"
	"time"
)

func blink() {
	rp2040_lora.LED.Set(!rp2040_lora.LED.Get())
	time.AfterFunc(time.Second, blink)
}

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)

	time.Sleep(3 * time.Second)
	println(node.NodeID)
	println(node.LongName)
	println(node.ShortName)

	for {
		err := main2()
		if err != nil {
			slog.Error("firmware fatal err", err)
		}
	}

}
func main2() error {
	slog.Info("firmware starting")
	rp2040_lora.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	blink() // bet this gets pretty crazy if main2 gets called multiple times
	err := node.Start()
	if err != nil {
		return err
	}
	slog.Info("firmware started")
	select {}
}
