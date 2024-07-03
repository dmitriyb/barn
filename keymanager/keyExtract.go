package keymanager

import (
	"encoding/binary"
	"errors"
)

func extractKeyAndPassword(randomData []byte) (actualKey, password []byte, err error) {
	actualKeySize := 1024 // 1 KB
	passwordSize := 128   // 128 bytes

	// Read the offsets
	if len(randomData) < 5 {
		return nil, nil, errors.New("insufficient data for offsets")
	}

	offset1 := int64(binary.BigEndian.Uint32(append([]byte{0}, randomData[:3]...)))
	offset2 := binary.BigEndian.Uint16(randomData[3:5])

	// Calculate positions
	actualKeyStart := 5 + offset1
	passwordStart := actualKeyStart + int64(actualKeySize) + int64(offset2)
	passwordEnd := passwordStart + int64(passwordSize)

	// Ensure offsets are within bounds
	if actualKeyStart+int64(actualKeySize) > int64(len(randomData)) || passwordEnd > int64(len(randomData)) {
		return nil, nil, errors.New("offsets exceed data length")
	}

	// Extract key and password
	actualKey = randomData[actualKeyStart : actualKeyStart+int64(actualKeySize)]
	password = randomData[passwordStart:passwordEnd]

	return actualKey, password, nil
}