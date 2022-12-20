package gui

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// drawBuilding draw all buildings like RedBase, BlueBase, NeutralRock, RedRock, BlueRock, ...
func drawBuilding(screen *ebiten.Image, t *core.Tank, building string, active bool, world *core.World) {
	xf := t.Pos().Xf
	yf := t.Pos().Yf
	cashRed, cashBlue := world.CashStat()

	// prepare image
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(xf-core.BlockRadius, yf-core.BlockRadius) // basic image 64x64
	op.Filter = ebiten.FilterLinear                             // Specify linear filter.

	// select image
	switch building {
	case core.NeutralRock, core.RedRock, core.BlueRock:
		// ROCKS
		screen.DrawImage(resources.Imgs.Rock, op)
	case core.RedBase:
		// BASE 1 (red)
		screen.DrawImage(resources.Imgs.Base1, op)
		writeCash(screen, xf, yf, cashRed)
	case core.BlueBase:
		// BASE 2 (blue)
		screen.DrawImage(resources.Imgs.Base2, op)
		writeCash(screen, xf, yf, cashBlue)
	default:
		// ERROR!!!
		screen.DrawImage(resources.Imgs.Error, op) // ERROR
	}

	// draw active building
	if active {
		ebitenutil.DrawCircle(screen, xf, yf, core.BlockRadius*1.1, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x33})
	}

	// draw text (CENTER)
	writeHealth(screen, xf, yf, t.Health())
}

// writeCash is for RedBase or BlueBase and write the $ cash.
func writeCash(screen *ebiten.Image, xf, yf float64, cash int) {
	txt := fmt.Sprintf("$%d", cash)
	xf = xf - (6 / 2 * float64(len(txt)))
	yf = yf + 15
	ebitenutil.DebugPrintAt(screen, txt, int(xf), int(yf))
}
