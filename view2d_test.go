package view2d

import (
	"fmt"
	"testing"

	"github.com/Konstantin8105/gog"
)

func Test(t *testing.T) {
	vf := OneCurve(0, []Curve{
		Line{
			p1: gog.Point{X: 0, Y: 0},
			p2: gog.Point{X: 1, Y: 0},
		},
		Line{
			p1: gog.Point{X: 0, Y: 1},
			p2: gog.Point{X: 0, Y: 0},
		},
	})
	fmt.Println("expect") // TODO
	fmt.Println(vf)
}
