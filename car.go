package main

import (
	"math"
)

type vec struct {
	x, y float64
}

// Car the car of the game
type Car struct {
	x, y, width, height, angle, angularVel, power, turn, drag, angDrag float64
	vel                                                                vec
}

// NewCar constructor for Car returns *Car
func NewCar(x, y, angle, w, h float64) *Car {
	return &Car{
		x:          x,
		y:          y,
		width:      w,
		height:     h,
		angle:      angle,
		angularVel: 0,
		power:      0.5,
		turn:       0.3,
		drag:       0.9,
		angDrag:    0.9,
		vel:        vec{0, 0}}
}

func (c *Car) drawAngle() float64 {
	return float64(int(c.angle)%360) * 2 * math.Pi / 360
}

func (c *Car) acc() {
	c.vel.x += math.Sin(c.angle*math.Pi/180) * c.power
	c.vel.y -= math.Cos(c.angle*math.Pi/180) * c.power
}

func (c *Car) br() {
	c.vel.x = c.vel.x * 0.9
	c.vel.y = c.vel.y * 0.9
}

func (c *Car) turnLeft() {
	c.angularVel -= c.turn
}
func (c *Car) turnRight() {
	c.angularVel += c.turn
}

func (c *Car) update() {
	c.x += c.vel.x
	c.y += c.vel.y
	c.vel.x *= c.drag
	c.vel.y *= c.drag
	c.angle += c.angularVel
	c.angularVel *= c.angDrag
}
