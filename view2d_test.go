package view2d

import (
	"math"
	"testing"

	"github.com/Konstantin8105/efmt"
	"github.com/Konstantin8105/gog"
	"github.com/Konstantin8105/pow"

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
		P1: gog.Point{X: 0, Y: 0},
		P2: gog.Point{X: w, Y: 0},
	}
	l2 := Line{
		P1: gog.Point{X: 0, Y: 0},
		P2: gog.Point{X: w * math.Cos(angleR), Y: w * math.Sin(angleR)},
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
		return a.l.P1.X, a.l.P1.Y
	}
	return a.l.P2.X, a.l.P2.Y
}

func record() {
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
		if err := p.Save(8*vg.Inch, 8*vg.Inch, gr.name+".png"); err != nil {
			panic(err)
		}
	}
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
			record() // store data
			debug = old
		}()
	}
	a(90) // action
	record()
}

func TestLine(t *testing.T) {
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

func TestArc(t *testing.T) {
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
		r1 = 1.0
		r2 = 3.0
	)
	var (
		p10 = gog.Point{+r1, 0.0}
		p11 = gog.Point{0.0, +r1}
		p12 = gog.Point{-r1, 0.0}
		a11 = Arc{p10, p11, p12}

		p20 = gog.Point{-r2, 0.0}
		p21 = gog.Point{0.0, +r2}
		p22 = gog.Point{+r2, 0.0}
		a21 = Arc{p20, p21, p22}

		cs = []Curve{a11, a21}
	)
	var vfs [][]float64
	for i := range cs {
		vf := OneCurve(cs[i], cs)
		vfs = append(vfs, vf)
		t.Logf("view factors: %.5f", vf)
		total := 0.0
		for i := range vf {
			total += vf[i]
		}
		t.Logf("total: %.5f", total)
	}
	// expect
	var (
		r   = r1 / r2 // ratio
		F12 = 1 - math.Acos(r)/math.Pi + 1/(math.Pi*r)*(math.Sqrt(1-r*r)+r-1)
	)
	// compare F12
	actF12 := vfs[0][1]
	if diff := math.Abs((actF12 - F12) / F12); 1e-2 < diff {
		t.Errorf("F12: {%.5f, %.5f} diff = %.5f", actF12, F12, diff)
	}
}

func TestConcentricCylinder(t *testing.T) {
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
		r1 = 1.0
		r2 = 3.0
	)
	var (
		c1 = Circle{
			Radius:        r1,
			VectorOutside: true,
		}
		c2 = Circle{
			Radius: r2,
		}

		cs = []Curve{c1, c2}
	)
	var vfs [][]float64
	for i := range cs {
		vf := OneCurve(cs[i], cs)
		vfs = append(vfs, vf)
		t.Logf("view factors: %.5f", vf)
		total := 0.0
		for i := range vf {
			total += vf[i]
		}
		t.Logf("total: %.5f", total)
	}
	// expect
	var (
		r   = r1 / r2 // ratio
		F11 = 0.0
		F12 = 1.0
		F21 = r
		F22 = 1.0 - r
	)
	{
		// compare F12
		actF11 := vfs[0][0]
		if diff := math.Abs((actF11 - F11)); 1e-2 < diff {
			t.Errorf("F11: {%.5f, %.5f} diff = %.5f", actF11, F11, diff)
		}
	}
	{
		// compare F12
		actF12 := vfs[0][1]
		if diff := math.Abs((actF12 - F12) / F12); 1e-2 < diff {
			t.Errorf("F12: {%.5f, %.5f} diff = %.5f", actF12, F12, diff)
		}
	}
	{
		// compare F21
		actF21 := vfs[1][0]
		if diff := math.Abs((actF21 - F21) / F21); 1e-2 < diff {
			t.Errorf("F21: {%.5f, %.5f} diff = %.5f", actF21, F21, diff)
		}
	}
	{
		// compare F22
		actF22 := vfs[1][1]
		if diff := math.Abs((actF22 - F22) / F22); 1e-2 < diff {
			t.Errorf("F22: {%.5f, %.5f} diff = %.5f", actF22, F22, diff)
		}
	}
}

func TestSingleRow(t *testing.T) {
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
	for S := 2.5; S <= 5.1; S += 0.5 {
		var (
			D = 1.0
			S = 3.0
		)
		var (
			c1 = Circle{
				Radius:        D / 2.0,
				VectorOutside: true,
			}

			c2 = Circle{Center: gog.Point{+1 * S, 0}, Radius: D / 2.0, VectorOutside: true}

			p3 = gog.Point{X: S, Y: D/2 + 0.00001}
			p4 = gog.Point{X: 0, Y: D/2 + 0.00001}
			l2 = Line{p3, p4}

			cs = []Curve{c2, c1, l2}
		)
		vf := OneCurve(l2, cs)
		t.Logf("view factors: %.5f", vf)
		total := 0.0
		for i := range vf {
			total += vf[i]
		}
		t.Logf("total: %.5f", total)
		// expect
		var (
			Fij = 1.0 - math.Sqrt(1-pow.E2(D/S)) + (D/S)*math.Atan(math.Sqrt(pow.E2(S/D)-1.0))
		)
		{
			// compare Fij
			actFij := vf[0] + vf[1] // view factor on tube
			if diff := math.Abs((actFij - Fij) / Fij); 1e-2 < diff {
				t.Errorf("Fij: {%.5f, %.5f} diff = %.5f", actFij, Fij, diff)
			}
		}
	}
}
