package jarviscrypto

import (
	"crypto/elliptic"
	"errors"
	"math/big"
)

var secp256k1 *elliptic.CurveParams
var secp256k1A *big.Int

func init() {
	secp256k1 = new(elliptic.CurveParams)
	secp256k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	secp256k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	secp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	secp256k1.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	secp256k1.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	secp256k1.BitSize = 256

	secp256k1A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
}

// addMod computes z = (x + y) % p.
func addMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	z = new(big.Int).Add(x, y)
	z.Mod(z, p)
	return z
}

// subMod computes z = (x - y) % p.
func subMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	z = new(big.Int).Sub(x, y)
	z.Mod(z, p)
	return z
}

// mulMod computes z = (x * y) % p.
func mulMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	n := new(big.Int).Set(x)
	z = big.NewInt(0)

	for i := 0; i < y.BitLen(); i++ {
		if y.Bit(i) == 1 {
			z = addMod(z, n, p)
		}
		n = addMod(n, n, p)
	}

	return z
}

// invMod computes z = (1/x) % p.
func invMod(x *big.Int, p *big.Int) (z *big.Int) {
	z = new(big.Int).ModInverse(x, p)
	return z
}

// expMod computes z = (x^e) % p.
func expMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	z = new(big.Int).Exp(x, y, p)
	return z
}

// sqrtMod computes z = sqrt(x) % p.
func sqrtMod(x *big.Int, p *big.Int) (z *big.Int) {
	/* assert that p % 4 == 3 */
	if new(big.Int).Mod(p, big.NewInt(4)).Cmp(big.NewInt(3)) != 0 {
		panic("p is not equal to 3 mod 4!")
	}

	/* z = sqrt(x) % p = x^((p+1)/4) % p */

	/* e = (p+1)/4 */
	e := new(big.Int).Add(p, big.NewInt(1))
	e = e.Rsh(e, 2)

	z = expMod(x, e, p)
	return z
}

// Decompress decompresses coordinate x and ylsb (y's least significant bit) into a Point P on EllipticCurve ec.
func decompress(x *big.Int, ylsb uint) (*big.Int, error) {
	/* y**2 = x**3 + a*x + b  % p */
	rhs := addMod(
		addMod(
			expMod(x, big.NewInt(3), secp256k1.P),
			mulMod(secp256k1A, x, secp256k1.P),
			secp256k1.P),
		secp256k1.B, secp256k1.P)

	/* y = sqrt(rhs) % p */
	y := sqrtMod(rhs, secp256k1.P)

	/* Use -y if opposite lsb is required */
	if y.Bit(0) != (ylsb & 0x1) {
		y = subMod(big.NewInt(0), y, secp256k1.P)
	}

	if !secp256k1.IsOnCurve(x, y) {
		return y, errors.New("Compressed (x, ylsb) not on curve")
	}

	return y, nil
}
