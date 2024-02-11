package main

import (
	"encoding/hex"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/meshnet-gophers/firmware/sx126x"
	"machine"
	"time"
	"tinygo.org/x/drivers/lora"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// 3 second loop to give serial monitor time to come up
	for i := 0; i < 3; i++ {
		machine.LED.Low()
		time.Sleep(500 * time.Millisecond)
		machine.LED.High()
		time.Sleep(500 * time.Millisecond)
	}

	println("\n# sx1262 test")
	println("# ----------------------")

	loraRadio, err := configureLoRa()
	if err != nil {
		panic(err)
	}

	detected := loraRadio.DetectDevice()
	if !detected {
		panic("sx1262 not detected")
	}

	// Meshtastic LongFast preset
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
		SyncWord:       0x24b4, // 0x2b sync word with some magic bit shifting to be compatible with older LoRa radios
		LoraTxPowerDBm: 0,      // disable radio for now
	}

	loraRadio.LoraConfig(loraConf)

	// Avoid duplicate packets
	dedupe := meshtastic.NewDeduplicator(nil, 10*time.Minute)
	println("main: Receiving Lora ")
	for {
		buf, err := loraRadio.Rx(0xffffff)
		if err != nil {
			println("RX Error: ", err)
		} else if buf != nil {
			//log.Println("Packet Received: len=", len(buf), string(buf))
			packet, headerFlags, err := ParsePacket(buf)
			if err != nil {
				panic(err)
			}
			_ = headerFlags
			// ignore duplicates of the packet
			if dedupe.Seen(packet.Sender, packet.PacketID) {
				continue
			}

			println(packet.Sender, packet.Destination, "0x"+hex.EncodeToString(packet.Payload))
			/*
				var pkt pb.MeshPacket

				pkt.PayloadVariant = &pb.MeshPacket_Decoded{Decoded: &pb.Data{
					Portnum:      0,
					Payload:      packet.Payload,
					WantResponse: false,
					Dest:         0,
					Source:       0,
					RequestId:    0,
					ReplyId:      0,
					Emoji:        0,
				}}
				println("packet constructed")
				out, err := json.MarshalIndent(&pkt, "", "  ")
				if err == nil {
					panic(err)
				}
				println("json done")
				println(string(out))
				println("printing done")

			*/
			// hex bytes of text message packet with message of "test": d4b66213862739e3
			/*
				// MeshPacket is not a representation of the packet on the wire, I don't think.
				var pkt pb.MeshPacket
				err = proto.Unmarshal(packet.Payload, &pkt)
				if err != nil {
					panic(err)
				}
				println("unmarhalled packet payload to Meshpacket")
			*/

		}
	}
}

func configureLoRa() (*sx126x.Device, error) {
	err := machine.SPI1.Configure(machine.SPIConfig{
		//Mode:      0,
		Frequency: 8 * 1e6,
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
	machine.LORA_RESET.Low()
	time.Sleep(100 * time.Nanosecond)
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
