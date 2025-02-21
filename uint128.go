// Package wide provides implementations of Int128 and Uint128 for Go. See README.md for more information.
package wide

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/Overclock-Validator/wide/internal/bits"
)

// Uint128 is a representation of an unsigned 128-bit integer
type Uint128 struct {
	Hi, Lo uint64
}

// String returns a hexadecimal representation of a Uint128
func (x Uint128) String() string {
	if x.Hi == 0 {
		return fmt.Sprintf("%#x", x.Lo) // ignore leading 0's
	}
	return fmt.Sprintf("%#x%016x", x.Hi, x.Lo)
}

// NewUint128 returns a Uint128 from the high and low 64 bits
func NewUint128(hi, lo uint64) Uint128 {
	return Uint128{Hi: hi, Lo: lo}
}

// Uint128FromBigInt returns a Uint128 from a big.Int
func Uint128FromBigInt(a *big.Int) (z Uint128) {
	z.Lo = a.Uint64()
	b := new(big.Int).Rsh(a, int64Size)
	z.Hi = b.Uint64()
	return z
}

// Uint128FromUint64 returns a Uint128 from a uint64
func Uint128FromUint64(x uint64) Uint128 {
	return Uint128{Hi: 0, Lo: x}
}

// RandUint128 returns a pseudo-random Uint128
func RandUint128() (z Uint128) {
	z.Hi = rand.Uint64()
	z.Lo = rand.Uint64()
	return z
}

// Add returns the sum of two Uint128's
func (x Uint128) Add(y Uint128) (z Uint128) {
	z.Hi = x.Hi + y.Hi
	z.Lo = x.Lo + y.Lo
	if z.Lo < x.Lo {
		z.Hi++
	}
	return z
}

// And returns the bitwise AND of two Uint128's
func (x Uint128) And(y Uint128) (z Uint128) {
	z.Hi = x.Hi & y.Hi
	z.Lo = x.Lo & y.Lo
	return z
}

// AndNot returns the bitwise AndNot of two Uint128's
func (x Uint128) AndNot(y Uint128) (z Uint128) {
	z.Hi = x.Hi &^ y.Hi
	z.Lo = x.Lo &^ y.Lo
	return z
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x Uint128) Cmp(y Uint128) int {
	switch {
	case x.Hi > y.Hi:
		return 1
	case x.Hi < y.Hi:
		return -1
	case x.Lo > y.Lo:
		return 1
	case x.Lo < y.Lo:
		return -1
	default:
		return 0
	}
}

// Dec returns the predecessor of a Uint128
func (x Uint128) Dec() (z Uint128) {
	z.Lo = x.Lo - 1
	if z.Lo > x.Lo {
		z.Hi = x.Hi - 1
		return z
	}
	z.Hi = x.Hi
	return z
}

// Div returns the quotient corresponding to the provided dividend and divisor
//
// Div panics on division by 0. It checks some common/faster cases before fully committing to long division. This can probably be further optimized by
// implementing a successive approximation algorithm, with an initial seed value determined by a 64-bit division of the most significant bits.
func (x Uint128) Div(d Uint128) (q Uint128) {
	q, _ = x.DivMod(d)
	return q
}

// DivMod returns the quotient and remainder corresponding to the provided dividend and divisor
//
// DivMod panics on division by 0. It checks some common/faster cases before fully committing to long division. This can probably be further optimized by
// implementing a successive approximation algorithm, with an initial seed value determined by a 64-bit division of the most significant bits.
func (x Uint128) DivMod(d Uint128) (q, r Uint128) {
	// Handle edge cases and some more common/faster cases
	switch {
	// Case 1: D = 0, divide by zero panic
	case d.Hi == 0 && d.Lo == 0:
		panic("runtime error: integer divide by zero")
	// Case 2: N < D
	case x.Hi < d.Hi || (x.Hi == d.Hi && x.Lo < d.Lo):
		r.Hi, r.Lo = x.Hi, x.Lo
		return q, r
	// Case 3: N >= D (per above), and D is large enough that N / D = 1
	case d.Hi > maxInt64 || (d.Hi == maxUint64 && d.Lo == maxUint64):
		q.Lo = 1
		r = x.Sub(d)
		return q, r
	}
	/*
		// Case 4: N and D have 64 leading zero bits
		case x.hi == 0 && d.hi == 0:
			q.lo = x.lo / d.lo
			r.lo = x.lo % d.lo
			return q, r
		// Case 5: N and D have 64 trailing zero bits
		case x.lo == 0 && d.lo == 0:
			q.lo = x.hi / d.hi
			// The following remainder calculation can probably be optimized further
			dq := d.Mul(q)
			r = x.Sub(dq)
			return q, r
	*/
	n := x.Len() - d.Len()
	if n < 0 {
		r = x
		return q, r
	}
	d = d.LShiftN(n)
	var i uint
	for i = 0; i <= n; i++ {
		q = q.LShift()
		if d.Lte(x) {
			q.Lo |= 1
			x = x.Sub(d)
		}
		d = d.RShift()
	}
	r = x
	return q, r
}

// Eq returns whether x is equal to y
func (x Uint128) Eq(y Uint128) bool {
	return x.Hi == y.Hi && x.Lo == y.Lo
}

// Gt returns whether x is greater than y
func (x Uint128) Gt(y Uint128) bool {
	switch {
	case x.Hi > y.Hi:
		return true
	case x.Hi < y.Hi:
		return false
	case x.Lo > y.Lo:
		return true
	default:
		return false
	}
}

// Gte returns whether x is greater than or equal to y
func (x Uint128) Gte(y Uint128) bool {
	switch {
	case x.Hi > y.Hi:
		return true
	case x.Hi < y.Hi:
		return false
	case x.Lo >= y.Lo:
		return true
	default:
		return false
	}
}

// Inc returns the successor of a Uint128
func (x Uint128) Inc() (z Uint128) {
	z.Lo = x.Lo + 1
	if z.Lo == 0 {
		z.Hi = x.Hi + 1
		return z
	}
	z.Hi = x.Hi
	return z
}

// IsInt64 checks if the Uint128 can be represented as an int64 without overflowing
func (x Uint128) IsInt64() bool {
	return x.Hi == 0 && x.Lo <= maxInt64
}

// IsUint64 checks if the Uint128 can be represented as a uint64 without overflowing
func (x Uint128) IsUint64() bool {
	return x.Hi == 0
}

// Int128 returns a Int128 representation of a Uint128
//
// This function overflows silently
func (x Uint128) Int128() (z Int128) {
	z.Hi = int64(x.Hi)
	z.Lo = x.Lo
	return z
}

// Int64 returns a representation of the Uint128 as the builtin int64
//
// This function overflows silently
func (x Uint128) Int64() int64 {
	return int64(x.Lo)
}

// Len returns the minimum number of bits required to represent x
//
// Edge cases:
//
//	Uint128{0, 0}.Len() -> 0
func (x Uint128) Len() uint {
	if x.Hi == 0 {
		return bits.Len64(x.Lo)
	}
	return bits.Len64(x.Hi) + int64Size
}

// LShift returns a Uint128 left-shifted by 1
func (x Uint128) LShift() (z Uint128) {
	z.Hi = x.Hi<<1 | x.Lo>>(int64Size-1)
	z.Lo = x.Lo << 1
	return z
}

// LShiftN returns a Uint128 left-shifted by a uint (i.e. x << n)
func (x Uint128) LShiftN(n uint) (z Uint128) {
	switch {
	case n >= int128Size:
		return z // z.hi, z.lo = 0, 0
	case n >= int64Size:
		z.Hi = x.Lo << (n - int64Size)
		z.Lo = 0
		return z
	default:
		z.Hi = x.Hi<<n | x.Lo>>(int64Size-n)
		z.Lo = x.Lo << n
		return z
	}
}

// Lt returns whether x is less than y
func (x Uint128) Lt(y Uint128) bool {
	switch {
	case x.Hi < y.Hi:
		return true
	case x.Hi > y.Hi:
		return false
	case x.Lo < y.Lo:
		return true
	default:
		return false
	}
}

// Lte returns whether x is less than or equal to y
func (x Uint128) Lte(y Uint128) bool {
	switch {
	case x.Hi < y.Hi:
		return true
	case x.Hi > y.Hi:
		return false
	case x.Lo <= y.Lo:
		return true
	default:
		return false
	}
}

// Mod returns the remainder corresponding to the provided dividend and divisor
//
// Mod panics on division by 0. It checks some common/faster cases before fully committing to long division. This can probably be further optimized by
// implementing a successive approximation algorithm, with an initial seed value determined by a 64-bit division of the most significant bits.
func (x Uint128) Mod(d Uint128) (r Uint128) {
	_, r = x.DivMod(d)
	return r
}

// Mul returns the product of two Uint128's
func (x Uint128) Mul(y Uint128) (z Uint128) {
	var i uint
	for i = 0; i < int64Size; i++ {
		if y.Lo&(1<<i) != 0 {
			z = z.Add(x.LShiftN(i))
		}
	}
	for i = 0; i < int64Size; i++ {
		if y.Hi&(1<<i) != 0 {
			z = z.Add(x.LShiftN(i + int64Size))
		}
	}
	return z
}

// Nand returns the bitwise NAND of two Uint128's
func (x Uint128) Nand(y Uint128) (z Uint128) {
	z.Hi = ^(x.Hi & y.Hi)
	z.Lo = ^(x.Lo & y.Lo)
	return z
}

// Neg returns the additive inverse of a Uint128
func (x Uint128) Neg() (z Uint128) {
	z.Hi = -x.Hi
	z.Lo = -x.Lo
	if z.Lo > 0 {
		z.Hi--
	}
	return z
}

// Nor returns the bitwise NOR of two Uint128's
func (x Uint128) Nor(y Uint128) (z Uint128) {
	z.Hi = ^(x.Hi | y.Hi)
	z.Lo = ^(x.Lo | y.Lo)
	return z
}

// Not returns the bitwise Not of a Uint128
func (x Uint128) Not() (z Uint128) {
	z.Hi = ^x.Hi
	z.Lo = ^x.Lo
	return z
}

// Or returns the bitwise OR of two Uint128's
func (x Uint128) Or(y Uint128) (z Uint128) {
	z.Hi = x.Hi | y.Hi
	z.Lo = x.Lo | y.Lo
	return z
}

// RShift returns a Uint128 right-shifted by 1
func (x Uint128) RShift() (z Uint128) {
	z.Hi = x.Hi >> 1
	z.Lo = x.Hi<<(int64Size-1) | x.Lo>>1
	return z
}

// RShiftN returns a Uint128 right-shifted by a uint (i.e. x >> n)
func (x Uint128) RShiftN(n uint) (z Uint128) {
	switch {
	case n >= int128Size:
		return z // zhi, zlo = 0, 0
	case n >= int64Size:
		z.Hi = 0
		z.Lo = x.Hi >> (n - int64Size)
		return z
	default:
		z.Hi = x.Hi >> n
		z.Lo = x.Lo>>n | x.Hi<<(int64Size-n)
		return z
	}
}

// RShift128 returns a Uint128 right-shifted by a Uint128 (i.e. x >> y)
func (x Uint128) RShift128(y Uint128) (z Uint128) {
	if y.Hi != 0 || y.Lo >= int128Size {
		return z
	}
	return x.RShiftN(uint(y.Lo))
}

// Sub returns the difference of two Uint128's
func (x Uint128) Sub(y Uint128) (z Uint128) {
	z.Hi = x.Hi - y.Hi
	z.Lo = x.Lo - y.Lo
	if z.Lo > x.Lo {
		z.Hi--
	}
	return z
}

// Uint64 returns a representation of the Uint128 as the builtin uint64
//
// This function overflows silently
func (x Uint128) Uint64() uint64 {
	return x.Lo
}

// Xor returns the bitwise XOR of two Uint128's
func (x Uint128) Xor(y Uint128) (z Uint128) {
	z.Hi = x.Hi ^ y.Hi
	z.Lo = x.Lo ^ y.Lo
	return z
}

// Xnor returns the bitwise XNOR of two Uint128's
func (x Uint128) Xnor(y Uint128) (z Uint128) {
	z.Hi = ^(x.Hi ^ y.Hi)
	z.Lo = ^(x.Lo ^ y.Lo)
	return z
}
