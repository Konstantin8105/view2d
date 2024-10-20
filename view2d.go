package view2d

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Konstantin8105/gog"
)

type Vector struct {
	start, finish gog.Point
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
	v.start.X = p1.X + (p1.X-p2.X)*rand
	v.start.Y = p1.Y + (p1.Y-p2.Y)*rand

	return
}

func (l Line) Box() (begin, finish gog.Point) {
	begin.X = min(p1.X, p2.X)
	begin.Y = min(p1.Y, p2.Y)
	finish.X = max(p1.X, p2.X)
	finish.Y = max(p1.Y, p2.Y)
	return
}

func Calc(curves []Curve) {
	for i := range curves {
		fmt.Println(OneCurve(curves[i], curves))
	}
	return
}

var Amount = 1000

func OneCurve(index int, curves []Curve) (viewFactor []float64) {
	vf := make([]curves, len(curves))
	present := curves[index]
	// calculate scale of geometry
	scale := 0.0
	{
		var begin, finish gog.Point
		for i := range curves {
			b, f := curves.Box()
			begin.X = min(begin.X, b.X, f.X)
			begin.Y = min(begin.Y, b.Y, f.Y)
			finish.X = max(finish.X, b.X, f.X)
			finish.Y = max(finish.Y, b.Y, f.Y)
		}
		scale = math.Sqrt(pow.E2(begin.X-finish.X) + pow.E2(begin.Y-finish.Y))
	}
	for iter := 0; iter <= Amount; iter++ {
		v := present.GetVector(rand.Float64())
		// scale vector
		v.Scale(scale)
		// find intersection between vector and curves

		// get first intersection

		// add to intersection view factor
	}
	for i := range viewFactor {
		viewFactor[i] = float64(vf[i]) / float64(Amount)
	}
	return
}
