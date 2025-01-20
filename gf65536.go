package gf65536

import "fmt"

const Default Field = 0x1002b

func Add(x, y uint16) uint16 {
	return Default.Add(x, y)
}

func Mul(x, y uint16) uint16 {
	return Default.Mul(x, y)
}

func Inv(x uint16) uint16 {
	return Default.Inv(x)
}

type Field uint64

func New(poly uint64) (Field, error) {
	if nBits(poly) != 17 {
		return 0, fmt.Errorf("polynomial must be of degree 16")
	} else if reducible(poly) {
		return 0, fmt.Errorf("polynomial must be irreducible")
	}
	return Field(poly), nil
}

func (f Field) Add(x, y uint16) uint16 {
	return x ^ y
}

func (f Field) Mul(x, y uint16) uint16 {
	return uint16(mul(uint64(f), uint64(x), uint64(y)))
}

func (f Field) Inv(x uint16) uint16 {
	return uint16(inv(uint64(f), uint64(x)))
}

func mul(p, x, y uint64) uint64 {
	_, rem := polyDiv(polyMul(x, y), p)
	return rem
}

func inv(poly, x uint64) uint64 {
	var (
		a, b uint64 = poly, x
		u, v uint64 = 0, 1
		s, t uint64 = 1, 0
	)

	for b != 0 {
		q, _ := polyDiv(a, b)
		a, b = b, a^polyMul(q, b)
		s, u = u, s^polyMul(q, u)
		t, v = v, t^polyMul(q, v)
	}

	return t
}

func polyMul(x, y uint64) uint64 {
	var z uint64
	for y > 0 {
		if y&1 == 1 {
			z ^= x
		}
		x <<= 1
		y >>= 1
	}
	return z
}

func polyDiv(p, q uint64) (uint64, uint64) {
	var (
		quot uint64
		np   = nBits(p)
		nq   = nBits(q)
	)

	for ; np >= nq; np-- {
		if p&(1<<(np-1)) != 0 {
			p ^= q << (np - nq)
			quot |= 1 << (np - nq)
		}
	}

	return quot, p
}

func reducible(p uint64) bool {
	factorMax := uint64(1) << (nBits(p)/2 + 1)
	for factor := uint64(2); factor < factorMax; factor++ {
		if _, rem := polyDiv(p, factor); rem == 0 {
			return true
		}
	}

	return false
}

func nBits(p uint64) uint64 {
	var n uint64
	for ; p > 0; p >>= 1 {
		n++
	}
	return n
}
