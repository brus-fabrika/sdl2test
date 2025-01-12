package shapes

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type AbrRenderer struct {
	*sdl.Renderer
}

type Vec3F struct {
	X, Y, Z float32
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
	Triangles []TriangleF
	Color     sdl.Color
}

func (m *Mesh) Clone() *Mesh {
	newMesh := &Mesh{
		Color: m.Color,
	}
	newMesh.Triangles = make([]TriangleF, len(m.Triangles))
	copy(newMesh.Triangles, m.Triangles)
	return newMesh
}

func (m *Mesh) LoadFromFile(fp string) error {
	readFile, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	vs := make([]Vec3F, 0)
	ts := make([]TriangleF, 0)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		s := strings.Split(line, " ")
		if s[0] == "#" {
			// skip comments
		} else if s[0] == "v" {
			// vertex
			f1, _ := strconv.ParseFloat(s[1], 32)
			f2, _ := strconv.ParseFloat(s[2], 32)
			f3, _ := strconv.ParseFloat(s[3], 32)

			v := Vec3F{X: float32(f1), Y: float32(f2), Z: float32(f3)}
			vs = append(vs, v)
		} else if s[0] == "vt" {
			// texture coordinate
		} else if s[0] == "vn" {
			// normal
		} else if s[0] == "f" {
			// face

			// keep it super simple for now
			v1, _ := strconv.ParseInt(strings.Split(s[1], "/")[0], 10, 64)
			v2, _ := strconv.ParseInt(strings.Split(s[2], "/")[0], 10, 64)
			v3, _ := strconv.ParseInt(strings.Split(s[3], "/")[0], 10, 64)

			ts = append(ts, TriangleF{A: vs[v1-1], B: vs[v2-1], C: vs[v3-1]})
			/*
				if len(ss) == 3 {
					// vertex/texture/normal
				} else if len(ss) == 2 {
					// vertex/texture
				} else if len(ss) == 1 {
					// vertex
				} else {
					// unknown
				}
			*/
		} else {
			// unknown
		}
	}

	m.Triangles = ts

	return nil
}

func Add(v1, v2 *Vec3F) Vec3F {
	return Vec3F{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

func (t *Vec3F) Add(v *Vec3F) {
	t.X += v.X
	t.Y += v.Y
	t.Z += v.Z
}

func Mul3F(v *Vec3F, s float32) Vec3F {
	return Vec3F{v.X * s, v.Y * s, v.Z * s}
}

func (t *Vec3F) Mul(s float32) {
	t.X *= s
	t.Y *= s
	t.Z *= s
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
		err = renderer.DrawLineF(t.A.X, t.A.Y, t.B.X, t.B.Y)
		err = renderer.DrawLineF(t.B.X, t.B.Y, t.C.X, t.C.Y)
		err = renderer.DrawLineF(t.C.X, t.C.Y, t.A.X, t.A.Y)
	}

	return err
}

func (renderer *AbrRenderer) DrawMeshInRect(m *Mesh, x float32, y float32, w float32, h float32) error {
	if err := renderer.SetDrawColor(m.Color.R, m.Color.G, m.Color.B, m.Color.A); err != nil {
		return err
	}

	var err error

	for _, t := range m.Triangles {
		if (t.A.X < x || t.A.Y < y || t.A.X > x+w || t.A.Y > y+h) && (t.B.X < x || t.B.Y < y || t.B.X > x+w || t.B.Y > y+h) && (t.C.X < x || t.C.Y < y || t.C.X > x+w || t.C.Y > y+h) {
			continue
		}
		err = renderer.DrawLineF(t.A.X, t.A.Y, t.B.X, t.B.Y)
		err = renderer.DrawLineF(t.B.X, t.B.Y, t.C.X, t.C.Y)
		err = renderer.DrawLineF(t.C.X, t.C.Y, t.A.X, t.A.Y)
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

func (renderer *AbrRenderer) MeshDebugPrint(m *Mesh) {
	fmt.Printf("Mesh:\n")
	for _, t := range m.Triangles {
		fmt.Printf("[{%f,%f,%f}, {%f,%f,%f}, {%f,%f,%f}]\n", t.A.X, t.A.Y, t.A.Z, t.B.X, t.B.Y, t.B.Z, t.C.X, t.C.Y, t.C.Z)
	}
}
