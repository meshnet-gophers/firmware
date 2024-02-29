package node

import (
	"context"
	"errors"
	"github.com/meshnet-gophers/firmware/hardware/mcu/waveshare/rp2040-lora"
	"github.com/meshnet-gophers/firmware/internal"
	_ "github.com/meshnet-gophers/firmware/node/plugins"
	"github.com/meshnet-gophers/firmware/plugin/registry"
	"github.com/meshnet-gophers/firmware/router"
	"github.com/meshnet-gophers/meshtastic-go"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"log/slog"
	"strconv"
	"tinygo.org/x/drivers/lora"
)

var (
	NodeID    string
	LongName  string
	ShortName string
	UserID    string
	ID        uint32
)

func init() {
	id, err := strconv.ParseUint(NodeID, 10, 32)
	if err != nil {
		panic(err)
	}
	ID = uint32(id)
	UserID = meshtastic.NodeID(ID).String()
}

type NodeInfo struct {
	pb.User
	ID uint32
}

func Start() error {
	// TODO: this is very much device specific. Need to make it easier to work across different devices
	slog.Debug("configure lora")
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
		LoraTxPowerDBm: 22,
	}
	loraRadio.LoraConfig(loraConf)
	slog.Debug("configure router")

	r := router.NewMeshRouter(context.TODO(), 2, loraRadio)
	slog.Debug("start router")

	r.Start()
	u := &pb.User{
		Id:         UserID,
		LongName:   LongName,
		ShortName:  ShortName,
		Macaddr:    nil,
		HwModel:    pb.HardwareModel_RP2040_LORA,
		IsLicensed: false,
		Role:       pb.Config_DeviceConfig_CLIENT,
	}
	slog.Debug("start plugins")

	registry.StartPlugins(r, u, ID)
	return nil
}
