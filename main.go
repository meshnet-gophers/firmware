package main

import (
	"github.com/meshnet-gophers/firmware/sx126x"
	"log"
	"machine"
	"time"
	"tinygo.org/x/drivers/lora"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

func main() {
	time.Sleep(1 * time.Second)

	println("\n# sx1262 test")
	println("# ----------------------")

	loraRadio, err := configureLoRa()
	if err != nil {
		panic(err)
	}
	for i := 0; ; i++ {
		detected := loraRadio.DetectDevice()
		if detected {
			println("DEVICE DETECTED")
			break
		}
		println(i, ":", "sx126x not detected.")
		time.Sleep(1 * time.Second)
	}

	loraConf := lora.Config{
		Freq:           906875000,
		Bw:             lora.Bandwidth_250_0,
		Sf:             lora.SpreadingFactor11,
		Cr:             lora.CodingRate4_5,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       16,
		Ldr:            lora.LowDataRateOptimizeOff,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       0x24b4,
		LoraTxPowerDBm: 1,
	}

	loraRadio.LoraConfig(loraConf)
	//loraRadio.SetDioIrqParams(sx126x.SX126X_IRQ_ALL, sx126x.SX126X_IRQ_ALL, 0x00, 0x00)

	//go func() {
	//	for {
	//		in := <-loraRadio.GetRadioEventChan()
	//		println("interrupt", in.EventType, in.EventData)
	//	}
	//}()

	println("main: Receiving Lora ")
	for {
		buf, err := loraRadio.Rx2() //0xffffff
		if err != nil {
			println("RX Error: ", err)
		} else if buf != nil {
			log.Println("Packet Received: len=", len(buf), string(buf))
		}
	}
}

func configureLoRa() (*sx126x.Device, error) {
	err := machine.SPI1.Configure(machine.SPIConfig{
		//Mode:      0,
		Frequency: 16 * 1e6,
		//SDO: machine.GPIO15,
		//SDI: machine.GPIO24,
		//SCK: machine.GPIO14,
	})
	if err != nil {
		return nil, err
	}
	machine.LORA_BUSY.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.LORA_RESET.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LORA_ANT_SW.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// trigger a radio reset
	//machine.LORA_RESET.Low()
	//time.Sleep(100 * time.Nanosecond)
	machine.LORA_RESET.High()

	loraRadio := sx126x.New(machine.SPI1)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)
	controller := sx126x.NewRadioControl(
		machine.LORA_CS, machine.LORA_BUSY, machine.LORA_DIO1,
		machine.LORA_ANT_SW, machine.NoPin, machine.NoPin)
	//controller.SetRfSwitchMode(sx126x.RFSWITCH_RX)
	err = loraRadio.SetRadioController(controller)
	if err != nil {
		return nil, err
	}
	//	loraRadio.ExecSetCommand(sx126x.SX126X_CMD_SET_DIO2_AS_RF_SWITCH_CTRL, []uint8{sx126x.SX126X_DIO2_AS_RF_SWITCH})

	return loraRadio, nil
}
