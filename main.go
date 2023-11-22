package main

import (
	"math"

	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600
	FRAMERATE     = 10.0
)

type Engine struct {
	Window   *sdl.Window
	Renderer *shapes.AbrRenderer
	Font     *ttf.Font
}

func (e *Engine) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	wnd, err := sdl.CreateWindow("SDL2 Test Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		return err
	}
	e.Window = wnd

	rend, err := sdl.CreateRenderer(e.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	e.Renderer = &shapes.AbrRenderer{Renderer: rend}

	if err := ttf.Init(); err != nil {
		return err
	}

	fnt, err := ttf.OpenFont("C:/Windows/Fonts/arial.ttf", 20)
	if err != nil {
		return err
	}
	e.Font = fnt

	return nil
}

func (e *Engine) Destroy() {
	println("Destroying...")

	if e.Font != nil {
		e.Font.Close()
	}
	if e.Renderer != nil {
		e.Renderer.Destroy()
	}
	if e.Window != nil {
		e.Window.Destroy()
	}

	ttf.Quit()
	sdl.Quit()
}

func (e *Engine) Rendertext(s string) error {
	text_surf, err := e.Font.RenderUTF8Solid(s, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return err
	}
	defer text_surf.Free()

	text_texture, err := e.Renderer.CreateTextureFromSurface(text_surf)
	if err != nil {
		return err
	}
	defer text_texture.Destroy()

	ws, hs, err := e.Font.SizeUTF8(s)

	text_rect := sdl.Rect{X: 0, Y: 0, W: int32(ws), H: int32(hs)}
	e.Renderer.Copy(text_texture, nil, &text_rect)

	return nil
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

	//	e.Renderer.Present()

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

		e.Renderer.Present()
		end := sdl.GetPerformanceCounter()

		elapsed := float64(end-start) / perfFreq * 1000.0

		if 1000.0/FRAMERATE > elapsed {
			sdl.Delay(uint32(math.Floor(1000.0/FRAMERATE - elapsed)))
		}

		println(frameCounter)

	}

	println("Exit")
}
