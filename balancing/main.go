package main

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"os"
	"sort"
)

type testCase struct {
	W1  string
	W2  string
	Out string
}

// main - the whole package simulate duels between two types of weapons.
// All variants (armor, damage, speed) are tested against each other.
func main() {
	tcs := []testCase{
		{W1: core.WeaponCannon, W2: core.WeaponCannon, Out: "C:\\Users\\user\\Desktop\\cvc.csv"},
		{W1: core.WeaponCannon, W2: core.WeaponRockets, Out: "C:\\Users\\user\\Desktop\\cvr.csv"},
		{W1: core.WeaponCannon, W2: core.WeaponArtillery, Out: "C:\\Users\\user\\Desktop\\cva.csv"},

		{W1: core.WeaponRockets, W2: core.WeaponCannon, Out: "C:\\Users\\user\\Desktop\\rvc.csv"},
		{W1: core.WeaponRockets, W2: core.WeaponRockets, Out: "C:\\Users\\user\\Desktop\\rvr.csv"},
		{W1: core.WeaponRockets, W2: core.WeaponArtillery, Out: "C:\\Users\\user\\Desktop\\rva.csv"},

		{W1: core.WeaponArtillery, W2: core.WeaponCannon, Out: "C:\\Users\\user\\Desktop\\avc.csv"},
		{W1: core.WeaponArtillery, W2: core.WeaponRockets, Out: "C:\\Users\\user\\Desktop\\avr.csv"},
		{W1: core.WeaponArtillery, W2: core.WeaponArtillery, Out: "C:\\Users\\user\\Desktop\\ava.csv"},
	}

	for _, tc := range tcs {
		println()
		println(tc.Out)
		WriteFile(tc.W1, tc.W2, tc.Out)
	}
}

// WriteFile simulate the test case and write the data to a file
func WriteFile(w1, w2, out string) {

	// simulate and sort (left points)
	results := MassSimulation(w1, w2)
	sort.SliceStable(results, func(i, j int) bool {
		ip := results[i].Points()
		ii := results[i].Iteration
		jp := results[j].Points()
		ji := results[j].Iteration
		if ip == jp {
			return ii < ji
		} else {
			return ip > jp
		}
	})

	// write data
	fh, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer func(fh *os.File) {
		_ = fh.Close()
	}(fh)

	_, _ = fh.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s\n", "Points", "A1", "D1", "S1", "A2", "D2", "S2", "Iteration"))
	for _, r := range results {
		_, err = fh.WriteString(fmt.Sprintf("%d;%d;%d;%d;%d;%d;%d;%d\n", r.Points(), r.A1, r.D1, r.S1, r.A2, r.D2, r.S2, r.Iteration))
		if err != nil {
			panic(err)
		}
	}
}
