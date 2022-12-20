package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"math/rand"
)

// FireWall fires at random positions in front of the tank.
// Most fun in combination with the rocket launcher.
func FireWall(t *core.Tank) {
	if t == nil || t.Weapon() == nil || t.Weapon().Type() == core.WeaponNone {
		return // EXIT
	}

	// get random angle: -35° to +35°
	angle := rand.Intn(70) - 35 + t.Angle()

	// fire
	t.Fire(angle, 9999)
}
