package maps

import "github.com/SchnorcherSepp/TankWars/core"

// InitTest is an empty field for weapon tests with unlimited cash.
func InitTest(world *core.World) (t1, a1, t2, a2 *core.Tank) {
	sw := world.ScreenWidth()
	sh := world.ScreenHeight()

	// set player cash
	//-----------------
	world.SetCash(1000000, 1000000)

	// build base 1
	//------------------------
	build(world, core.RedBase, 200, sh/2)

	// build base 2 (mirror)
	//------------------------
	build(world, core.BlueBase, sw-200, sh/2)

	// return
	return
}
