package jarviscrypto

import (
	"errors"
	"math/big"

	"github.com/fomichev/secp256k1"
)

// var secp256k1A *big.Int

// func init() {
// 	secp256k1A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
// }

/*** Modular Arithmetic ***/

/* NOTE: Returning a new z each time below is very space inefficient, but the
 * alternate accumulator based design makes the point arithmetic functions look
 * absolutely hideous. I may still change this in the future. */

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

func decompressY(x *big.Int, ybit uint) (*big.Int, error) {
	c := secp256k1.SECP256K1().Params()

	// y^2 = x^3 + b
	// y   = sqrt(x^3 + b)
	var y, x3b big.Int
	x3b.Mul(x, x)
	x3b.Mul(&x3b, x)
	x3b.Add(&x3b, c.B)
	x3b.Mod(&x3b, c.P)
	y.ModSqrt(&x3b, c.P)

	if y.Bit(0) != ybit {
		y.Sub(c.P, &y)
	}
	if y.Bit(0) != ybit {
		return nil, errors.New("incorrectly encoded X and Y bit")
	}

	// /* y**2 = x**3 + a*x + b  % p */
	// rhs := addMod(
	// 	addMod(
	// 		expMod(x, big.NewInt(3), c.P),
	// 		mulMod(secp256k1A, x, c.P),
	// 		c.P),
	// 	c.B, c.P)

	// /* y = sqrt(rhs) % p */
	// y := sqrtMod(rhs, c.P)

	// /* Use -y if opposite lsb is required */
	// if y.Bit(0) != (ybit & 0x1) {
	// 	y = subMod(big.NewInt(0), y, c.P)
	// }

	// if !c.IsOnCurve(x, y) {
	// 	return nil, errors.New("IsOnCurve X and Y fail")
	// }

	return &y, nil
}
