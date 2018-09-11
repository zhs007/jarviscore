package jarviscrypto

import (
	"bytes"
	"testing"
)

func TestBase58(t *testing.T) {
	var b58Vectors = []struct {
		bytes   []byte
		encoded string
	}{
		{hex2bytes("4e19"), "6wi"},
		{hex2bytes("3ab7"), "5UA"},
		{hex2bytes("ae0ddc9b"), "5T3W5p"},
		{hex2bytes("65e0b4c9"), "3c3E6L"},
		{hex2bytes("25793686e9f25b6b"), "7GYJp3ZThFG"},
		{hex2bytes("94b9ac084a0d65f5"), "RspedB5CMo2"},
	}

	/* Test base-58 encoding */
	for i := 0; i < len(b58Vectors); i++ {
		got := Base58Encode(b58Vectors[i].bytes)
		if got != b58Vectors[i].encoded {
			t.Fatalf("b58encode(%v): got %s, expected %s", b58Vectors[i].bytes, got, b58Vectors[i].encoded)
		}
	}
	t.Log("success b58encode() on valid vectors")

	/* Test base-58 decoding */
	for i := 0; i < len(b58Vectors); i++ {
		got, err := Base58Decode(b58Vectors[i].encoded)
		if err != nil {
			t.Fatalf("b58decode(%s): got error %v, expected %v", b58Vectors[i].encoded, err, b58Vectors[i].bytes)
		}
		if bytes.Compare(got, b58Vectors[i].bytes) != 0 {
			t.Fatalf("b58decode(%s): got %v, expected %v", b58Vectors[i].encoded, got, b58Vectors[i].bytes)
		}
	}
	t.Log("success b58decode() on valid vectors")

	/* Test base-58 decoding of invalid strings */
	b58InvalidVectors := []string{
		"5T3IW5p", // Invalid character I
		"6Owi",    // Invalid character O
	}

	for i := 0; i < len(b58InvalidVectors); i++ {
		got, err := Base58Decode(b58InvalidVectors[i])
		if err == nil {
			t.Fatalf("b58decode(%s): got %v, expected error", b58InvalidVectors[i], got)
		}
		t.Logf("b58decode(%s): got expected err %v", b58InvalidVectors[i], err)
	}
	t.Log("success b58decode() on invalid vectors")
}

func TestBase58Check(t *testing.T) {
	var b58CheckVectors = []struct {
		ver     uint8
		bytes   []byte
		encoded string
	}{
		{0x00, hex2bytes("010966776006953D5567439E5E39F86A0D273BEE"), "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"},
		{0x00, hex2bytes("000000006006953D5567439E5E39F86A0D273BEE"), "111112LbMksD9tCRVsyW67atmDssDkHHG"},
		{0x80, hex2bytes("0C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D"), "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"},
	}

	/* Test base-58 check encoding */
	for i := 0; i < len(b58CheckVectors); i++ {
		got := Base58CheckEncode(b58CheckVectors[i].ver, b58CheckVectors[i].bytes)
		if got != b58CheckVectors[i].encoded {
			t.Fatalf("b58checkencode(0x%02x, %v): got %s, expected %s", b58CheckVectors[i].ver, b58CheckVectors[i].bytes, got, b58CheckVectors[i].encoded)
		}
	}
	t.Log("success b58checkencode() on valid vectors")

	/* Test base-58 check decoding */
	for i := 0; i < len(b58CheckVectors); i++ {
		ver, got, err := Base58CheckDecode(b58CheckVectors[i].encoded)
		if err != nil {
			t.Fatalf("b58checkdecode(%s): got error %v, expected ver %v, bytes %v", b58CheckVectors[i].encoded, err, b58CheckVectors[i].ver, b58CheckVectors[i].bytes)
		}
		if ver != b58CheckVectors[i].ver || bytes.Compare(got, b58CheckVectors[i].bytes) != 0 {
			t.Fatalf("b58checkdecode(%s): got ver %v, bytes %v, expected ver %v, bytes %v", b58CheckVectors[i].encoded, ver, got, b58CheckVectors[i].ver, b58CheckVectors[i].bytes)
		}
	}
	t.Log("success b58checkdecode() on valid vectors")

	/* Test base-58 check decoding of invalid strings */
	b58CheckInvalidVectors := []string{
		"5T3IW5p", // Invalid base58
		"6wi",     // Missing checksum
		"6UwLL9Risc3QfPqBUvKofHmBQ7wMtjzm", // Invalid checksum
	}

	for i := 0; i < len(b58CheckInvalidVectors); i++ {
		ver, got, err := Base58CheckDecode(b58CheckInvalidVectors[i])
		if err == nil {
			t.Fatalf("b58checkdecode(%s): got ver %v, bytes %v, expected error", b58CheckInvalidVectors[i], ver, got)
		}
		t.Logf("b58checkdecode(%s): got expected err %v", b58CheckInvalidVectors[i], err)
	}
	t.Log("success b58checkdecode() on invalid vectors")
}
