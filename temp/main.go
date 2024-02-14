package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	pb "github.com/meshnet-gophers/protobufs/meshtastic"

	"time"
)

type Packet struct {
	Sender   uint32
	PacketID uint32
	Payload  []byte
}

func main() {
	packetPayload, _ := hex.DecodeString("6c804c8fc044afd9")
	p := Packet{
		Sender:   2062296788,
		PacketID: 463280907,
		Payload:  packetPayload,
	}

	// Sleep for a bit to allow the serial port to come up
	time.Sleep(3 * time.Second)
	println("decrypting")
	var DefaultKey = []byte{0xd4, 0xf1, 0xbb, 0x3a, 0x20, 0x29, 0x07, 0x59, 0xf0, 0xbc, 0xff, 0xab, 0xcf, 0x4e, 0x69, 0x01}
	decrypted, err := XOR(p.Payload, DefaultKey, p.PacketID, p.Sender)
	if err != nil {
		println("error decrypting", err.Error())
		return
	}
	println("unmarshalling")
	var meshPacket pb.Data
	err = meshPacket.UnmarshalVT(decrypted)
	if err != nil {
		println("error unmarshal decrypted to Data:", err.Error())
		return
	}
	println("message content:", string(meshPacket.Payload))

	select {}
}

// CreateNonce creates a 128-bit nonce.
// It takes a uint32 packetId, converts it to a uint64, and a uint32 fromNode.
// The nonce is concatenated as [64-bit packetId][32-bit fromNode][32-bit block counter].
func CreateNonce(packetId uint32, fromNode uint32) ([]byte, error) {
	// Expand packetId to 64 bits
	packetId64 := uint64(packetId)

	// Initialize block counter (32-bit, starts at zero)
	blockCounter := uint32(0)

	// Create a buffer for the nonce
	buf := new(bytes.Buffer)

	// Write packetId, fromNode, and block counter to the buffer
	err := binary.Write(buf, binary.LittleEndian, packetId64)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, fromNode)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, blockCounter)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// XOR encrypts or decrypts text with the specified key. It requires the packetID and sending node ID for the AES IV
func XOR(text []byte, key []byte, packetID, fromNode uint32) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("key length must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. It's common to include it at
	// the beginning of the text. In CTR mode, the IV size is equal to the block size.
	//if len(text) < aes.BlockSize {
	//	return nil, fmt.Errorf("text too short")
	//}
	iv, err := CreateNonce(packetID, fromNode)
	if err != nil {
		return nil, err
	}
	//text = text[aes.BlockSize:]

	// CTR mode is the same for both encryption and decryption, so we use
	// the NewCTR function rather than NewCBCDecrypter.
	stream := cipher.NewCTR(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	plaintext := make([]byte, len(text))
	stream.XORKeyStream(plaintext, text)

	return plaintext, nil
}
