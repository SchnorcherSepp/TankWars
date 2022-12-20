package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"testing"
)

func TestMoveTo(t *testing.T) {
	w := core.NewWorld(1000, 1000)

	nt, _ := core.NewTank(w, core.RedTank, 5, 15, core.WeaponCannon)
	nt.SetPosition(core.NewPosition(500, 500), core.North)
	w.AddTank(nt)
	rock, _ := core.NewTank(w, core.NeutralRock, 5, 15, core.WeaponNone)
	rock.SetPosition(core.NewPosition(600, 600), core.Northwest)
	w.AddTank(rock)

	// test nil
	MoveTo(nil, core.Position{})

	// move
	for i := 0; i < 110; i++ {
		MoveTo(nt, core.NewPosition(700, 700))
		w.Update()
	}

	// check
	if nt.Blocked() != true || nt.Pos().X < 520 || nt.Pos().Y < 520 {
		t.Error("wrong value", nt.Blocked(), nt.Pos().X, nt.Pos().Y)
	}

	// move without rock
	rock.Remove() // remove rock
	nt.Forward()  // reset 'block'
	w.Update()

	for i := 0; i < 100; i++ {
		MoveTo(nt, core.NewPosition(700, 700))
		w.Update()
	}

	// check
	if nt.Blocked() != false || nt.Pos().X < 670 || nt.Pos().Y < 670 {
		t.Error("wrong value", nt.Blocked(), nt.Pos().X, nt.Pos().Y)
	}
}
