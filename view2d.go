package view2d

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Konstantin8105/gog"
	"github.com/Konstantin8105/pow"
)

type Vector struct {
	start, finish gog.Point
}

func (v *Vector) Scale(scale float64) {
	v.finish.X = v.start.X + (v.finish.X-v.start.X)*scale
	v.finish.Y = v.start.Y + (v.finish.Y-v.start.Y)*scale
}

func (v *Vector) Rotate(rand float64) {
	if rand < 0 || 1 < rand {
		panic(fmt.Errorf("not valid random value: %.5f", rand))
	}
	// TODO rotate at random angle from -90 ... +90 degree
}

type Curve interface {
	GetVector(rand float64) (v Vector)
	Box() (begin, finish gog.Point)
}

var _ Curve = new(Line)

type Line struct {
	p1, p2 gog.Point
}

func (l Line) GetVector(rand float64) (v Vector) {
	if rand < 0 || 1 < rand {
		panic(fmt.Errorf("not valid random value: %.5f", rand))
	}
	v.start.X = l.p1.X + (l.p2.X-l.p1.X)*rand
	v.start.Y = l.p1.Y + (l.p2.Y-l.p1.Y)*rand

	// TODO v.finish perpendicular vector 1 size
	return
}

func (l Line) Box() (begin, finish gog.Point) {
	begin.X = min(l.p1.X, l.p2.X)
	begin.Y = min(l.p1.Y, l.p2.Y)
	finish.X = max(l.p1.X, l.p2.X)
	finish.Y = max(l.p1.Y, l.p2.Y)
	return
}

func Calc(curves []Curve) {
	for i := range curves {
		fmt.Println(OneCurve(curves[i], curves))
	}
	return
}

var Amount int64 = 100

func OneCurve(present Curve, curves []Curve) (viewFactor []float64) {
	vf := make([]int64, len(curves))
	// calculate scale of geometry
	scale := 0.0
	{
		var begin, finish gog.Point
		for i := range curves {
			b, f := curves[i].Box()
			begin.X = min(begin.X, b.X, f.X)
			begin.Y = min(begin.Y, b.Y, f.Y)
			finish.X = max(finish.X, b.X, f.X)
			finish.Y = max(finish.Y, b.Y, f.Y)
		}
		scale = math.Sqrt(pow.E2(begin.X-finish.X) + pow.E2(begin.Y-finish.Y))
	}
	for iter := int64(0); iter <= Amount; iter++ {
		v := present.GetVector(rand.Float64())
		// scale vector
		v.Scale(scale)
		// rotate vector
		v.Rotate(rand.Float64())
		// TODO find intersection between vector and curves

		// TODO get first intersection

		// TODO add to intersection view factor
	}
	for i := range viewFactor {
		viewFactor[i] = float64(vf[i]) / float64(Amount)
	}
	return
}
