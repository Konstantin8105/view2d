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
	start, finish gog.Point
}

func (r *Ray) Scale(scale float64) {
	r.finish.X = r.start.X + (r.finish.X-r.start.X)*scale
	r.finish.Y = r.start.Y + (r.finish.Y-r.start.Y)*scale
}

// rand value from 0...1
func (r *Ray) Rotate(randomValue float64) {
	if randomValue < 0 || 1 < randomValue {
		panic(fmt.Errorf("not valid random value: %.5f", randomValue))
	}
	// rotate at random angle from -pi/2 ... +pi/2
	angle := math.Asin(1 - 2*randomValue)
	r.finish = gog.Rotate(r.start.X, r.start.Y, angle, r.finish)
}

type Curve interface {
	GetVector(rand float64) (r Ray)
	Box() (begin, finish gog.Point)
}

var _ Curve = new(Line)

type Line struct {
	p1, p2 gog.Point
}

func (l Line) GetVector(rand float64) (r Ray) {
	if rand < 0 || 1 < rand {
		panic(fmt.Errorf("not valid random value: %.5f", rand))
	}
	r.start.X = l.p1.X + (l.p2.X-l.p1.X)*rand
	r.start.Y = l.p1.Y + (l.p2.Y-l.p1.Y)*rand
	r.finish.X = r.start.X + (l.p2.X - l.p1.X)
	r.finish.Y = r.start.Y + (l.p2.Y - l.p1.Y)
	// rotate at 90 degree
	xc, yc := r.start.X, r.start.Y
	res := gog.Rotate(xc, yc, math.Pi/2.0, r.finish)
	r.finish = res
	// avoid scaling to lenght of vector to 1.0 and it is ok
	return
}

func (l Line) Box() (begin, finish gog.Point) {
	begin.X = min(l.p1.X, l.p2.X)
	begin.Y = min(l.p1.Y, l.p2.Y)
	finish.X = max(l.p1.X, l.p2.X)
	finish.Y = max(l.p1.Y, l.p2.Y)
	return
}

var Amount int64 = 100000

var (
	debug    bool
	intersec []Line
	miss     []Line
)

func OneCurve(present Curve, curves []Curve) (viewFactor []float64) {
	if debug {
		intersec = []Line{}
		miss = []Line{}
	}
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
		scale *= 2.1
	}
	// calculation
	var mut sync.Mutex
	counter := int64(0)
	run := func(steps int64) {
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
			distance := math.MaxFloat64
			for i := range curves {
				switch c := curves[i].(type) {
				case Line:
					// TODO : if gog.Orientation(c.p1, c.p2, v.start) != gog.ClockwisePoints {
					// TODO : 	continue
					// TODO : }
					pi, stA, stB := gog.LineLine(
						v.start, v.finish,
						c.p1, c.p2,
					)
					_ = stA
					_ = stB
					if len(pi) == 0 {
						if debug {
							miss = append(miss, Line{v.start, v.finish})
						}
						continue
					}
					d := gog.Distance(v.start, pi[0])
					if 1e-6 < math.Abs(d) && d < distance {
						index = i
						distance = d
						if debug {
							intersec = append(intersec, Line{v.start, pi[0]})
						}
					}
				default:
					panic(fmt.Errorf("not implemented: %#v", v))
				}
			}
			// get first intersection
			mut.Lock()
			counter++
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
				go func() {
					run(dstep)
					wg.Done()
				}()
			} else {
				go func() {
					run(amount)
					wg.Done()
				}()
			}
		}
		wg.Wait()
	} else {
		run(Amount)
	}
	if counter != Amount {
		panic("not valid amount counter")
	}

	viewFactor = make([]float64, len(curves))
	for i := range curves {
		viewFactor[i] = float64(vf[i]) / float64(Amount)
	}
	return
}
