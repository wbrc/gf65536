# gf65536 - arithmetic over GF(2^16)

> **Attention**: This implementation is not hardened against side-channel
> attacks. Be cautious when using it in security-critical applications.

This package provides addition, multiplication, and the multiplicative inverse
over the galois field (finite field) GF(2^16). Exponentiation and discrete
logarithm are not implemented.

The default irreducible polynomial is `x^16 + x^5 + x^3 + x + 1` (0x1002b), 
which is the least weight irreducible polynomial of degree 16.

## Example
```go
f, err := gf65536.New(0x1002d)
if err != nil {
	panic(err)
}

// x = 5/2 + 1000*45001 + 60000*2001
t1 := f.Mul(5, f.Inv(2))
t2 := f.Mul(1000, 45001)
t3 := f.Mul(60000, 2001)
x := f.Add(t1, f.Add(t2, t3))
fmt.Printf("x = %d\n", x)
```