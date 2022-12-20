package core

import (
	"math"
)

// Position stores x and y as an int and as a float value.
// Update() work with the float value and round to int.
type Position struct {
	X  int
	Xf float64
	Y  int
	Yf float64
}

// NewPosition return a position and set x,y float + int
func NewPosition(x, y int) Position {
	p := Position{}
	p.Set(x, y)
	return p
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// Set x,y float and int
func (p *Position) Set(x, y int) {
	p.X = x
	p.Xf = float64(x)
	p.Y = y
	p.Yf = float64(y)
}

//---------------- UPDATE --------------------------------------------------------------------------------------------//

// Update calculate the new position.
// Internally, float values are used. The int values are rounded.
// If the int values and the float values are differ, the int value is used.
func (p *Position) Update(angle, speed int) {

	// check diff
	if math.Abs(float64(p.X)-p.Xf) > 1 {
		p.Xf = float64(p.X)
		println("WARNING: reset X:", p.X, p.Xf)
	}
	if math.Abs(float64(p.Y)-p.Yf) > 1 {
		p.Yf = float64(p.Y)
		println("WARNING: reset Y:", p.Y, p.Yf)
	}

	// move
	r := float64(angle-90) * math.Pi / 180
	p.Xf += MovePerTick * float64(speed) * math.Cos(r)
	p.Yf += MovePerTick * float64(speed) * math.Sin(r)

	// update
	p.X = int(math.Round(p.Xf))
	p.Y = int(math.Round(p.Yf))
}
