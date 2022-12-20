package core

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

// World is the game and holds all active objects on the map.
type World struct {
	xWidth  int // world dimension X (with 64*64 blocks)
	yHeight int // world dimension Y (with 64*64 blocks)

	iteration   uint64
	tanks       []*Tank
	projectiles []*Projectile

	freeze   bool    // disable the Update() routine if true
	cashRed  float64 // is increased by Update() as long as a red base exists
	cashBlue float64 // is increased by Update() as long as a blue base exists
}

// NewWorld create a new world.
// The attributes XWidth and YHeight are blocks (64x64).
// see WorldXWidth and WorldYHeight.
func NewWorld(XWidth, YHeight int) *World {
	return &World{
		xWidth:  XWidth,
		yHeight: YHeight,

		tanks:       make([]*Tank, 0),
		projectiles: make([]*Projectile, 0),
	}
}

//---------------- GETTER --------------------------------------------------------------------------------------------//

// IsFrozen returns the freeze state.
// see Freeze()
func (w *World) IsFrozen() bool {
	return w.freeze
}

// XWidth return the block width
func (w *World) XWidth() int {
	return w.xWidth
}

// YHeight return the block height
func (w *World) YHeight() int {
	return w.yHeight
}

// Iteration returns the current calculation round.
func (w *World) Iteration() uint64 {
	return w.iteration
}

// ScreenWidth is the GUI width (= XWidth * 64)
func (w *World) ScreenWidth() int {
	return w.xWidth * BlockSize
}

// ScreenHeight is the GUI width (= YHeight * 64)
func (w *World) ScreenHeight() int {
	return w.yHeight * BlockSize
}

// Tanks returns the tank list (= all destructible objects)
func (w *World) Tanks() []*Tank {
	return w.tanks
}

// Projectiles returns all flying bullets.
func (w *World) Projectiles() []*Projectile {
	return w.projectiles
}

// CashStat returns the current cash amount of the players.
// The amount is slowly generated by the home base.
//
// see Update().
func (w *World) CashStat() (cashRed, cashBlue int) {
	return int(w.cashRed), int(w.cashBlue)
}

// UnitCount returns the sum of all units.
// This number is used for the victory condition.
func (w *World) UnitCount() (red, blue int) {
	for _, t := range w.tanks {
		if t != nil && (t.owner == RedTank || t.owner == RedBase) {
			red++
		} else if t != nil && (t.owner == BlueTank || t.owner == BlueBase) {
			blue++
		}
	}
	return
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// Freeze disable the Update() routine if true.
func (w *World) Freeze(status bool) {
	w.freeze = status
}

// AddTank adds a new tank to the world.
// Use Tank.SetPosition() to set the correct position.
func (w *World) AddTank(tank *Tank) {
	w.tanks = append(w.tanks, tank)
}

// SetCash overwrites the current value of both players.
func (w *World) SetCash(cashRed, cashBlue int) {
	w.cashRed, w.cashBlue = float64(cashRed), float64(cashBlue)
}

// BuyTank buy a tank and place it near the home base.
func (w *World) BuyTank(tank *Tank) error {

	// pay credits
	// (and get base string)
	//-----------------------
	var base string
	if tank != nil && tank.owner == RedTank && w.cashRed >= TankBudget {
		w.cashRed -= TankBudget
		base = RedBase
	} else if tank != nil && tank.owner == BlueTank && w.cashBlue >= TankBudget {
		w.cashBlue -= TankBudget
		base = BlueBase
	} else {
		// ERROR EXIT
		tank.Remove() // may be nil
		return errors.New("not enough tank budget or unknown owner")
	}

	// find home base pos
	//--------------------
	var homePos = NewPosition(-1000000, -1000000)
	for _, o := range w.tanks {
		if o != nil && o.owner == base {
			homePos = o.Pos()
			break
		}
	}
	if homePos.X < -100 && homePos.Y < -100 {
		// ERROR EXIT
		tank.Remove()
		return errors.New("home base not found")
	}

	// random spawn new tank
	//-----------------------
	rand.Seed(time.Now().UnixNano())
	for try := 1; try < 1000; try++ {

		// generate random spawn point
		rndX := homePos.X + rand.Intn(4*BlockSize) - 2*BlockSize
		rndY := homePos.Y + rand.Intn(4*BlockSize) - 2*BlockSize
		spawnPos := NewPosition(rndX, rndY)

		// check collisions
		tank.SetPosition(spawnPos, South)
		tank.Forward()
		tank.Update()

		if tank.Blocked() {
			// ERROR: retry
			continue
		} else {
			// success
			tank.Stop()
			w.AddTank(tank)
			return nil // success EXIT
		}
	}

	// ERROR: can't spawn tank
	tank.Remove()
	return errors.New("not enough space to spawn")
}

// Clear all tanks (objects) with prefix from world.
// Kill a player like 'red' or 'blue'.
func (w *World) Clear(prefix string) {
	list := make([]*Tank, 0, len(w.tanks))

	for _, t := range w.tanks {
		if !strings.HasPrefix(t.owner, prefix) {
			list = append(list, t)
		} else {
			t.health = 0 // kill object
		}
	}

	w.tanks = list
}

//---------------- UPDATE --------------------------------------------------------------------------------------------//

// UpdateN calls Update() n-times.
func (w *World) UpdateN(n int) {
	for i := 0; i < n; i++ {
		w.Update()
	}
}

// Update is called 30 times (see GameSpeed) per second.
// The method also calls Update() of all tanks and all projectiles.
func (w *World) Update() {
	if w.freeze {
		return // no updates
	}

	// update all tanks
	for _, t := range w.tanks {
		t.Update()
	}

	// update all projectiles
	for _, p := range w.projectiles {
		p.Update()
	}

	// increased cash for player
	for _, t := range w.tanks {
		if t != nil && t.owner == RedBase {
			w.cashRed += 100.0 / 3600.0
		} else if t != nil && t.owner == BlueBase {
			w.cashBlue += 100.0 / 3600.0
		}
	}

	// finish this iteration
	w.iteration++
}
