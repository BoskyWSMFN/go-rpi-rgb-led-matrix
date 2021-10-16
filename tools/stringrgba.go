package tools

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/enotofil/cyrfont"
)

// DrawString - simple way, requires delimiter "\n" for new line.
func (tk *ToolKit) DrawString(message string, col color.Color, face *basicfont.Face) error {
	if face == nil {
		face = cyrfont.Face9x15
		face.Advance = 12
	}

	row := 1
	toWrite := strings.Split(message, "\n")
	img := image.NewRGBA(tk.Canvas.Bounds())
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
	}

	for _, line := range toWrite {
		point := fixed.Point26_6{
			X: fixed.Int26_6(64),
			Y: fixed.Int26_6(face.Advance * 64 * row),
		}

		d.Dot = point
		d.DrawString(line)

		row++
	}

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), img, image.Point{}, draw.Src)

	return tk.Canvas.Render()
}
