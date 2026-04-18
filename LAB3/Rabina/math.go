package main

import (
	"math/big"
)

func FastExp(a, z, n *big.Int) *big.Int {
	res := big.NewInt(1)
	base := new(big.Int).Mod(a, n)
	exp := new(big.Int).Set(z)

	for exp.Sign() > 0 {
		if exp.Bit(0) == 1 {
			res.Mul(res, base).Mod(res, n)
		}
		base.Mul(base, base).Mod(base, n)
		exp.Rsh(exp, 1)
	}
	return res
}

func ExtendedEuclid(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	if a.Sign() == 0 {
		return b, big.NewInt(0), big.NewInt(1)
	}
	gcd, x1, y1 := ExtendedEuclid(new(big.Int).Mod(b, a), a)
	x := new(big.Int).Sub(y1, new(big.Int).Mul(new(big.Int).Div(b, a), x1))
	y := x1
	return gcd, x, y
}

func RabinDecryptByte(c, b, n, p, q *big.Int) byte {
	four := big.NewInt(4)
	D := new(big.Int).Mul(b, b)
	D.Add(D, new(big.Int).Mul(four, c))
	D.Mod(D, n)

	pExp := new(big.Int).Add(p, big.NewInt(1))
	pExp.Div(pExp, big.NewInt(4))
	mp := FastExp(D, pExp, p)

	qExp := new(big.Int).Add(q, big.NewInt(1))
	qExp.Div(qExp, big.NewInt(4))
	mq := FastExp(D, qExp, q)

	_, yp, yq := ExtendedEuclid(p, q)

	t1 := new(big.Int).Mul(yp, p)
	t1.Mul(t1, mq)
	
	t2 := new(big.Int).Mul(yq, q)
	t2.Mul(t2, mp)

	roots := make([]*big.Int, 4)
	roots[0] = new(big.Int).Add(t1, t2)
	roots[0].Mod(roots[0], n)
	
	roots[1] = new(big.Int).Sub(n, roots[0])
	
	roots[2] = new(big.Int).Sub(t1, t2)
	roots[2].Mod(roots[2], n)
	
	roots[3] = new(big.Int).Sub(n, roots[2])

	inv2 := new(big.Int).ModInverse(big.NewInt(2), n)

	for _, d := range roots {
		m := new(big.Int).Sub(d, b)
		m.Mul(m, inv2)
		m.Mod(m, n)

		if m.Cmp(big.NewInt(256)) < 0 && m.Sign() >= 0 {
			return byte(m.Uint64())
		}
	}
	return 0
}