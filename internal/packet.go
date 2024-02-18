package internal

import (
	"encoding/binary"
	"errors"
)

const (
	DestinationOffset = 0
	SenderOffset      = 4
	PacketIDOffset    = 8
	FlagsOffset       = 12
	ChannelHashOffset = 13
	PaddingOffset     = 14
	DataOffset        = 16
)

type Packet struct {
	Destination uint32
	Sender      uint32
	PacketID    uint32
	Flags       PacketHeaderFlags
	ChannelHash byte
	Payload     []byte
}

// PacketHeaderFlags represents the flags found in the packet header.
type PacketHeaderFlags struct {
	HopLimit uint8
	WantAck  bool
	ViaMQTT  bool
}

const (
	HeaderSize = 16
)

// MarshalPacket takes a Packet struct and serializes it into a byte slice.
func MarshalPacket(p *Packet) ([]byte, error) {
	// Ensure packet payload does not exceed maximum allowed size
	if len(p.Payload) > 237 {
		return nil, errors.New("payload exceeds maximum allowed size")
	}

	packet := make([]byte, HeaderSize+len(p.Payload))
	binary.LittleEndian.PutUint32(packet[DestinationOffset:], p.Destination)
	binary.LittleEndian.PutUint32(packet[SenderOffset:], p.Sender)
	binary.LittleEndian.PutUint32(packet[PacketIDOffset:], p.PacketID)
	packet[FlagsOffset] = marshalHeaderFlags(p.Flags)
	packet[ChannelHashOffset] = p.ChannelHash
	copy(packet[DataOffset:], p.Payload)

	return packet, nil
}

// marshalHeaderFlags serializes the packet header flags into a single byte.
func marshalHeaderFlags(flags PacketHeaderFlags) byte {
	var flagsByte byte
	flagsByte |= flags.HopLimit & 0x07 // First 3 bits
	if flags.WantAck {
		flagsByte |= 0x08 // Fourth bit
	}
	if flags.ViaMQTT {
		flagsByte |= 0x10 // Fifth bit
	}
	// Bits 6-8 are currently unused.
	return flagsByte
}

// ParsePacket takes a byte slice and parses the packet header and payload.
func ParsePacket(packet []byte) (*Packet, error) {
	if len(packet) < DataOffset {
		return nil, errors.New("packet is too short to contain a valid header")
	}

	headerFlags := parseHeaderFlags(packet[FlagsOffset])

	return &Packet{
		Destination: binary.LittleEndian.Uint32(packet[DestinationOffset:]),
		Sender:      binary.LittleEndian.Uint32(packet[SenderOffset:]),
		PacketID:    binary.LittleEndian.Uint32(packet[PacketIDOffset:]),
		Flags:       headerFlags,
		ChannelHash: packet[ChannelHashOffset],
		Payload:     packet[DataOffset:],
	}, nil
}

// parseHeaderFlags parses the flags from the packet header.
func parseHeaderFlags(flagsByte byte) PacketHeaderFlags {
	return PacketHeaderFlags{
		HopLimit: flagsByte & 0x07,        // First 3 bits
		WantAck:  (flagsByte & 0x08) != 0, // Fourth bit
		ViaMQTT:  (flagsByte & 0x10) != 0, // Fifth bit
		// Bits 6-8 are currently unused.
	}
}
