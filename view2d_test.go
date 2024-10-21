package view2d

import (
	"fmt"
	"math"
	"testing"

	"github.com/Konstantin8105/efmt"
	"github.com/Konstantin8105/gog"
)

func a(angle float64) (F12, f float64) {
	var (
		angleR = angle * math.Pi / 180 // degree to radian
		w      = 1.0
	)
	l1 := Line{
		p1: gog.Point{X: 0, Y: 0},
		p2: gog.Point{X: w, Y: 0},
	}
	l2 := Line{
		p1: gog.Point{X: 0, Y: 0},
		p2: gog.Point{X: w * math.Cos(angleR), Y: w * math.Sin(angleR)},
	}
	vf := OneCurve(l1, []Curve{l1, l2})
	// expect: F12 = 1 - sin(1/2*angle)
	F12 = 1.0 - math.Sin(0.5*angleR)
	f = vf[1]
	return
}

func TestRecorder(t *testing.T) {
	{
		old := Amount
		Amount = 100
		defer func() {
			Amount = old
		}()
	}
	{
		old := debug
		debug = true
		defer func() {
			debug = old
		}()
	}
	F12, f := a(60)
	fmt.Println(miss)
	fmt.Println(intersec)
	fmt.Println("result", efmt.Sprint(F12), efmt.Sprint(f))
}

func TestL(t *testing.T) {
	for angle := 5.0; angle < 180; angle += 5 {
		F12, f := a(angle)
		fmt.Println(efmt.Sprint(angle), efmt.Sprint(F12), efmt.Sprint(f))
	}
	// TODO verification tests
}

func TestTriangle(t *testing.T) {
	var (
		p0 = gog.Point{X: 0, Y: 0}
		p1 = gog.Point{X: 1, Y: 0}
		p2 = gog.Point{X: 0, Y: 1}

		l0 = Line{p0, p1}
		l1 = Line{p1, p2}
		l2 = Line{p2, p0}
		// 		l0 = Line{p1, p0}
		// 		l1 = Line{p2, p1}
		// 		l2 = Line{p0, p2}
	)
	Calc([]Curve{l0, l1, l2})
}
