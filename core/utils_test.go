package core

import "testing"

func TestRelativeAngle(t *testing.T) {

	if a := RelativeAngle(NewPosition(100, 100), NewPosition(100, 50)); a != North {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(150, 50)); a != Northeast {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(150, 100)); a != East {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(150, 150)); a != Southeast {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(100, 150)); a != South {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(50, 150)); a != Southwest {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(50, 100)); a != West {
		t.Error("wrong value", a)
	}
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(50, 50)); a != Northwest {
		t.Error("wrong value", a)
	}
	// other
	if a := RelativeAngle(NewPosition(100, 100), NewPosition(140, 50)); a != 38 {
		t.Error("wrong value", a)
	}
}

func TestCalcPosFromAngle(t *testing.T) {
	pos := CalcPosFromAngle(NewPosition(100, 200), North, 50)
	if pos.X != 100 || pos.Y != 150 {
		t.Error("wrong position", pos.Xf, pos.Yf)
	}

	pos = CalcPosFromAngle(NewPosition(100, 200), East, 50)
	if pos.X != 150 || pos.Y != 200 {
		t.Error("wrong position", pos.Xf, pos.Yf)
	}

	pos = CalcPosFromAngle(NewPosition(100, 200), Southwest, 50)
	if pos.X != 65 || pos.Y != 235 {
		t.Error("wrong position", pos.Xf, pos.Yf)
	}
}
