package main

import (
	"encoding/hex"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/meshnet-gophers/firmware/hardware"
	pb "github.com/meshnet-gophers/protobufs/meshtastic"
	"machine"
	"time"
	"tinygo.org/x/drivers/lora"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

func main() {
	//	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	hardware.LED.High()
	main2()
}
func main2() {

	go func() {
		hardware.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
		for {
			hardware.LED.High()
			time.Sleep(1 * time.Second)
			hardware.LED.Low()
			time.Sleep(1 * time.Second)
		}
	}()

	// sleep for 3 seconds to let serial monitor connect
	for i := 0; i < 3; i++ {
		println("sleep cycle", i)
		time.Sleep(1 * time.Second)
	}

	println("\n# sx1262 test")
	println("# ----------------------")

	loraRadio, err := hardware.ConfigureLoRa()
	if err != nil {
		println(err)
		return
	}
	_ = loraRadio

	detected := loraRadio.DetectDevice()
	if !detected {
		println("sx1262 not detected")
		return
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

	_ = loraConf
	loraRadio.LoraConfig(loraConf)

	dedupe := meshtastic.NewDeduplicator(nil, 10*time.Minute)
	println("main: Receiving Lora ")

	// importing this package makes the whole program not run
	thing := new(pb.NodeInfo)
	thing.Num = 69420
	println(thing.Num)
	// removing the three lines above (and the pb import) allows program to run

	for {
		//buf, err := loraRadio.Rx(0xffffff)
		var buf []byte
		var err error
		if err != nil {
			println("RX Error: ", err)
		} else if buf != nil {
			//log.Println("Packet Received: len=", len(buf), string(buf))
			packet, err := ParsePacket(buf)
			if err != nil {
				println("error parsing packet:", err.Error())
				continue
			}
			// ignore duplicates of the packet
			if dedupe.Seen(packet.Sender, packet.PacketID) {
				continue
			}
			println("Packet received:", hex.EncodeToString(buf))
			println("sender, destination, packet ID, hop limit, channel, want ack, via mqtt")
			println(packet.Sender, packet.Destination, packet.PacketID, packet.Flags.HopLimit, packet.ChannelHash, packet.Flags.WantAck, packet.Flags.ViaMQTT)
			println("payload:", hex.EncodeToString(packet.Payload))
			println()

			continue /*
				data, err := decrypt(packet)
				if err != nil {
					println("error decrypting:", err.Error())
					continue
				}
				if data.Portnum == pb.PortNum_TEXT_MESSAGE_APP {
					println("MESSAGE:", string(data.Payload))
				}
			*/
		}
	}
}

//	func decrypt(packet Packet) (*pb.Data, error) {
//		return nil, errors.New("uh")
//
//			decrypted, err := XOR(packet.Payload, DefaultKey, packet.PacketID, packet.Sender)
//			if err != nil {
//				//log.Error("failed decrypting packet", "error", err)
//				return nil, err
//			}
//			println("unmarshalling")
//			var meshPacket *pb.Data
//			err = meshPacket.UnmarshalVT(decrypted)
//			if err != nil {
//				return nil, err
//			}
//			return meshPacket, nil
//
// }
