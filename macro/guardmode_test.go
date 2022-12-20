package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestGuardMode(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	// prepare world
	w := core.NewWorld(333, 666)
	red, _ := core.NewTank(w, core.RedTank, 11, 22, core.WeaponRockets)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)
	blue, _ := core.NewTank(w, core.BlueTank, 11, 22, core.WeaponCannon)
	blue.SetPosition(core.NewPosition(200, 100), core.North)
	w.AddTank(blue)
	blue2, _ := core.NewTank(w, core.BlueTank, 11, 22, core.WeaponCannon)
	blue2.SetPosition(core.NewPosition(100, 200), core.West)
	w.AddTank(blue2)

	// check
	if red.Health() != 100 || blue.Health() != 100 {
		t.Error("wrong value")
	}

	// test
	for i := 0; i < 900; i++ {
		GuardMode(nil) // test nil
		GuardMode(red)
		GuardMode(blue)
		GuardMode(blue2)
		w.Update()
	}

	// check
	if red.Health() == 100 || blue.Health() == 100 {
		t.Error("wrong value")
	}
}
