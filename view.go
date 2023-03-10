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
  Children      []*Rect
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

func (r *Rect) contains(x, y float64) bool {
  return x >= r.TL.X &&
    x <= r.BR.X &&
    y <= r.TL.Y &&
    y >= r.BR.Y
}

type Direction int

const (
  DIR_HORIZONTAL Direction = iota
  DIR_VERTICAL
)

func (r *Rect) Split(dir Direction, ratio float64) (*Rect, *Rect) {
  var first, second *Rect
  switch dir {
  case DIR_HORIZONTAL:
    //  TL--------------+
    //  |      top      |
    //  +---------------+
    //  |     bottom    |
    //  +---------------BR
    height := r.TL.Y - r.BR.Y

    // top
    first = NewRect(
      r.TL,
      pixel.V(r.BR.X, r.TL.Y-height*ratio),
    )

    // bottom
    second = NewRect(
      pixel.V(r.TL.X, r.TL.Y-height*ratio),
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

    // left
    first = NewRect(
      r.TL,
      pixel.V(r.TL.X+width*ratio, r.BR.Y),
    )

    // right
    second = NewRect(
      pixel.V(r.TL.X+width*ratio, r.TL.Y),
      r.BR,
    )
  }
  r.Children = append(r.Children, first, second)
  return first, second
}

func GetRectAtPosition(r *Rect, v pixel.Vec) *Rect {
  if r.Children == nil {
    return r
  }
  for _, c := range r.Children {
    if c.contains(v.X, v.Y) {
      return GetRectAtPosition(c, v)
    }
  }
  return nil
}

func (r *Rect) DrawFill(win *pixelgl.Window, imd *imdraw.IMDraw) {
  imd.Push(r.TL, r.BR)
  imd.Rectangle(0)
  imd.Draw(win)
}
