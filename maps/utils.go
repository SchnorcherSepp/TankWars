package maps

import "github.com/SchnorcherSepp/TankWars/core"

// build a building on position x,y.
// A building is a tank with a specific owner string.
// The GUI function drawBuilding() draw this differently.
//
// see core.NeutralRock, core.RedRock, core.BlueRock, core.RedBase and core.BlueBase.
func build(world *core.World, building string, x, y int) *core.Tank {

	// create tank (building)
	b, err := core.NewTank(world, building, core.TankMaxArmor, core.TankMinDamage, core.WeaponNone)
	if err != nil {
		panic(err) // can't happen
	}

	// set in world
	b.SetPosition(core.NewPosition(x, y), core.North)
	world.AddTank(b)

	// return
	return b
}
