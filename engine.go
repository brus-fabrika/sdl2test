package main

import (
	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

func (e *Engine) RenderText(s string) error {
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
