package main

import "github.com/veandco/go-sdl2/sdl"

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("SDL2 Test Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

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

	rend, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer rend.Destroy()

	rect2 := sdl.Rect{X: 350, Y: 250, W: 100, H: 100}
	if err := rend.SetDrawColor(255, 255, 255, 255); err != nil {
		println("Error setting draw colour: ", sdl.GetError().Error())
		panic(err)
	}

	if err := rend.DrawRect(&rect2); err != nil {
		println("Error draw rectangle: ", sdl.GetError())
		panic(err)
	}

	rend.Present()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
