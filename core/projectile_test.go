package core

import (
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestProjectile_Getter(t *testing.T) {
	pos := NewPosition(13, 17)
	p := NewProjectile(nil, nil, pos, 7, 33, 9, 900, 12, true)

	if p.Parent() != nil {
		t.Error("wrong value")
	}
	if p.Pos().X != 13 || p.Pos().Y != 17 {
		t.Error("wrong value")
	}
	if p.StartPos().X != 13 || p.StartPos().Y != 17 {
		t.Error("wrong value")
	}
	if p.Angle() != 7 {
		t.Error("wrong value")
	}
	if p.Distance() != 33 {
		t.Error("wrong value")
	}
	if p.Speed() != 9 {
		t.Error("wrong value")
	}
	if p.Damage() != 900 {
		t.Error("wrong value")
	}
	if p.AoERadius() != 12 {
		t.Error("wrong value")
	}
	if p.Collision() != true {
		t.Error("wrong value")
	}
}

func TestProjectile_Pos(t *testing.T) {
	pos := NewPosition(13, 17)
	p := NewProjectile(nil, nil, pos, East, 33, 1000, 900, 12, true)
	p.Update()

	if p.Pos().X != 13+1000*MovePerTick || p.Pos().Y != 17 {
		t.Errorf("wrong value: %#v", p.Pos())
	}
	if p.StartPos().X != 13 || p.StartPos().Y != 17 {
		t.Errorf("wrong value: %#v", p.StartPos())
	}
	if p.EndPos().X != 46 || p.EndPos().Y != 17 {
		t.Errorf("wrong value: %#v", p.EndPos())
	}
}

func TestProjectile_Exploded(t *testing.T) {
	p := NewProjectile(nil, nil, NewPosition(0, 0), 0, 100, 100, 100, 100, true)
	if p.Exploded() != false {
		t.Error("wrong value")
	}
	p.exploded = 1
	if p.Exploded() != true {
		t.Error("wrong value")
	}
}

func TestProjectile_Update(t *testing.T) {
	resources.MuteSound = true // mute for tests

	// create projectile
	w := NewWorld(1000, 1000)
	w.UpdateN(500)

	pj := NewProjectile(w, nil, NewPosition(100, 100), East, 1000, 1600, 99, 0, true)

	// create target
	ta, err := NewTank(w, "test", 10, 20, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	ta.SetPosition(NewPosition(200, 100), West)
	w.AddTank(ta)

	// CHECK: Exploded
	pj.Update()
	if pj.Exploded() {
		t.Error("wrong value")
	}
	pj.Update()
	if !pj.Exploded() {
		t.Error("wrong value")
	}
	pj.Update()
	if !pj.Exploded() {
		t.Error("wrong value")
	}

	// CHECK: remove
	w.projectiles = append(w.projectiles, pj)
	w.projectiles = append(w.projectiles, NewProjectile(nil, nil, Position{}, 0, 0, 0, 0, 0, false))
	if len(w.projectiles) != 2 {
		t.Error("wrong value")
	}
	for i := 0; i < 1000; i++ {
		pj.Update()
	}
	if len(w.projectiles) != 1 {
		t.Error("wrong value")
	}

	// CHECK: max distance (cannon)
	pj = NewProjectile(w, nil, Position{}, South, 100, 300, 0, 0, true)
	w.projectiles = append(w.projectiles, pj)
	if len(w.projectiles) != 2 {
		t.Errorf("wrong value: %d", len(w.projectiles))
	}
	for i := 0; i < 22; i++ {
		pj.Update()
	}
	if len(w.projectiles) != 1 {
		t.Errorf("wrong value: %d, %#v", len(w.projectiles), pj)
	}

	// CHECK: max distance (explosion)
	pj = NewProjectile(w, nil, Position{}, South, 100, 300, 0, 0, false)
	w.projectiles = append(w.projectiles, pj)
	if len(w.projectiles) != 2 || pj.Exploded() {
		t.Errorf("wrong value: %d", len(w.projectiles))
	}
	for i := 0; i < 22; i++ {
		pj.Update()
	}
	if len(w.projectiles) != 2 || !pj.Exploded() {
		t.Errorf("wrong value: %d, %#v", len(w.projectiles), pj)
	}
	for i := 0; i < 2200; i++ {
		pj.Update()
	}
	if len(w.projectiles) != 1 {
		t.Errorf("wrong value: %d, %#v", len(w.projectiles), pj)
	}
}
