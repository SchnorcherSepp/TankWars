package main

// Result holds all data of a single battle.
type Result struct {
	// left
	A1 int
	D1 int
	S1 int
	W1 string
	// right
	A2 int
	D2 int
	S2 int
	W2 string
	// result
	Winner    string
	LeftHP    int
	RightHP   int
	Iteration uint64
}

// Points rated how high the winner won.
func (r Result) Points() int {
	if r.Winner == "none" {
		return r.LeftHP - r.RightHP
	}
	if r.LeftHP > 0 {
		return r.LeftHP
	} else {
		return -r.RightHP
	}
}
