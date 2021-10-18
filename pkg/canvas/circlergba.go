package canvas

import (
	"image"
	"image/color"
	"image/draw"
)

type circle struct {
	p image.Point
	r int
	c color.Color
}

func (c *circle) RGBA() (r, g, b, a uint32) {
	return c.c.RGBA()
}

func (c *circle) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return c.c
	}
	return color.RGBA{}
}

// DrawCircle - simple way.
func (tk *ToolKit) DrawCircle(x, y, r int, col color.Color) error {
	p := image.Point{X: x, Y: y}
	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), &circle{p, r, col}, image.Point{}, draw.Src)

	return tk.Canvas.Render()
}
