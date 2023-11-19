package main

import (
	"math"

	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

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

	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()

	font, err := ttf.OpenFont("C:/Windows/Fonts/arial.ttf", 20)
	if err != nil {
		panic(sdl.GetError())
	}
	defer font.Close()

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
	abrRender := &shapes.AbrRenderer{Renderer: rend}
	defer abrRender.Destroy()

	rect3 := shapes.RectangleShape{Rect: sdl.Rect{X: 300, Y: 200, W: 200, H: 200}, Color: sdl.Color{R: 255, G: 0, B: 255, A: 255}}
	abrRender.DrawRectangleShape(&rect3)

	c := shapes.MakeCircleShape(300, 200, 50, 4, sdl.Color{R: 0, G: 0, B: 255, A: 255})
	c.SetOrigin(50, 50)
	abrRender.DrawCircleShape(c)

	c.Rotate(math.Pi / 3)
	c.SetColor(sdl.Color{R: 0, G: 255, B: 0, A: 255})
	abrRender.DrawCircleShape(c)

	ps := sdl.Point{X: 350, Y: 250}
	pe := sdl.Point{X: 400, Y: 250}
	pe2 := sdl.Point{X: 400, Y: 300}

	abrRender.DrawLine(ps.X, ps.Y, pe.X, pe.Y)

	b := gfx.ThickLineRGBA(abrRender.Renderer, ps.X, ps.Y, pe2.X, pe2.Y, 10, 0, 0, 255, 255)
	if b != true {
		panic(sdl.GetError())
	}

	str := "Hello, World!"

	text_surf, err := font.RenderUTF8Solid(str, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		panic(err)
	}
	defer text_surf.Free()

	text_texture, err := abrRender.CreateTextureFromSurface(text_surf)
	if err != nil {
		panic(err)
	}
	defer text_texture.Destroy()

	ws, hs, err := font.SizeUTF8(str)

	text_rect := sdl.Rect{X: 0, Y: 0, W: int32(ws), H: int32(hs)}
	abrRender.Copy(text_texture, nil, &text_rect)

	abrRender.Present()

	frameCounter := 0

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

		abrRender.SetDrawColor(0, 0, 0, 255)
		abrRender.Clear()
		abrRender.Present()

		abrRender.SetDrawColor(250, 0, 0, 255)
		shapes.RotateSdlPointBase(&pe, ps, math.Pi/3)
		abrRender.DrawLine(ps.X, ps.Y, pe.X, pe.Y)

		c.Rotate(math.Pi / 3)
		abrRender.DrawCircleShape(c)
		abrRender.Present()

		end := sdl.GetPerformanceCounter()

		elapsed := float64(end-start) / perfFreq * 1000.0

		if 1000.0/30.0 > elapsed {
			sdl.Delay(uint32(math.Floor(1000.0/60.0 - elapsed)))
		}

		println(frameCounter)

	}

	println("Exit")
}
