package jarviscrypto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
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
		jarvisbase.Fatal("jarviscrypto.IntToHex", zap.Error(err))
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

// CheckWIF checks that string wif is a valid Wallet Import Format or Wallet Import Format Compressed string. If it is not, err is populated with the reason.
func CheckWIF(wif string) (valid bool, err error) {
	/* See https://en.bitcoin.it/wiki/Wallet_import_format */

	/* Base58 Check Decode the WIF string */
	ver, privbytes, err := Base58CheckDecode(wif)
	if err != nil {
		return false, err
	}

	/* Check that the version byte is 0x80 */
	if ver != 0x80 {
		return false, fmt.Errorf("Invalid WIF version 0x%02x, expected 0x80", ver)
	}

	/* Check that private key bytes length is 32 or 33 */
	if len(privbytes) != 32 && len(privbytes) != 33 {
		return false, fmt.Errorf("Invalid private key bytes length %d, expected 32 or 33", len(privbytes))
	}

	/* If the private key bytes length is 33, check that suffix byte is 0x01 (for compression) */
	if len(privbytes) == 33 && privbytes[len(privbytes)-1] != 0x01 {
		return false, fmt.Errorf("Invalid private key bytes, unknown suffix byte 0x%02x", privbytes[len(privbytes)-1])
	}

	return true, nil
}
