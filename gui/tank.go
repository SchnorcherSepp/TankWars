package gui

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"strings"
)

// drawTank draw all tanks (RedTank and BlueTank).
func drawTank(screen *ebiten.Image, t *core.Tank, owner string, active, rangeCircles, debug bool) {
	xf := t.Pos().Xf
	yf := t.Pos().Yf

	// prepare image
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(xf-core.BlockRadius, yf-core.BlockRadius) // basic image 64x64
	op.Filter = ebiten.FilterLinear                             // Specify linear filter.

	// select image
	if strings.HasPrefix(owner, core.RedTank) {
		switch t.Angle() {
		case core.North:
			screen.DrawImage(resources.Imgs.Tank1North, op)
		case core.Northeast:
			screen.DrawImage(resources.Imgs.Tank1Northeast, op)
		case core.East:
			screen.DrawImage(resources.Imgs.Tank1East, op)
		case core.Southeast:
			screen.DrawImage(resources.Imgs.Tank1Southeast, op)
		case core.South:
			screen.DrawImage(resources.Imgs.Tank1South, op)
		case core.Southwest:
			screen.DrawImage(resources.Imgs.Tank1Southwest, op)
		case core.West:
			screen.DrawImage(resources.Imgs.Tank1West, op)
		case core.Northwest:
			screen.DrawImage(resources.Imgs.Tank1Northwest, op)
		default:
			// ERROR!!!
			screen.DrawImage(resources.Imgs.Error, op) // ERROR
		}
	} else {
		switch t.Angle() {
		case core.North:
			screen.DrawImage(resources.Imgs.Tank2North, op)
		case core.Northeast:
			screen.DrawImage(resources.Imgs.Tank2Northeast, op)
		case core.East:
			screen.DrawImage(resources.Imgs.Tank2East, op)
		case core.Southeast:
			screen.DrawImage(resources.Imgs.Tank2Southeast, op)
		case core.South:
			screen.DrawImage(resources.Imgs.Tank2South, op)
		case core.Southwest:
			screen.DrawImage(resources.Imgs.Tank2Southwest, op)
		case core.West:
			screen.DrawImage(resources.Imgs.Tank2West, op)
		case core.Northwest:
			screen.DrawImage(resources.Imgs.Tank2Northwest, op)
		default:
			// ERROR!!!
			screen.DrawImage(resources.Imgs.Error, op) // ERROR
		}
	}

	// draw weapon range
	if rangeCircles && t.Weapon() != nil && t.Weapon().Type() != core.WeaponNone {
		// get player color
		clr := color.RGBA{R: 0, G: 0, B: 0xff, A: 0x22} // blue
		if strings.HasPrefix(owner, core.RedTank) {
			clr = color.RGBA{R: 0xff, G: 0, B: 0, A: 0x22} // red
		}

		// draw
		drawWeaponRange(screen, xf, yf, t.Weapon().Range(), clr, t.Weapon().AnyFireAngle())
	}

	// draw active tank
	if active {
		ebitenutil.DrawCircle(screen, xf, yf, core.BlockRadius*1.1, color.RGBA{R: 0, G: 0xff, B: 0, A: 0x33})
	}

	// write health (CENTER)
	writeHealth(screen, xf, yf, t.Health())

	// write tank text
	if t.Weapon() != nil {
		writeTankLine1(screen, xf, yf, t) // write LINE 1
		writeTankLine2(screen, xf, yf, t) // write LINE 2
		if debug {
			writeTankLine3(screen, xf, yf, t) // write LINE 3
		}
	}
}

// drawWeaponRange draw the projectile range for artillery and cannons
func drawWeaponRange(screen *ebiten.Image, xf, yf float64, wRange int, clr color.Color, anyFireAngle bool) {
	if anyFireAngle {
		// simple circle
		ebitenutil.DrawCircle(screen, xf, yf, float64(wRange), clr)

	} else {
		// tanks can not fire at any angle!
		var path vector.Path
		var x = float32(xf)
		var y = float32(yf)
		var wr = float32(wRange)
		// horizontal
		path.MoveTo(x-wr, y-core.BallRadius)
		path.LineTo(x+wr, y-core.BallRadius)
		path.LineTo(x+wr, y+core.BallRadius)
		path.LineTo(x-wr, y+core.BallRadius)
		// vertical
		path.MoveTo(x-core.BallRadius, y-wr)
		path.LineTo(x+core.BallRadius, y-wr)
		path.LineTo(x+core.BallRadius, y+wr)
		path.LineTo(x-core.BallRadius, y+wr)
		// top-left to bottom-right
		path.MoveTo(x-wr*0.7-7, y-wr*0.7+7)
		path.LineTo(x-wr*0.7+7, y-wr*0.7-7)
		path.LineTo(x+wr*0.7+7, y+wr*0.7-7)
		path.LineTo(x+wr*0.7-7, y+wr*0.7+7)
		// bottom-left to top-right
		path.MoveTo(x-wr*0.7+7, y+wr*0.7+7)
		path.LineTo(x-wr*0.7-7, y+wr*0.7-7)
		path.LineTo(x+wr*0.7-7, y-wr*0.7-7)
		path.LineTo(x+wr*0.7+7, y-wr*0.7+7)

		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		drawVerticesForUtil(screen, vs, is, clr)
	}
}

// drawVerticesForUtil is a helper function for drawWeaponRange()
func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(color.White)
	whiteSubImage := whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	op := &ebiten.DrawTrianglesOptions{}
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

// writeHealth write the health number in the center of all objects (tanks and buildings)
func writeHealth(screen *ebiten.Image, xf, yf float64, health int) {
	txt := fmt.Sprintf("%d", health)
	xf = xf - (6 / 2 * float64(len(txt)))
	yf = yf - 8
	ebitenutil.DebugPrintAt(screen, txt, int(xf), int(yf))
}

// writeTankLine1 write the first line of tanks (stats).
func writeTankLine1(screen *ebiten.Image, xf, yf float64, t *core.Tank) {
	txt := fmt.Sprintf("D:%d | A:%d | S:%d", t.Weapon().Damage(), t.Armor(), t.Speed())
	xf = xf - (6 / 2 * float64(len(txt)))
	yf = yf + 30
	ebitenutil.DebugPrintAt(screen, txt, int(xf), int(yf))
}

// writeTankLine2 write the second line of tanks (type, mode, ...)
func writeTankLine2(screen *ebiten.Image, xf, yf float64, t *core.Tank) {
	// generate text
	_, status := t.Status()
	txt := fmt.Sprintf("%s: %s", t.Weapon().Type(), status)

	// add blocked
	if t.Blocked() {
		txt += " (Blocked)"
	}

	// add macro
	if t.ActiveMacro() {
		txt += " (M)"
	}

	// write text
	xf = xf - (6 / 2 * float64(len(txt)))
	yf = yf + 42
	ebitenutil.DebugPrintAt(screen, txt, int(xf), int(yf))
}

// writeTankLine3 write the third line of tanks (debug)
func writeTankLine3(screen *ebiten.Image, xf, yf float64, t *core.Tank) {
	// generate text
	txt := fmt.Sprintf("R:%d, M:%d, F:%d", t.LastRotate(), t.Weapon().LastMove(), t.Weapon().LastFire())

	// write text
	xf = xf - (6 / 2 * float64(len(txt)))
	yf = yf + 54
	ebitenutil.DebugPrintAt(screen, txt, int(xf), int(yf))
}
