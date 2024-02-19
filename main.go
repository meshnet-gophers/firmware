package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/meshnet-gophers/firmware/hardware/mcu/waveshare/rp2040-lora"
	"github.com/meshnet-gophers/firmware/internal"
	"machine"
	"math"
	"time"

	dedup "github.com/crypto-smoke/meshtastic-go/dedupe"
	pb "github.com/meshnet-gophers/firmware/meshtastic"
	"github.com/meshnet-gophers/firmware/sx126x"
	"tinygo.org/x/drivers/lora"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

func blink() {
	rp2040_lora.LED.Set(!rp2040_lora.LED.Get())
	time.AfterFunc(time.Second, blink)
}

// For testing, injecting a value at link time will cause this message to be sent.
// tinygo build -target=waveshare-rp2040-lora -ldflags '-X main.SendMsg="hello from linker"' -o firmware.uf2
var SendMsg string

type DispatchFunc func(string, *internal.Packet, *pb.Data) error

type NamedKey struct {
	name string
	key  []byte
	hash byte
}

func NewNamedKey(name string, key []byte) NamedKey {
	hash, err := internal.ChannelHash(name, key)
	if err != nil {
		panic(err.Error())
	}
	return NamedKey{
		name: name,
		key:  key,
		hash: uint8(hash),
	}
}

type MeshNode struct {
	radio       *sx126x.Device
	dedup       *dedup.PacketDeduplicator
	handlers    map[pb.PortNum]DispatchFunc
	repeatAfter func(*internal.Packet) time.Duration
	keys        []NamedKey
}

func (m *MeshNode) sendTest(msg string) (*internal.Packet, error) {
	pktBytes, err := hex.DecodeString("ffffffffd426ec7a66f7be31035e528374b5a62151")
	if err != nil {
		return nil, err
	}
	txtPkt, err := internal.ParsePacket(pktBytes)
	if err != nil {
		return nil, err
	}
	_, data, err := m.decrypt(txtPkt)
	if err != nil {
		return nil, err
	}

	txtPkt.PacketID, err = machine.GetRNG()
	if err != nil {
		return nil, err
	}
	txtPkt.Sender = math.MaxUint32 - 1

	data.Payload = []byte(msg)

	txtPkt, err = m.encrypt(internal.DefaultKey, txtPkt, data)
	if err != nil {
		return nil, err
	}

	out, err := internal.MarshalPacket(txtPkt)
	if err != nil {
		return nil, err
	}

	if err := m.radio.Tx(out, 1000*10); err != nil {
		return nil, err
	}
	return txtPkt, nil
}

func (m *MeshNode) recvLoop() {
	for {
		buf, err := m.radio.Rx(0xffffff)
		if err != nil {
			println("RX Error: ", err)
			continue
		}

		//log.Println("Packet Received: len=", len(buf), string(buf))
		packet, err := internal.ParsePacket(buf)
		if err != nil {
			println("error parsing packet:", err.Error())
			continue
		}
		// ignore duplicates of the packet
		if m.dedup.Seen(packet.Sender, packet.PacketID) {
			continue
		}

		// Schedule retransmit
		if packet.Flags.HopLimit > 0 {
			if delay := m.repeatAfter(packet); delay >= 0 {
				time.AfterFunc(delay, func() {
					if err := m.repeat(packet); err != nil {
						println("error retransmitting: ", err.Error())
					} else {
						println("retransmatted")
					}
				})
			}
		}

		println("Packet received:", hex.EncodeToString(buf))
		println("sender, destination, packet ID, hop limit, channel, want ack, via mqtt")
		println(packet.Sender, packet.Destination, packet.PacketID, packet.Flags.HopLimit, packet.ChannelHash, packet.Flags.WantAck, packet.Flags.ViaMQTT)
		println("payload:", hex.EncodeToString(packet.Payload))
		println()
		keyName, data, err := m.decrypt(packet)
		if err != nil {
			println("error decrypting:", err.Error())
			continue
		}
		if f, ok := m.handlers[data.Portnum]; ok {
			if err := f(keyName, packet, data); err != nil {
				println("error processing incoming packet: ", err.Error())
			}
		} else {
			println("no handler for ", data.Portnum)
		}
	}
}

func (m *MeshNode) repeat(packet *internal.Packet) error {
	println("repeat goes here")
	packet.Flags.HopLimit = packet.Flags.HopLimit - 1
	_, err := internal.MarshalPacket(packet)
	if err != nil {
		return err
	}
	return nil

	// return m.radio.Tx(out, 1000*10)
}

func NewNamedKeyB64(name, k64 string) NamedKey {
	dec, err := base64.StdEncoding.DecodeString(k64)
	if err != nil {
		panic("error reading private key: " + err.Error())
	}
	return NewNamedKey(name, dec)
}

func main() {
	rp2040_lora.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	blink()

	// sleep for 3 seconds to let serial monitor connect
	for i := 0; i < 3; i++ {
		println("sleep cycle", i)
		time.Sleep(1 * time.Second)
	}

	println("\n# sx1262 test")
	println("# ----------------------")

	loraRadio, err := rp2040_lora.ConfigureLoRa()
	if err != nil {
		println("failed configuring radio:", err.Error())
	}

	if !loraRadio.DetectDevice() {
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

	println("main: Receiving Lora ")
	node := &MeshNode{
		radio: loraRadio,
		dedup: dedupe,
		keys: []NamedKey{
			NewNamedKey("LongFast", internal.DefaultKey),
			NewNamedKey("Quiet", internal.DefaultKey),
			NewNamedKeyB64("Private", "5haeCewpmA5LqXB/xxrwGp3bwtMnhuxhwzpi+MvVaVo="),
		},
		repeatAfter: func(*internal.Packet) time.Duration { return time.Second },
		handlers: map[pb.PortNum]DispatchFunc{
			pb.PortNum_TEXT_MESSAGE_APP: func(kname string, packet *internal.Packet, data *pb.Data) error {
				println("MESSAGE on:", kname, string(data.Payload))
				return nil
			},
			pb.PortNum_NODEINFO_APP: func(kname string, packet *internal.Packet, data *pb.Data) error {
				u := new(pb.User)
				if err := u.UnmarshalVT(data.Payload); err != nil {
					println("failed unmarshalling user:", err.Error())
					return err
				}
				println("nodeinfo on:", kname, u.LongName)
				return nil
			},
		}}

	for _, nk := range node.keys {
		println("key", nk.name, "has a hash of", nk.hash)
	}

	if SendMsg != "" {
		pkt, err := node.sendTest(SendMsg)
		if err != nil {
			println("error in sendTest(): ", err.Error())
		} else {
			println("transmitted...")
			dedupe.Seen(pkt.Sender, pkt.PacketID)
		}
	}

	node.recvLoop()
}

func (m *MeshNode) decrypt(packet *internal.Packet) (string, *pb.Data, error) {
	for _, namedKey := range m.keys {
		if packet.ChannelHash != namedKey.hash {
			println("wrong channel hash for", namedKey.name, ", got", packet.ChannelHash, ", want", namedKey.hash)
			continue
		}
		decrypted, err := internal.XOR(packet.Payload, namedKey.key, packet.PacketID, packet.Sender)
		if err != nil {
			continue
		}
		meshPacket := new(pb.Data)
		err = meshPacket.UnmarshalVT(decrypted)
		if err != nil {
			return "", nil, err
		}
		return namedKey.name, meshPacket, nil
	}
	return "", nil, errors.New("unable to decrypt")
}

func (m *MeshNode) encrypt(key []byte, packet *internal.Packet, data *pb.Data) (*internal.Packet, error) {
	d, err := data.MarshalVT()
	if err != nil {
		return nil, err
	}
	encrypted, err := internal.XOR(d, key, packet.PacketID, packet.Sender)
	if err != nil {
		//log.Error("failed decrypting packet", "error", err)
		return nil, err
	}
	packet.Payload = encrypted
	return packet, nil
}
