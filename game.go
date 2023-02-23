package gameoflife

import(
    "time"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
)

const (
    State_Menu int = iota
    State_Game
)

const width = 1920
const height = 1080

const res = 8
const board_width = width / res
const board_height = height / res

func Run() {
    win := initWindow("Game of life (in GO)", width, height)

    board := CreateEmptyBoard(board_width, board_height)
    board.InitializeRandom(0.2)

    state := State_Menu

    for !win.Closed() {
        win.Clear(colornames.Skyblue)

        switch (state) {
            case State_Menu: {
                pos := win.MousePosition()
                hovering := pos.X > 200 && pos.X < width - 200 && pos.Y > 150 && pos.Y < 250

                drawTitleScreen(win, hovering)

                if win.JustPressed(pixelgl.MouseButtonLeft) {
                    if hovering {
                        state = State_Game
                    }
                }
            }

            case State_Game: {
                drawBoard(board, win)
                board.NextGen()
                time.Sleep(50 * time.Millisecond)
            }
        }

        win.Update()
        if win.Pressed(pixelgl.KeyQ) {
            return
        }
    }
}

func initWindow(title string, width float64, height float64) *pixelgl.Window {
    cfg := pixelgl.WindowConfig{
        Title: title,
        Bounds: pixel.R(0, 0, width, height),
        VSync:  true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }
    return win
}
