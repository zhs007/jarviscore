package jarviscrypto

import (
	"errors"
	"math/big"

	"github.com/fomichev/secp256k1"
)

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
	return &y, nil
}
