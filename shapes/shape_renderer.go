package shapes

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type AbrRenderer struct {
	*sdl.Renderer
}

type Vec3F struct {
	X, Y, Z float64
}

type Vec3 struct {
	X, Y, Z int32
}

type TriangleF struct {
	A, B, C Vec3F
}

type Triangle struct {
	A, B, C Vec3
}

type Mesh struct {
	Triangles []Triangle
	Color     sdl.Color
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

func (renderer *AbrRenderer) DrawMesh(m *Mesh) error {
	if err := renderer.SetDrawColor(m.Color.R, m.Color.G, m.Color.B, m.Color.A); err != nil {
		return err
	}

	var err error

	for _, t := range m.Triangles {
		err = renderer.DrawLine(t.A.X, t.A.Y, t.B.X, t.B.Y)
		err = renderer.DrawLine(t.B.X, t.B.Y, t.C.X, t.C.Y)
		err = renderer.DrawLine(t.C.X, t.C.Y, t.A.X, t.A.Y)
	}

	return err
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
