package gf65536

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func Test_polyMul(t *testing.T) {
	type args struct {
		x uint64
		y uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "(x^5+x^2+1)*1",
			args: args{0b100101, 1},
			want: 0b100101,
		},
		{
			name: "(x^5+x^2+1)*0",
			args: args{0b100101, 0},
			want: 0,
		},
		{
			name: "(x^5+x^2+1)*(x^3+x^2)",
			args: args{0b100101, 0b1100},
			want: 0b110111100,
		},
		{
			name: "(x^7+x^6+1)*(x^6+x^5+x+1)",
			args: args{0b11000001, 0b1100011},
			want: 0b10100100100011,
		},
		{
			name: "(x^6+x^5+x+1)*(x^7+x^6+1)",
			args: args{0b1100011, 0b11000001},
			want: 0b10100100100011,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := polyMul(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("polyMul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nBits(t *testing.T) {
	tests := []struct {
		x uint64
		y uint64
	}{
		{x: 0b0, y: 0},
		{x: 0b1, y: 1},
		{x: 0b11, y: 2},
		{x: 0b101, y: 3},
		{x: 0b1100101, y: 7},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%0b", tt.x), func(t *testing.T) {
			if got := nBits(tt.x); got != tt.y {
				t.Errorf("nBits() = %v, want %v", got, tt.y)
			}
		})
	}
}

func Test_polyDiv(t *testing.T) {
	type args struct {
		p uint64
		q uint64
	}
	tests := []struct {
		name  string
		args  args
		wantQ uint64
		wantR uint64
	}{
		{
			name:  "(x+1)/x",
			args:  args{0b11, 0b10},
			wantQ: 1,
			wantR: 1,
		},
		{
			name:  "(x^7+x^6+x^5+x^2+x+1)/1",
			args:  args{0b11100111, 1},
			wantQ: 0b11100111,
			wantR: 0,
		},
		{
			name:  "(x^7+x^6+x^5+x^2+x+1)/(x^3+x^2+1)",
			args:  args{0b11100111, 0b1101},
			wantQ: 0b10100,
			wantR: 0b11,
		},
		{
			name:  "(x^11+x^7+x^2+x^1)/(x^8+x^7+x^6+x^3+x^2+x)",
			args:  args{0b100010000110, 0b111001110},
			wantQ: 0b1101,
			wantR: 0,
		},
		{
			name:  "(x^2+1)/(x^5+x^4+x+1)",
			args:  args{0b101, 0b110011},
			wantQ: 0,
			wantR: 0b101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, r := polyDiv(tt.args.p, tt.args.q)
			if q != tt.wantQ {
				t.Errorf("polyDiv() got = %v, want %v", q, tt.wantQ)
			}
			if r != tt.wantR {
				t.Errorf("polyDiv() got1 = %v, want %v", r, tt.wantR)
			}
		})
	}
}

func Test_reducible(t *testing.T) {
	tests := []struct {
		name string
		p    uint64
		want bool
	}{
		{
			name: "reducible",
			p:    0b110111100,
			want: true,
		},
		{
			name: "reducible",
			p:    0b10100100100011,
			want: true,
		},
		{
			name: "irreducible",
			p:    0b10000011,
			want: false,
		},
		{
			name: "irreducible",
			p:    0b100111001,
			want: false,
		},
		{
			name: "irreducible",
			p:    0b111111001,
			want: false,
		},
		{
			name: "irreducible",
			p:    0b10001000000001011,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reducible(tt.p); got != tt.want {
				t.Errorf("reducible(%0b) = %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}

func Test_inv(t *testing.T) {
	// gen degree 8 irreducible polynomial and test all inverses of GF(2^8)
	var poly uint64 = 0b100000000
	for ; reducible(poly); poly++ {
	}

	for x := uint64(1); x < 256; x++ {
		inv := inv(poly, x)
		if mul(poly, x, inv) != 1 {
			t.Errorf("inv(%0b, %d) = %0b, want 1", poly, x, inv)
		}
	}
}

func TestNew(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		_, err := New(0xff)
		if err == nil {
			t.Errorf("New(0xff) should return an error")
		}
	})

	t.Run("reducible", func(t *testing.T) {
		_, err := New(0x10020)
		if err == nil {
			t.Errorf("New(0x10020) should return an error")
		}
	})

	t.Run("ok", func(t *testing.T) {
		_, err := New(0x1002b)
		if err != nil {
			t.Errorf("New(0x1002b) should not return an error")
		}
	})
}

func TestDefault(t *testing.T) {
	// a + b*c = d
	// a = d + b*c
	// b = (d + a) * inv(c)
	// c = (d + a) * inv(b)

	var a, b, c, d uint16
	for range 1000 {
		a = uint16(rand.IntN(1<<16-2) + 1)
		b = uint16(rand.IntN(1<<16-2) + 1)
		c = uint16(rand.IntN(1<<16-2) + 1)
		d = Add(a, Mul(b, c))

		aa := Add(d, Mul(b, c))
		if a != aa {
			t.Errorf("a = %d, want %d", a, aa)
		}

		bb := Mul(Add(d, a), Inv(c))
		if b != bb {
			t.Errorf("b = %d, want %d", b, bb)
		}

		cc := Mul(Add(d, a), Inv(b))
		if c != cc {
			t.Errorf("c = %d, want %d", c, cc)
		}
	}
}
