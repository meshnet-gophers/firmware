package internal

import (
	"bytes"
	"testing"
)

func TestMarshalAndParsePacket(t *testing.T) {
	originalPacket := &Packet{
		Destination: 0x11223344,
		Sender:      0x55667788,
		PacketID:    0x99AABBCC,
		Flags: PacketHeaderFlags{
			HopLimit: 3,
			WantAck:  true,
			ViaMQTT:  false,
		},
		ChannelHash: 0xDD,
		Payload:     []byte("Hello, World!"),
	}

	marshaled, err := MarshalPacket(originalPacket)
	if err != nil {
		t.Fatalf("Error marshaling packet: %v", err)
	}

	parsedPacket, err := ParsePacket(marshaled)
	if err != nil {
		t.Fatalf("Error parsing packet: %v", err)
	}

	if originalPacket.Destination != parsedPacket.Destination ||
		originalPacket.Sender != parsedPacket.Sender ||
		originalPacket.PacketID != parsedPacket.PacketID ||
		originalPacket.Flags.HopLimit != parsedPacket.Flags.HopLimit ||
		originalPacket.Flags.WantAck != parsedPacket.Flags.WantAck ||
		originalPacket.Flags.ViaMQTT != parsedPacket.Flags.ViaMQTT ||
		originalPacket.ChannelHash != parsedPacket.ChannelHash ||
		!bytes.Equal(originalPacket.Payload, parsedPacket.Payload) {
		t.Errorf("Parsed packet does not match original")
	}
}

func TestMarshalPacketPayloadTooLarge(t *testing.T) {
	packet := &Packet{
		Destination: 0x11223344,
		Sender:      0x55667788,
		PacketID:    0x99AABBCC,
		Flags: PacketHeaderFlags{
			HopLimit: 3,
			WantAck:  true,
			ViaMQTT:  false,
		},
		ChannelHash: 0xDD,
		Payload:     make([]byte, 238), // One byte too large
	}

	_, err := MarshalPacket(packet)
	if err == nil {
		t.Errorf("Expected an error for payload size exceeding maximum allowed, but got nil")
	}
}

func TestParsePacketTooShort(t *testing.T) {
	packetData := make([]byte, DataOffset-1) // One byte too short
	_, err := ParsePacket(packetData)
	if err == nil {
		t.Errorf("Expected an error for packet data being too short, but got nil")
	}
}
