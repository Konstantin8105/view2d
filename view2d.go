package view2d

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/Konstantin8105/gog"
	"github.com/Konstantin8105/pow"
)

type Ray struct {
	Line
}

func (r *Ray) Scale(scale float64) {
	r.P2.X = r.P1.X + (r.P2.X-r.P1.X)*scale
	r.P2.Y = r.P1.Y + (r.P2.Y-r.P1.Y)*scale
}

// rand value from 0...1
func (r *Ray) Rotate(randomValue float64) {
	if randomValue < 0 || 1 < randomValue {
		panic(fmt.Errorf("not valid random value: %.5f", randomValue))
	}
	// rotate at random angle from -pi/2 ... +pi/2
	angle := math.Asin(1 - 2*randomValue)
	r.P2 = gog.Rotate(r.P1.X, r.P1.Y, angle, r.P2)
}

///////////////////////////////////////////////////////////////////////////////

type Curve interface {
	GetVector(rand float64) (r Ray)
	Box() (begin, finish gog.Point)
}

///////////////////////////////////////////////////////////////////////////////

var _ Curve = new(Line)

type Line struct {
	P1, P2 gog.Point
}

func (l Line) GetVector(rand float64) (r Ray) {
	if rand < 0 || 1 < rand {
		panic(fmt.Errorf("not valid random value: %.5f", rand))
	}
	r.P1.X = l.P1.X + (l.P2.X-l.P1.X)*rand
	r.P1.Y = l.P1.Y + (l.P2.Y-l.P1.Y)*rand
	r.P2.X = r.P1.X + (l.P2.X - l.P1.X)
	r.P2.Y = r.P1.Y + (l.P2.Y - l.P1.Y)
	// rotate at 90 degree
	xc, yc := r.P1.X, r.P1.Y
	res := gog.Rotate(xc, yc, math.Pi/2.0, r.P2)
	r.P2 = res
	// change vector size to 1.0
	dist := gog.Distance(r.P1, r.P2)
	r.P2.X = r.P1.X + (r.P2.X-r.P1.X)*dist/1.0
	r.P2.Y = r.P1.Y + (r.P2.Y-r.P1.Y)*dist/1.0
	return
}

func (l Line) Box() (begin, finish gog.Point) {
	begin.X = min(l.P1.X, l.P2.X)
	begin.Y = min(l.P1.Y, l.P2.Y)
	finish.X = max(l.P1.X, l.P2.X)
	finish.Y = max(l.P1.Y, l.P2.Y)
	return
}

///////////////////////////////////////////////////////////////////////////////

var _ Curve = new(Arc)

type Arc struct {
	P1, P2, P3 gog.Point
}

func (a Arc) GetVector(rand float64) (r Ray) {
	// find random point on arc
	isClock := gog.Orientation(a.P1, a.P2, a.P3) == gog.ClockwisePoints
	xc, yc, radius := gog.Arc(a.P1, a.P2, a.P3)
	a1 := math.Atan2(a.P1.Y-yc, a.P1.X-xc) // begin angle
	a3 := math.Atan2(a.P3.Y-yc, a.P3.X-xc) // end angle
	if isClock {
		a1, a3 = a3, a1 // angle by clock
	}
	var fullAngle float64
	if a1 < a3 {
		fullAngle = a3 - a1
	} else {
		fullAngle = 2*math.Pi - (a1 - a3)
	}
	angle := fullAngle * rand // random angle
	// create vector
	r.P1.X = xc + radius*math.Cos(a1+angle)
	r.P1.Y = yc + radius*math.Sin(a1+angle)
	r.P2.X, r.P2.Y = xc, yc // end of ray at the center
	// change vector size to 1.0
	dist := gog.Distance(r.P1, r.P2)
	if !isClock {
		dist = -dist
	}
	r.P2.X = r.P1.X + (r.P2.X-r.P1.X)*dist/1.0
	r.P2.Y = r.P1.Y + (r.P2.Y-r.P1.Y)*dist/1.0
	return
}

func (a Arc) Box() (begin, finish gog.Point) {
	xc, yc, r := gog.Arc(a.P1, a.P2, a.P3)
	begin.X = min(xc-r, xc+r)
	begin.Y = min(yc-r, yc+r)
	finish.X = max(xc-r, xc+r)
	finish.Y = max(yc-r, yc+r)
	return
}

///////////////////////////////////////////////////////////////////////////////

var _ Curve = new(Circle)

type Circle struct {
	Center        gog.Point
	Radius        float64
	VectorOutside bool
}

func (c Circle) GetVector(rand float64) (r Ray) {
	// generate angle
	angle := 2 * math.Pi * rand
	// create vector
	xc, yc, radius := c.Center.X, c.Center.Y, c.Radius
	r.P1.X = xc + radius*math.Cos(angle)
	r.P1.Y = yc + radius*math.Sin(angle)
	r.P2.X, r.P2.Y = xc, yc // end of ray at the center
	// change vector size to 1.0
	dist := gog.Distance(r.P1, r.P2)
	if c.VectorOutside {
		dist = -dist
	}
	r.P2.X = r.P1.X + (r.P2.X-r.P1.X)*dist/1.0
	r.P2.Y = r.P1.Y + (r.P2.Y-r.P1.Y)*dist/1.0
	return
}

func (c Circle) Box() (begin, finish gog.Point) {
	xc, yc, r := c.Center.X, c.Center.Y, c.Radius
	begin.X = min(xc-r, xc+r)
	begin.Y = min(yc-r, yc+r)
	finish.X = max(xc-r, xc+r)
	finish.Y = max(yc-r, yc+r)
	return
}

///////////////////////////////////////////////////////////////////////////////

var Amount int64 = 100000

var (
	debug    bool
	intersec []Line
	miss     []Line
)

func intersect(c Curve, v Ray) (pi []gog.Point) {
	switch c := c.(type) {
	case Line:
		pi, _, _ = gog.LineLine(
			v.P1, v.P2,
			c.P1, c.P2,
		)
	case Arc:
		pi, _, _ = gog.LineArc(
			v.P1, v.P2,
			c.P1, c.P2, c.P3,
		)
	case Circle:
		arc1 := Arc{
			P1: gog.Point{X: c.Center.X + c.Radius, Y: c.Center.Y},
			P2: gog.Point{X: c.Center.X, Y: c.Center.Y + c.Radius},
			P3: gog.Point{X: c.Center.X - c.Radius, Y: c.Center.Y},
		}
		arc2 := Arc{
			P1: gog.Point{X: c.Center.X - c.Radius, Y: c.Center.Y},
			P2: gog.Point{X: c.Center.X, Y: c.Center.Y - c.Radius},
			P3: gog.Point{X: c.Center.X + c.Radius, Y: c.Center.Y},
		}
		if c.VectorOutside {
			arc1.P1, arc1.P3 = arc1.P3, arc1.P1
			arc2.P1, arc2.P3 = arc2.P3, arc2.P1
		}
		pi = append(intersect(arc1, v), intersect(arc2, v)...)
	default:
		panic(fmt.Errorf("not implemented: %#v", c))
	}
	return
}

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
		// for garantee large of all point system, them multiply
		// by coefficient more 1.0
		scale *= 1.1
	}
	// calculation
	var mut sync.Mutex
	run := func(cpu int, steps int64) {
		for iter := int64(0); iter < steps; iter++ {
			v := present.GetVector(rand.Float64())
			// scale vector
			v.Scale(scale)
			// rotate vector
			v.Rotate(rand.Float64())
			// find intersection between vector and curves
			// present curve to all
			// and store minimal distance
			index := -1
			var pint gog.Point // interssection point
			distance := math.MaxFloat64
			for i := range curves {
				pis := intersect(curves[i], v)
				for p := range pis {
					if d := gog.Distance(v.P1, pis[p]); 1e-6 < math.Abs(d) && d < distance {
						pint = pis[p]
						distance = d
						index = i
					}
				}
			}
			if index < 0 {
				if debug {
					miss = append(miss, Line{v.P1, v.P2})
				}
				continue
			} else if debug {
				intersec = append(intersec, Line{v.P1, pint})
			}
			// get first intersection
			mut.Lock()
			if 0 <= index {
				vf[index]++
			}
			mut.Unlock()
		}
	}

	// parallel calculation
	cpus := runtime.NumCPU()
	if int64(cpus)*10 < Amount {
		var wg sync.WaitGroup
		wg.Add(cpus)
		amount := Amount
		dstep := Amount / int64(cpus)
		for i := 0; i < cpus; i++ {
			if i != cpus-1 {
				amount -= dstep
				go func(cpu int) {
					run(cpu, dstep)
					wg.Done()
				}(i)
			} else {
				go func(cpu int) {
					run(cpu, amount)
					wg.Done()
				}(i)
			}
		}
		wg.Wait()
	} else {
		run(0, Amount)
	}

	viewFactor = make([]float64, len(curves))
	for i := range curves {
		viewFactor[i] = float64(vf[i]) / float64(Amount)
	}
	return
}
