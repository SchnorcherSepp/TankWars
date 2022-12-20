package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"math/rand"
)

// AttackMove moves the tank in the aligned direction.
// If the tank encounters an enemy, it stops and opens fire.
// If the enemy is destroyed, the tank continues to move.
func AttackMove(t *core.Tank, filter ...string) {
	if t == nil || t.Weapon() == nil || t.Weapon().Type() == core.WeaponNone {
		return // EXIT
	}

	// macro
	list := core.PossibleTargets(t, filter...)
	if len(list) > 0 {
		t.Stop()
		GuardMode(t, filter...)
	} else {
		if t.Blocked() {
			// random left/right
			if rand.Intn(2) == 1 {
				t.Left()
			} else {
				t.Right()
			}
		}
		t.Forward()
	}
}
