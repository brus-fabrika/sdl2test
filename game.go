package main

import (
	"math"
	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
)

const (
	SCREEN_WIDTH        = 600
	SCREEN_HEIGHT       = 600
	FRAMERATE           = 30.0
	USE_FIXED_FRAMERATE = true
)


type RotateDirection int

const (
	NO_ROTATE RotateDirection	= 0
	X_ROTATE = 1
	Y_ROTATE = 2
	Z_ROTATE = 4
	ALL_ROTATE = X_ROTATE | Y_ROTATE | Z_ROTATE
)


type PerspectiveParameters struct {
	Theta float32
	Zoom  float32
	Zbase float32
	Znear float32
	Zfar  float32
}

type GameState struct {
	ctrlButton bool
	running bool
	paused bool
}

type Game struct {
	e			*Engine
	m			shapes.Mesh
	perspective PerspectiveParameters

	rotTheta	float32
	state		GameState
}

func (g *Game) Init(e *Engine) error {
	g.e = e
	g.e.SetEngineListener(g)

	// perspective parameters
	theta := math.Pi / 3
	g.perspective = PerspectiveParameters{
		Theta: math.Pi / 3,
		Zoom:  float32(1.0 / math.Tan(theta/2)),
		Zbase: 6.8,
		Znear: 0.1,
		Zfar:  100.0,
	}

	g.rotTheta = float32(0.0)

	g.state.ctrlButton = false

	g.m.Color = COLOR_PURPLE
	err := g.m.LoadFromFile("./obj/sphere.obj")
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) OnPaused(){}
func (g *Game) OnUnpaused(){}
func (g *Game) OnExit() {
	println("Quit")
}

func (g *Game) OnPreRender(){
	// mesh render pipeline
	t := g.m.Clone()
	for i := range t.Triangles {
		tr := &t.Triangles[i]
		
		RotateTriangle(tr, float64(g.rotTheta), ALL_ROTATE)

		// check if current triangle is visible if looking from (0,0,0)
		
		// for now let's calculate normale vector ourselves
		ab := shapes.Vec3F{X: tr.B.X - tr.A.X, Y: tr.B.Y - tr.A.Y, Z: tr.B.Z - tr.A.Z}
		ac := shapes.Vec3F{X: tr.C.X - tr.A.X, Y: tr.C.Y - tr.A.Y, Z: tr.C.Z - tr.A.Z}
		nX := ab.Y*ac.Z - ab.Z*ac.Y
		nY := ab.Z*ac.X - ab.X*ac.Z
		nZ := ab.X*ac.Y - ab.Y*ac.X
		nmod := float32(math.Sqrt(float64(nX*nX + nY*nY + nZ*nZ)))
		nv := shapes.Vec3F{X: nX/nmod, Y: nY/nmod, Z: nZ/nmod}

		t.Normales[i] = nv
		t.Visibles[i] = nv.Z < 0

		ProjectTriangle(tr, g.perspective)
	}

	g.e.Renderer.DrawMesh(t)
	g.e.Renderer.DrawRectangleShape(
		&shapes.RectangleShape{
			Rect: sdl.Rect{X: SCREEN_WIDTH / 4, Y: SCREEN_HEIGHT / 4, W: SCREEN_WIDTH / 2, H: SCREEN_HEIGHT / 2},
			Color: COLOR_GRAY,
	})

	DrawAxes(g.e, &g.perspective)

	g.rotTheta += 0.01
}

func (g *Game) OnRender(){}
func (g *Game) OnPostRender(){}

func (g *Game) OnKeyEvent(event sdl.Event) {
	g.handleKeyboardEvent(event)
}

func (g *Game) OnMouseWheelEvent(event sdl.Event) {
	g.handleMouseWheelEvent(event)	
}

func (g *Game) LoadMeshFromFile(f string) error {
	m := shapes.Mesh{
		Color: COLOR_PURPLE,
	}

	return m.LoadFromFile(f)
}

func (g *Game) handleMouseWheelEvent(event sdl.Event) {
	if event == nil {
		return
	}
	g.perspective.Zoom += 0.1 * float32(event.(*sdl.MouseWheelEvent).Y)
}

func (g *Game) handleKeyboardEvent(event sdl.Event/*, gameState *GameState*/) {
	if event == nil {
		return
	}

	if event.(*sdl.KeyboardEvent).Type == sdl.KEYDOWN {
		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_SPACE:
			/*gameState.paused = !gameState.paused
			if !gameState.paused {
				//frameCounter = 0
				//currentFrameCounter = 0
			}*/
		case sdl.K_UP:
			g.perspective.Zbase += 0.1
		case sdl.K_DOWN:
			g.perspective.Zbase -= 0.1
		case sdl.K_LCTRL:
			g.state.ctrlButton = true
		case sdl.K_RCTRL:
			g.state.ctrlButton = true
		case sdl.K_EQUALS: // zoom in
			if g.state.ctrlButton {
				g.perspective.Zoom += 0.1
			}
		case sdl.K_MINUS: // zoom out
			if g.state.ctrlButton {
				g.perspective.Zoom -= 0.1
			}
		}
	}
	if event.(*sdl.KeyboardEvent).Type == sdl.KEYUP {
		switch event.(*sdl.KeyboardEvent).Keysym.Sym {
		case sdl.K_LCTRL:
			g.state.ctrlButton = false
		case sdl.K_RCTRL:
			g.state.ctrlButton = false
		}
	}
}

func RotateTriangle(t *shapes.TriangleF, theta float64, rotate RotateDirection) {
	if rotate == NO_ROTATE {
		return
	}

	sinRotTheta, cosRotTheta := math.Sincos(theta)

	if rotate & X_ROTATE != 0 {
		// X-rotate
		y := t.A.Y*float32(cosRotTheta) - t.A.Z*float32(sinRotTheta)
		z := t.A.Y*float32(sinRotTheta) + t.A.Z*float32(cosRotTheta)

		t.A.Y = y
		t.A.Z = z

		y = t.B.Y*float32(cosRotTheta) - t.B.Z*float32(sinRotTheta)
		z = t.B.Y*float32(sinRotTheta) + t.B.Z*float32(cosRotTheta)

		t.B.Y = y
		t.B.Z = z

		y = t.C.Y*float32(cosRotTheta) - t.C.Z*float32(sinRotTheta)
		z = t.C.Y*float32(sinRotTheta) + t.C.Z*float32(cosRotTheta)

		t.C.Y = y
		t.C.Z = z
	}

	if rotate & Y_ROTATE != 0 {
		// Y-rotate
		x := t.A.X*float32(cosRotTheta) + t.A.Z*float32(sinRotTheta)
		z := -t.A.X*float32(sinRotTheta) + t.A.Z*float32(cosRotTheta)
		t.A.X = x
		t.A.Z = z

		x = t.B.X*float32(cosRotTheta) + t.B.Z*float32(sinRotTheta)
		z = -t.B.X*float32(sinRotTheta) + t.B.Z*float32(cosRotTheta)
		t.B.X = x
		t.B.Z = z

		x = t.C.X*float32(cosRotTheta) + t.C.Z*float32(sinRotTheta)
		z = -t.C.X*float32(sinRotTheta) + t.C.Z*float32(cosRotTheta)
		t.C.X = x
		t.C.Z = z
	}

	if rotate & Z_ROTATE != 0 {
		//Z-rotate
		x := t.A.X*float32(cosRotTheta) - t.A.Y*float32(sinRotTheta)
		y := t.A.X*float32(sinRotTheta) + t.A.Y*float32(cosRotTheta)
		t.A.X = x
		t.A.Y = y

		x = t.B.X*float32(cosRotTheta) - t.B.Y*float32(sinRotTheta)
		y = t.B.X*float32(sinRotTheta) + t.B.Y*float32(cosRotTheta)
		t.B.X = x
		t.B.Y = y

		x = t.C.X*float32(cosRotTheta) - t.C.Y*float32(sinRotTheta)
		y = t.C.X*float32(sinRotTheta) + t.C.Y*float32(cosRotTheta)
		t.C.X = x
		t.C.Y = y
	}
}

func ProjectTriangle(t *shapes.TriangleF, p PerspectiveParameters) {
	// z-translate
	t.A.Add(&shapes.Vec3F{X: 0.0, Y: 0.0, Z: p.Zbase})
	t.B.Add(&shapes.Vec3F{X: 0.0, Y: 0.0, Z: p.Zbase})
	t.C.Add(&shapes.Vec3F{X: 0.0, Y: 0.0, Z: p.Zbase})

	// perspective projection
	t.A.Mul(p.Zoom / t.A.Z)
	t.B.Mul(p.Zoom / t.B.Z)
	t.C.Mul(p.Zoom / t.C.Z)

	// scale into view
	t.A.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
	t.B.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
	t.C.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})

	t.A.Mul(0.5)
	t.B.Mul(0.5)
	t.C.Mul(0.5)

	t.A.X *= SCREEN_WIDTH
	t.A.Y *= SCREEN_HEIGHT
	t.B.X *= SCREEN_WIDTH
	t.B.Y *= SCREEN_HEIGHT
	t.C.X *= SCREEN_WIDTH
	t.C.Y *= SCREEN_HEIGHT
}

func DrawAxes(e *Engine, perspective *PerspectiveParameters) {
	d0_1 := shapes.Vec3F{X: 0.0, Y: 0.0, Z: 1.0}
	d1_1 := shapes.Vec3F{X: 1.0, Y: 0.0, Z: 1.0}
	d2_1 := shapes.Vec3F{X: 0.0, Y: 1.0, Z: 1.0}
	d3_1 := shapes.Vec3F{X: 0.0, Y: 0.0, Z: 2.0}

	// scale axes into view
	d0 := d0_1
	d1 := d1_1
	d2 := d2_1
	d3 := d3_1

	// perspective projection
	d0.Mul(perspective.Zoom / d0.Z)
	d1.Mul(perspective.Zoom / d1.Z)
	d2.Mul(perspective.Zoom / d2.Z)
	d3.Mul(perspective.Zoom / d3.Z)

	d0.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
	d1.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
	d2.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
	d3.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})

	d0.Mul(0.5)
	d1.Mul(0.5)
	d2.Mul(0.5)
	d3.Mul(0.5)

	d0.X *= SCREEN_WIDTH
	d0.Y *= SCREEN_HEIGHT
	d1.X *= SCREEN_WIDTH
	d1.Y *= SCREEN_HEIGHT
	d2.X *= SCREEN_WIDTH
	d2.Y *= SCREEN_HEIGHT
	d3.X *= SCREEN_WIDTH
	d3.Y *= SCREEN_HEIGHT
	// draw axes
	gfx.ThickLineRGBA(e.Renderer.Renderer, int32(d0.X), int32(d0.Y), int32(d1.X), int32(d1.Y), 2, 127, 127, 127, 255)
	gfx.ThickLineRGBA(e.Renderer.Renderer, int32(d0.X), int32(d0.Y), int32(d2.X), int32(d2.Y), 2, 127, 127, 127, 255)
	gfx.ThickLineRGBA(e.Renderer.Renderer, int32(d0.X), int32(d0.Y), int32(d3.X), int32(d3.Y), 2, 127, 127, 127, 255)
}
