package beacon

import (
	"github.com/meshnet-gophers/firmware/internal"
	"github.com/meshnet-gophers/firmware/plugin/registry"
	"github.com/meshnet-gophers/firmware/router"
	"github.com/meshnet-gophers/meshtastic-go"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"log/slog"
	"machine"
	"math"
	"time"
)

func init() {
	registry.RegisterPlugin(new(Beacon))
}

type Beacon struct {
	r      *router.MeshRouter
	u      *pb.User
	nodeID uint32
}

func (b *Beacon) Name() string {
	return "beacon"
}

func (b *Beacon) Init(r *router.MeshRouter, u *pb.User, nodeID uint32) {
	b.r = r
	b.u = u
	b.nodeID = nodeID
	go func() {
		for {
			select {
			case packet := <-r.ReceivePacket():
				slog.Info("packet received",
					"from", packet.From, "to", packet.To,
					"id", packet.Id, "hop limit", packet.HopLimit,
					"channel", packet.Channel, "want ack", packet.WantAck, "via mqtt", packet.ViaMqtt)
			}
		}
	}()
	info, err := b.buildNodeInfo()
	if err != nil {
		slog.Error("error building nodeinfo", "err", err)
		return
	}
	r.SendPacket(info)
	slog.Info("nodeinfo sent")
	go func() {
		for i := 0; ; i++ {
			select {
			case <-time.NewTimer(30 * time.Second).C:
				info, err := b.buildNodeInfo()
				if err != nil {
					slog.Error("error building nodeinfo", "err", err)
					continue
				}
				r.SendPacket(info)
				slog.Info("nodeinfo sent")
			}
		}
	}()

}

// creates a node info packet and encrypts it for LongFast
func (b *Beacon) buildNodeInfo() (*pb.MeshPacket, error) {
	id := meshtastic.NodeID(b.nodeID)
	nodeInfo := &pb.User{
		Id:         id.String(),
		LongName:   b.u.GetLongName(),
		ShortName:  b.u.GetShortName(),
		HwModel:    b.u.GetHwModel(),
		IsLicensed: b.u.GetIsLicensed(),
		Role:       b.u.GetRole(),
	}
	slog.Debug("nodeinfo", "id", nodeInfo.GetId(), "long", nodeInfo.GetLongName(), "short", nodeInfo.GetShortName())
	user, err := nodeInfo.MarshalVT()
	if err != nil {
		return nil, err
	}
	packetPayload := &pb.Data{
		Portnum:      pb.PortNum_NODEINFO_APP,
		Payload:      user,
		WantResponse: true,
	}
	payload, _ := packetPayload.MarshalVT()
	rng, err := machine.GetRNG()
	if err != nil {
		return nil, err

	}

	payload, err = internal.XOR(payload, internal.DefaultKey, rng, b.nodeID)
	if err != nil {
		return nil, err

	}
	pkt := &pb.MeshPacket{
		From:           b.nodeID,
		To:             math.MaxUint32,
		Channel:        8,
		Id:             rng,
		HopLimit:       7,
		WantAck:        false,
		PayloadVariant: &pb.MeshPacket_Encrypted{Encrypted: payload},
	}

	return pkt, nil
}
