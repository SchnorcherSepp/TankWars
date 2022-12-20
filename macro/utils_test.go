package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"testing"
)

func TestRotationsToTarget(t *testing.T) {
	// test nil
	if RotationsToTarget(nil, core.Target{}.RelativeAngle) != 0 {
		t.Error("wrong value")
	}

	// test normal
	w := core.NewWorld(1000, 1000)

	me, _ := core.NewTank(w, "red", 5, 15, core.WeaponCannon)
	me.SetPosition(core.NewPosition(100, 100), core.North)
	w.AddTank(me)

	tank, _ := core.NewTank(w, "North", 5, 15, core.WeaponArtillery)
	tank.SetPosition(core.NewPosition(100, 60), core.North)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "Northeast", 5, 15, core.WeaponArtillery)
	tank.SetPosition(core.NewPosition(140, 53), core.Northeast)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "East", 5, 15, core.WeaponNone)
	tank.SetPosition(core.NewPosition(159, 100), core.East)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "Southeast", 5, 15, core.WeaponNone)
	tank.SetPosition(core.NewPosition(150, 161), core.Southeast)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "South", 5, 15, core.WeaponCannon)
	tank.SetPosition(core.NewPosition(93, 150), core.South)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "Southwest", 5, 15, core.WeaponCannon)
	tank.SetPosition(core.NewPosition(55, 150), core.Southwest)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "West", 5, 15, core.WeaponCannon)
	tank.SetPosition(core.NewPosition(50, 101), core.West)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "Northwest", 5, 15, core.WeaponCannon)
	tank.SetPosition(core.NewPosition(50, 50), core.Northwest)
	w.AddTank(tank)

	tank, _ = core.NewTank(w, "other", 5, 15, core.WeaponCannon)
	tank.SetPosition(core.NewPosition(400, 240), core.North)
	w.AddTank(tank)

	//--------------------------------------------------------------

	// TEST 1
	list := core.PossibleTargets(me, "")
	if len(list) != 8 {
		t.Fatal("wrong value", len(list))
	}
	for i, test := range []int{0, 1, -1, 2, -2, -3, 3, -4} {
		rs := RotationsToTarget(me, list[i].RelativeAngle)
		if rs != test {
			t.Error("wrong value", rs, list[i].Tank.Owner(), list[i].Distance, list[i].RelativeAngle)
		}
	}

	// TEST 2
	me.SetPosition(core.NewPosition(100, 100), core.Southeast)
	list = core.PossibleTargets(me, "")
	if len(list) != 8 {
		t.Fatal("wrong value", len(list))
	}
	for i, test := range []int{0, -1, 1, 2, -2, 3, -3, 4} {
		rs := RotationsToTarget(me, list[i].RelativeAngle)
		if rs != test {
			t.Error("wrong value", rs, list[i].Tank.Owner(), list[i].Distance, list[i].RelativeAngle)
		}
	}
}
