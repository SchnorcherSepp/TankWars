package gui

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/remote"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strconv"
)

// Update updates an img by one tick. The given argument represents a screen image.
//
// Update updates only the img logic and Draw draws the screen.
//
// In the first frame, it is ensured that Update is called at least once before Draw. You can use Update
// to initialize the img state.
//
// After the first frame, Update might not be called or might be called once
// or more for one frame. The frequency is determined by the current TPS (tick-per-second).
func (g *Game) Update() error {

	// switch inputs
	switch g.inputMode {
	case "shop":
		controlsShop(g)
	default:
		controlsGame(g)
	}

	// UPDATE and RETURN
	g.world.UpdateN(g.speed)
	return nil
}

//--------------------------------------------------------------------------------------------------------------------//

var shopState int
var shopType string // 'C', 'R' or 'A'
var shopArmor int   // 1, 2, 3, ...
var shopDamage int  // 1, 2, 3, ...
var shopErr error

// controlsGame provides the shop controls
func controlsShop(g *Game) {

	switch shopState {
	case 0: // select type
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			shopType = core.WeaponCannon
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			shopType = core.WeaponRockets
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			shopType = core.WeaponArtillery
			shopState++
		}
	case 1: // select armor
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			shopArmor = 5
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
			shopArmor = 15
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
			shopArmor = 25
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
			shopArmor = 35
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
			shopArmor = 45
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
			shopArmor = 55
			shopState++
		}
	case 2: // select damage
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			shopDamage = 15
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
			shopDamage = 25
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
			shopDamage = 35
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
			shopDamage = 45
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
			shopDamage = 55
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
			shopDamage = 65
			shopState++
		} else if inpututil.IsKeyJustPressed(ebiten.Key7) {
			shopDamage = 70
			shopState++
		}
	case 3: // buy
		// find owner
		var owner = core.BlueTank
		if g.activeTank != nil && g.activeTank.Owner() == core.RedBase {
			owner = core.RedTank
		}
		// build tank
		var tank *core.Tank
		tank, shopErr = core.NewTank(g.world, owner, shopArmor, shopDamage, shopType)
		if shopErr == nil {
			shopErr = g.world.BuyTank(tank)
		}
		// fin
		shopState++
		if shopErr == nil {
			g.inputMode = "" // close shop
			shopState = 0
		}
	}

	// toggle KEY: close shop
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		// switch input
		g.inputMode = ""
		// reset shop
		shopState = 0
	}
}

//--------------------------------------------------------------------------------------------------------------------//

// controlsGame provides the default game controls
func controlsGame(g *Game) {

	// LEFT: select active tank
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		for _, t := range g.world.Tanks() {
			if core.IsCollided(t.Pos(), core.BlockRadius, core.NewPosition(cx, cy), core.BallRadius) {
				g.activeTank = t
				break
			}
			g.activeTank = nil
		}
	}

	// RIGHT: fire at
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.activeTank != nil {
		cx, cy := ebiten.CursorPosition()
		pos := core.NewPosition(cx, cy)
		g.activeTank.Stop()
		g.activeTank.FireAt(pos)
	}

	// KEY: send command
	if g.activeTank != nil && (g.activeTank.Owner() == core.RedTank || g.activeTank.Owner() == core.BlueTank) {
		// movement
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.activeTank.Left()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.activeTank.Right()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.activeTank.Forward()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.activeTank.Backward()
		}
		if ebiten.IsKeyPressed(ebiten.KeyControlRight) {
			g.activeTank.Stop()
		}

		// macros
		if ebiten.IsKeyPressed(ebiten.KeyQ) { // remove
			remote.SetMacro(g.world, "", g.activeTank.ID(), core.MacroReset)
		}
		if ebiten.IsKeyPressed(ebiten.Key1) { // GuardMode
			remote.SetMacro(g.world, "", g.activeTank.ID(), core.MacroGuardMode)
		}
		if ebiten.IsKeyPressed(ebiten.Key2) { // FireAndManeuver
			remote.SetMacro(g.world, "", g.activeTank.ID(), core.MacroFireAndManeuver)
		}
		if ebiten.IsKeyPressed(ebiten.Key3) { // AttackMove
			remote.SetMacro(g.world, "", g.activeTank.ID(), core.MacroAttackMove)
		}
		if ebiten.IsKeyPressed(ebiten.Key4) { // FireWall
			remote.SetMacro(g.world, "", g.activeTank.ID(), core.MacroFireWall)
		}
		if ebiten.IsKeyPressed(ebiten.Key5) { // FireWall
			cx, cy := ebiten.CursorPosition()
			g.activeTank.Forward()
			remote.SetMacroMoveTo(g.world, "", g.activeTank.ID(), strconv.Itoa(cx), strconv.Itoa(cy))
		}
	}

	// toggle KEY R: range circles
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.rangeCircles = !g.rangeCircles
	}
	// toggle KEY D: debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}
	// toggle KEY: help
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.help = !g.help
	}
	// toggle KEY: open shop
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.activeTank != nil && (g.activeTank.Owner() == core.RedBase || g.activeTank.Owner() == core.BlueBase) {
			g.inputMode = "shop" // switch input
		}
	}
}
