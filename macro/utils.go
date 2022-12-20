package macro

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"math"
)

// RotationsToTarget returns the number of rotation steps.
// Negative numbers require rotation to the left.
// Positive numbers require rotation to the right.
func RotationsToTarget(t *core.Tank, relativeAngle int) int {
	// no tank
	if t == nil {
		return 0
	}

	// add own angle to relative angle
	// and convert [0 ... 360] to [-180 ... +180]
	ra := (t.Angle() + 180 - relativeAngle + 360) % 360 // modulo
	ra -= 180

	// invert (-1 is left, +1 is right)
	ra *= -1

	// compare absolut values
	return int(math.Round(float64(ra) / 45))
}
