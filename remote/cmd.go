package remote

import (
	"errors"
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/macro"
	"strconv"
)

//---------------- GETTER --------------------------------------------------------------------------------------------//

// MyName returns the active player of this connection (RedTank or BlueTank)
func MyName(owner string) string {
	return owner
}

// GameStatus returns a json with all world data.
func GameStatus(w *core.World) string {
	if w == nil {
		return "err: invalid world status"
	}

	jw := NewJsonWorld(w)
	return jw.Get()
}

// TankStatus returns a json with all data of a requested tank.
func TankStatus(w *core.World, tankID string) string {
	// get tank
	t, err := id2Tank(w, "", tankID)
	if err != nil {
		return err.Error()
	}

	// return
	jt := NewJsonTank(t)
	return jt.Get()
}

// CloseTargets returns all objects in the world that are theoretical in weapon range.
// The weapon type is irrelevant (WeaponCannon or WeaponArtillery) and the angle of the tank is ignored.
// The list is sorted by distance (from the closest to the farthest).
func CloseTargets(w *core.World, tankID string, filter ...string) string {
	// get tank
	t, err := id2Tank(w, "", tankID)
	if err != nil {
		return err.Error()
	}

	// return
	list := core.CloseTargets(t, filter...)
	ts := NewJsonTargets(list)
	return ts.Get()
}

// PossibleTargets extends CloseTargets.
// It only returns objects that can actually be attacked,depending on the weapon type.
// However, it may be necessary for the battle tank to change its angle.
// The list is sorted by the rotation required to reach the target.
func PossibleTargets(w *core.World, tankID string, filter ...string) string {
	// get tank
	t, err := id2Tank(w, "", tankID)
	if err != nil {
		return err.Error()
	}

	// return
	list := core.PossibleTargets(t, filter...)
	ts := NewJsonTargets(list)
	return ts.Get()
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// BuyTank buy a new tank and place it near the home base.
func BuyTank(w *core.World, owner, armor, damage, weapon string) string {
	// convert input
	a, err := strconv.Atoi(armor)
	if err != nil {
		return "err: armor: " + err.Error()
	}
	d, err := strconv.Atoi(damage)
	if err != nil {
		return "err: damage: " + err.Error()
	}

	// build tank
	t, err := core.NewTank(w, owner, a, d, weapon)
	if err != nil {
		return "err: " + err.Error()
	}

	// return
	err = w.BuyTank(t)
	if err != nil {
		return "err: " + err.Error()
	} else {
		return fmt.Sprintf("ok %s", t.ID())
	}
}

// Fire creates a new projectile.
// The attributes fireAngle and distance determine the direction and distance of the shot.
// Cannons can fire in vehicle angle only.
// The distance is limited by the weapon range.
func Fire(w *core.World, owner, tankID, angle, distance string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// convert input
	a, err := strconv.Atoi(angle)
	if err != nil {
		return "err: angle: " + err.Error()
	}
	d, err := strconv.Atoi(distance)
	if err != nil {
		return "err: distance: " + err.Error()
	}

	// return
	ok, txt := t.Fire(a, d)
	if ok {
		return "ok"
	} else {
		return "err: " + txt
	}
}

// FireAt is a wrapper for Fire() and convert the position to fireAngle and distance.
func FireAt(w *core.World, owner, tankID, x, y string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// convert input
	xInt, err := strconv.Atoi(x)
	if err != nil {
		return "err: X: " + err.Error()
	}
	yInt, err := strconv.Atoi(y)
	if err != nil {
		return "err: Y: " + err.Error()
	}

	// return
	ok, txt := t.FireAt(core.NewPosition(xInt, yInt))
	if ok {
		return "ok"
	} else {
		return "err: " + txt
	}
}

// Forward send the tank forward.
func Forward(w *core.World, owner, tankID string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// return
	t.Forward()
	return "ok"
}

// Backward send the tank back.
func Backward(w *core.World, owner, tankID string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// return
	t.Backward()
	return "ok"
}

// Stop the movement.
// Weapons can only build up when the tank is stationary
func Stop(w *core.World, owner, tankID string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// return
	t.Stop()
	return "ok"
}

// Left turn the tank direction 45° left.
func Left(w *core.World, owner, tankID string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// return
	ok, txt := t.Left()
	if ok {
		return "ok"
	} else {
		return "err: " + txt
	}
}

// Right turn the tank direction 45° right.
func Right(w *core.World, owner, tankID string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// return
	ok, txt := t.Right()
	if ok {
		return "ok"
	} else {
		return "err: " + txt
	}
}

// SetMacroMoveTo sets a special macro with a position that is called with every update.
func SetMacroMoveTo(w *core.World, owner, tankID, x, y string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// convert input
	xInt, err := strconv.Atoi(x)
	if err != nil {
		return "err: X: " + err.Error()
	}
	yInt, err := strconv.Atoi(y)
	if err != nil {
		return "err: Y: " + err.Error()
	}

	// set macro
	f := func(t *core.Tank) {
		macro.MoveTo(t, core.NewPosition(xInt, yInt))
	}
	t.SetMacro(f)

	// return
	return "ok"
}

// SetMacro sets a macro that is called with every update.
// (see MacroAttackMove, MacroFireWall, MacroFireAndManeuver, MacroGuardMode and MacroReset)
func SetMacro(w *core.World, owner, tankID, mco string) string {
	// get tank
	t, err := id2Tank(w, owner, tankID)
	if err != nil {
		return err.Error()
	}

	// convert input
	switch mco {
	case core.MacroAttackMove:
		var filters = genFilters(t)
		f := func(t *core.Tank) {
			macro.AttackMove(t, filters...)
		}
		t.SetMacro(f)
		return "ok"

	case core.MacroFireAndManeuver:
		f := func(t *core.Tank) {
			macro.FireAndManeuver(t)
		}
		t.SetMacro(f)
		return "ok"

	case core.MacroFireWall:
		f := func(t *core.Tank) {
			macro.FireWall(t)
		}
		t.SetMacro(f)
		return "ok"

	case core.MacroGuardMode:
		var filters = genFilters(t)
		f := func(t *core.Tank) {
			macro.GuardMode(t, filters...)
		}
		t.SetMacro(f)
		return "ok"

	case "", core.MacroReset, "reset", "null", "remove", "disable":
		t.SetMacro(nil)
		return "ok: disable macro"
	default:
		t.SetMacro(nil)
		return "err: macro not found"
	}
}

//---------------- HELPER --------------------------------------------------------------------------------------------//

// id2Tank is a helper function and find a tank by id.
// This function enforce the owner-control. Use "" to disable the owner check.
func id2Tank(w *core.World, owner, id string) (*core.Tank, error) {
	// find tank
	if w != nil {
		for _, t := range w.Tanks() {
			if t != nil && t.ID() == id {
				// found tank -> check owner
				if owner == "" || owner == t.Owner() {
					return t, nil // no owner oo correct owner
				} else {
					return nil, errors.New("err: no access to other players units") // wrong owner
				}
			}
		}
	}

	// err: no tank
	return nil, errors.New("err: tank not found")
}

// genFilters is a helper function and generate filter based on the owner.
func genFilters(t *core.Tank) []string {
	filters := make([]string, 0, 6)

	if t != nil {
		owner := t.Owner() // RedTank or BlueTank

		// own tanks and own base
		filters = append(filters, owner)
		filters = append(filters, owner+"_base")

		// all rocks
		filters = append(filters, core.NeutralRock)
		filters = append(filters, core.RedRock)
		filters = append(filters, core.BlueRock)
	}

	return filters
}
