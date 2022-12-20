package macro

import "github.com/SchnorcherSepp/TankWars/core"

// MoveTo uses Forward(), Left() and Right() to reach the given position.
// If the tank becomes Blocked(), the algorithm will be paused and must be reset manually with Forward().
func MoveTo(t *core.Tank, to core.Position) {
	if t == nil {
		return // EXIT
	}

	// tank position
	me := t.Pos()

	// movement start/stop
	d := core.Distance(me, to)
	if d > core.BlockRadius {
		if !t.Moving() && !t.Blocked() {
			t.Forward()
		}
	} else {
		if t.Moving() {
			t.Stop()
			return // EXIT!!
		}
	}

	// direction left/right
	ra := core.RelativeAngle(me, to)
	r := RotationsToTarget(t, ra)
	if r < 0 {
		t.Left()
	}
	if r > 0 {
		t.Right()
	}
}
