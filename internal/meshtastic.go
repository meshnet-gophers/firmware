package internal

// Converts an sx127x style uint8 sync word to a sx126x uint16 variant with integreated control bits
func ConvertSyncWord(syncWord, controlBits uint8) uint16 {
	firstByte := (syncWord & 0xF0) | ((controlBits & 0xF0) >> 4)
	secondByte := ((syncWord & 0x0F) << 4) | (controlBits & 0x0F)
	return (uint16(firstByte) << 8) | uint16(secondByte)
}
