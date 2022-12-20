package core

import "testing"

func Test_NewPosition(t *testing.T) {
	p := NewPosition(1, -1)

	if p.X != 1 || p.Xf != 1 {
		t.Error("wrong x value")
	}
	if p.Y != -1 || p.Yf != -1 {
		t.Error("wrong y value")
	}
}

func TestPosition_int_float(t *testing.T) {

	// check: int != float  (set float)
	p := Position{Xf: 100, Yf: 300}
	p.Update(0, 0)
	if p.X != 0 || p.Xf != 0 || p.Y != 0 || p.Yf != 0 {
		t.Error("wrong value")
	}

	// check: int != float  (set int)
	p = Position{X: 100, Y: 300}
	p.Update(0, 0)
	if p.X != 100 || p.Xf != 100 || p.Y != 300 || p.Yf != 300 {
		t.Error("wrong value")
	}

}

func TestPosition_Update(t *testing.T) {
	p := NewPosition(13, 17)
	p.Update(North, 33)
	if p.X != 13 || p.Xf != 13 || p.Y != 17-1 || p.Yf != 17-33*MovePerTick {
		t.Errorf("wrong value: %#v", p)
	}

	p = NewPosition(13, 17)
	p.Update(East, 33)
	if p.X != 13+1 || p.Xf != 13+33*MovePerTick || p.Y != 17 || p.Yf != 17 {
		t.Errorf("wrong value: %#v", p)
	}
}
