package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestAttackMove(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	// prepare world
	w := core.NewWorld(333, 666)
	red, _ := core.NewTank(w, core.RedTank, 55, 20, core.WeaponCannon)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)
	blue, _ := core.NewTank(w, core.BlueTank, 55, 20, core.WeaponCannon)
	blue.SetPosition(core.NewPosition(600, 100), core.West)
	w.AddTank(blue)
	blue2, _ := core.NewTank(w, core.BlueTank, 55, 20, core.WeaponCannon)
	blue2.SetPosition(core.NewPosition(1200, 10), core.North)
	w.AddTank(blue2)

	// check
	if red.Health() != 100 || blue.Health() != 100 { // fire
		t.Error("wrong value")
	}
	if red.Pos().X != 100 || blue.Pos().X != 600 { // move
		t.Error("wrong value", red.Pos().Xf, blue.Pos().Xf)
	}

	// test
	for i := 0; i < 900; i++ {
		AttackMove(nil) // test nil
		AttackMove(red)
		AttackMove(blue)
		AttackMove(blue2)
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
