package main

import (
	"math"
	"runtime"
	"strconv"

	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH        = 600
	SCREEN_HEIGHT       = 600
	FRAMERATE           = 30.0
	USE_FIXED_FRAMERATE = false
)

type PerspectiveParameters struct {
	Theta float32
	Zoom  float32
	Zbase float32
	Znear float32
	Zfar  float32
}

func main() {

	//	surface, err := window.GetSurface()
	//	if err != nil {
	//		panic(err)
	//	}
	//	surface.FillRect(nil, 0)
	//
	//	rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	//	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	//	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	//	surface.FillRect(&rect, pixel)

	e := Engine{}
	if err := e.Init(); err != nil {
		e.Destroy()
		panic(err)
	}
	defer e.Destroy()

	m := shapes.Mesh{
		Color: sdl.Color{R: 255, G: 0, B: 255, A: 255},
	}

	err := m.LoadFromFile("obj/sphere.obj")
	if err != nil {
		panic(err)
	}

	frameCounter := int64(0)
	currentFrameCounter := int64(0)

	perfFreq := float64(sdl.GetPerformanceFrequency())
	globStart := sdl.GetPerformanceCounter()
	globElapsed := float64(0)
	curFps := float64(FRAMERATE)

	d0_1 := shapes.Vec3F{X: 0.0, Y: 0.0, Z: 1.0}
	d1_1 := shapes.Vec3F{X: 1.0, Y: 0.0, Z: 1.0}
	d2_1 := shapes.Vec3F{X: 0.0, Y: 1.0, Z: 1.0}
	d3_1 := shapes.Vec3F{X: 0.0, Y: 0.0, Z: 2.0}

	// perspective parameters
	theta := math.Pi / 3
	perspective := PerspectiveParameters{
		Theta: math.Pi / 3,
		Zoom:  float32(1.0 / math.Tan(theta/2)),
		Zbase: 6.8,
		Znear: 0.1,
		Zfar:  100.0,
	}

	rotTheta := float32(0.0)

	ctrlButton := false
	running := true
	paused := false

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)

	for running {

		frameCounter++
		currentFrameCounter++

		start := sdl.GetPerformanceCounter()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				if event.(*sdl.KeyboardEvent).Type == sdl.KEYDOWN {
					switch event.(*sdl.KeyboardEvent).Keysym.Sym {
					case sdl.K_ESCAPE:
						println("Quit")
						running = false
					case sdl.K_SPACE:
						paused = !paused
						if !paused {
							frameCounter = 0
							currentFrameCounter = 0
						}
					case sdl.K_UP:
						perspective.Zbase += 0.1
					case sdl.K_DOWN:
						perspective.Zbase -= 0.1
					case sdl.K_LCTRL:
						ctrlButton = true
					case sdl.K_RCTRL:
						ctrlButton = true
					case sdl.K_EQUALS: // zoom in
						if ctrlButton {
							perspective.Zoom += 0.1
						}
					case sdl.K_MINUS: // zoom out
						if ctrlButton {
							perspective.Zoom -= 0.1
						}
					}
				}
				if event.(*sdl.KeyboardEvent).Type == sdl.KEYUP {
					switch event.(*sdl.KeyboardEvent).Keysym.Sym {
					case sdl.K_LCTRL:
						ctrlButton = false
					case sdl.K_RCTRL:
						ctrlButton = false
					}
				}
			case *sdl.MouseWheelEvent:
				perspective.Zoom += 0.1 * float32(event.(*sdl.MouseWheelEvent).Y)
			}
		}

		if paused {
			continue
		}

		e.Renderer.SetDrawColor(0, 0, 0, 255)
		e.Renderer.Clear()

		e.RenderText("Frame: "+strconv.FormatInt(frameCounter, 10), 0, 0)
		e.RenderText("FPS: "+strconv.FormatFloat(curFps, 'f', 2, 64), 0, 16)
		e.RenderText("Zbase: "+strconv.FormatFloat(float64(perspective.Zbase), 'f', 2, 64), 0, 32)
		e.RenderText("Zoom: "+strconv.FormatFloat(float64(perspective.Zoom), 'f', 2, 64), 0, 48)

		e.RenderText("Heap Sys: "+strconv.FormatFloat(float64(mem.HeapSys)/1024/1024, 'f', 2, 64), 20*16, 0)
		e.RenderText("Heap Alloc: "+strconv.FormatFloat(float64(mem.HeapAlloc)/1024/1024, 'f', 2, 64), 20*16, 16)
		e.RenderText("Heap In Use: "+strconv.FormatFloat(float64(mem.HeapInuse)/1024/1024, 'f', 2, 64), 20*16, 32)
		e.RenderText("Heap Idle: "+strconv.FormatFloat(float64(mem.HeapIdle)/1024/1024, 'f', 2, 64), 20*16, 48)

		// cube render pipeline
		t := m.Clone()
		for i := range t.Triangles {
			RotateTriangle(&t.Triangles[i], float64(rotTheta))
			ProjectTriangle(&t.Triangles[i], perspective)
		}

		//e.Renderer.MeshDebugPrint(t)

		e.Renderer.DrawMesh(t)
		e.Renderer.DrawRectangleShape(&shapes.RectangleShape{Rect: sdl.Rect{X: SCREEN_WIDTH / 4, Y: SCREEN_HEIGHT / 4, W: SCREEN_WIDTH / 2, H: SCREEN_HEIGHT / 2}, Color: sdl.Color{127, 127, 127, 0}})
		//e.Renderer.DrawMeshInRect(t, SCREEN_WIDTH/4, SCREEN_HEIGHT/4, SCREEN_WIDTH/2, SCREEN_HEIGHT/2)
		//e.Renderer.DrawMeshInRect(t, 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)
		// end pipeline

		// scale into view

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

		e.Renderer.Present()
		end := sdl.GetPerformanceCounter()

		elapsed := float64(end-start) / perfFreq * 1000.0

		globElapsed = float64(end-globStart) / perfFreq * 1000.0
		if globElapsed > 10.0*1000.0 {
			globStart = sdl.GetPerformanceCounter()
			curFps = float64(currentFrameCounter) / globElapsed * 1000.0
			currentFrameCounter = 0

			runtime.ReadMemStats(&mem)
		}

		if USE_FIXED_FRAMERATE {
			if 1000.0/FRAMERATE > elapsed {
				sdl.Delay(uint32(math.Floor(1000.0/FRAMERATE - elapsed)))
			}
		}

		rotTheta += 0.01

	}
}

func RotateTriangle(t *shapes.TriangleF, theta float64) {
	sinRotTheta, cosRotTheta := math.Sincos(theta)
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

	// Y-rotate
	x := t.A.X*float32(cosRotTheta) + t.A.Z*float32(sinRotTheta)
	z = -t.A.X*float32(sinRotTheta) + t.A.Z*float32(cosRotTheta)
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

	//Z-rotate
	x = t.A.X*float32(cosRotTheta) - t.A.Y*float32(sinRotTheta)
	y = t.A.X*float32(sinRotTheta) + t.A.Y*float32(cosRotTheta)
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
