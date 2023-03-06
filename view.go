package gameoflife

import (
  "math"

  "github.com/faiface/pixel"
  "github.com/faiface/pixel/imdraw"
  "github.com/faiface/pixel/pixelgl"
)

type Rect struct {
  //  TL--+
  //  |   |
  //  +---BR
  TL, BR        pixel.Vec
  Width, Height float64
}

func NewRect(tl, br pixel.Vec) *Rect {
  r := Rect{}
  r.TL = tl
  r.BR = br
  r.Width = br.X - tl.X
  r.Height = tl.Y - br.Y
  return &r
}

// Get factor by which to scale 'other' to make it fit into
// this rectangle, preserving its aspect ratio
// FIXME: Scaling is broken in some weird aspect ratios
func (r *Rect) getScalingFactor(other *Rect) float64 {
  if other.Width <= r.Width && other.Height <= r.Height {
    return 1.0
  }

  scale := math.Min(r.Width/other.Width, r.Height/other.Height)
  return scale
}

type Direction int

const (
  DIR_HORIZONTAL Direction = iota
  DIR_VERTICAL
)

type Split interface {
  Draw(win *pixelgl.Window)
  Split(dir Direction)
}

func (r *Rect) Split(dir Direction) (*Rect, *Rect) {
  var first, second *Rect
  switch dir {
  case DIR_HORIZONTAL:
    //  TL--------------+
    //  |      top      |
    //  +---------------+
    //  |     bottom    |
    //  +---------------BR
    height := r.TL.Y - r.BR.Y
    first = NewRect(
      r.TL,
      pixel.V(r.BR.X, r.TL.Y-height/2),
    )
    second = NewRect(
      pixel.V(r.TL.X, r.TL.Y-height/2),
      r.BR,
    )

  case DIR_VERTICAL:
    //  TL------+-------+
    //  |   l   |   r   |
    //  |   e   |   i   |
    //  |   f   |   g   |
    //  |   t   |   h   |
    //  |       |   t   |
    //  +-------+-------BR
    width := r.BR.X - r.TL.X
    first = NewRect(
      r.TL,
      pixel.V(r.TL.X+width/2, r.BR.Y),
    )
    second = NewRect(
      pixel.V(r.TL.X+width/2, r.TL.Y),
      r.BR,
    )
  }
  return first, second
}

func (r *Rect) DrawFill(win *pixelgl.Window, imd *imdraw.IMDraw) {
  imd.Push(r.TL, r.BR)
  imd.Rectangle(0)
  imd.Draw(win)
}
