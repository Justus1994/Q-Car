package main

import "math"

/*** https://github.com/hajimehoshi/ebiten/blob/master/examples/raycasting/main.go **/
type line struct {
	X1, Y1, X2, Y2 float64
}

func (l *line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

func newRay(x, y, length, angle float64) line {
	return line{
		X1: x,
		Y1: y,
		X2: x + length*math.Cos(angle),
		Y2: y + length*math.Sin(angle),
	}
}

// intersection calculates the intersection of given two lines.
func intersection(l1, l2 line) (float64, float64, bool) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
	uNum := -((l1.X1-l1.X2)*(l1.Y1-l2.Y1) - (l1.Y1-l1.Y2)*(l1.X1-l2.X1))

	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom
	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}

	x := l1.X1 + t*(l1.X2-l1.X1)
	y := l1.Y1 + t*(l1.Y2-l1.Y1)
	return x, y, true
}

func carBoundings(x, y, w, h, angle float64) []line {

	top := newRay(x, y, math.Sqrt(math.Pow(h/2, 2)+math.Pow(w/2, 2)), angle-60*math.Pi/180)

	bot := newRay(x, y, math.Sqrt(math.Pow(h/2, 2)+math.Pow(w/2, 2)), angle-120*math.Pi/180)

	left := newRay(x, y, math.Sqrt(math.Pow(h/2, 2)+math.Pow(w/2, 2)), angle+60*math.Pi/180)
	right := newRay(x, y, math.Sqrt(math.Pow(h/2, 2)+math.Pow(w/2, 2)), angle+120*math.Pi/180)
	//left := newRay(x, y, h, angle)
	//right := newRay(top.X2, top.Y2, h, angle)
	//bottom := newRay(left.X2, left.Y2, w, angle)

	return []line{top, bot, left, right}

}
