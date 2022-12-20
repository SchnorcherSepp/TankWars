package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
)

// GuardMode makes the tank wait and attack anything that approaches.
// Cannons can change their angle but cannot move.
func GuardMode(t *core.Tank, filter ...string) {
	if t == nil || t.Weapon() == nil || t.Weapon().Type() == core.WeaponNone {
		return // EXIT
	}

	// macro
	list := core.PossibleTargets(t, filter...)
	if len(list) > 0 {

		// fire at targets
		if t.Weapon().AnyFireAngle() {
			// Artillery
			t.FireAt(list[0].Tank.Pos())

		} else {
			// Cannon
			rs := RotationsToTarget(t, list[0].RelativeAngle)
			if rs < 0 {
				t.Left()
			} else if rs > 0 {
				t.Right()
			} else {
				t.FireAt(list[0].Tank.Pos())
			}
		}
	}
}
