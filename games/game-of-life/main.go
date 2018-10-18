package main

import (
    "fmt"
    "image/color"
    "runtime"

    "github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
    appName string = "Web-GOF"
    screenWidth, screenHeight int = 640, 640

    fieldWidth, fieldHeight int = 64, 64
    colWidth = 10
    rowHeight = 10
)

func cellId(x, y int) int {
    x %= fieldWidth;
    y %= fieldHeight;
    if x < 0 { x += fieldWidth }
    if y < 0 { y += fieldHeight }
    return x + y * fieldWidth
}

type MapPreset int

const (
    EmptyMapPreset MapPreset = iota
    StableSquaresPreset
)

type Game struct {
    paused bool
    hideGrid bool

    field [fieldWidth*fieldHeight]bool
    tempField [fieldWidth*fieldHeight]bool
}

func NewGame() *Game {
    g := &Game{paused: true, hideGrid: false}
    g.clearMap(StableSquaresPreset)

    g.field[cellId(3, 2)] = true

    return g
}

func (g *Game) Update(screen *ebiten.Image) error {
    g.updateField()

    if ebiten.IsDrawingSkipped() {
        return nil
    }

    screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

    g.drawField(screen)

    ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
    return nil
}

func (g *Game) updateField() {
    for col := 0; col < fieldWidth; col++ {
        for row := 0; row < fieldHeight; row++ {
            nc := g.countNeighbors(col, row)

            if g.field[col + row * fieldWidth] {
                g.tempField[col + row * fieldWidth] = (nc == 2 || nc == 3)
            } else {
                g.tempField[col + row * fieldWidth] = (nc == 3)
            }
        }
    }

    for i := 0; i < len(g.field); i++ {
        g.field[i] = g.tempField[i]
    }
}

func (g *Game) countNeighbors(x, y int) (count uint) {
    if g.field[cellId(x, y + 1)] { count += 1 } //Top
    if g.field[cellId(x, y - 1)] { count += 1 } //Down
    if g.field[cellId(x + 1, y)] { count += 1 } //Right
    if g.field[cellId(x - 1, y)] { count += 1 } //Left

    if g.field[cellId(x + 1, y + 1)] { count += 1 } //Top right
    if g.field[cellId(x - 1, y + 1)] { count += 1 } //Top left
    if g.field[cellId(x + 1, y - 1)] { count += 1 } //Down right
    if g.field[cellId(x - 1, y - 1)] { count += 1 } //Down left

    return
}

func (g *Game) drawField(screen *ebiten.Image) {
    pixels := make([]byte, screenWidth*screenHeight*4)
    for col := 0; col < fieldWidth; col++ {
        for row := 0; row < fieldHeight; row++ {
            if g.field[col + row * fieldWidth] {
                for y := 0; y < rowHeight; y++ {
                    for x := 0; x < colWidth; x++ {
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 0] = 0
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 1] = 0
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 2] = 0
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 3] = 0
                    }
                }
            } else {
                for y := 0; y < rowHeight; y++ {
                    for x := 0; x < colWidth; x++ {
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 0] = 255
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 1] = 255
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 2] = 255
                        pixels[4*((row*rowHeight+y)*screenWidth + (col*colWidth+x)) + 3] = 255
                    }
                }
            }
        }
    }
    screen.ReplacePixels(pixels)
}

func (g *Game) clearMap(preset MapPreset) {
    switch preset {
    case EmptyMapPreset:
        g.setEmptyPreset()
    case StableSquaresPreset:
        g.setStableSquaresPreset()
    default:
        panic("Reached default")
    }
}

func (g *Game) setEmptyPreset() {
    for i := 0; i < len(g.field); i++ {
        g.field[i] = false
    }
}

func (g *Game) setStableSquaresPreset() {
    for col := 0; col < fieldWidth; col++ {
        for row := 0; row < fieldHeight; row++ {
            g.field[col + row * fieldWidth] = !(col % 3 == 0 || row % 3 == 0)
        }
    }
}

func main() {
    game := NewGame()

    if runtime.GOARCH == "js" || runtime.GOOS == "js" {
        ebiten.SetFullscreen(true);
    }

    if err := ebiten.Run(game.Update, screenWidth, screenHeight, 1, appName); err != nil {
        panic(err)
    }
}
