package gameoflife

import (
  "time"

  "github.com/faiface/pixel"
  "github.com/faiface/pixel/imdraw"
  "github.com/faiface/pixel/pixelgl"
  "golang.org/x/image/colornames"
)

const FRAME_RATE = 30

type BoardDrawer struct {
  Game               *Game
  Width, Height, Res float64
}

func (bd *BoardDrawer) Run() {
  cfg := pixelgl.WindowConfig{
    Title:  "Game of life (in GO)",
    Bounds: pixel.R(0, 0, bd.Width, bd.Height),
    VSync:  true,
  }
  win, err := pixelgl.NewWindow(cfg)
  if err != nil {
    panic(err)
  }

  for !win.Closed() {
    // TODO: Implement click handling for splits
    // bd.handleClick(win)
    win.Clear(colornames.Skyblue)
    mainWindow := Rect{TL: pixel.V(0, bd.Height), BR: pixel.V(bd.Width, 0)}
    top, bottom := mainWindow.Split(DIR_HORIZONTAL, 0.7)
    bd.drawBoardToRect(win, top)

    imd := imdraw.New(nil)
    imd.Color = colornames.Salmon
    bottom.DrawFill(win, imd)
    win.Update()
    time.Sleep(1000 / FRAME_RATE * time.Millisecond)
  }
}

func (bd *BoardDrawer) handleClick(win *pixelgl.Window) {
  if win.JustPressed(pixelgl.MouseButtonLeft) {
    mouse_pos := win.MousePosition()

    var change Change
    change.X = int(mouse_pos.X / bd.Res)
    change.Y = int(mouse_pos.Y / bd.Res)
    change.Alive = !bd.Game.Board.IsAlive(change.X, change.Y)
    change.Gen = bd.Game.Board.Gen

    bd.Game.Changes <- change

    changes := make([]Change, 1)
    changes[0] = change
    go bd.Game.Client.SendChanges(changes)
  }
}

func (bd *BoardDrawer) drawBoard(win *pixelgl.Window) {
  imd := imdraw.New(nil)
  imd.Color = colornames.Black

  for y := 0; y < bd.Game.Board.Height; y++ {
    for x := 0; x < bd.Game.Board.Width; x++ {
      if bd.Game.Board.IsAlive(x, y) {
        imd.Push(pixel.V(float64(x)*bd.Res, float64(y)*bd.Res), pixel.V(float64(x+1)*bd.Res, float64(y+1)*bd.Res))
        imd.Rectangle(0)
      }
    }
  }

  imd.Draw(win)
}

func (bd *BoardDrawer) drawBoardToRect(win *pixelgl.Window, rect *Rect) {
  border := NewRect(pixel.V(0, bd.Height), pixel.V(bd.Width, 0))
  scale := rect.getScalingFactor(border)

  // TODO: fix scaling
  newScaleX := (bd.Width * scale) / float64(bd.Game.Board.Width)
  newScaleY := (bd.Height * scale) / float64(bd.Game.Board.Height)

  imd := imdraw.New(nil)
  imd.Color = colornames.Black

  for y := 0; y < bd.Game.Board.Height; y++ {
    for x := 0; x < bd.Game.Board.Width; x++ {
      if bd.Game.Board.IsAlive(x, y) {
        px := rect.TL.X + float64(x)*newScaleX
        py := rect.TL.Y - float64(y)*newScaleY
        px2 := rect.TL.X + float64(x+1)*newScaleX
        py2 := rect.TL.Y - float64(y+1)*newScaleY

        imd.Push(pixel.V(px, py), pixel.V(px2, py2))
        imd.Rectangle(0)
      }
    }
  }

  imd.Draw(win)
}
