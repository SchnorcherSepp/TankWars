package gui

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

// interface check: ebiten.Game
var _ ebiten.Game = (*Game)(nil)

// Game is the GUI
type Game struct {
	world        *core.World
	speed        int
	xWidth       int
	yHeight      int
	screenWidth  int
	screenHeight int

	activeTank *core.Tank
	inputMode  string

	// toggle
	rangeCircles bool
	help         bool
	debug        bool
}

// RunGame starts a GUI window and displays the specified world.
// The callUpdate option activates the core.WorldMap Update() call with 60 Ticks per second.
// Do not activate this option if the update is done externally.
//
// This call is blocking.
func RunGame(title string, world *core.World, speed int, mute bool) error {
	resources.MuteSound = mute

	// config img
	game := &Game{
		world:        world,
		speed:        speed,
		xWidth:       world.XWidth(),       // world dimension X
		yHeight:      world.YHeight(),      // world dimension Y
		screenWidth:  world.ScreenWidth(),  // basic image 64x64
		screenHeight: world.ScreenHeight(), // basic image 64x64
	}

	// config window
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowIcon([]image.Image{resources.Imgs.Logo})
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(core.GameSpeed) // default: 60 ticks per second

	// run (BLOCKING)
	return ebiten.RunGame(game)
}

//--------------------------------------------------------------------------------------------------------------------//

// Layout accepts a native outside size in device-independent pixels and returns the img logical screen
// size.
//
// On desktops, the outside is a window or a monitor (fullscreen mode). On browsers, the outside is a body
// element. On mobiles, the outside is the view's size.
//
// Even though the outside size and the screen size differ, the rendering scale is automatically adjusted to
// fit with the outside.
//
// Layout is called almost every frame.
//
// It is ensured that Layout is invoked before Update is called in the first frame.
//
// If Layout returns non-positive numbers, the caller can panic.
//
// You can return a fixed screen size if you don't care, or you can also return a calculated screen size
// adjusted with the given outside size.
func (g *Game) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}

//---------------- DRAW ----------------------------------------------------------------------------------------------//

// Draw draws the img screen by one frame.
//
// The give argument represents a screen image. The updated content is adopted as the img screen.
func (g *Game) Draw(screen *ebiten.Image) {

	// draw background
	drawBackground(screen, g.screenWidth, g.screenHeight)

	// draw ALL buildings or tanks
	for _, t := range g.world.Tanks() {
		switch t.Owner() {
		case core.NeutralRock, core.RedRock, core.BlueRock, core.RedBase, core.BlueBase:
			// BUILDING
			drawBuilding(screen, t, t.Owner(), t == g.activeTank, g.world)
		default:
			// TANK (default)
			drawTank(screen, t, t.Owner(), t == g.activeTank, g.rangeCircles, g.debug)
		}
	}

	// draw ALL projectiles
	for _, p := range g.world.Projectiles() {
		drawProjectiles(screen, p)
	}

	// write global text
	writeGlobalText(screen, g)

	// draw shop
	if g.inputMode == "shop" {
		drawShop(screen, g.screenWidth, g.screenHeight, g.activeTank, g.world)
	}

	// draw victory
	scoreRed, scoreBlue := g.world.UnitCount()
	drawVictory(screen, g.screenWidth, g.screenHeight, scoreRed, scoreBlue)
}

// drawBackground fills the entire background with grass. Call this drawing first.
func drawBackground(screen *ebiten.Image, screenWidth, screenHeight int) {
	xWidth := screenWidth / core.BlockSize
	yHeight := screenHeight / core.BlockSize

	for xCol := 0; xCol < xWidth; xCol++ {
		for yRow := 0; yRow < yHeight; yRow++ {
			op := new(ebiten.DrawImageOptions)
			op.GeoM.Translate(float64(xCol*core.BlockSize), float64(yRow*core.BlockSize)) // basic image 64x64
			op.Filter = ebiten.FilterLinear                                               // Specify linear filter.
			screen.DrawImage(resources.Imgs.Grass, op)
		}
	}
}

// drawVictory draw the victory image
func drawVictory(screen *ebiten.Image, screenWidth, screenHeight, scoreRed, scoreBlue int) {
	// image: 285 x 120
	x := float64(screenWidth-285) / 2
	y := float64(screenHeight-120) / 2

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(x, y)
	op.Filter = ebiten.FilterLinear

	if scoreRed == 0 && scoreBlue == 0 {
		screen.DrawImage(resources.Imgs.VictoryDraw, op)
	} else if scoreRed == 0 {
		screen.DrawImage(resources.Imgs.VictoryBlue, op)
	} else if scoreBlue == 0 {
		screen.DrawImage(resources.Imgs.VictoryRed, op)
	}
}

// writeGlobalText write the global text top left.
func writeGlobalText(screen *ebiten.Image, g *Game) {
	s := "\n"

	if g.help {
		s += "  - Right click to select\n"
		s += "  - Left click to fire at pos\n"
		s += "  - Arrow keys to navigate\n"
		s += "  - Ctrl right key to stop\n"
		s += "  - 'R' key toggle range view\n"
		s += "  - 'D' key toggle debug mode\n"
		s += "  - 'S' open shop (selected base)\n"
		s += "  - 'Q' remove macro\n"
		s += "  - '1' set GuardMode macro\n"
		s += "  - '2' set FireAndManeuver macro\n"
		s += "  - '3' set AttackMove macro\n"
		s += "  - '4' set FireWall macro\n"
		s += "  - '5' set MoveTo(cursor) macro\n"
		s += "\n"
	} else {
		s += "   Press 'H' for help\n"
		s += "\n"
	}

	s += "   player units: %d red; %d blue"
	if g.debug {
		s += fmt.Sprintf("\n   round: %d", g.world.Iteration())
	}

	r, b := g.world.UnitCount()
	txt := fmt.Sprintf(s, r, b)
	ebitenutil.DebugPrint(screen, txt)
}
