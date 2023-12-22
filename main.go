package main

import (
	"fmt"
	"math/big"
	"math/bits"
)

type Point struct {
	X, Y *big.Int
}

func (p Point) IsEqual(p1 Point) bool {
	return p.X.Cmp(p1.X) == 0 && p.Y.Cmp(p1.Y) == 0
}

func (p Point) Minus(mod *big.Int) Point {
	nP := p.Copy()

	nP.Y.Mul(nP.Y, minusOne)
	nP.Y.Mod(nP.Y, mod)

	return nP
}

func (p Point) Copy() Point {
	return Point{X: new(big.Int).Set(p.X), Y: new(big.Int).Set(p.Y)}
}

func (p Point) IsInf() bool {
	return p.X.Cmp(minusOne) == 0 && p.Y.Cmp(minusOne) == 0
}

func (p Point) String() string {
	if p.IsInf() {
		return "O Inf"
	}

	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

type Curve struct {
	A, B, Mod *big.Int
}

func main() {
	//mod := big.NewInt(31)
	mod := big.NewInt(631)
	//a := big.NewInt(0)
	a := big.NewInt(30)
	//b := big.NewInt(11)
	b := big.NewInt(34)

	c := Curve{
		A:   a,
		B:   b,
		Mod: mod,
	}

	//P, Q := Point{X: big.NewInt(2), Y: big.NewInt(9)}, Point{X: big.NewInt(3), Y: big.NewInt(10)}
	P, Q := Point{X: big.NewInt(36), Y: big.NewInt(571)}, Point{X: big.NewInt(420), Y: big.NewInt(48)}
	R := Point{X: big.NewInt(0), Y: big.NewInt(36)}

	minusR := R.Minus(c.Mod)
	PMinusR := Add(P, minusR, c)

	fmt.Println(Wail(5, P, Add(Q, R, c), c))
	fmt.Println(Wail(5, P, R, c))
	fmt.Println(Wail(5, Q, minusR, c))
	fmt.Println(Wail(5, Q, PMinusR, c))
}

var two = big.NewInt(2)
var three = big.NewInt(3)
var minusOne = big.NewInt(-1)

func Add(p1, p2 Point, c Curve) Point {
	if p1.IsInf() {
		return p2.Copy()
	}

	if p1.X.Cmp(p2.X) == 0 {
		if p1.Y.Cmp(p2.Y) == 0 {
			return double(p1, c.A, c.Mod)
		}

		return inf()
	}

	return add(p1, p2, c.Mod)
}

func add(p1, p2 Point, mod *big.Int) Point {
	p := Point{X: new(big.Int), Y: new(big.Int)}

	p.X = otherAngleCoef(p1, p2, mod)
	p.X.Exp(p.X, two, mod)
	p.X.Sub(p.X, p1.X)
	p.X.Mod(p.X, mod)
	p.X.Sub(p.X, p2.X)
	p.X.Mod(p.X, mod)

	p.Y = otherAngleCoef(p1, p2, mod)
	p.Y.Mod(p.Y, mod)

	tmpY := new(big.Int).Sub(p.X, p1.X)
	tmpY = tmpY.Mod(tmpY, mod)
	p.Y.Mul(p.Y, tmpY)
	p.Y.Mul(p.Y, minusOne)
	p.Y.Sub(p.Y, p1.Y)

	p.Y.Mod(p.Y, mod)

	return p
}

func double(p1 Point, a, mod *big.Int) Point {
	p := Point{X: new(big.Int), Y: new(big.Int)}

	tmp := selfAngleCoef(p1, a, mod)

	p.X.Exp(tmp, two, mod)

	p.X.Sub(p.X, p1.X)
	p.X.Sub(p.X, p1.X)
	p.X.Mod(p.X, mod)

	tmpY := new(big.Int).Sub(p.X, p1.X)
	p.Y.Mul(tmp, tmpY)
	p.Y.Mod(p.Y, mod)
	p.Y.Mul(p.Y, minusOne)

	p.Y.Sub(p.Y, p1.Y)
	p.Y.Mod(p.Y, mod)

	return p
}

func inf() Point {
	return Point{X: big.NewInt(-1), Y: big.NewInt(-1)}
}

func AngleCoef(p1, p2 Point, c Curve) {
	if p1.IsEqual(p2) {

	}

}

func selfAngleCoef(p1 Point, a, mod *big.Int) *big.Int {
	tmpXMul, tmpXDiv := new(big.Int), new(big.Int)

	tmpXMul.Exp(p1.X, two, mod)
	tmpXMul.Mul(tmpXMul, three)
	tmpXMul.Add(tmpXMul, a)

	tmpXDiv.Mul(p1.Y, two)
	tmpXDiv.ModInverse(tmpXDiv, mod)

	coef := new(big.Int).Mul(tmpXMul, tmpXDiv)
	coef.Mod(coef, mod)

	return coef
}

func otherAngleCoef(p1, p2 Point, mod *big.Int) *big.Int {
	tmpXMul, tmpXDiv := new(big.Int), new(big.Int)
	tmpXMul.Sub(p1.Y, p2.Y)
	tmpXMul.Mod(tmpXMul, mod)

	tmpXDiv.Sub(p1.X, p2.X)
	tmpXDiv.Mod(tmpXDiv, mod)

	tmpXDiv.ModInverse(tmpXDiv, mod)

	coef := new(big.Int).Mul(tmpXMul, tmpXDiv)
	coef.Mod(coef, mod)

	return coef
}

func Wail(n uint, basePoint, calcPoint Point, c Curve) *big.Int {
	mask := uint(1)

	f := big.NewInt(1)
	z := basePoint.Copy()

	for i := bits.Len(n) - 2; i >= 0; i-- {
		bitVal := (n & (mask << i)) >> i

		f.Mul(f, f)
		zzCoef := selfAngleCoef(z, c.A, c.Mod)

		lv := LineValueAtPoint(z, zzCoef, c.Mod)(calcPoint)
		z2 := Add(z, z, c)
		vv := VerticalLineValueAtPoint(z2, c.Mod)(calcPoint)

		f.Mul(f, lv)
		f.Mul(f, vv.ModInverse(vv, c.Mod))
		f.Mod(f, c.Mod)
		z = z2

		if bitVal == 1 {
			f.Mul(f, f)
			pzCoef := otherAngleCoef(basePoint, z, c.Mod)
			lv := LineValueAtPoint(z, pzCoef, c.Mod)(calcPoint)
			pz := Add(basePoint, z, c)
			vv := VerticalLineValueAtPoint(pz, c.Mod)(calcPoint)

			f.Mul(f, lv)
			f.Mul(f, vv.ModInverse(vv, c.Mod))
			f.Mod(f, c.Mod)
			z = pz
		}
	}

	return f
}

func LineValueAtPoint(p1 Point, coef, mod *big.Int) func(p Point) *big.Int {
	return func(p Point) *big.Int {
		val := new(big.Int)

		val.Mul(coef, p1.X)

		val.Sub(val, p1.Y)
		val.Mod(val, mod)

		val.Add(val, p.Y)
		val.Sub(val, new(big.Int).Mul(coef, p.X))
		val.Mod(val, mod)

		return val
	}
}

func VerticalLineValueAtPoint(p1 Point, mod *big.Int) func(p Point) *big.Int {
	return func(p Point) *big.Int {
		val := new(big.Int).Sub(p.X, p1.X)
		val.Mod(val, mod)

		return val
	}
}
