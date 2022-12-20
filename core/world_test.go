package core

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestNewWorld(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	w := NewWorld(111, 222)
	w.UpdateN(500)

	// get dimensions
	if w.XWidth() != 111 {
		t.Error("wrong value")
	}
	if w.YHeight() != 222 {
		t.Error("wrong value")
	}
	if w.ScreenWidth() != 111*BlockSize {
		t.Error("wrong value")
	}
	if w.ScreenHeight() != 222*BlockSize {
		t.Error("wrong value")
	}

	// check lists
	tank, err := NewTank(w, "test", 11, 22, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(tank)

	// one tank, no ball
	if len(w.Tanks()) != 1 {
		t.Error("wrong value")
	}
	if len(w.Projectiles()) != 0 {
		t.Error("wrong value")
	}

	// sleep and fire
	w.UpdateN(int(tank.weapon.prepTime))
	if suc, txt := tank.Fire(0, 999); !suc || txt != StatusReady {
		t.Error("wrong value", suc, txt)
	}

	// one tank and one ball
	if len(w.Tanks()) != 1 {
		t.Error("wrong value")
	}
	if len(w.Projectiles()) != 1 {
		t.Error("wrong value")
	}

	// update
	if w.Iteration() != tank.weapon.prepTime+500 {
		t.Error("wrong value")
	}
	w.Update()
	if w.Iteration() != tank.weapon.prepTime+500+1 {
		t.Error("wrong value")
	}
}

func TestWorld_UpdateN(t *testing.T) {
	w := NewWorld(222, 333)
	w.UpdateN(500)

	if w.iteration != 500 {
		t.Error("wrong value")
	}
	w.UpdateN(0)
	if w.iteration != 500 {
		t.Error("wrong value")
	}
	w.UpdateN(1)
	if w.iteration != 501 {
		t.Error("wrong value")
	}
	w.UpdateN(9)
	if w.iteration != 510 {
		t.Error("wrong value")
	}
}

func TestWorld_Update(t *testing.T) {
	w := NewWorld(222, 333)
	w.UpdateN(500)

	if w.iteration != 500 || w.freeze != false {
		t.Error("wrong value")
	}
	w.Update()
	if w.iteration != 501 || w.freeze != false {
		t.Error("wrong value")
	}
	w.Freeze(true)
	w.Update()
	w.Update()
	w.Update()
	if w.iteration != 501 || w.IsFrozen() != true {
		t.Error("wrong value")
	}
}

func TestWorld_BuyTank(t *testing.T) {
	// get world
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	w.SetCash(1000000, 1000000)

	// set red bases
	o, _ := NewTank(w, RedBase, 30, 30, WeaponNone)
	o.SetPosition(NewPosition(100, 100), North)
	w.AddTank(o)

	// test nil
	if err := w.BuyTank(nil); err == nil {
		t.Fatal("wrong value")
	}

	// test wrong owner
	nt, _ := NewTank(w, "orange", 22, 22, WeaponCannon)
	if err := w.BuyTank(nt); err == nil {
		t.Fatal("wrong value")
	}

	// test blue (not exist)
	nt, _ = NewTank(w, "blue", 22, 22, WeaponCannon)
	if err := w.BuyTank(nt); err == nil {
		t.Fatal("wrong value")
	}

	// test red (success)
	nt, _ = NewTank(w, "red", 22, 22, WeaponCannon)
	if err := w.BuyTank(nt); err != nil {
		t.Fatal("wrong value", err)
	}
	if len(w.Tanks()) != 2 {
		t.Fatal("wrong value")
	}

	// add blue base
	o, _ = NewTank(w, BlueBase, 30, 30, WeaponNone)
	o.SetPosition(NewPosition(800, 800), North)
	w.AddTank(o)

	// test blue (success)
	nt, _ = NewTank(w, "blue", 22, 22, WeaponCannon)
	if err := w.BuyTank(nt); err != nil {
		t.Fatal("wrong value", err)
	}
	if len(w.Tanks()) != 4 {
		t.Fatal("wrong value")
	}

	// test no cash
	w.SetCash(1000000, 99)
	nt, _ = NewTank(w, "blue", 22, 22, WeaponCannon)
	if err := w.BuyTank(nt); err == nil {
		t.Fatal("wrong value")
	}

	// test: not enough space!
	for i := 0; i < 1000; i++ {
		nt, _ = NewTank(w, "red", 22, 22, WeaponCannon)
		if err := w.BuyTank(nt); err != nil {
			fmt.Printf("%v\n", err)
			break
		}
	}
	if len(w.Tanks()) < 10 || len(w.Tanks()) > 30 {
		t.Fatal("wrong value", len(w.Tanks()))
	}
}

func TestWorld_UnitCount(t *testing.T) {
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	// init
	if r, b := w.UnitCount(); r != 0 || b != 0 {
		t.Error("wrong value", r, b)
	}

	// other
	nt, err := NewTank(w, "orange", 22, 33, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(nt)
	if r, b := w.UnitCount(); r != 0 || b != 0 {
		t.Error("wrong value", r, b)
	}

	// rock
	nt, err = NewTank(w, RedRock, 22, 33, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(nt)
	if r, b := w.UnitCount(); r != 0 || b != 0 {
		t.Error("wrong value", r, b)
	}

	// RED
	nt, err = NewTank(w, RedTank, 22, 33, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(nt)
	if r, b := w.UnitCount(); r != 1 || b != 0 {
		t.Error("wrong value", r, b)
	}

	// BLUE
	nt, err = NewTank(w, BlueTank, 22, 33, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(nt)
	if r, b := w.UnitCount(); r != 1 || b != 1 {
		t.Error("wrong value", r, b)
	}

	// BASE (blue)
	nt, err = NewTank(w, BlueBase, 22, 33, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(nt)
	if r, b := w.UnitCount(); r != 1 || b != 2 {
		t.Error("wrong value", r, b)
	}

	// BASE (red)
	nt, err = NewTank(w, RedBase, 22, 33, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(nt)
	if r, b := w.UnitCount(); r != 2 || b != 2 {
		t.Error("wrong value", r, b)
	}
}

func TestWorld_CashStat(t *testing.T) {
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	o, _ := NewTank(w, RedBase, 30, 30, WeaponNone)
	w.AddTank(o)
	o, _ = NewTank(w, BlueBase, 30, 30, WeaponNone)
	w.AddTank(o)

	// get and update
	if r, b := w.CashStat(); r != 0 || b != 0 {
		t.Error("wrong value", r, b)
	}
	for i := 0; i < 10000; i++ {
		w.Update()
	}
	if r, b := w.CashStat(); r != 277 || b != 277 {
		t.Error("wrong value", r, b)
	}

	// set
	w.SetCash(333, 444)
	if r, b := w.CashStat(); r != 333 || b != 444 {
		t.Error("wrong value", r, b)
	}
}

func TestWorld_Clear(t *testing.T) {
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	tank, _ := NewTank(w, "blue_stone", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "blue_base", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "blue", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "blue1", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "red_stone", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "red_base", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "red", 5, 15, WeaponNone)
	w.AddTank(tank)
	tank, _ = NewTank(w, "red1", 5, 15, WeaponNone)
	w.AddTank(tank)

	if len(w.tanks) != 8 {
		t.Error("wrong value")
	}
	w.Clear("green")
	if len(w.tanks) != 8 {
		t.Error("wrong value")
	}
	w.Clear("blue_")
	if len(w.tanks) != 6 {
		t.Error("wrong value")
	}
	w.Clear("")
	if len(w.tanks) != 0 {
		t.Error("wrong value")
	}
}
