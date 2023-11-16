package shapes

import "github.com/veandco/go-sdl2/sdl"

type AbrRenderer struct {
	*sdl.Renderer
}

func (renderer *AbrRenderer) DrawShape(rect *RectangleShape) error {
	if err := renderer.SetDrawColor(rect.Color.R, rect.Color.G, rect.Color.B, rect.Color.A); err != nil {
		return err
	}
	if err := renderer.DrawRect(&rect.Rect); err != nil {
		return err
	}
	return nil
}
