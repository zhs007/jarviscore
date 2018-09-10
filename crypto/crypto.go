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

	"golang.org/x/crypto/ripemd160"
)

// // var secp256k1N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

const privKeyBytesLen = 32

// const version = byte(0x00)
// const addressChecksumLen = 4

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
	privkey, err := ecdsa.GenerateKey(secp256k1, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	return PrivateKey{priKey: privkey}
}

// func newKeyPair() ([]byte, []byte) {
// 	// curve := elliptic.P256()
// 	private, err := ecdsa.GenerateKey(secp256k1, rand.Reader)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	d := private.D.Bytes()

// 	paddedd := append(bytes.Repeat([]byte{0x00}, privKeyBytesLen-len(d)), d...)
// 	// b := make([]byte, 0, privKeyBytesLen)
// 	// priKet := paddedAppend(privKeyBytesLen, b, d)
// 	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

// 	return paddedd, pubKey
// }

// derive derives a Bitcoin public key from a Bitcoin private key.
func (key *PrivateKey) derive() {
	/* See Certicom's SEC1 3.2.1, pg.23 */

	/* Derive public key from Q = d*G */
	x, y := secp256k1.ScalarBaseMult(key.priKey.D.Bytes())

	/* Check that Q is on the curve */
	if !secp256k1.IsOnCurve(x, y) {
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
	key.priKey.PublicKey.Curve = secp256k1

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
	key.pubKey.Curve = secp256k1

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
		y, err := decompress(x, uint(b[0]&0x1))
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

// func newKeyPair() ([]byte, []byte) {
// 	curve := elliptic.P256()
// 	private, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	d := private.D.Bytes()
// 	b := make([]byte, 0, privKeyBytesLen)
// 	priKet := paddedAppend(privKeyBytesLen, b, d)
// 	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

// 	return priKet, pubKey
// }

// // paddedAppend appends the src byte slice to dst, returning the new slice.
// // If the length of the source is smaller than the passed size, leading zero
// // bytes are appended to the dst slice before appending src.
// func paddedAppend(size uint, dst, src []byte) []byte {
// 	for i := 0; i < int(size)-len(src); i++ {
// 		dst = append(dst, 0)
// 	}
// 	return append(dst, src...)
// }

// // Wallet stores private and public keys
// type Wallet struct {
// 	PrivateKey ecdsa.PrivateKey
// 	PublicKey  []byte
// }

// // NewWallet creates and returns a Wallet
// func NewWallet() *Wallet {
// 	private, public := newKeyPair()
// 	wallet := Wallet{private, public}

// 	return &wallet
// }

// // GetAddress returns wallet address
// func (w Wallet) GetAddress() []byte {
// 	pubKeyHash := HashPubKey(w.PublicKey)

// 	versionedPayload := append([]byte{version}, pubKeyHash...)
// 	checksum := checksum(versionedPayload)

// 	fullPayload := append(versionedPayload, checksum...)
// 	address := Base58Encode(fullPayload)

// 	return address
// }

// // HashPubKey hashes public key
// func HashPubKey(pubKey []byte) []byte {
// 	publicSHA256 := sha256.Sum256(pubKey)

// 	RIPEMD160Hasher := ripemd160.New()
// 	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

// 	return publicRIPEMD160
// }

// // ValidateAddress check if address if valid
// func ValidateAddress(address string) bool {
// 	pubKeyHash := Base58Decode([]byte(address))
// 	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
// 	version := pubKeyHash[0]
// 	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
// 	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

// 	return bytes.Compare(actualChecksum, targetChecksum) == 0
// }

// // Checksum generates a checksum for a public key
// func checksum(payload []byte) []byte {
// 	firstSHA := sha256.Sum256(payload)
// 	secondSHA := sha256.Sum256(firstSHA[:])

// 	return secondSHA[:addressChecksumLen]
// }

// func newKeyPair() (ecdsa.PrivateKey, []byte) {
// 	curve := elliptic.P256()
// 	private, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

// 	return *private, pubKey
// }
