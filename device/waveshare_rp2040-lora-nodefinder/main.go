//go:build waveshare_rp2040_lora

package main

import (
	"context"
	"errors"
	"github.com/meshnet-gophers/firmware/hardware/mcu/waveshare/rp2040-lora"
	"github.com/meshnet-gophers/firmware/internal"
	"github.com/meshnet-gophers/firmware/router"
	"github.com/meshnet-gophers/meshtastic-go"
	"log/slog"
	"machine"
	"math"
	"strconv"
	"time"

	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
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
// tinygo build -ldflags '-X main.SendMsg="hello from linker"' -o firmware.uf2
var (
	NodeID        string
	ParsedNodeID  uint32
	NodeName      string
	NodeShortName string
)

func init() {
	if NodeID != "" {
		out, err := strconv.ParseUint(NodeID, 10, 32)
		if err != nil {
			panic(err)
		}
		ParsedNodeID = uint32(out)
	} else {
		ParsedNodeID = math.MaxUint32
	}

	if NodeName == "" {
		NodeName = "Meshtastic"
	}

	if NodeShortName == "" {
		NodeShortName = "TEST"
	}
}
func main() {
	time.Sleep(3 * time.Second)
	for {
		err := main2()
		if err != nil {
			time.Sleep(1 * time.Second)
			slog.Error("firmware fatal err", err)
		}
	}

}
func main2() error {
	slog.Info("Firmware starting")
	rp2040_lora.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	//blink() // bet this gets pretty crazy if main2 gets called multiple times
	loraRadio, err := rp2040_lora.ConfigureLoRa()
	if err != nil {
		return errors.Join(err, errors.New("failed configuring LoRa radio"))
	}
	if !loraRadio.DetectDevice() {
		return errors.New("LoRa radio not detected")
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
		SyncWord:       internal.ConvertSyncWord(0x2b, 0x44),
		LoraTxPowerDBm: 1,
	}
	loraRadio.LoraConfig(loraConf)
	r := router.NewMeshRouter(context.TODO(), 2, loraRadio)

	go func() {
		for {
			select {
			case packet := <-r.ReceivePacket():
				slog.Info("packet received",
					"from", packet.From, "to", packet.To,
					"id", packet.Id, "hop limit", packet.HopLimit,
					"channel", packet.Channel, "want ack", packet.WantAck, "via mqtt", packet.ViaMqtt)

				//println("forwarding packet")
				//r.SendPacket(packet)
			}
		}
	}()
	r.Start()

	r.SendPacket(buildNodeInfo())
	for i := 0; ; i++ {
		select {
		case <-time.NewTimer(30 * time.Second).C:
			r.SendPacket(buildNodeInfo())
			slog.Info("nodeinfo sent")
		}
	}
}

// creates a node info packet and encrypts it for LongFast
func buildNodeInfo() *pb.MeshPacket {
	id := meshtastic.NodeID(ParsedNodeID)
	nodeInfo := &pb.User{
		Id:         id.String(),
		LongName:   NodeName,
		ShortName:  NodeShortName,
		HwModel:    pb.HardwareModel_RP2040_LORA,
		IsLicensed: false,
		Role:       pb.Config_DeviceConfig_CLIENT,
	}
	user, err := nodeInfo.MarshalVT()
	if err != nil {
		panic(err)
	}
	packetPayload := &pb.Data{
		Portnum:      pb.PortNum_NODEINFO_APP,
		Payload:      user,
		WantResponse: true,
	}
	payload, _ := packetPayload.MarshalVT()
	rng, err := machine.GetRNG()
	if err != nil {
		panic(err)
	}

	payload, err = internal.XOR(payload, internal.DefaultKey, rng, ParsedNodeID)
	if err != nil {
		panic(err)
	}
	pkt := &pb.MeshPacket{
		From:           ParsedNodeID,
		To:             math.MaxUint32,
		Channel:        8,
		Id:             rng,
		HopLimit:       7,
		WantAck:        false,
		PayloadVariant: &pb.MeshPacket_Encrypted{Encrypted: payload},
	}

	return pkt
}
