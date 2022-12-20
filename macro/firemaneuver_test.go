package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestFireAndManeuver(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	// prepare world
	w := core.NewWorld(333, 666)
	red, _ := core.NewTank(w, core.RedTank, 55, 20, core.WeaponCannon)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)
	blue, _ := core.NewTank(w, core.BlueTank, 55, 20, core.WeaponCannon)
	blue.SetPosition(core.NewPosition(600, 100), core.West)
	w.AddTank(blue)

	// check
	if red.Health() != 100 || blue.Health() != 100 { // fire
		t.Error("wrong value")
	}
	if red.Pos().X != 100 || blue.Pos().X != 600 { // move
		t.Error("wrong value", red.Pos().Xf, blue.Pos().Xf)
	}

	// test
	for i := 0; i < 900; i++ {
		FireAndManeuver(nil) // test nil
		FireAndManeuver(red)
		FireAndManeuver(blue)
		w.Update()
	}

	// check
	if red.Health() == 100 || blue.Health() == 100 { // fire
		t.Error("wrong value", red.Health(), blue.Health())
	}
	if red.Pos().X == 100 || blue.Pos().X == 600 { // move
		t.Error("wrong value", red.Pos().Xf, blue.Pos().Xf)
	}
}
