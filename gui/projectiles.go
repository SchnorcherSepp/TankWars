package gui

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

// drawProjectiles draw all projectiles and the explosions.
func drawProjectiles(screen *ebiten.Image, p *core.Projectile) {
	x := p.Pos().Xf
	y := p.Pos().Yf
	start := p.StartPos()
	end := p.EndPos()

	// draw line
	ebitenutil.DrawLine(screen, start.Xf, start.Yf, x, y, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x55})
	ebitenutil.DrawLine(screen, x, y, end.Xf, end.Yf, color.RGBA{R: 0xff, G: 0, B: 0, A: 0x44})

	// draw ball
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(x-core.BallRadius, y-core.BallRadius) // ball image 20x19
	op.Filter = ebiten.FilterLinear                         // Specify linear filter.
	screen.DrawImage(resources.Imgs.Ball, op)

	// draw explosion
	if p.Exploded() {
		// draw circle
		ebitenutil.DrawCircle(screen, x, y, float64(p.AoERadius()), color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x55})
		// draw image
		op := new(ebiten.DrawImageOptions)
		op.GeoM.Translate(x-core.BlockRadius, y-core.BlockRadius) // basic image 64x64
		op.Filter = ebiten.FilterLinear                           // Specify linear filter.
		screen.DrawImage(resources.Imgs.Explosion, op)
	}
}
