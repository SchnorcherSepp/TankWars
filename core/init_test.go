package core

import (
	"reflect"
	"testing"
)

func TestInitialization_Projectile(t *testing.T) {
	// origin
	nw := NewWorld(81, 82)
	nt, _ := NewTank(nw, "test parent", 22, 33, WeaponRockets)

	o := &Projectile{
		world:     nw,
		parent:    nt,
		pos:       NewPosition(91, 92),
		startPos:  NewPosition(93, 94),
		endPos:    NewPosition(95, 96),
		angle:     1,
		distance:  2,
		speed:     3,
		damage:    4,
		aoeRadius: 5,
		collision: true,
		exploded:  6,
	}

	// TestInitialization
	clone := new(Projectile)
	clone.TestInitialization(o.world, o.parent, o.pos, o.startPos, o.endPos, o.angle, o.distance, o.speed, o.damage, o.aoeRadius, o.collision, o.exploded)

	// compare
	if !reflect.DeepEqual(o, clone) {
		t.Error("not equal")
	}
}

func TestInitialization_Weapon(t *testing.T) {
	// origin
	nw := NewWorld(81, 82)
	nt, _ := NewTank(nw, "test parent", 22, 33, WeaponRockets)

	o := &Weapon{
		world:         nw,
		parent:        nt,
		typ:           "test type",
		rng:           1,
		prepTime:      2,
		reloadTime:    3,
		projSpeed:     4,
		damage:        5,
		aoeRadius:     6,
		projCollision: true,
		anyFireAngle:  true,
		lastMove:      7,
		lastFire:      8,
	}

	// TestInitialization
	clone := new(Weapon)
	clone.TestInitialization(o.world, o.parent, o.typ, o.rng, o.prepTime, o.reloadTime, o.projSpeed, o.damage, o.aoeRadius, o.projCollision, o.anyFireAngle, o.lastMove, o.lastFire)

	// compare
	if !reflect.DeepEqual(o, clone) {
		t.Error("not equal")
	}
}

func TestInitialization_Tank(t *testing.T) {
	// origin
	nw := NewWorld(81, 82)
	we := NewWeaponArtillery(nw, new(Tank), 99)

	o := &Tank{
		world:      nw,
		id:         "test id",
		owner:      "test owner",
		weapon:     we,
		health:     1,
		armor:      2,
		speed:      3,
		pos:        NewPosition(11, 22),
		command:    4,
		angle:      5,
		isBlocked:  true,
		lastRotate: 6,
		macro:      nil,
	}

	// TestInitialization
	clone := new(Tank)
	clone.TestInitialization(o.world, o.id, o.owner, o.weapon, o.health, o.armor, o.speed, o.pos, o.command, o.angle, o.isBlocked, o.lastRotate, o.macro)

	// compare
	if !reflect.DeepEqual(o, clone) {
		t.Error("not equal")
	}
}

func TestInitialization_World(t *testing.T) {
	// origin
	o := &World{
		xWidth:      1,
		yHeight:     2,
		iteration:   3,
		tanks:       make([]*Tank, 1),
		projectiles: make([]*Projectile, 2),
		freeze:      true,
		cashRed:     4,
		cashBlue:    5,
	}

	// TestInitialization
	clone := new(World)
	clone.TestInitialization(o.xWidth, o.yHeight, o.iteration, o.tanks, o.projectiles, o.freeze, o.cashRed, o.cashBlue)

	// compare
	if !reflect.DeepEqual(o, clone) {
		t.Error("not equal")
	}
}
