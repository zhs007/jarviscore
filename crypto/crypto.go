// https://github.com/vsergeev/btckeygenie

package jarviscrypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"github.com/fomichev/secp256k1"
	"golang.org/x/crypto/ripemd160"
)

const privKeyBytesLen = 32

// PrivateKey -
type PrivateKey struct {
	priKey *ecdsa.PrivateKey
}

// PublicKey -
type PublicKey struct {
	pubKey *ecdsa.PublicKey
}

// GenerateKey -
func GenerateKey() PrivateKey {
	privkey, err := ecdsa.GenerateKey(secp256k1.SECP256K1(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	return PrivateKey{priKey: privkey}
}

// derive derives a Bitcoin public key from a Bitcoin private key.
func (key *PrivateKey) derive() {
	/* See Certicom's SEC1 3.2.1, pg.23 */

	/* Derive public key from Q = d*G */
	x, y := key.priKey.Curve.ScalarBaseMult(key.priKey.D.Bytes())

	/* Check that Q is on the curve */
	if !key.priKey.Curve.IsOnCurve(x, y) {
		panic("Catastrophic math logic failure in public key derivation")
	}

	key.priKey.PublicKey.X = x
	key.priKey.PublicKey.Y = y
}

// ToPrivateBytes -
func (key *PrivateKey) ToPrivateBytes() []byte {
	priKeyBytes := key.priKey.D.Bytes()

	priKeyBytes = append(bytes.Repeat([]byte{0x00}, privKeyBytesLen-len(priKeyBytes)), priKeyBytes...)

	return priKeyBytes
}

// FromBytes - converts a 32-byte byte slice to a Bitcoin private key and derives the corresponding Bitcoin public key.
func (key *PrivateKey) FromBytes(b []byte) (err error) {
	key.priKey = new(ecdsa.PrivateKey)
	key.priKey.PublicKey.Curve = secp256k1.SECP256K1()

	if len(b) != 32 {
		return fmt.Errorf("Invalid private key bytes length %d, expected 32", len(b))
	}

	key.priKey.D = new(big.Int).SetBytes(b)

	/* Derive public key from private key */
	key.derive()

	return nil
}

// ToWIFC - prikey to WIFC
func (key *PrivateKey) ToWIFC() string {
	privkey := key.ToPrivateBytes()
	/* Append 0x01 to tell Bitcoin wallet to use compressed public keys */
	privkey = append(privkey, []byte{0x01}...)

	/* Convert bytes to base-58 check encoded string with version 0x80 */
	wifc := Base58CheckEncode(0x80, privkey)

	return wifc
}

// FromWIFC converts a Wallet Import Format string to a Bitcoin private key and derives the corresponding Bitcoin public key.
func (key *PrivateKey) FromWIFC(wif string) (err error) {
	/* See https://en.bitcoin.it/wiki/Wallet_import_format */

	/* Base58 Check Decode the WIF string */
	ver, privbytes, err := Base58CheckDecode(wif)
	if err != nil {
		return err
	}

	/* Check that the version byte is 0x80 */
	if ver != 0x80 {
		return fmt.Errorf("Invalid WIF version 0x%02x, expected 0x80", ver)
	}

	/* If the private key bytes length is 33, check that suffix byte is 0x01 (for compression) and strip it off */
	if len(privbytes) == 33 {
		if privbytes[len(privbytes)-1] != 0x01 {
			return fmt.Errorf("Invalid private key, unknown suffix byte 0x%02x", privbytes[len(privbytes)-1])
		}
		privbytes = privbytes[0:32]
	} else {
		return fmt.Errorf("Invalid private key, len fail")
	}

	/* Convert from bytes to a private key */
	err = key.FromBytes(privbytes)
	if err != nil {
		return err
	}

	/* Derive public key from private key */
	key.derive()

	return nil
}

// ToPublicBytes - converts a Bitcoin public key to a 33-byte byte slice with point compression.
func (key *PrivateKey) ToPublicBytes() (b []byte) {

	/* See Certicom SEC1 2.3.3, pg. 10 */
	x := key.priKey.PublicKey.X.Bytes()

	/* Pad X to 32-bytes */
	paddedx := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if key.priKey.PublicKey.Y.Bit(0) == 0 {
		return append([]byte{0x02}, paddedx...)
	}
	return append([]byte{0x03}, paddedx...)
}

// ToAddress - converts a Bitcoin public key to a compressed Bitcoin address string.
func (key *PrivateKey) ToAddress() (address string) {
	/* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */

	/* Convert the public key to bytes */
	pubbytes := key.ToPublicBytes()

	/* SHA256 Hash */
	sha256h := sha256.New()
	sha256h.Reset()
	sha256h.Write(pubbytes)
	pubhash1 := sha256h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160h := ripemd160.New()
	ripemd160h.Reset()
	ripemd160h.Write(pubhash1)
	pubhash2 := ripemd160h.Sum(nil)

	/* Convert hash bytes to base58 check encoded sequence */
	address = Base58CheckEncode(0x00, pubhash2)

	return address
}

// FromBytes - converts a byte slice (either with or without point compression) to a Bitcoin public key.
func (key *PublicKey) FromBytes(b []byte) (err error) {
	key.pubKey = new(ecdsa.PublicKey)
	key.pubKey.Curve = secp256k1.SECP256K1()

	/* See Certicom SEC1 2.3.4, pg. 11 */

	if len(b) < 33 {
		return fmt.Errorf("Invalid public key bytes length %d, expected at least 33", len(b))
	}

	if b[0] == 0x02 || b[0] == 0x03 {
		/* Compressed public key */

		if len(b) != 33 {
			return fmt.Errorf("Invalid public key bytes length %d, expected 33", len(b))
		}

		x := new(big.Int).SetBytes(b[1:33])
		y, err := decompressY(x, uint(b[0]&0x1))
		if err != nil {
			return fmt.Errorf("Invalid compressed public key bytes, decompression error: %v", err)
		}

		key.pubKey.X = x
		key.pubKey.Y = y

	} else {
		return fmt.Errorf("Invalid public key prefix byte 0x%02x, expected 0x02, 0x03", b[0])
	}

	return nil
}

// ToBytes - converts a Bitcoin public key to a 33-byte byte slice with point compression.
func (key *PublicKey) ToBytes() (b []byte) {

	/* See Certicom SEC1 2.3.3, pg. 10 */
	x := key.pubKey.X.Bytes()

	/* Pad X to 32-bytes */
	paddedx := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if key.pubKey.Y.Bit(0) == 0 {
		return append([]byte{0x02}, paddedx...)
	}
	return append([]byte{0x03}, paddedx...)
}

// ToAddress - converts a Bitcoin public key to a compressed Bitcoin address string.
func (key *PublicKey) ToAddress() (address string) {
	/* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */

	/* Convert the public key to bytes */
	pubbytes := key.ToBytes()

	/* SHA256 Hash */
	sha256h := sha256.New()
	sha256h.Reset()
	sha256h.Write(pubbytes)
	pubhash1 := sha256h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160h := ripemd160.New()
	ripemd160h.Reset()
	ripemd160h.Write(pubhash1)
	pubhash2 := ripemd160h.Sum(nil)

	/* Convert hash bytes to base58 check encoded sequence */
	address = Base58CheckEncode(0x00, pubhash2)

	return address
}
