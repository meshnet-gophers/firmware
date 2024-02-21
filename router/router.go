package router

import (
	"context"
	"github.com/meshnet-gophers/firmware/hardware/lora"
	"github.com/meshnet-gophers/firmware/internal"
	meshtastic "github.com/meshnet-gophers/firmware/meshtastic"
	dedup "github.com/meshnet-gophers/meshtastic-go/dedupe"
	"math"
	"time"
)

// MeshRouter is the main struct for the packet router.
type MeshRouter struct {
	incoming chan *meshtastic.MeshPacket // Channel for receiving packets to be transmitted.
	outgoing chan *meshtastic.MeshPacket // Channel for emitting received packets for processing.
	ctx      context.Context             // Context for managing cancellation and cleanup.
	cancel   context.CancelFunc          // Function to cancel the context.
	radio    lora.Radio                  // LoRa radio to receive packets from and transmit packets on
	dedup    *dedup.PacketDeduplicator
}

// NewMeshRouter creates a new MeshRouter with initialized channels and context.
func NewMeshRouter(ctx context.Context, bufferSize int, radio lora.Radio) *MeshRouter {
	ctx, cancel := context.WithCancel(ctx)
	return &MeshRouter{
		incoming: make(chan *meshtastic.MeshPacket, bufferSize),
		outgoing: make(chan *meshtastic.MeshPacket, bufferSize),
		ctx:      ctx,
		cancel:   cancel,
		radio:    radio,
		dedup:    dedup.NewDeduplicator(10 * time.Minute),
	}
}

// Start starts the packet routing process with context handling for graceful shutdown.
func (r *MeshRouter) Start() {
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return // Exit the goroutine when context is canceled.

			case packet := <-r.incoming: // packets to transmit
				var payload []byte
				var err error
				switch p := packet.GetPayloadVariant().(type) {
				case *meshtastic.MeshPacket_Decoded:
					payload, err = p.Decoded.MarshalVT()
					if err != nil {
						println("error marshalling mesh packet payload:", err.Error())
						continue
					}
				case *meshtastic.MeshPacket_Encrypted:
					payload = p.Encrypted
				}

				pkt := internal.Packet{
					Destination: packet.To,
					Sender:      packet.From,
					PacketID:    packet.Id,
					Flags:       internal.PacketHeaderFlags{},
					ChannelHash: byte(packet.Channel),
					Payload:     payload,
				}
				pktBytes, err := internal.MarshalPacket(&pkt)
				if err != nil {
					println("error marshalling packet:", err.Error())
					continue
				}
				timeout := 10 * time.Second
				err = r.radio.Tx(pktBytes, uint32(timeout.Milliseconds()))
				if err != nil {
					println("error transmitting packet:", err.Error())
					continue
				}
				println("packet transmitted")

			default: // if nothing to tx, listen for a packet for 1 second
				// TODO: I don't really like this, we should find a way to listen indefinitely but be able to break out of
				// it should we get a packet to TX. We probably need to take over the sx126x.Device and implement a lot
				// of this logic in there.
				buf, err := r.radio.Rx(1000)
				if err != nil {
					println("RX Error: ", err)
					continue
				}
				if buf == nil {
					// no packet received
					continue
				}
				// grab these MeshPacket fields before another packet comes in
				RssiPk, SnrPkt, SignalRssiPkt := r.radio.GetPacketStatus()
				rssi := -(float64(RssiPk) / 2)
				rssi = math.Round(rssi*100) / 100

				snr := float64(SnrPkt) / 4
				snr = math.Round(snr*100) / 100

				signalRSSI := -(float64(SignalRssiPkt) / 2)
				signalRSSI = math.Round(signalRSSI*100) / 100

				packet, err := internal.ParsePacket(buf)
				if err != nil {
					println("error parsing packet:", err.Error())
					println("packet len", len(buf))
					continue
				}
				// ignore duplicates of the packet
				if r.dedup.Seen(packet.Sender, packet.PacketID) {
					continue
				}

				meshPacket := &meshtastic.MeshPacket{
					From:           packet.Sender,
					To:             packet.Destination,
					Channel:        uint32(packet.ChannelHash), // I think all packets are transmitted OTA as an encrypted protobuf
					PayloadVariant: &meshtastic.MeshPacket_Encrypted{Encrypted: packet.Payload},
					Id:             packet.PacketID,
					RxTime:         uint32(time.Now().Unix()),
					RxSnr:          float32(snr),
					HopLimit:       uint32(packet.Flags.HopLimit),
					WantAck:        packet.Flags.WantAck,
					Priority:       0, // not sure what to do with this (TODO)
					RxRssi:         int32(rssi),
					Delayed:        0,
					ViaMqtt:        packet.Flags.ViaMQTT,
				}
				select {
				case r.outgoing <- meshPacket: // Attempt to send the packet.
				case <-r.ctx.Done(): // Do nothing if context is already cancelled.
				}
			}
		}
	}()
}

// SendPacket allows consumers to send a MeshPacket to be transmitted, respecting context cancellation.
func (r *MeshRouter) SendPacket(packet *meshtastic.MeshPacket) {
	select {
	case r.incoming <- packet: // Attempt to send the packet.
	case <-r.ctx.Done(): // Do nothing if context is already cancelled.
	}
}

// ReceivePacket provides a channel on which consumers can receive incoming MeshPackets, respecting context cancellation.
func (r *MeshRouter) ReceivePacket() <-chan *meshtastic.MeshPacket {
	return r.outgoing
}

// Close cleans up resources used by MeshRouter, now respecting the cancellation context.
func (r *MeshRouter) Close() {
	r.cancel() // Cancel the context to clean up and exit any ongoing operations gracefully.
	close(r.incoming)
	close(r.outgoing)
}
