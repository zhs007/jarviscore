package jarviscrypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

// BASE58TABLE -
const BASE58TABLE = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Base58Encode - encodes a byte slice b into a base-58 encoded string.
func Base58Encode(b []byte) (s string) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */

	/* Convert big endian bytes to big int */
	x := new(big.Int).SetBytes(b)

	/* Initialize */
	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	/* Convert big int to string */
	for x.Cmp(zero) > 0 {
		/* x, r = (x / 58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		s = string(BASE58TABLE[r.Int64()]) + s
	}

	return s
}

// Base58Decode - decodes a base-58 encoded string into a byte slice b.
func Base58Decode(s string) (b []byte, err error) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */

	/* Initialize */
	x := big.NewInt(0)
	m := big.NewInt(58)

	/* Convert string to big int */
	for i := 0; i < len(s); i++ {
		b58index := strings.IndexByte(BASE58TABLE, s[i])
		if b58index == -1 {
			return nil, fmt.Errorf("Invalid base-58 character encountered: '%c', index %d", s[i], i)
		}
		b58value := big.NewInt(int64(b58index))
		x.Mul(x, m)
		x.Add(x, b58value)
	}

	/* Convert big int to big endian bytes */
	b = x.Bytes()

	return b, nil
}

// Base58CheckEncode - encodes version ver and byte slice b into a base-58 check encoded string.
func Base58CheckEncode(ver uint8, b []byte) (s string) {
	/* Prepend version */
	bcpy := append([]byte{ver}, b...)

	/* Create a new SHA256 context */
	sha256h := sha256.New()

	/* SHA256 Hash #1 */
	sha256h.Reset()
	sha256h.Write(bcpy)
	hash1 := sha256h.Sum(nil)

	/* SHA256 Hash #2 */
	sha256h.Reset()
	sha256h.Write(hash1)
	hash2 := sha256h.Sum(nil)

	/* Append first four bytes of hash */
	bcpy = append(bcpy, hash2[0:4]...)

	/* Encode base58 string */
	s = Base58Encode(bcpy)

	/* For number of leading 0's in bytes, prepend 1 */
	for _, v := range bcpy {
		if v != 0 {
			break
		}
		s = "1" + s
	}

	return s
}

// Base58CheckDecode - decodes base-58 check encoded string s into a version ver and byte slice b.
func Base58CheckDecode(s string) (ver uint8, b []byte, err error) {
	/* Decode base58 string */
	b, err = Base58Decode(s)
	if err != nil {
		return 0, nil, err
	}

	/* Add leading zero bytes */
	for i := 0; i < len(s); i++ {
		if s[i] != '1' {
			break
		}
		b = append([]byte{0x00}, b...)
	}

	/* Verify checksum */
	if len(b) < 5 {
		return 0, nil, fmt.Errorf("Invalid base-58 check string: missing checksum")
	}

	/* Create a new SHA256 context */
	sha256h := sha256.New()

	/* SHA256 Hash #1 */
	sha256h.Reset()
	sha256h.Write(b[:len(b)-4])
	hash1 := sha256h.Sum(nil)

	/* SHA256 Hash #2 */
	sha256h.Reset()
	sha256h.Write(hash1)
	hash2 := sha256h.Sum(nil)

	/* Compare checksum */
	if bytes.Compare(hash2[0:4], b[len(b)-4:]) != 0 {
		return 0, nil, fmt.Errorf("Invalid base-58 check string: invalid checksum")
	}

	/* Strip checksum bytes */
	b = b[:len(b)-4]

	/* Extract and strip version */
	ver = b[0]
	b = b[1:]

	return ver, b, nil
}
