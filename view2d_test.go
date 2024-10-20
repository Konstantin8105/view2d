package view2d

import (
	"math"
	"testing"

	"github.com/Konstantin8105/efmt"
	"github.com/Konstantin8105/gog"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

type arr struct {
	l Line
}

func (a arr) Len() int { return 2 }
func (a arr) XY(i int) (x, y float64) {
	if i == 0 {
		return a.l.p1.X, a.l.p1.Y
	}
	return a.l.p2.X, a.l.p2.Y
}

func TestRecorder(t *testing.T) {
	{
		old := Amount
		size := int64(80)
		if size < Amount {
			Amount = size
		}
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
	a(90) // action

	for _, gr := range []struct {
		name string
		data []Line
	}{
		{"miss", miss},
		{"intersec", intersec},
	} {
		p := plot.New()

		p.Title.Text = "View factor"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		var v []interface{}
		v = append(v, gr.name)
		for i := range gr.data {
			a := arr{l: gr.data[i]}
			v = append(v, a)
		}
		err := plotutil.AddLinePoints(p, v...)
		if err != nil {
			panic(err)
		}
		// Save the plot to a PNG file.
		if err := p.Save(4*vg.Inch, 4*vg.Inch, gr.name+".png"); err != nil {
			panic(err)
		}
	}
}

func TestL(t *testing.T) {
	{
		old := Amount
		size := int64(100000)
		if size < Amount {
			Amount = size
		}
		defer func() {
			Amount = old
		}()
	}
	for angle := 5.0; angle < 100; angle += 5 {
		F12, f := a(angle)
		if diff := math.Abs((F12 - f) / F12); 1.0/100.0 < diff { // 1%
			t.Errorf("angle = %.2f diff = %.5f", angle, diff)
		}
		t.Logf("%s %s %s", efmt.Sprint(angle), efmt.Sprint(F12), efmt.Sprint(f))
	}
}

func TestTriangle(t *testing.T) {
	var (
		p0 = gog.Point{X: 0, Y: 0}
		p1 = gog.Point{X: 1, Y: 0}
		p2 = gog.Point{X: 0, Y: 1}

		l0 = Line{p0, p1}
		l1 = Line{p1, p2}
		l2 = Line{p2, p0}
		cs = []Curve{l0, l1, l2}
	)
	for i, l := range []Line{l0, l1, l2} {
		vf := OneCurve(l, cs)
		total := 0.0
		for i := range vf {
			total += vf[i]
		}
		if diff := math.Abs((total - 1) / 1.0); 1e-2 < diff {
			t.Errorf("i = %d diff = %.5f", i, diff)
		}
	}
}

func TestVerification(t *testing.T) {
	{
		old := Amount
		size := int64(100000)
		if size < Amount {
			Amount = size
		}
		defer func() {
			Amount = old
		}()
	}
	var (
		p0  = gog.Point{0.00, 0.00}
		p1  = gog.Point{0.35, 0.00}
		A12 = Line{p0, p1}

		p2  = gog.Point{0.8, 0.00}
		A11 = Line{p1, p2}

		p3  = gog.Point{1.00, 0.00}
		A10 = Line{p2, p3}

		p4 = gog.Point{1, 0.5}
		A9 = Line{p3, p4}

		p5 = gog.Point{0.6, 0.5}
		A8 = Line{p4, p5}

		p6 = gog.Point{1., 0.5001}
		A7 = Line{p5, p6}

		p7 = gog.Point{1, 1}
		A6 = Line{p6, p7}

		p8 = gog.Point{0, 1}
		A5 = Line{p7, p8}

		p9 = gog.Point{0, 0.5001}
		A4 = Line{p8, p9}

		p10 = gog.Point{0.4, 0.5}
		A2  = Line{p9, p10}

		p11 = gog.Point{0, 0.5}
		A3  = Line{p10, p11}

		A1 = Line{p11, p0}

		cs = []Curve{A1, A2, A3, A4, A5, A6, A7, A8, A9, A10, A11, A12}
	)
	vf := OneCurve(A5, cs)
	t.Logf("view factors: %.5f", vf)
	total := 0.0
	for i := range vf {
		total += vf[i]
	}
	t.Logf("total: %.5f", total)
}
