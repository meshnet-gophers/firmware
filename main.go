package main

import (
	"encoding/hex"
	"machine"
	"math"
	"time"

	dedup "github.com/crypto-smoke/meshtastic-go/dedupe"
	"github.com/meshnet-gophers/firmware/hardware"
	meshtastic "github.com/meshnet-gophers/firmware/meshtastic"
	pb "github.com/meshnet-gophers/firmware/meshtastic"
	"github.com/meshnet-gophers/firmware/sx126x"
	"tinygo.org/x/drivers/lora"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

func blink() {
	hardware.LED.Set(!hardware.LED.Get())
	time.AfterFunc(time.Second, blink)
}

func sendTest(loraRadio *sx126x.Device, msg string) (*Packet, error) {
	pktBytes, err := hex.DecodeString("ffffffffd426ec7a66f7be31035e528374b5a62151")
	if err != nil {
		return nil, err
	}
	txtPkt, err := ParsePacket(pktBytes)
	if err != nil {
		return nil, err
	}
	data, err := decrypt(txtPkt)
	if err != nil {
		return nil, err
	}

	txtPkt.PacketID, err = machine.GetRNG()
	if err != nil {
		return nil, err
	}
	txtPkt.Sender = math.MaxUint32 - 1

	data.Payload = []byte(msg)

	txtPkt, err = encrypt(txtPkt, data)
	if err != nil {
		return nil, err
	}

	out, err := MarshalPacket(txtPkt)
	if err != nil {
		return nil, err
	}

	if err := loraRadio.Tx(out, 1000*10); err != nil {
		return nil, err
	}
	return txtPkt, nil
}

// For testing, injecting a value at link time will cause this message to be sent.
// tinygo build -target=waveshare-rp2040-lora -ldflags '-X main.SendMsg="hello from linker"' -o firmware.uf2
var SendMsg string

func main() {
	hardware.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	blink()

	// sleep for 3 seconds to let serial monitor connect
	for i := 0; i < 3; i++ {
		println("sleep cycle", i)
		time.Sleep(1 * time.Second)
	}

	println("\n# sx1262 test")
	println("# ----------------------")

	loraRadio, err := hardware.ConfigureLoRa()
	if err != nil {
		println("failed configuring radio:", err.Error())
	}

	detected := loraRadio.DetectDevice()
	if !detected {
		println("sx1262 not detected")
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
	dedupe := dedup.NewDeduplicator(10 * time.Minute)

	if SendMsg != "" {
		pkt, err := sendTest(loraRadio, SendMsg)
		if err != nil {
			println("error in sendTest(): ", err.Error())
		} else {
			println("transmitted...")
			dedupe.Seen(pkt.Sender, pkt.PacketID)
		}
	}

	println("main: Receiving Lora ")
	for {
		buf, err := loraRadio.Rx(0xffffff)
		if err != nil {
			println("RX Error: ", err)
			continue
		}

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
		data, err := decrypt(packet)
		if err != nil {
			println("error decrypting:", err.Error())
			continue
		}
		switch data.Portnum {
		case pb.PortNum_TEXT_MESSAGE_APP:
			println("MESSAGE:", string(data.Payload))
			data.Payload = []byte("ok")
			payload, err := data.MarshalVT()
			if err != nil {
				println("error marshalling payload:", err.Error())
				continue
			}
			packet.PacketID++
			packet.Sender = math.MaxUint32 - 1

			packet.Payload = payload

			out, err := MarshalPacket(packet)
			if err != nil {
				println("error marshalling packet:", err.Error())
				continue
			}

			if err := loraRadio.Tx(out, 1000*10); err != nil {
				println("error transmitting packet:", err.Error())
				continue
			}
		case pb.PortNum_NODEINFO_APP:
			u := new(meshtastic.User)
			if err := u.UnmarshalVT(data.Payload); err != nil {
				println("failed unmarshalling user:", err.Error())
				continue
			}
			println("nodeinfo:", u.LongName)

			data.WantResponse = true
		}
	}
}

func decrypt(packet *Packet) (*pb.Data, error) {
	decrypted, err := XOR(packet.Payload, DefaultKey, packet.PacketID, packet.Sender)
	if err != nil {
		//log.Error("failed decrypting packet", "error", err)
		return nil, err
	}
	meshPacket := new(pb.Data)
	err = meshPacket.UnmarshalVT(decrypted)
	if err != nil {
		return nil, err
	}
	return meshPacket, nil

}

func encrypt(packet *Packet, data *pb.Data) (*Packet, error) {
	d, err := data.MarshalVT()
	if err != nil {
		return nil, err
	}
	encrypted, err := XOR(d, DefaultKey, packet.PacketID, packet.Sender)
	if err != nil {
		//log.Error("failed decrypting packet", "error", err)
		return nil, err
	}
	packet.Payload = encrypted
	return packet, nil
}
