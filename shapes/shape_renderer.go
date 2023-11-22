package shapes

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type AbrRenderer struct {
	*sdl.Renderer
}

func (renderer *AbrRenderer) DrawRectangleShape(rect *RectangleShape) error {
	if err := renderer.SetDrawColor(rect.Color.R, rect.Color.G, rect.Color.B, rect.Color.A); err != nil {
		return err
	}
	if err := renderer.DrawRect(&rect.Rect); err != nil {
		return err
	}
	return nil
}

func (renderer *AbrRenderer) DrawCircleShape(c *CircleShape) error {
	if err := renderer.SetDrawColor(c.Color.R, c.Color.G, c.Color.B, c.Color.A); err != nil {
		return err
	}
	if err := renderer.DrawLines(c.Points); err != nil {
		return err
	}
	p := c.Points[len(c.Points)-1]
	if err := renderer.DrawLine(p.X, p.Y, c.Points[0].X, c.Points[0].Y); err != nil {
		return err
	}
	return nil
}

func RotateSdlPoint(p *sdl.Point, angle float64) *sdl.Point {
	if p == nil {
		return nil
	}
	sina, cosa := math.Sincos(angle)
	newX := int32(math.Round(float64(p.X)*cosa - float64(p.Y)*sina))
	newY := int32(math.Round(float64(p.X)*sina + float64(p.Y)*cosa))
	p.X = newX
	p.Y = newY
	return p
}

func RotateSdlPointBase(p *sdl.Point, bp sdl.Point, angle float64) *sdl.Point {
	if p == nil {
		return nil
	}
	sina, cosa := math.Sincos(angle)
	newX := bp.X + int32(math.Round(float64(p.X-bp.X)*cosa-float64(p.Y-bp.Y)*sina))
	newY := bp.Y + int32(math.Round(float64(p.X-bp.X)*sina+float64(p.Y-bp.Y)*cosa))

	p.X = newX
	p.Y = newY
	return p
}

func RotatePoint(x, y, cx, cy, angle float64) (float64, float64) {
	sin, cos := math.Sincos(angle)
	x -= cx
	y -= cy
	return x*cos - y*sin + cx, x*sin + y*cos + cy
}
