package canvas

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

// DrawString - simple way, requires delimiter "\n" for a new line.
func (tk *ToolKit) DrawString(message string, indent int, col color.Color, face *basicfont.Face) error {
	if face == nil {
		face = cyrfont.Scaled9x15(1)
	}

	face.Advance = face.Advance + 3

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
			X: fixed.Int26_6(64 * indent),
			Y: fixed.Int26_6(face.Advance * 64 * row),
		}

		d.Dot = point
		d.DrawString(line)

		row++
	}

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), img, image.Point{}, draw.Over)

	return tk.Canvas.Render()
}
