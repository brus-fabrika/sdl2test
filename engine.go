package main

import (
	"math"
	"runtime"
	"strconv"

	"github.com/brus-fabrika/sdl2test/shapes"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	COLOR_RED		= sdl.Color{R: 255, G:   0, B:   0, A: 255}
	COLOR_GREEN		= sdl.Color{R:   0, G: 255, B:   0, A: 255}
	COLOR_BLUE		= sdl.Color{R:   0, G:   0, B: 255, A: 255}
	COLOR_PURPLE	= sdl.Color{R: 255, G:   0, B: 255, A: 255}
	COLOR_WHITE		= sdl.Color{R: 255, G: 255, B: 255, A: 255}
	COLOR_GRAY		= sdl.Color{R: 127, G: 127, B: 127, A: 255}
	COLOR_BLACK		= sdl.Color{R:   0, G:   0, B:   0, A: 255}
)

type EngineListener interface {
	OnPaused()
	OnUnpaused()
	OnExit()
	OnPreRender()
	OnRender()
	OnPostRender()
	OnKeyEvent(event sdl.Event)
	OnMouseWheelEvent(event sdl.Event)
}

type EngineState int

const (
	STOPPED		= iota
	RUNNING
	PAUSED
)

type Engine struct {
	Window   	*sdl.Window
	Renderer 	*shapes.AbrRenderer
	Font     	*ttf.Font

	state		EngineState
	listener	EngineListener
}

func (e *Engine) OnPaused() {
	if e.listener != nil {
		e.listener.OnPaused()
	}
}

func (e *Engine) OnUnpaused() {
	if e.listener != nil {
		e.listener.OnUnpaused()
	}
}

func (e *Engine) OnPreRender() {
	if e.listener != nil {
		e.listener.OnPreRender()
	}
}

func (e *Engine) OnRender() {
	if e.listener != nil {
		e.listener.OnRender()
	}
}

func (e *Engine) OnPostRender() {
	if e.listener != nil {
		e.listener.OnPostRender()
	}
}

func (e *Engine) OnExit() {
	if e.listener != nil {
		e.listener.OnExit()
	}
}

func (e *Engine) OnKeyEvent(event sdl.Event) {
	if e.listener != nil {
		e.listener.OnKeyEvent(event)
	}
}

func (e *Engine) OnMouseWheelEvent(event sdl.Event) {
	if e.listener != nil {
		e.listener.OnMouseWheelEvent(event)
	}
}

func (e *Engine) Run() {
	frameCounter := int64(0)
	currentFrameCounter := int64(0)
	perfFreq := float64(sdl.GetPerformanceFrequency())
	globStart := sdl.GetPerformanceCounter()
	globElapsed := float64(0)
	curFps := float64(FRAMERATE)

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)

	e.state = RUNNING
	for e.state == RUNNING {
		// handle events
		e.processEvents()
		
		// pre-render
		frameCounter++
		currentFrameCounter++

		start := sdl.GetPerformanceCounter()

		e.Renderer.SetDrawColor(0, 0, 0, 255)
		e.Renderer.Clear()

		e.OnPreRender()
		
		// render
		e.RenderText("Frame: "+strconv.FormatInt(frameCounter, 10), 0, 0)
		e.RenderText("FPS: "+strconv.FormatFloat(curFps, 'f', 2, 64), 0, 16)
		//e.RenderText("Zbase: "+strconv.FormatFloat(float64(perspective.Zbase), 'f', 2, 64), 0, 32)
		//e.RenderText("Zoom: "+strconv.FormatFloat(float64(perspective.Zoom), 'f', 2, 64), 0, 48)

		//e.RenderText("Heap Sys: "+strconv.FormatFloat(float64(mem.HeapSys)/1024/1024, 'f', 2, 64), 20*16, 0)
		//e.RenderText("Heap Alloc: "+strconv.FormatFloat(float64(mem.HeapAlloc)/1024/1024, 'f', 2, 64), 20*16, 16)
		//e.RenderText("Heap In Use: "+strconv.FormatFloat(float64(mem.HeapInuse)/1024/1024, 'f', 2, 64), 20*16, 32)
		//e.RenderText("Heap Idle: "+strconv.FormatFloat(float64(mem.HeapIdle)/1024/1024, 'f', 2, 64), 20*16, 48)

		e.Renderer.Present()

		// postrender
		e.OnPostRender()

		end := sdl.GetPerformanceCounter()

		elapsed := float64(end-start) / perfFreq * 1000.0

		globElapsed = float64(end - globStart) / perfFreq * 1000.0
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
	}
}

func (e *Engine) IsRunning() bool {
	return e.state == RUNNING
}

func (e *Engine) SetEngineListener(l EngineListener) {
	e.listener = l
}

func (e *Engine) processEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			e.state = STOPPED
			e.OnExit()
		case *sdl.KeyboardEvent:
			switch event.(*sdl.KeyboardEvent).Keysym.Sym {
			case sdl.K_ESCAPE:
				e.state = STOPPED
				e.OnExit()
			default:
				e.OnKeyEvent(event) // pass to the listener(s)
			}
		case *sdl.MouseWheelEvent:
				e.OnMouseWheelEvent(event)
		}
	}
}

func (e *Engine) Init() error {
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

	fnt, err := ttf.OpenFont("/home/abrus/.fonts/FiraCodeNerdFont-Regular.ttf", 16)
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

func (e *Engine) RenderText(s string, x, y int32) error {
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

	text_rect := sdl.Rect{X: x, Y: y, W: int32(ws), H: int32(hs)}
	e.Renderer.Copy(text_texture, nil, &text_rect)

	return nil
}
