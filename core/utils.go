package core

import "math"

// IsCollided return true if two objects touch.
func IsCollided(pos1 Position, rad1 int, pos2 Position, rad2 int) bool {
	return Distance(pos1, pos2) < float64(rad1+rad2)
}

// Distance returns the distance between the centers of two objects.
func Distance(a, b Position) float64 {
	xd := math.Abs(a.Xf - b.Xf)
	yd := math.Abs(a.Yf - b.Yf)
	return Length(xd, yd)
}

// Length calculate the vector length.
//
//	Sqrt(x*x + y*y)
func Length(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

// CalcPosFromAngle returns the new coordinates with given angle and length.
// The angle is defined with North is 0° (see North, East, South, West, ...).
func CalcPosFromAngle(start Position, angle, length int) Position {
	// convert to radian
	angle = (angle - 90 + 360) % 360           // correct angle (North is 0°)
	angleRad := float64(angle) * math.Pi / 180 // convert to radian

	// calc new X, Y
	x := start.Xf + float64(length)*math.Cos(angleRad)
	y := start.Yf + float64(length)*math.Sin(angleRad)

	// build new position
	end := NewPosition(int(math.Round(x)), int(math.Round(y)))
	end.Xf = x
	end.Yf = y

	// return
	return end
}

//--------------------------------------------------------------------------------------------------------------------//

// CheckBorders returns true if the object collide with the borders.
func CheckBorders(p Position, r, screenWidth, screenHeight int) bool {

	// top
	if p.Y-r < 0 {
		return true
	}
	// bottom
	if p.Y+r > screenHeight {
		return true
	}
	// left
	if p.X-r < 0 {
		return true
	}
	// right
	if p.X+r > screenWidth {
		return true
	}

	// no problem
	return false
}

//--------------------------------------------------------------------------------------------------------------------//

// RelativeAngle returns the relative position of OTHER as an angle.
// Return [0 ... 359]
//
//	North     = 0
//	Northeast = 45
//	East      = 90
//	Southeast = 135
//	South     = 180
//	Southwest = 225
//	West      = 270
//	Northwest = 315
func RelativeAngle(me Position, other Position) int {
	xd := me.Xf - other.Xf
	yd := me.Yf - other.Yf

	// calc angle
	// and correct North to 0
	a2 := math.Atan2(yd, xd) / math.Pi * 180
	a2 -= 90 // correct North

	// modulo 360
	ret := int(a2) + 360
	ret %= 360

	// return
	return ret
}
