package jarviscrypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"log"
)

// var secp256k1N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

const privKeyBytesLen = 32
const version = byte(0x00)
const addressChecksumLen = 4

func newKeyPair() ([]byte, []byte) {
	// curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(secp256k1, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	d := private.D.Bytes()
	// b := make([]byte, 0, privKeyBytesLen)
	// priKet := paddedAppend(privKeyBytesLen, b, d)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return d, pubKey
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
