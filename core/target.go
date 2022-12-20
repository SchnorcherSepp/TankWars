package core

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Target provide the tank, the distance and the relative angle to the tank.
// It is used by CloseTargets() and PossibleTargets().
// see RelativeAngle().
type Target struct {
	Tank          *Tank
	Distance      int
	RelativeAngle int
}

//--------------------------------------------------------------------------------------------------------------------//

// CloseTargets returns all objects in the world that are theoretical in weapon range.
// The weapon type is irrelevant (WeaponCannon or WeaponArtillery) and the angle of the tank is ignored.
// The list is sorted by distance (from the closest to the farthest).
func CloseTargets(t *Tank, filter ...string) []Target {
	// no tank or no weapon
	if t == nil || t.weapon == nil {
		return make([]Target, 0)
	}

	// check all world objects
	list := make([]Target, 0, 16)
	for _, ot := range t.world.tanks {
		if ot == nil || t == ot {
			continue
		}

		// filter
		filterThis := false
		for _, f := range filter {
			if len(f) > 0 && strings.HasPrefix(ot.owner, f) {
				filterThis = true
				break
			}
		}
		if filterThis {
			continue // skip this target
		}

		// calc distance and InRange
		dist := Distance(t.pos, ot.pos)
		inRange := dist < float64(t.weapon.rng+BlockRadius)

		// add other tank to list
		if inRange {
			list = append(list, Target{
				Tank:          ot,
				Distance:      int(dist),
				RelativeAngle: RelativeAngle(t.pos, ot.pos),
			})
		}
	}

	// sort list and return
	sort.Slice(list, func(i, j int) bool {
		return list[i].Distance < list[j].Distance
	})
	return list
}

// PossibleTargets extends CloseTargets.
// It only returns objects that can actually be attacked,depending on the weapon type.
// However, it may be necessary for the battle tank to change its angle.
// The list is sorted by the rotation required to reach the target.
func PossibleTargets(t *Tank, filter ...string) []Target {
	// no tank or no weapon
	if t == nil || t.weapon == nil {
		return make([]Target, 0)
	}

	// get targets in weapon range
	targetsInRange := CloseTargets(t, filter...)

	// Artillery -> RETURN
	if t.weapon.anyFireAngle {
		return targetsInRange // FIN!
	}

	// ------ for cannons it's more complicated ------ //

	// Cannon: horizontal and vertical
	cannonTargets := make([]Target, 0, len(targetsInRange))
	for _, ot := range targetsInRange {
		possibleX := ot.Tank.pos.X+(BlockRadius-1) > t.pos.X-(BallRadius-1) && ot.Tank.pos.X-(BlockRadius-1) < t.pos.X+(BallRadius-1)
		possibleY := ot.Tank.pos.Y+(BlockRadius-1) > t.pos.Y-(BallRadius-1) && ot.Tank.pos.Y-(BlockRadius-1) < t.pos.Y+(BallRadius-1)
		if possibleX || possibleY {
			cannonTargets = append(cannonTargets, ot)
		}
	}
	// Cannon: diagonal
	for _, ot := range targetsInRange {
		dX := math.Abs(float64(ot.Tank.pos.X - t.pos.X))
		dY := math.Abs(float64(ot.Tank.pos.Y - t.pos.Y))
		dd := math.Abs(dX - dY)
		if dd < BlockRadius+BallRadius+10 {
			cannonTargets = append(cannonTargets, ot)
		}
	}

	// deduplicate list: a tank can be a target horizontal/vertical and a target diagonal at the same time!
	{
		allKeys := make(map[string]bool)
		newList := make([]Target, 0, len(cannonTargets))
		for _, target := range cannonTargets {
			if target.Tank != nil && target.Tank.id != "" {
				id := target.Tank.id
				if _, ok := allKeys[id]; !ok {
					allKeys[id] = true
					newList = append(newList, target)
				}
			} else {
				fmt.Printf("WARNING: Remove tank without ID from possible target list!\n")
			}
		}
		cannonTargets = newList
	}

	// sort list: minimize rotations
	sort.SliceStable(cannonTargets, func(i, j int) bool {
		// consider own angle
		iAD := (cannonTargets[i].RelativeAngle + 360 - t.angle) % 360 // modulo
		jAD := (cannonTargets[j].RelativeAngle + 360 - t.angle) % 360 // modulo

		// convert [0 ... 360] to [-180 ... +180]
		iAD -= 180
		jAD -= 180

		// compare absolut values
		return math.Abs(float64(iAD)) > math.Abs(float64(jAD))
	})

	// Cannon -> RETURN
	return cannonTargets
}
