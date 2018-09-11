package jarviscrypto

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/fomichev/secp256k1"
)

// var mysecp256k1 *elliptic.CurveParams
// var mysecp256k1A *big.Int

// func init() {
// 	mysecp256k1 = new(elliptic.CurveParams)
// 	mysecp256k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
// 	mysecp256k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
// 	mysecp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
// 	mysecp256k1.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
// 	mysecp256k1.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
// 	mysecp256k1.BitSize = 256

// 	mysecp256k1A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
// }

var keyPairVectors = []struct {
	// wif                    string
	wifc              string
	privbytes         []byte
	addresscompressed string
	// address_uncompressed   string
	pubbytescompressed []byte
	// pub_bytes_uncompressed []byte
	D *big.Int
	X *big.Int
	Y *big.Int
}{
	{
		// "5J1F7GHadZG3sCCKHCwg8Jvys9xUbFsjLnGec4H125Ny1V9nR6V",
		"Kx45GeUBSMPReYQwgXiKhG9FzNXrnCeutJp4yjTd5kKxCitadm3C",
		hex2bytes("18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"),
		"1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs",
		// "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
		hex2bytes("0250863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B2352"),
		// hex2bytes("0450863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B23522CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582BA6"),
		hex2int("18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725"),
		hex2int("50863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B2352"),
		hex2int("2CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582BA6"),
	},
	{
		// "5JbDYniwPgAn3YqPUkVvrCQdJsjjFx2rV2EYeg5CAH3wNncziMm",
		"Kze2PJp755t9pFWaDUzgg9MHFtwbWyuBQgSnhHTWFwqy14NafA1S",
		hex2bytes("660527765029F5F1BC6DFD5821A7FF336C10EDA391E19BB4517DB4E23E5B112F"),
		"1ChaLikBC5E2uTCA7GZh9vaMQMuRt7h1yq",
		// "17FBpEDgirwQJTvHT6ZgSirWSCbdTB9f76",
		hex2bytes("03A83B8DE893467D3A88D959C0EB4032D9CE3BF80F175D4D9E75892A3EBB8AB7E5"),
		// hex2bytes("04A83B8DE893467D3A88D959C0EB4032D9CE3BF80F175D4D9E75892A3EBB8AB7E5370F723328C24B7A97FE34063BA68F253FB08F8645D7C8B9A4FF98E3C29E7F0D"),
		hex2int("660527765029F5F1BC6DFD5821A7FF336C10EDA391E19BB4517DB4E23E5B112F"),
		hex2int("A83B8DE893467D3A88D959C0EB4032D9CE3BF80F175D4D9E75892A3EBB8AB7E5"),
		hex2int("370F723328C24B7A97FE34063BA68F253FB08F8645D7C8B9A4FF98E3C29E7F0D"),
	},
	{
		// "5KPaskZdrcPmrH3AFdpMF7FFBcYigwdrEfpBN9K5Ch4Ch6Bort4",
		"L4AgX1H3fyDWxVnXqbVzMGZsbqu11J9eKLKEkKgmYbo8bbs4K9Sq",
		hex2bytes("CF4DBE1ABCB061DB64CC87404AB736B6A56E8CDD40E9846144582240C5366758"),
		"1DP4edYeSPAF5UkXomAFKhsXwKq59r26aY",
		// "1K1EJ6Zob7mr6Wye9mF1pVaU4tpDhrYMKJ",
		hex2bytes("03F680556678E25084A82FA39E1B1DFD0944F7E69FDDAA4E03CE934BD6B291DCA0"),
		// hex2bytes("04F680556678E25084A82FA39E1B1DFD0944F7E69FDDAA4E03CE934BD6B291DCA052C10B721D34447E173721FB0151C68DE1106BADB089FB661523B8302A9097F5"),
		hex2int("CF4DBE1ABCB061DB64CC87404AB736B6A56E8CDD40E9846144582240C5366758"),
		hex2int("F680556678E25084A82FA39E1B1DFD0944F7E69FDDAA4E03CE934BD6B291DCA0"),
		hex2int("52C10B721D34447E173721FB0151C68DE1106BADB089FB661523B8302A9097F5"),
	},
	{
		// "5KTzSQJFWc48YdgxXJPb7BhnHu98TUd6C8CDNw6D2dq8fVfC5G8",
		"L4W8X7q93fipJcwN4jkhYJEzub8survHv6kVdojz6DZSpYUmJYkM",
		hex2bytes("D94F024E82D787FB38369BEA7478AA61308DC2F7080ADDF69919A881490CFF48"),
		"1Hicf8AisGTeFqhNuSTw5m5UsYbHRxDxfj",
		// "1CqhvePnxy5ZdvuunhZ7KzaqJVrNfXAk5E",
		hex2bytes("02692B035A2BB89C503E68A732596491524808BC2CC6A95061CD2CDE5151B34CD8"),
		// hex2bytes("04692B035A2BB89C503E68A732596491524808BC2CC6A95061CD2CDE5151B34CD8B9FBFB401C7BDA0C77F161ADE0AA54688412E591DAF2E3A652DB00A533645B24"),
		hex2int("D94F024E82D787FB38369BEA7478AA61308DC2F7080ADDF69919A881490CFF48"),
		hex2int("692B035A2BB89C503E68A732596491524808BC2CC6A95061CD2CDE5151B34CD8"),
		hex2int("B9FBFB401C7BDA0C77F161ADE0AA54688412E591DAF2E3A652DB00A533645B24"),
	},
	{
		// "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsrgA9tXshp",
		"KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgd9M7rFUFqJ5Vvujp",
		hex2bytes("0000000000000000000000000000000000000000000000000000000000000400"),
		"1G73bvYR97QGVb8bfeX2TqvSKietBDybQC",
		// "17imJe7o4mpq2MMfZ328evDJQfbt6ShvxA",
		hex2bytes("03241FEBB8E23CBD77D664A18F66AD6240AAEC6ECDC813B088D5B901B2E285131F"),
		// hex2bytes("04241FEBB8E23CBD77D664A18F66AD6240AAEC6ECDC813B088D5B901B2E285131F513378D9FF94F8D3D6C420BD13981DF8CD50FD0FBD0CB5AFABB3E66F2750026D"),
		hex2int("0000000000000000000000000000000000000000000000000000000000000400"),
		hex2int("241FEBB8E23CBD77D664A18F66AD6240AAEC6ECDC813B088D5B901B2E285131F"),
		hex2int("513378D9FF94F8D3D6C420BD13981DF8CD50FD0FBD0CB5AFABB3E66F2750026D"),
	},
}

var invalidPublicKeyBytesVectors = [][]byte{
	hex2bytes("0250863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B23"),                                                                 /* Short compressed */
	hex2bytes("0450863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B23522CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582B"), /* Short uncompressed */
	hex2bytes("03A83B8DFF93467D3A88D959C0EB4032FFFF3BF80F175D4D9E75892A3EBB8FF7E5"),                                                               /* Invalid compressed */
	hex2bytes("02c8e337cee51ae9af3c0ef923705a0cb1b76f7e8463b3d3060a1c8d795f9630fd"),                                                               /* Invalid compressed */
}

var wifInvalidVectors = []string{
	"5T3IW5p", // Invalid base58
	"6wi",     // Missing checksum
	"6Mcb23muAxyXaSMhmB6B1mqkvLdWhtuFZmnZsxDczHRraMcNG",    // Invalid checksum
	"huzKTSifqNioknFPsoA7uc359rRHJQHRg42uiKn6P8Rnv5qxV5",   // Invalid version byte
	"yPoVP5njSzmEVK4VJGRWWAwqnwCyLPRcMm5XyrKgYUpeXtGyM",    // Invalid private key byte length
	"Kx45GeUBSMPReYQwgXiKhG9FzNXrnCeutJp4yjTd5kKxCj6CAKu3", // Invalid private key suffix byte
}

func TestPrivateKeyDerive(t *testing.T) {
	priv := PrivateKey{priKey: &ecdsa.PrivateKey{}}
	priv.priKey.Curve = secp256k1.SECP256K1()

	for i := 0; i < len(keyPairVectors); i++ {
		priv.priKey.D = keyPairVectors[i].D

		/* Derive public key from private key */
		priv.derive()

		if priv.priKey.PublicKey.X.Cmp(keyPairVectors[i].X) != 0 || priv.priKey.PublicKey.Y.Cmp(keyPairVectors[i].Y) != 0 {
			t.Fatalf("derived public key does not match test vector on index %d", i)
		}
	}
	t.Log("success PrivateKey derive()")
}

func TestCheckWIF(t *testing.T) {
	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		// got, err := CheckWIF(keyPairVectors[i].wif)
		// if got == false {
		// 	t.Fatalf("CheckWIF(%s): got false, error %v, expected true", keyPairVectors[i].wif, err)
		// }
		got, err := CheckWIF(keyPairVectors[i].wifc)
		if got == false {
			t.Fatalf("CheckWIF(%s): got false, error %v, expected true", keyPairVectors[i].wifc, err)
		}
	}
	t.Log("success CheckWIF() on valid vectors")

	/* Check invalid vectors */
	for i := 0; i < len(wifInvalidVectors); i++ {
		got, err := CheckWIF(wifInvalidVectors[i])
		if got == true {
			t.Fatalf("CheckWIF(%s): got true, expected false", wifInvalidVectors[i])
		}
		t.Logf("CheckWIF(%s): got false, err %v", wifInvalidVectors[i], err)
	}
	t.Log("success CheckWIF() on invalid vectors")
}

func TestPrivateKeyFromBytes(t *testing.T) {
	var priv PrivateKey

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		err := priv.FromBytes(keyPairVectors[i].privbytes)
		if err != nil {
			t.Fatalf("priv.FromBytes(D): got error %v, expected success on index %d", err, i)
		}
		if priv.priKey.D.Cmp(keyPairVectors[i].D) != 0 {
			t.Fatalf("private key does not match test vector on index %d", i)
		}
		if priv.priKey.X.Cmp(keyPairVectors[i].X) != 0 || priv.priKey.Y.Cmp(keyPairVectors[i].Y) != 0 {
			t.Fatalf("public key does not match test vector on index %d", i)
		}
	}
	t.Log("success PrivateKey FromBytes() on valid vectors")

	/* Invalid short private key */
	err := priv.FromBytes(keyPairVectors[0].privbytes[0:31])
	if err == nil {
		t.Fatalf("priv.FromBytes(D): got success, expected error")
	}
	/* Invalid long private key */
	err = priv.FromBytes(append(keyPairVectors[0].privbytes, []byte{0x00}...))
	if err == nil {
		t.Fatalf("priv.FromBytes(D): got success, expected error")
	}

	t.Log("success PrivateKey FromBytes() on invaild vectors")
}

func TestPrivateKeyToBytes(t *testing.T) {
	priv := NewPrivateKey()

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		priv.priKey.D = keyPairVectors[i].D
		b := priv.ToPrivateBytes()
		if bytes.Compare(keyPairVectors[i].privbytes, b) != 0 {
			t.Fatalf("private key bytes do not match test vector in index %d", i)
		}
	}
	t.Log("success PrivateKey ToBytes()")
}

func TestPrivateKeyFromWIF(t *testing.T) {
	priv := NewPrivateKey()

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		// /* Import WIF */
		// err := priv.FromWIF(keyPairVectors[i].wif)
		// if err != nil {
		// 	t.Fatalf("priv.FromWIF(%s): got error %v, expected success", keyPairVectors[i].wif, err)
		// }
		// if priv.D.Cmp(keyPairVectors[i].D) != 0 {
		// 	t.Fatalf("private key does not match test vector on index %d", i)
		// }
		// if priv.X.Cmp(keyPairVectors[i].X) != 0 || priv.Y.Cmp(keyPairVectors[i].Y) != 0 {
		// 	t.Fatalf("public key does not match test vector on index %d", i)
		// }

		/* Import WIFC */
		err := priv.FromWIFC(keyPairVectors[i].wifc)
		if err != nil {
			t.Fatalf("priv.FromWIF(%s): got error %v, expected success", keyPairVectors[i].wifc, err)
		}
		if priv.priKey.D.Cmp(keyPairVectors[i].D) != 0 {
			t.Fatalf("private key does not match test vector on index %d", i)
		}
		if priv.priKey.X.Cmp(keyPairVectors[i].X) != 0 || priv.priKey.Y.Cmp(keyPairVectors[i].Y) != 0 {
			t.Fatalf("public key does not match test vector on index %d", i)
		}
	}
	t.Log("success PrivateKey FromWIF() on valid vectors")

	/* Check invalid vectors */
	for i := 0; i < len(wifInvalidVectors); i++ {
		err := priv.FromWIFC(wifInvalidVectors[i])
		if err == nil {
			t.Fatalf("priv.FromWIF(%s): got success, expected error", wifInvalidVectors[i])
		}
		t.Logf("priv.FromWIF(%s): got err %v", wifInvalidVectors[i], err)
	}
	t.Log("success PrivateKey FromWIF() on invalid vectors")
}

func TestPrivateKeyToWIF(t *testing.T) {
	priv := NewPrivateKey()

	// /* Check valid vectors */
	// for i := 0; i < len(keyPairVectors); i++ {
	// 	/* Export WIF */
	// 	priv.D = keyPairVectors[i].D
	// 	wif := priv.ToWIF()
	// 	if wif != keyPairVectors[i].wif {
	// 		t.Fatalf("priv.ToWIF() %s != expected %s", wif, keyPairVectors[i].wif)
	// 	}
	// }
	// t.Log("success PrivateKey ToWIF()")

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		/* Export WIFC */
		priv.priKey.D = keyPairVectors[i].D
		wifc := priv.ToWIFC()
		if wifc != keyPairVectors[i].wifc {
			t.Fatalf("priv.ToWIFC() %s != expected %s", wifc, keyPairVectors[i].wifc)
		}
	}
	t.Log("success PrivateKey ToWIFC()")
}

func TestPublicKeyToBytes(t *testing.T) {
	pub := NewPublicKey()

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		pub.pubKey.X = keyPairVectors[i].X
		pub.pubKey.Y = keyPairVectors[i].Y

		bytescompressed := pub.ToBytes()
		if bytes.Compare(keyPairVectors[i].pubbytescompressed, bytescompressed) != 0 {
			t.Fatalf("public key compressed bytes do not match test vectors on index %d", i)
		}

		// bytesuncompressed := pub.ToBytesUncompressed()
		// if bytes.Compare(keyPairVectors[i].pub_bytes_uncompressed, bytes_uncompressed) != 0 {
		// 	t.Fatalf("public key uncompressed bytes do not match test vectors on index %d", i)
		// }
	}
	t.Log("success PublicKey ToBytes() and ToBytesUncompressed()")
}

func TestPublicKeyFromBytes(t *testing.T) {
	pub := NewPublicKey()

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		err := pub.FromBytes(keyPairVectors[i].pubbytescompressed)
		if err != nil {
			t.Fatalf("pub.FromBytes(): got error %v, expected success on index %d", err, i)
		}
		if pub.pubKey.X.Cmp(keyPairVectors[i].X) != 0 || pub.pubKey.Y.Cmp(keyPairVectors[i].Y) != 0 {
			t.Fatalf("public key does not match test vectors on index %d", i)
		}

		// err = pub.FromBytes(keyPairVectors[i].pub_bytes_uncompressed)
		// if err != nil {
		// 	t.Fatalf("pub.FromBytes(): got error %v, expected success on index %d", err, i)
		// }
		// if pub.X.Cmp(keyPairVectors[i].X) != 0 || pub.Y.Cmp(keyPairVectors[i].Y) != 0 {
		// 	t.Fatalf("public key does not match test vectors on index %d", i)
		// }
	}
	t.Log("success PublicKey FromBytes() on valid vectors")

	/* Check invalid vectors */
	for i := 0; i < len(invalidPublicKeyBytesVectors); i++ {
		err := pub.FromBytes(invalidPublicKeyBytesVectors[i])
		if err == nil {
			t.Fatal("pub.FromBytes(): got success, expected error")
		}
		t.Logf("pub.FromBytes(): got error %v", err)
	}
	t.Log("success PublicKey FromBytes() on invalid vectors")
}

func TestToAddress(t *testing.T) {
	pub := NewPublicKey()

	/* Check valid vectors */
	for i := 0; i < len(keyPairVectors); i++ {
		pub.pubKey.X = keyPairVectors[i].X
		pub.pubKey.Y = keyPairVectors[i].Y

		addresscompressed := pub.ToAddress()
		if addresscompressed != keyPairVectors[i].addresscompressed {
			t.Fatalf("public key compressed address does not match test vectors on index %d", i)
		}

		// address_uncompressed := pub.ToAddressUncompressed()
		// if address_uncompressed != keyPairVectors[i].address_uncompressed {
		// 	t.Fatalf("public key uncompressed address does not match test vectors on index %d", i)
		// }
	}
	t.Log("success PublicKey ToAddress() and ToAddressUncompressed()")
}
