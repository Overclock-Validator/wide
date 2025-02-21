package wide

import (
	"testing"
)

func TestStringUint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected string
	}{
		{Uint128{Hi: 0x0, Lo: 0x0}, "0x0"},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, "0xffffffffffffffffffffffffffffffff"},
		{Uint128{Hi: 0xdeadbeef, Lo: 0xbaadf00d}, "0xdeadbeef00000000baadf00d"},
	}
	for _, test := range tests {
		result := test.inp.String()
		if result != test.expected {
			t.Errorf("Expected %+v.String() == %s, got: %s", test.inp, test.expected, result)
		}
	}
}

func TestAddUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 3}},
		{Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 3}},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 2, Lo: 0}, Uint128{Hi: 3, Lo: 0}},
		{Uint128{Hi: 2, Lo: 0}, Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 3, Lo: 0}},
		{Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: 0}, Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: maxUint64, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.Add(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Add(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestAndUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}},
	}
	for _, test := range tests {
		result := test.op1.And(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.And(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestAndNotUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.AndNot(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.AndNot(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestCmpUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected int
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, 0},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, +1},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, -1},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, 0},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}, +1},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 1, Lo: 0}, -1},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 1, Lo: 0}, 0},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: maxUint64}, +1},
		{Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 1, Lo: 0}, -1},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64 - 1}, +1},
		{Uint128{Hi: maxUint64, Lo: maxUint64 - 1}, Uint128{Hi: maxUint64, Lo: maxUint64}, -1},
	}
	for _, test := range tests {
		result := test.op1.Cmp(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.Cmp(%s) == %v, got: %v", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestDivUint128(t *testing.T) {
	tests := []struct {
		expected Uint128
		op2      Uint128
		op1      Uint128
	}{
		{Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 0, Lo: 5}, Uint128{Hi: 0, Lo: 15}},
		{Uint128{Hi: 0, Lo: 5}, Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 0, Lo: 15}},
		{Uint128{Hi: 3, Lo: 0}, Uint128{Hi: 0, Lo: 5}, Uint128{Hi: 15, Lo: 0}},
		{Uint128{Hi: 5, Lo: 0}, Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 15, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1 << 63}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 1 << 63}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0xFFFFFFFFFFFFFFFF}, Uint128{Hi: 0, Lo: 0xFFFFFFFFFFFFFFFF}, Uint128{Hi: 0xFFFFFFFFFFFFFFFE, Lo: 1}},
	}
	for _, test := range tests {
		result := test.op1.Div(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Div(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestDivModUint128(t *testing.T) {
	tests := []struct {
		op1       Uint128
		op2       Uint128
		expected1 Uint128
		expected2 Uint128
	}{
		// Edge cases
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 3, Lo: 0}, Uint128{Hi: 2, Lo: 0}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 65535}, Uint128{Hi: 0, Lo: 0x1000100010001}, Uint128{Hi: 0, Lo: 1}},
		// Randomly generated tests, so that we aren't exclusively looking at edge cases
		{Uint128{Hi: 18078839827447545686, Lo: 14690264446208930433}, Uint128{Hi: 8998388556054393762, Lo: 1806737508591697479}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 82062715338758162, Lo: 11076789429025535475}},
		{Uint128{Hi: 16644964954028498953, Lo: 6044532115116883114}, Uint128{Hi: 201749283954463444, Lo: 14360937540127383586}, Uint128{Hi: 0, Lo: 82}, Uint128{Hi: 101523669762496481, Lo: 9039274542082732486}},
		{Uint128{Hi: 2093133179559761320, Lo: 16602062711200797661}, Uint128{Hi: 16864581198376902610, Lo: 9459801703427125641}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 2093133179559761320, Lo: 16602062711200797661}},
		{Uint128{Hi: 15749230925143681054, Lo: 5141327176156114442}, Uint128{Hi: 7370181155092257276, Lo: 6109534453883242647}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 1008868614959166501, Lo: 11369002342099180764}},
		{Uint128{Hi: 9688725833947941869, Lo: 5829725676053634587}, Uint128{Hi: 10771428793017394942, Lo: 11481679973018363035}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 9688725833947941869, Lo: 5829725676053634587}},
		{Uint128{Hi: 13115430490891893466, Lo: 8603977268700753405}, Uint128{Hi: 9528846791087986823, Lo: 10534857072914627768}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 3586583699803906642, Lo: 16515864269495677253}},
		{Uint128{Hi: 4464682195163803491, Lo: 7489173075960054176}, Uint128{Hi: 16981508686037311992, Lo: 10249255444110602208}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 4464682195163803491, Lo: 7489173075960054176}},
		{Uint128{Hi: 14969323870889724866, Lo: 2190849366522262425}, Uint128{Hi: 15744153244977445381, Lo: 6216957372856731963}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 14969323870889724866, Lo: 2190849366522262425}},
		{Uint128{Hi: 4242347467606309934, Lo: 617379611895620583}, Uint128{Hi: 10940397127107242540, Lo: 2302117645027716276}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 4242347467606309934, Lo: 617379611895620583}},
		{Uint128{Hi: 13878862832891726409, Lo: 12247025815474997093}, Uint128{Hi: 15933048822954309307, Lo: 10203778256787916680}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 13878862832891726409, Lo: 12247025815474997093}},
		{Uint128{Hi: 838346584253750531, Lo: 8728622130014530564}, Uint128{Hi: 7623255793094361036, Lo: 8997606623242799596}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 838346584253750531, Lo: 8728622130014530564}},
		{Uint128{Hi: 16329672898340907698, Lo: 1812476799316563631}, Uint128{Hi: 15241696152186921136, Lo: 805856988913712022}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1087976746153986562, Lo: 1006619810402851609}},
		{Uint128{Hi: 6779770950567198933, Lo: 5234941809987140509}, Uint128{Hi: 12958508852484540729, Lo: 6724476107841768838}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 6779770950567198933, Lo: 5234941809987140509}},
		{Uint128{Hi: 14498257309946160406, Lo: 18084320459228438784}, Uint128{Hi: 13094649420129182662, Lo: 13044860498710942198}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1403607889816977744, Lo: 5039459960517496586}},
		{Uint128{Hi: 15511746663743037887, Lo: 9414436614834301257}, Uint128{Hi: 241543290869542064, Lo: 4264400748188406476}, Uint128{Hi: 0, Lo: 64}, Uint128{Hi: 52976048092345776, Lo: 13193949836419561033}},
		{Uint128{Hi: 8084744342537816804, Lo: 1423207416017422773}, Uint128{Hi: 15656244447055551121, Lo: 11594742766098937249}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 8084744342537816804, Lo: 1423207416017422773}},
		{Uint128{Hi: 4319800846968343225, Lo: 14534725286450089249}, Uint128{Hi: 13287988548133458958, Lo: 5713573001805029812}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 4319800846968343225, Lo: 14534725286450089249}},
		{Uint128{Hi: 13508470227635185554, Lo: 12357508788529651644}, Uint128{Hi: 10774097062273294804, Lo: 7568828230680467196}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 2734373165361890750, Lo: 4788680557849184448}},
		{Uint128{Hi: 14855998319421094055, Lo: 4486542090348589603}, Uint128{Hi: 6678463923086748523, Lo: 1254263323666630076}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 1499070473247597009, Lo: 1978015443015329451}},
		{Uint128{Hi: 10551012307773401926, Lo: 7564802811090251473}, Uint128{Hi: 6463539642689907276, Lo: 18401146206923932121}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 4087472665083494649, Lo: 7610400677875870968}},
	}
	for _, test := range tests {
		result1, result2 := test.op1.DivMod(test.op2)
		if result1.Lo != test.expected1.Lo || result1.Hi != test.expected1.Hi || result2.Lo != test.expected2.Lo || result2.Hi != test.expected2.Hi {
			t.Errorf("Expected %s.DivMod(%s) == %s, %s got: %s, %s", test.op1, test.op2, test.expected1, test.expected2, result1, result2)
		}
	}
}

func TestDivByZeroUint128(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Division by 0 did not panic")
		}
	}()
	Uint128{Hi: 1, Lo: 1}.DivMod(Uint128{Hi: 0, Lo: 0})
}

func TestEqUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 1}, Uint128{Hi: 1, Lo: 1}, true},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1}, false},
	}
	for _, test := range tests {
		result := test.op1.Eq(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.Eq(%s) == %v, got: %v", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestGtUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 1}, Uint128{Hi: 1, Lo: 1}, false},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: maxUint64}, true},
	}
	for _, test := range tests {
		result := test.op1.Gt(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.Gt(%s) == %v, got: %v", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestGteUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 1}, Uint128{Hi: 1, Lo: 1}, true},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: maxUint64}, true},
	}
	for _, test := range tests {
		result := test.op1.Gte(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.Gte(%s) == %v, got: %v", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestIsInt64Uint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: maxInt64}, true},
		{Uint128{Hi: 0, Lo: maxInt64 + 1}, false},
		{Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, false},
	}
	for _, test := range tests {
		result := test.inp.IsInt64()
		if test.expected != result {
			t.Errorf("Expected %s.IsInt64() == %v, got: %v", test.inp, test.expected, result)
		}
	}
}

func TestIsUint64Uint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: maxUint64}, true},
		{Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: maxUint64, Lo: 0}, false},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, false},
	}
	for _, test := range tests {
		result := test.inp.IsUint64()
		if test.expected != result {
			t.Errorf("Expected %s.IsUint64() == %v, got: %v", test.inp, test.expected, result)
		}
	}
}

func TestInt64Uint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected int64
	}{
		{Uint128{Hi: 0, Lo: 0}, 0},
		{Uint128{Hi: 0, Lo: maxInt64}, maxInt64},
		{Uint128{Hi: 0, Lo: maxInt64 + 1}, minInt64},
		{Uint128{Hi: 0, Lo: maxUint64}, -1},
		{Uint128{Hi: 1, Lo: 0}, 0},
		{Uint128{Hi: maxUint64, Lo: 0}, 0},
		{Uint128{Hi: maxUint64, Lo: maxInt64}, maxInt64},
	}
	for _, test := range tests {
		result := test.inp.Int64()
		if test.expected != result {
			t.Errorf("Expected %s.Int64() == %v, got: %v", test.inp, test.expected, result)
		}
	}
}

func TestLenUint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected uint
	}{
		{Uint128{Hi: 0, Lo: 0}, 0},
		{Uint128{Hi: 0, Lo: 1}, 1},
		{Uint128{Hi: 0, Lo: 2}, 2},
		{Uint128{Hi: 0, Lo: 3}, 2},
		{Uint128{Hi: 0, Lo: 4}, 3},
		{Uint128{Hi: 0, Lo: maxUint64}, 64},
		{Uint128{Hi: 1, Lo: 0}, 65},
		{Uint128{Hi: 2, Lo: 0}, 66},
		{Uint128{Hi: 3, Lo: 0}, 66},
		{Uint128{Hi: 4, Lo: 0}, 67},
		{Uint128{Hi: maxUint64, Lo: 0}, 128},
	}
	for _, test := range tests {
		result := test.inp.Len()
		if test.expected != result {
			t.Errorf("Expected %s.Len() == %v, got: %v", test.inp, test.expected, result)
		}
	}
}

func TestLShiftUint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1 << 1}},
		{Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 1, Lo: maxUint64 - 1}},
		{Uint128{Hi: maxUint64, Lo: 0}, Uint128{Hi: maxUint64 - 1, Lo: 0}},
	}
	for _, test := range tests {
		result := test.inp.LShift()
		if test.expected != result {
			t.Errorf("Expected %s.LShift() == %s, got: %s", test.inp, test.expected, result)
		}
	}
}

func TestLShiftNUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      uint
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, 0, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, 1, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, 2, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, 0, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 0, Lo: 1}, 1, Uint128{Hi: 0, Lo: 2}},
		{Uint128{Hi: 0, Lo: 1}, 2, Uint128{Hi: 0, Lo: 4}},
		{Uint128{Hi: 0, Lo: 1}, 63, Uint128{Hi: 0, Lo: 1 << 63}},
		{Uint128{Hi: 0, Lo: 1}, 64, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, 127, Uint128{Hi: 1 << 63, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, 0, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, 1, Uint128{Hi: 2, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, 2, Uint128{Hi: 4, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, 63, Uint128{Hi: 1 << 63, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, 64, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.LShiftN(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.LShiftN(%v) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestLtUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 1, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 1}, Uint128{Hi: 1, Lo: 1}, false},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: maxUint64}, false},
		{Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 1, Lo: 0}, true},
	}
	for _, test := range tests {
		result := test.op1.Lt(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.Lt(%s) == %v, got: %v", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestLteUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected bool
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 1}, Uint128{Hi: 1, Lo: 1}, true},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 1}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}, false},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}, true},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1}, false},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: maxUint64}, false},
		{Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 1, Lo: 0}, true},
	}
	for _, test := range tests {
		result := test.op1.Lte(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.Lte(%s) == %v, got: %v", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestMulUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 0, Lo: 5}, Uint128{Hi: 0, Lo: 15}},
		{Uint128{Hi: 0, Lo: 5}, Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 0, Lo: 15}},
		{Uint128{Hi: 3, Lo: 0}, Uint128{Hi: 0, Lo: 5}, Uint128{Hi: 15, Lo: 0}},
		{Uint128{Hi: 5, Lo: 0}, Uint128{Hi: 0, Lo: 3}, Uint128{Hi: 15, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1<<32 + 1<<31}, Uint128{Hi: 0, Lo: 1<<32 + 1<<31}, Uint128{Hi: 2, Lo: 1 << 62}},
		{Uint128{Hi: 0, Lo: 1 << 63}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 1 << 63}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0xFFFFFFFFFFFFFFFF}, Uint128{Hi: 0, Lo: 0xFFFFFFFFFFFFFFFF}, Uint128{Hi: 0xFFFFFFFFFFFFFFFE, Lo: 1}},
	}
	for _, test := range tests {
		result := test.op1.Mul(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Mul(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestNandUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.Nand(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Nand(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestNegUint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 1 << 63, Lo: 1}, Uint128{Hi: 1<<63 - 1, Lo: 1<<64 - 1}},
		{Uint128{Hi: 1<<63 - 1, Lo: 1<<64 - 1}, Uint128{Hi: 1 << 63, Lo: 1}},
		{Uint128{Hi: 1 << 63, Lo: 0}, Uint128{Hi: 1 << 63, Lo: 0}}, // most negative number has no positive counterpart
	}
	for _, test := range tests {
		result := test.inp.Neg()
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Neg() == %s, got: %s", test.inp, test.expected, result)
		}
	}
}

func TestNorUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.Nor(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Nor(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestNotUint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.inp.Not()
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Not() == %s, got: %s", test.inp, test.expected, result)
		}
	}
}

func TestOrUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}},
	}
	for _, test := range tests {
		result := test.op1.Or(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Or(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestRShiftUint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 0, Lo: 4}, Uint128{Hi: 0, Lo: 2}},
		{Uint128{Hi: 1, Lo: maxUint64 - 1}, Uint128{Hi: 0, Lo: maxUint64}},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 1 << 63}},
	}
	for _, test := range tests {
		result := test.inp.RShift()
		if test.expected != result {
			t.Errorf("Expected %s.RShift() == %s, got: %s", test.inp, test.expected, result)
		}
	}
}

func TestRShiftNUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      uint
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, 0, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, 1, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, 0, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 0, Lo: 1}, 1, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 2}, 1, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 0, Lo: 4}, 2, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 0, Lo: 1 << 63}, 63, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 1, Lo: 0}, 1, Uint128{Hi: 0, Lo: 1 << 63}},
		{Uint128{Hi: 2, Lo: 0}, 1, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 4, Lo: 0}, 1, Uint128{Hi: 2, Lo: 0}},
		{Uint128{Hi: 4, Lo: 0}, 2, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, 64, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: 1 << 63, Lo: 0}, 127, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, 127, Uint128{Hi: 0, Lo: 1}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, 128, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.RShiftN(test.op2)
		if test.expected != result {
			t.Errorf("Expected %s.RShiftN(%v) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestSubUint128(t *testing.T) {
	tests := []struct {
		expected Uint128
		op2      Uint128
		op1      Uint128
	}{
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 3}},
		{Uint128{Hi: 0, Lo: 2}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 3}},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 2, Lo: 0}, Uint128{Hi: 3, Lo: 0}},
		{Uint128{Hi: 2, Lo: 0}, Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 3, Lo: 0}},
		{Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: maxUint64}, Uint128{Hi: 1, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: 0}, Uint128{Hi: 1, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 1, Lo: 0}, Uint128{Hi: maxUint64, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 1}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 1}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.Sub(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Sub(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestUint64Uint128(t *testing.T) {
	tests := []struct {
		inp      Uint128
		expected uint64
	}{
		{Uint128{Hi: 0, Lo: 0}, 0},
		{Uint128{Hi: 0, Lo: maxInt64}, maxInt64},
		{Uint128{Hi: 0, Lo: maxInt64 + 1}, maxInt64 + 1},
		{Uint128{Hi: 0, Lo: maxUint64}, maxUint64},
		{Uint128{Hi: 1, Lo: 0}, 0},
		{Uint128{Hi: 1, Lo: maxUint64}, maxUint64},
		{Uint128{Hi: maxUint64, Lo: 0}, 0},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, maxUint64},
	}
	for _, test := range tests {
		result := test.inp.Uint64()
		if test.expected != result {
			t.Errorf("Expected %s.Uint64() == %v, got: %v", test.inp, test.expected, result)
		}
	}
}

func TestXorUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
	}
	for _, test := range tests {
		result := test.op1.Xor(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Xor(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}

func TestXnorUint128(t *testing.T) {
	tests := []struct {
		op1      Uint128
		op2      Uint128
		expected Uint128
	}{
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}},
		{Uint128{Hi: 0, Lo: 0}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: 0, Lo: 0}, Uint128{Hi: 0, Lo: 0}},
		{Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}, Uint128{Hi: maxUint64, Lo: maxUint64}},
	}
	for _, test := range tests {
		result := test.op1.Xnor(test.op2)
		if result.Lo != test.expected.Lo || result.Hi != test.expected.Hi {
			t.Errorf("Expected %s.Xnor(%s) == %s, got: %s", test.op1, test.op2, test.expected, result)
		}
	}
}
