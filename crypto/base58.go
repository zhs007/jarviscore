package jarviscrypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Encode encodes the given version and data to a base58check encoded string
func Base58CheckEncode(version string, data []byte) (string, error) {
	prefix, err := hex.DecodeString(version)
	if err != nil {
		return "", err
	}

	dataBytes := append(prefix, data...)

	// Performing SHA256 twice
	sha256hash := sha256.New()
	sha256hash.Write(dataBytes)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	checksum := hash[:4]
	dataBytes = append(dataBytes, checksum...)

	// For all the "00" versions or any prepended zeros as base58 removes them
	zeroCount := 0
	for i := 0; i < len(dataBytes); i++ {
		if dataBytes[i] == 0 {
			zeroCount++
		} else {
			break
		}
	}

	// Performing base58 encoding
	encoded := Base58Encode(dataBytes)

	for i := 0; i < zeroCount; i++ {
		encoded = "1" + encoded
	}

	return encoded, nil
}

// Decode decodes the given base58check encoded string and returns the version prepended decoded string
func Base58CheckDecode(encoded string) ([]byte, error) {
	zeroCount := 0
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == 49 {
			zeroCount++
		} else {
			break
		}
	}

	dataBytes, err := Base58Decode(encoded)
	if err != nil {
		return nil, err
	}
	data, checksum := dataBytes[:len(dataBytes)-4], dataBytes[len(dataBytes)-4:]

	for i := 0; i < zeroCount; i++ {
		data = append([]byte{0}, data...)
	}

	// Performing SHA256 twice to validate checksum
	sha256hash := sha256.New()
	sha256hash.Write(data)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	if !reflect.DeepEqual(checksum, hash[:4]) {
		return nil, errors.New("Data and checksum don't match")
	}

	return data, nil
}

func Base58Encode(data []byte) string {
	var encoded string
	decimalData := new(big.Int)
	decimalData.SetBytes(data)
	divisor, zero := big.NewInt(58), big.NewInt(0)

	for decimalData.Cmp(zero) > 0 {
		mod := new(big.Int)
		decimalData.DivMod(decimalData, divisor, mod)
		encoded = string(alphabet[mod.Int64()]) + encoded
	}

	return encoded
}

func Base58Decode(data string) ([]byte, error) {
	decimalData := new(big.Int)
	alphabetBytes := []byte(alphabet)
	multiplier := big.NewInt(58)

	for _, value := range data {
		pos := bytes.IndexByte(alphabetBytes, byte(value))
		if pos == -1 {
			return nil, errors.New("Character not found in alphabet")
		}
		decimalData.Mul(decimalData, multiplier)
		decimalData.Add(decimalData, big.NewInt(int64(pos)))
	}

	return decimalData.Bytes(), nil
}
