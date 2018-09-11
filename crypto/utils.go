package jarviscrypto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/big"
)

func hex2bytes(hexstring string) (b []byte) {
	b, _ = hex.DecodeString(hexstring)
	return b
}

func hex2int(hexstring string) (v *big.Int) {
	v, _ = new(big.Int).SetString(hexstring, 16)
	return v
}

// IntToHex - converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes - reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}
