package shapes

import "github.com/veandco/go-sdl2/sdl"

type RectangleShape struct {
	Rect  sdl.Rect
	Color sdl.Color
}

func (r *RectangleShape) setSize(width, height int32) {
	r.Rect.W = width
	r.Rect.H = height
}
