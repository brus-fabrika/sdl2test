package shapes

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type CircleShape struct {
	C      sdl.Point
	Points []sdl.Point
	Color  sdl.Color
	origin sdl.Point
}

func MakeCircleShape(x, y, radius int32, n int32, color sdl.Color) *CircleShape {
	circle := CircleShape{C: sdl.Point{X: x, Y: y}, origin: sdl.Point{X: x, Y: y}, Color: color}
	circle.Points = make([]sdl.Point, n)

	for i, p := range circle.Points {
		p.X = circle.C.X + radius + int32(float64(radius)*math.Cos(float64(i)*2.0*math.Pi/float64(n)))
		p.Y = circle.C.Y + radius + int32(float64(radius)*math.Sin(float64(i)*2.0*math.Pi/float64(n)))
		circle.Points[i] = p
	}

	return &circle
}

func (c *CircleShape) Rotate(a float64) {
	for i := range c.Points {
		RotateSdlPointBase(&c.Points[i], c.origin, a)
	}
}

func (r *CircleShape) SetColor(c sdl.Color) {
	r.Color = c
}

func (r *CircleShape) SetOrigin(x, y int32) {
	r.origin.X += x
	r.origin.Y += y
}
