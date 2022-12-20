package core

import (
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"math"
	"testing"
)

func TestWeapon_Status_PreparationTime(t *testing.T) {
	w := NewWorld(333, 444)
	w.UpdateN(500)

	// spawn new weapons
	cannon := NewWeaponCannon(w, nil, 999)
	artillery := NewWeaponArtillery(w, nil, 999)

	// PreparationTime: not ready
	if rdy, txt := cannon.Status(); rdy || txt != StatusPreparing {
		t.Error("[#1] cannon is ready too soon:", txt)
	}
	if rdy, txt := artillery.Status(); rdy || txt != StatusPreparing {
		t.Error("[#1] artillery is ready too soon:", txt)
	}

	w.UpdateN(int(cannon.prepTime))

	// PreparationTime: cannon rdy
	if rdy, txt := cannon.Status(); !rdy {
		t.Error("[#2] cannon is not ready:", txt)
	}
	if rdy, _ := artillery.Status(); rdy {
		t.Error("[#2] artillery is ready too soon")
	}

	w.UpdateN(int(artillery.prepTime))

	// PreparationTime: all ready
	if rdy, txt := cannon.Status(); !rdy {
		t.Error("[#3] cannon is not ready:", txt)
	}
	if rdy, txt := artillery.Status(); !rdy {
		t.Error("[#3] artillery is not ready:", txt)
	}

	cannon.Update(true)
	artillery.Update(true)

	// PreparationTime: not ready
	if rdy, txt := cannon.Status(); rdy || txt != StatusPreparing {
		t.Error("[#4] cannon is ready too soon")
	}
	if rdy, txt := artillery.Status(); rdy || txt != StatusPreparing {
		t.Error("[#4] artillery is ready too soon")
	}
}

func TestWeapon_Status_ReloadTime(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	w := NewWorld(333, 444)
	w.UpdateN(500)

	// spawn new weapons; set no movement
	cannon := NewWeaponCannon(w, nil, 999)
	artillery := NewWeaponArtillery(w, nil, 999)
	// reset lastMove
	w.UpdateN(int(artillery.prepTime))
	// FIRE!!
	cannon.Fire(Position{}, 0, 0, 100)
	artillery.Fire(Position{}, 0, 0, 100)

	// ReloadTime: not ready
	if rdy, txt := cannon.Status(); rdy || txt != StatusReloading {
		t.Error("[#1] cannon is ready too soon")
	}
	if rdy, txt := artillery.Status(); rdy || txt != StatusReloading {
		t.Error("[#1] artillery is ready too soon")
	}

	w.UpdateN(int(cannon.reloadTime))

	// ReloadTime: cannon rdy
	if rdy, txt := cannon.Status(); !rdy {
		t.Error("[#2] cannon is not ready:", txt)
	}
	if rdy, _ := artillery.Status(); rdy {
		t.Error("[#2] artillery is ready too soon")
	}

	w.UpdateN(int(artillery.reloadTime))

	// ReloadTime: all ready
	if rdy, txt := cannon.Status(); !rdy {
		t.Error("[#3] cannon is not ready:", txt)
	}
	if rdy, txt := artillery.Status(); !rdy {
		t.Error("[#3] artillery is not ready:", txt)
	}

	cannon.Fire(Position{}, 0, 0, 100)
	artillery.Fire(Position{}, 0, 0, 100)

	// ReloadTime: not ready
	if rdy, _ := cannon.Status(); rdy {
		t.Error("[#4] cannon is ready too soon")
	}
	if rdy, _ := artillery.Status(); rdy {
		t.Error("[#4] artillery is ready too soon")
	}
}

func TestWeapon_Getter(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	world := NewWorld(1111, 2222)
	world.UpdateN(100)

	tank, err := NewTank(world, "Owner", 11, 22, WeaponCannon)
	if err != nil {
		panic(err)
	}
	w := NewWeaponCannon(world, tank, 33)

	if w.Type() != WeaponCannon {
		t.Error("wrong value", w.Type())
	}
	if w.Range() != 338 {
		t.Error("wrong value", w.Range())
	}
	if w.PreparationTime() != 24 {
		t.Error("wrong value", w.PreparationTime())
	}
	if w.ReloadTime() != 90 {
		t.Error("wrong value", w.ReloadTime())
	}
	if w.ProjectileSpeed() != 600 {
		t.Error("wrong value", w.ProjectileSpeed())
	}
	if w.Damage() != int(math.Round(float64(33)*1.5)) {
		t.Error("wrong value", w.Damage())
	}
	if w.AoERadius() != 0 {
		t.Error("wrong value", w.AoERadius())
	}
	if w.ProjectileCollision() != true {
		t.Error("wrong value", w.ProjectileCollision())
	}
	if w.AnyFireAngle() != false {
		t.Error("wrong value", w.AnyFireAngle())
	}
	if w.LastMove() != 100 {
		t.Error("wrong value", w.LastMove())
	}
	if w.LastFire() != 0 {
		t.Error("wrong value", w.LastFire())
	}
	if rdy, txt := w.Status(); rdy || txt != StatusPreparing {
		t.Error("wrong value", rdy, txt)
	}
}

func TestWeapon_Status(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	world := NewWorld(1111, 2222)
	world.UpdateN(100)

	tank, err := NewTank(world, "Owner", 11, 22, WeaponCannon)
	if err != nil {
		panic(err)
	}
	w := NewWeaponCannon(world, tank, 33)

	// preparing
	w.Fire(tank.pos, tank.angle, 0, 10)
	if rdy, txt := w.Status(); rdy || txt != StatusPreparing {
		t.Error("wrong value", rdy, txt)
	}

	// reset (ready)
	world.UpdateN(int(w.prepTime))
	if rdy, txt := w.Status(); !rdy || txt != StatusReady {
		t.Error("wrong value", rdy, txt)
	}

	// fire
	w.Fire(tank.pos, tank.angle, 0, 99999)
	if rdy, txt := w.Status(); rdy || txt != StatusReloading {
		t.Error("wrong value", rdy, txt)
	}

	// ready
	world.UpdateN(int(w.reloadTime))
	if rdy, txt := w.Status(); !rdy || txt != StatusReady {
		t.Error("wrong value", rdy, txt)
	}

	// moving
	tank.Backward()
	if rdy, txt := w.Status(); rdy || txt != StatusMoving {
		t.Error("wrong value", rdy, txt)
	}
}
