package gameoflife

import (
        "fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
        "golang.org/x/image/font/basicfont"
)

func drawRect(win *pixelgl.Window,
          color pixel.RGBA,
          x1 float64, y1 float64,
          x2 float64, y2 float64,
          thickness int) {
      imd := imdraw.New(nil)
      imd.Color = color
      imd.Push(pixel.V(x1, y1))
      imd.Push(pixel.V(x2, y2))
      imd.Rectangle(0)
      imd.Draw(win)
}

func drawBoard(board *Board, win *pixelgl.Window) {
    imd := imdraw.New(nil)
    imd.Color = colornames.Black

    for y, row := range board.Board() {
        for x := range row {
            if board.IsAlive(x, y) {
                imd.Push(pixel.V(float64(x * res), float64(y * res)), pixel.V(float64((x + 1) * res), float64((y + 1) * res)))
                imd.Rectangle(0)
            }
        }
    }

    imd.Draw(win)
}

func drawTitleScreen(win *pixelgl.Window, drawShadow bool) {
      const wCenter = width / 2
      const hCenter = height / 2

      atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
      textScaling := 5

      title := "GOme of Life"
      textCenterW := (len(title) * 7 * textScaling) / 2
      textCenterH := (13 * textScaling) / 2
      titleText := text.New(pixel.V(
        float64(wCenter -  textCenterW),
        float64(hCenter - (textCenterH) + (height * 0.3))), atlas)
        fmt.Fprintln(titleText, title)

      shadowText := text.New(pixel.V(
        float64((wCenter -  textCenterW) + 3),
        float64((hCenter - (textCenterH) + (height * 0.3))) - 3), atlas)
        shadowText.Color = colornames.Black
        fmt.Fprintln(shadowText, title)

        shadowText.Draw(win, pixel.IM.Scaled(shadowText.Orig, float64(textScaling)))
        titleText.Draw(win, pixel.IM.Scaled(titleText.Orig, float64(textScaling)))

      if drawShadow {
        ofs := 5.0
        drawRect(win, pixel.RGB(0.4, 0.4, 0.4), 200 + ofs, 150 - ofs, (width - 200) + ofs, 250 - ofs, 0)
      }
      drawRect(win, pixel.RGB(0.5, 0.5, 0.5), 200, 150, width - 200, 250, 0)

      buttonLabel := "Play"
      buttonCenterW := (len(buttonLabel) * 7 * textScaling) / 2
      buttonCenterH := (13 * textScaling) / 2
      buttonText := text.New(pixel.V(
        float64(wCenter - buttonCenterW),
        float64(215 - buttonCenterH)), atlas)
        fmt.Fprintln(buttonText, buttonLabel)
        buttonText.Draw(win, pixel.IM.Scaled(buttonText.Orig, float64(textScaling)))
}
