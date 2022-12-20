package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestFireWall(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	// prepare world
	w := core.NewWorld(333, 666)

	red, _ := core.NewTank(w, core.RedTank, 5, 70, core.WeaponRockets)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)

	// test
	for i := 0; i < 688; i++ {
		FireWall(nil) // test nil
		FireWall(red)
		w.Update()
	}

	// check
	c := len(w.Projectiles())
	if c < 1 {
		t.Error("wrong value", c)
	}

}
