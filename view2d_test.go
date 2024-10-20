package view2d

import (
	"fmt"
	"testing"

	"github.com/Konstantin8105/gog"
)

func Test(t *testing.T) {
	l1 := Line{
		p1: gog.Point{X: 0, Y: 0},
		p2: gog.Point{X: 1, Y: 0},
	}
	l2 := Line{
		p1: gog.Point{X: 0, Y: 1},
		p2: gog.Point{X: 0, Y: 0},
	}
	vf := OneCurve(l1, []Curve{l1, l2})
	fmt.Println("expect") // TODO
	fmt.Println(vf)
	// TODO verification tests
}
