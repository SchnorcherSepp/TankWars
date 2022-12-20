package maps

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"math/rand"
)

// InitRandomWorld generates a flat world with random tanks every were
func InitRandomWorld(w *core.World, seed int64, macro func(t *core.Tank)) {
	rnd := rand.New(rand.NewSource(seed))

	wpList := []string{core.WeaponCannon, core.WeaponRockets, core.WeaponArtillery, core.WeaponNone}
	aList := []int{core.North, core.Northeast, core.East, core.Southeast, core.South, core.Southwest, core.West, core.Northwest}
	oList := []string{core.RedTank, core.BlueTank}

	// init random world
	for i := 0; i < 200; i++ {
		// random values
		rndArmor := rnd.Intn(core.TankMaxArmor+core.TankMinArmor) - core.TankMinArmor
		rndDamage := rnd.Intn(core.TankMaxDamage+core.TankMinDamage) - core.TankMinDamage
		rndOwner := oList[rnd.Intn(len(oList))]
		rndWeapon := wpList[rnd.Intn(len(wpList))]
		rndAngle := aList[rnd.Intn(len(aList))]
		rndPos := core.NewPosition(rnd.Intn(w.ScreenWidth()-100), rnd.Intn(w.ScreenHeight()-100))

		// add tank to world
		nt, err := core.NewTank(w, rndOwner, rndArmor, rndDamage, rndWeapon)
		if err == nil {
			nt.SetPosition(rndPos, rndAngle)
			w.AddTank(nt)
			nt.SetMacro(macro)
		}
	}
}
