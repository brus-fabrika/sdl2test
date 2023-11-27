package main

import (
	"math"
	"strconv"

	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 600
	SCREEN_HEIGHT = 600
	FRAMERATE     = 30.0
)

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
	f := float32(1.0 / math.Tan(theta/2))
	//Znear := 0.1
	//Zfar := 100.0
	Zbase := float32(3.0)

	rotTheta := float32(0.0)

	ctrlButton := false

	running := true
	for running {

		frameCounter++
		currentFrameCounter++

		start := sdl.GetPerformanceCounter()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				if event.(*sdl.KeyboardEvent).Type == sdl.KEYDOWN {
					switch event.(*sdl.KeyboardEvent).Keysym.Sym {
					case sdl.K_ESCAPE:
						println("Quit")
						running = false
						break
					case sdl.K_UP:
						Zbase += 0.1
						break
					case sdl.K_DOWN:
						Zbase -= 0.1
						break
					case sdl.K_LCTRL:
						ctrlButton = true
						break
					case sdl.K_RCTRL:
						ctrlButton = true
						break
					case sdl.K_EQUALS: // zoom in
						if ctrlButton {
							f += 0.1
						}
						break
					case sdl.K_MINUS: // zoom out
						if ctrlButton {
							f -= 0.1
						}
						break
					}
				}
				if event.(*sdl.KeyboardEvent).Type == sdl.KEYUP {
					switch event.(*sdl.KeyboardEvent).Keysym.Sym {
					case sdl.K_LCTRL:
						ctrlButton = false
						break
					case sdl.K_RCTRL:
						ctrlButton = false
						break
					}
				}
			}
		}

		e.Renderer.SetDrawColor(0, 0, 0, 255)
		e.Renderer.Clear()
		e.RenderText("Frame: "+strconv.FormatInt(frameCounter, 10), 0, 0)
		e.RenderText("FPS: "+strconv.FormatFloat(curFps, 'f', 2, 64), 0, 16)

		// cube render pipeline
		t := m.Clone()
		for i := range t.Triangles {
			sinRotTheta, cosRotTheta := math.Sincos(float64(rotTheta))
			// X-rotate
			y := t.Triangles[i].A.Y*float32(cosRotTheta) - t.Triangles[i].A.Z*float32(sinRotTheta)
			z := t.Triangles[i].A.Y*float32(sinRotTheta) + t.Triangles[i].A.Z*float32(cosRotTheta)

			t.Triangles[i].A.Y = y
			t.Triangles[i].A.Z = z

			y = t.Triangles[i].B.Y*float32(cosRotTheta) - t.Triangles[i].B.Z*float32(sinRotTheta)
			z = t.Triangles[i].B.Y*float32(sinRotTheta) + t.Triangles[i].B.Z*float32(cosRotTheta)

			t.Triangles[i].B.Y = y
			t.Triangles[i].B.Z = z

			y = t.Triangles[i].C.Y*float32(cosRotTheta) - t.Triangles[i].C.Z*float32(sinRotTheta)
			z = t.Triangles[i].C.Y*float32(sinRotTheta) + t.Triangles[i].C.Z*float32(cosRotTheta)

			t.Triangles[i].C.Y = y
			t.Triangles[i].C.Z = z

			// Y-rotate
			x := t.Triangles[i].A.X*float32(cosRotTheta) + t.Triangles[i].A.Z*float32(sinRotTheta)
			z = -t.Triangles[i].A.X*float32(sinRotTheta) + t.Triangles[i].A.Z*float32(cosRotTheta)
			t.Triangles[i].A.X = x
			t.Triangles[i].A.Z = z

			x = t.Triangles[i].B.X*float32(cosRotTheta) + t.Triangles[i].B.Z*float32(sinRotTheta)
			z = -t.Triangles[i].B.X*float32(sinRotTheta) + t.Triangles[i].B.Z*float32(cosRotTheta)
			t.Triangles[i].B.X = x
			t.Triangles[i].B.Z = z

			x = t.Triangles[i].C.X*float32(cosRotTheta) + t.Triangles[i].C.Z*float32(sinRotTheta)
			z = -t.Triangles[i].C.X*float32(sinRotTheta) + t.Triangles[i].C.Z*float32(cosRotTheta)
			t.Triangles[i].C.X = x
			t.Triangles[i].C.Z = z

			//Z-rotate
			x = t.Triangles[i].A.X*float32(cosRotTheta) - t.Triangles[i].A.Y*float32(sinRotTheta)
			y = t.Triangles[i].A.X*float32(sinRotTheta) + t.Triangles[i].A.Y*float32(cosRotTheta)
			t.Triangles[i].A.X = x
			t.Triangles[i].A.Y = y

			x = t.Triangles[i].B.X*float32(cosRotTheta) - t.Triangles[i].B.Y*float32(sinRotTheta)
			y = t.Triangles[i].B.X*float32(sinRotTheta) + t.Triangles[i].B.Y*float32(cosRotTheta)
			t.Triangles[i].B.X = x
			t.Triangles[i].B.Y = y

			x = t.Triangles[i].C.X*float32(cosRotTheta) - t.Triangles[i].C.Y*float32(sinRotTheta)
			y = t.Triangles[i].C.X*float32(sinRotTheta) + t.Triangles[i].C.Y*float32(cosRotTheta)
			t.Triangles[i].C.X = x
			t.Triangles[i].C.Y = y

			// z-translate
			t.Triangles[i].A.Add(&shapes.Vec3F{X: 0.0, Y: 0.0, Z: Zbase})
			t.Triangles[i].B.Add(&shapes.Vec3F{X: 0.0, Y: 0.0, Z: Zbase})
			t.Triangles[i].C.Add(&shapes.Vec3F{X: 0.0, Y: 0.0, Z: Zbase})

			// perspective projection
			t.Triangles[i].A.Mul(f / t.Triangles[i].A.Z)
			t.Triangles[i].B.Mul(f / t.Triangles[i].B.Z)
			t.Triangles[i].C.Mul(f / t.Triangles[i].C.Z)

			// scale into view
			t.Triangles[i].A.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
			t.Triangles[i].B.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})
			t.Triangles[i].C.Add(&shapes.Vec3F{X: 1.0, Y: 1.0, Z: 0.0})

			t.Triangles[i].A.Mul(0.5)
			t.Triangles[i].B.Mul(0.5)
			t.Triangles[i].C.Mul(0.5)

			t.Triangles[i].A.X *= SCREEN_WIDTH
			t.Triangles[i].A.Y *= SCREEN_HEIGHT
			t.Triangles[i].B.X *= SCREEN_WIDTH
			t.Triangles[i].B.Y *= SCREEN_HEIGHT
			t.Triangles[i].C.X *= SCREEN_WIDTH
			t.Triangles[i].C.Y *= SCREEN_HEIGHT
		}

		e.Renderer.DrawMesh(t)
		// end pipeline

		// scale into view

		// scale axes into view
		d0 := d0_1
		d1 := d1_1
		d2 := d2_1
		d3 := d3_1

		// perspective projection
		d0.Mul(f / d0.Z)
		d1.Mul(f / d1.Z)
		d2.Mul(f / d2.Z)
		d3.Mul(f / d3.Z)

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
		gfx.ThickLineRGBA(e.Renderer.Renderer, int32(d0.X), int32(d0.Y), int32(d1.X), int32(d1.Y), 2, 0, 127, 0, 255)
		gfx.ThickLineRGBA(e.Renderer.Renderer, int32(d0.X), int32(d0.Y), int32(d2.X), int32(d2.Y), 2, 0, 127, 0, 255)
		gfx.ThickLineRGBA(e.Renderer.Renderer, int32(d0.X), int32(d0.Y), int32(d3.X), int32(d3.Y), 2, 0, 127, 0, 255)

		e.Renderer.Present()
		end := sdl.GetPerformanceCounter()

		elapsed := float64(end-start) / perfFreq * 1000.0

		globElapsed = float64(end-globStart) / perfFreq * 1000.0
		if globElapsed > 10.0*1000.0 {
			globStart = sdl.GetPerformanceCounter()
			curFps = float64(currentFrameCounter) / globElapsed * 1000.0
			currentFrameCounter = 0
		}

		if 1000.0/FRAMERATE > elapsed {
			sdl.Delay(uint32(math.Floor(1000.0/FRAMERATE - elapsed)))
		}

		rotTheta += 0.01

	}
}
