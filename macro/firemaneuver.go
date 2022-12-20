package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
)

// FireAndManeuver moves the tank in the aligned direction.
// Whenever the weapon is loaded, the tank will stop and fire.
// While reloading, the tank keeps moving.
func FireAndManeuver(t *core.Tank) {
	if t == nil || t.Weapon() == nil || t.Weapon().Type() == core.WeaponNone {
		return // EXIT
	}

	// macro
	rdy, txt := t.Status()
	if rdy { // ready -> fire
		t.Fire(t.Angle(), 99999)
	} else if txt == core.StatusReloading { // move while reloading
		t.Forward()
	} else if txt == core.StatusMoving { // prepare for fire
		t.Stop()
	}
}
