package main

import (
	"math"
	"strconv"

	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600
	FRAMERATE     = 10.0
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
		Triangles: []shapes.Triangle{
			shapes.Triangle{shapes.Vec3{100, 100, 0}, shapes.Vec3{200, 100, 0}, shapes.Vec3{300, 300, 0}},
		},
		Color: sdl.Color{R: 255, G: 0, B: 255, A: 255},
	}

	frameCounter := int64(0)

	perfFreq := float64(sdl.GetPerformanceFrequency())

	running := true
	for running {

		frameCounter++

		start := sdl.GetPerformanceCounter()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		e.Renderer.SetDrawColor(0, 0, 0, 255)
		e.Renderer.Clear()
		e.RenderText("Frame: " + strconv.FormatInt(frameCounter, 10))

		e.Renderer.DrawMesh(&m)

		e.Renderer.Present()
		end := sdl.GetPerformanceCounter()

		elapsed := float64(end-start) / perfFreq * 1000.0

		if 1000.0/FRAMERATE > elapsed {
			sdl.Delay(uint32(math.Floor(1000.0/FRAMERATE - elapsed)))
		}

	}
}
