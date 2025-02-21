package wide

import (
	"fmt"
	"math/big"
)

// Int128 is a representation of a signed 128-bit integer
type Int128 struct {
	Hi int64
	Lo uint64
}

// String returns a hexadecimal (string) representation of an Int128
func (x Int128) String() string {
	switch {
	case x.Hi < 0:
		x = x.Neg()
		if x.Hi == 0 {
			return fmt.Sprintf("-%#x", x.Lo) // ignore leading 0's
		}
		return fmt.Sprintf("-%#x%016x", uint64(x.Hi), x.Lo)
	case x.Hi > 0:
		return fmt.Sprintf("%#x%016x", uint64(x.Hi), x.Lo)
	default:
		return fmt.Sprintf("%#x", x.Lo) // ignore leading 0's
	}
}

// NewInt128 returns an Int128 from the high and low 64 bits
func NewInt128(hi int64, lo uint64) Int128 {
	return Int128{Hi: hi, Lo: lo}
}

// Int128FromBigInt returns an Int128 from a big.Int
func Int128FromBigInt(a *big.Int) (z Int128) {
	var y Uint128
	neg := false
	if a.Sign() == -1 {
		a = new(big.Int).Neg(a)
		neg = true
	}
	y.Lo = a.Uint64()
	b := new(big.Int).Rsh(a, int64Size)
	y.Hi = b.Uint64()
	z = y.Int128()
	if neg {
		return z.Neg()
	}
	return z
}

// Int128FromInt64 returns an Int128 from an int64
func Int128FromInt64(x int64) Int128 {
	if x >= 0 {
		return Int128{Hi: 0, Lo: uint64(x)}
	}
	return Int128{Hi: -1, Lo: uint64(x)}
}

// Abs returns the absolute value of an Int128's
func (x Int128) Abs() Int128 {
	if x.Hi < 0 {
		return x.Neg()
	}
	return x
}

// Add returns the sum of two Int128's
func (x Int128) Add(y Int128) (z Int128) {
	z.Hi = x.Hi + y.Hi
	z.Lo = x.Lo + y.Lo
	if z.Lo < x.Lo {
		z.Hi++
	}
	return z
}

// And returns the bitwise AND of two Int128's
func (x Int128) And(y Int128) (z Int128) {
	z.Hi = x.Hi & y.Hi
	z.Lo = x.Lo & y.Lo
	return z
}

// AndNot returns the bitwise AndNot of two Int128's
func (x Int128) AndNot(y Int128) (z Int128) {
	z.Hi = x.Hi &^ y.Hi
	z.Lo = x.Lo &^ y.Lo
	return z
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x Int128) Cmp(y Int128) int {
	switch {
	case x.Hi > y.Hi:
		return 1
	case x.Hi < y.Hi:
		return -1
	case x.Lo > y.Lo:
		return 1
	case x.Lo < y.Lo:
		return -1
	}
	return 0
}

// CmpAbs compares |x| and |y| and returns:
//
//	-1 if |x| <  |y|
//	 0 if |x| == |y|
//	+1 if |x| >  |y|
func (x Int128) CmpAbs(y Int128) int {
	x = x.Abs()
	y = y.Abs()
	return x.Cmp(y)
}

// Dec returns the predecessor of an Int128
func (x Int128) Dec() (z Int128) {
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
func (x Int128) Div(d Int128) (q Int128) {
	q, _ = x.DivMod(d)
	return q
}

// DivMod returns the quotient and remainder corresponding to the provided dividend and divisor
//
// DivMod panics on division by 0. It checks some common/faster cases before fully committing to long division. This can probably be further optimized by
// implementing a successive approximation algorithm, with an initial seed value determined by a 64-bit division of the most significant bits.
func (x Int128) DivMod(d Int128) (q, r Int128) {
	var zero Int128
	qSign, rSign := +1, +1
	if x.Lt(zero) {
		qSign, rSign = -1, -1
		x = x.Neg()
	}
	if d.Lt(zero) {
		qSign = -qSign
		d = d.Neg()
	}
	qAbs, rAbs := x.Uint128().DivMod(d.Uint128())
	q, r = qAbs.Int128(), rAbs.Int128()
	if qSign < 0 {
		q = q.Neg()
	}
	if rSign < 0 {
		r = r.Neg()
	}
	return q, r
}

// Eq returns whether x is equal to y
func (x Int128) Eq(y Int128) bool {
	return x.Hi == y.Hi && x.Lo == y.Lo
}

// Gt returns whether x is greater than y
func (x Int128) Gt(y Int128) bool {
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
func (x Int128) Gte(y Int128) bool {
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

// Inc returns the successor of an Int128
func (x Int128) Inc() (z Int128) {
	z.Lo = x.Lo + 1
	if z.Lo == 0 {
		z.Hi = x.Hi + 1
		return z
	}
	z.Hi = x.Hi
	return z
}

// IsInt64 checks if the Int128 can be represented as an int64
func (x Int128) IsInt64() bool {
	switch x.Hi {
	case 0:
		return x.Lo <= maxInt64
	case -1:
		return x.Lo > maxInt64
	default:
		return false
	}
}

// IsNeg returns whether or not the Int128 is negative
func (x Int128) IsNeg() bool {
	if x.Hi < 0 {
		return true
	}
	return false
}

// IsPos returns whether or not the Int128 is positive
func (x Int128) IsPos() bool {
	switch {
	case x.Hi < 0:
		return false
	case x.Lo > 0:
		return true
	default: // x is zero
		return false
	}
}

// IsUint64 checks if the Int128 can be represented as a uint64 without wrapping
func (x Int128) IsUint64() bool {
	return x.Hi == 0
}

// Int64 returns a representation of the Int128 as the builtin int64
//
// This function overflows silently
func (x Int128) Int64() int64 {
	return int64(x.Lo)
}

// LShift returns an Int128 left-shifted by 1
func (x Int128) LShift() (z Int128) {
	z.Hi = int64(uint64(x.Hi)<<1 | x.Lo>>(int64Size-1))
	z.Lo = x.Lo << 1
	return z
}

// LShiftN returns an Int128 left-shifted by a uint (i.e. x << n)
func (x Int128) LShiftN(n uint) (z Int128) {
	switch {
	case n >= int128Size:
		return z // z.hi, z.lo = 0, 0
	case n >= int64Size:
		z.Hi = int64(x.Lo << (n - int64Size))
		z.Lo = 0
		return z
	default:
		z.Hi = int64(uint64(x.Hi)<<n | x.Lo>>(int64Size-n))
		z.Lo = x.Lo << n
		return z
	}
}

// lShiftNActual returns an Int128 left-shifted by a uint (i.e. x << n)
//
// Unlike LShiftN, it operates on the actual 2's complement representation
func (x Int128) lShiftNActual(n uint) (z Int128) {
	switch {
	case n >= int128Size:
		return z // z.hi, z.lo = 0, 0
	case n >= int64Size:
		z.Hi = int64(x.Lo << (n - int64Size))
		z.Lo = 0
		return z
	default:
		z.Hi = int64(uint64(x.Hi)<<n | x.Lo>>(int64Size-n))
		z.Lo = x.Lo << n
		return z
	}
}

// Lt returns whether x is less than y
func (x Int128) Lt(y Int128) bool {
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
func (x Int128) Lte(y Int128) bool {
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
func (x Int128) Mod(d Int128) (r Int128) {
	_, r = x.DivMod(d)
	return r
}

// Mul returns the product of two Int128's
func (x Int128) Mul(y Int128) (z Int128) {
	var i uint
	yhi := uint64(y.Hi)
	for i = 0; i < int64Size; i++ {
		if y.Lo&(1<<i) != 0 {
			z = z.Add(x.lShiftNActual(i))
		}
	}
	for i = 0; i < int64Size; i++ {
		if yhi&(1<<i) != 0 {
			z = z.Add(x.lShiftNActual(i + int64Size))
		}
	}
	return z
}

// Nand returns the bitwise NAND of two Int128's
func (x Int128) Nand(y Int128) (z Int128) {
	z.Hi = ^(x.Hi & y.Hi)
	z.Lo = ^(x.Lo & y.Lo)
	return z
}

// Neg returns the additive inverse of an Int128
func (x Int128) Neg() (z Int128) {
	z.Hi = -x.Hi
	z.Lo = -x.Lo
	if z.Lo > 0 {
		z.Hi--
	}
	return z
}

// Nor returns the bitwise NOR of two Int128's
func (x Int128) Nor(y Int128) (z Int128) {
	z.Hi = ^(x.Hi | y.Hi)
	z.Lo = ^(x.Lo | y.Lo)
	return z
}

// Not returns the bitwise Not of an Int128
func (x Int128) Not() (z Int128) {
	z.Hi = ^x.Hi
	z.Lo = ^x.Lo
	return z
}

// Or returns the bitwise OR of two Int128's
func (x Int128) Or(y Int128) (z Int128) {
	z.Hi = x.Hi | y.Hi
	z.Lo = x.Lo | y.Lo
	return z
}

// RShift returns an Int128 right-shifted by 1
func (x Int128) RShift() (z Int128) {
	xhi := uint64(x.Hi)
	z.Hi = int64(xhi >> 1)
	z.Lo = x.Lo>>1 | xhi<<(int64Size-1)
	return z
}

// RShiftN returns an Int128 right-shifted by a uint (i.e. x >> n)
//
// Could probably be made faster with sign extension
func (x Int128) RShiftN(n uint) (z Int128) {
	neg := false
	if x.Hi < 0 {
		x = x.Neg()
		neg = true
	}
	switch {
	case n >= int128Size:
		return z // z.hi, z.lo = 0, 0
	case n >= int64Size:
		z.Hi = 0
		z.Lo = uint64(x.Hi) >> (n - int64Size)
	default:
		z.Hi = x.Hi >> n
		z.Lo = x.Lo>>n | uint64(x.Hi)<<(int64Size-n)
	}
	if neg {
		return z.Neg()
	}
	return z
}

// RShift128 returns an Int128 right-shifted by a Uint128 (i.e. x >> y)
func (x Int128) RShift128(y Uint128) (z Int128) {
	if y.Hi != 0 || y.Lo >= int128Size {
		return x.RShiftN(int128Size)
	}
	return x.RShiftN(uint(y.Lo))
}

// Sign returns the sign of an Int128
func (x Int128) Sign() int {
	switch {
	case x.Hi > 0:
		return 1
	case x.Hi < 0:
		return -1
	case x.Lo > 0:
		return 1
	}
	return 0
}

// Sub returns the difference of two Int128's
func (x Int128) Sub(y Int128) (z Int128) {
	z.Hi = x.Hi - y.Hi
	z.Lo = x.Lo - y.Lo
	if z.Lo > x.Lo {
		z.Hi--
	}
	return z
}

// Uint128 returns a Uint128 representation of an Int128
//
// This function overflows silently
func (x Int128) Uint128() (z Uint128) {
	z.Hi = uint64(x.Hi)
	z.Lo = x.Lo
	return z
}

// Uint64 returns a representation of the Int128 as the builtin uint64
//
// This function overflows silently
func (x Int128) Uint64() uint64 {
	return x.Lo
}

// Xor returns the bitwise XOR of two Int128's
func (x Int128) Xor(y Int128) (z Int128) {
	z.Hi = x.Hi ^ y.Hi
	z.Lo = x.Lo ^ y.Lo
	return z
}

// Xnor returns the bitwise XNOR of two Int128's
func (x Int128) Xnor(y Int128) (z Int128) {
	z.Hi = ^(x.Hi ^ y.Hi)
	z.Lo = ^(x.Lo ^ y.Lo)
	return z
}
