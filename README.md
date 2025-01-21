# gf65536 - arithmetic over GF(2^16)

**Attention**: This implementation is not hardened against side-channel attacks.
Be cautious when using it in security-critical applications.

This package provides addition, multiplication, and the multiplicative inverse
over the galois field (finite field) GF(2^16). Exponentiation and discrete
logarithm are not implemented.

The default irreducible polynomial is `x^16 + x^5 + x^3 + x + 1` (0x1002b), 
which is the least weight irreducible polynomial of degree 16.
