package core

import (
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestTargetsInRange(t *testing.T) {
	// test nil
	if len(CloseTargets(nil, "")) != 0 {
		t.Error("wrong value")
	}

	// test normal
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	me, _ := NewTank(w, "blue", 5, 15, WeaponCannon)
	me.SetPosition(NewPosition(100, 100), South)
	w.AddTank(me)

	tank, _ := NewTank(w, "tank1", 5, 15, WeaponArtillery)
	tank.SetPosition(NewPosition(10, 30), North)
	w.AddTank(tank)

	tank, _ = NewTank(w, "tank2", 5, 15, WeaponNone)
	tank.SetPosition(NewPosition(170, 100), North)
	w.AddTank(tank)

	tank, _ = NewTank(w, "tank3", 5, 15, WeaponNone)
	tank.SetPosition(NewPosition(50, 200), North)
	w.AddTank(tank)

	tank, _ = NewTank(w, "tank4", 5, 15, WeaponCannon)
	tank.SetPosition(NewPosition(800, 800), North)
	w.AddTank(tank)

	//--------------------------------------------------------------

	// TEST 1
	list := CloseTargets(me, "")
	if len(list) != 3 {
		t.Fatal("wrong value")
	}
	i := 0
	if list[i].Tank.owner != "tank2" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}
	i = 1
	if list[i].Tank.owner != "tank3" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}
	i = 2
	if list[i].Tank.owner != "tank1" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}

	// TEST 2
	tank.SetPosition(NewPosition(164, 100), West)
	list = CloseTargets(me, "")
	if len(list) != 4 {
		t.Fatal("wrong value")
	}
	i = 0
	if list[i].Tank.owner != "tank4" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}
	i = 1
	if list[i].Tank.owner != "tank2" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}
	i = 2
	if list[i].Tank.owner != "tank3" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}
	i = 3
	if list[i].Tank.owner != "tank1" {
		t.Error("wrong value", list[i].Tank.owner, list[i].Distance)
	}

}

func TestTargetsInRange_filter(t *testing.T) {
	// test normal
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	me, _ := NewTank(w, "blue", 5, 15, WeaponCannon)
	me.SetPosition(NewPosition(100, 100), South)
	w.AddTank(me)

	tank, _ := NewTank(w, "tank1", 5, 15, WeaponArtillery)
	tank.SetPosition(NewPosition(10, 30), North)
	w.AddTank(tank)

	tank, _ = NewTank(w, "tank2", 5, 15, WeaponNone)
	tank.SetPosition(NewPosition(170, 100), North)
	w.AddTank(tank)

	// TEST
	if list := CloseTargets(me, ""); len(list) != 2 {
		t.Error("wrong value")
	}
	if list := CloseTargets(me, "tank1"); len(list) != 1 {
		t.Error("wrong value")
	}
	if list := CloseTargets(me, "tank1", "tank2"); len(list) != 0 {
		t.Error("wrong value")
	}
}

func TestPossibleTargets_WrongID(t *testing.T) {
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	me, _ := NewTank(w, "red", 5, 15, WeaponCannon)
	me.SetPosition(NewPosition(100, 100), North)
	w.AddTank(me)

	tank, _ := NewTank(w, "North", 5, 15, WeaponNone)
	tank.SetPosition(NewPosition(100, 50), North)
	w.AddTank(tank)

	// TEST 1
	list := PossibleTargets(me, "")
	if len(list) != 1 {
		t.Fatal("wrong value", len(list))
	}

	// TEST 2
	tank.id = ""
	me.id = ""
	list = PossibleTargets(me, "")
	if len(list) != 0 {
		t.Fatal("wrong value", len(list))
	}
}

func TestPossibleTargets(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	// test nil
	if len(PossibleTargets(nil, "")) != 0 {
		t.Error("wrong value")
	}

	// test normal
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	me, _ := NewTank(w, "red", 5, 15, WeaponCannon)
	me.SetPosition(NewPosition(100, 100), North)
	w.AddTank(me)

	tank, _ := NewTank(w, "North", 5, 15, WeaponArtillery)
	tank.SetPosition(NewPosition(100, 50), North)
	w.AddTank(tank)

	tank, _ = NewTank(w, "Northeast", 5, 15, WeaponArtillery)
	tank.SetPosition(NewPosition(150, 50), Northeast)
	w.AddTank(tank)

	tank, _ = NewTank(w, "East", 5, 15, WeaponNone)
	tank.SetPosition(NewPosition(150, 100), East)
	w.AddTank(tank)

	tank, _ = NewTank(w, "Southeast", 5, 15, WeaponNone)
	tank.SetPosition(NewPosition(150, 150), Southeast)
	w.AddTank(tank)

	tank, _ = NewTank(w, "South", 5, 15, WeaponCannon)
	tank.SetPosition(NewPosition(100, 150), South)
	w.AddTank(tank)

	tank, _ = NewTank(w, "Southwest", 5, 15, WeaponCannon)
	tank.SetPosition(NewPosition(50, 150), Southwest)
	w.AddTank(tank)

	tank, _ = NewTank(w, "West", 5, 15, WeaponCannon)
	tank.SetPosition(NewPosition(50, 100), West)
	w.AddTank(tank)

	tank, _ = NewTank(w, "Northwest", 5, 15, WeaponCannon)
	tank.SetPosition(NewPosition(50, 50), Northwest)
	w.AddTank(tank)

	tank, _ = NewTank(w, "other", 5, 15, WeaponCannon)
	tank.SetPosition(NewPosition(400, 240), North)
	w.AddTank(tank)

	//--------------------------------------------------------------

	// TEST 1
	me.weapon.anyFireAngle = true
	list := PossibleTargets(me, "")
	if len(list) != 9 {
		t.Fatal("wrong value", len(list))
	}

	// TEST 2
	me.weapon.anyFireAngle = false
	list = PossibleTargets(me, "")
	if len(list) != 8 {
		t.Fatal("wrong value", len(list))
	}
	for i, test := range []string{"North", "Northeast", "Northwest", "East", "West", "Southeast", "Southwest", "South"} {
		if list[i].Tank.owner != test {
			t.Error("wrong value", list[i].Tank.owner, list[i].Distance, list[i].RelativeAngle)
		}
	}

	// TEST 3
	me.angle = Southwest
	list = PossibleTargets(me, "")
	if len(list) != 8 {
		t.Fatal("wrong value", len(list))
	}
	for i, test := range []string{"Southwest", "South", "West", "Southeast", "Northwest", "North", "East", "Northeast"} {
		if list[i].Tank.owner != test {
			t.Error("wrong value", list[i].Tank.owner, list[i].Distance, list[i].RelativeAngle)
		}
	}

}
