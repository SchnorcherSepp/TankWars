package maps

import "github.com/SchnorcherSepp/TankWars/core"

// InitOpenField is an open field for big tank battles.
// There is very little coverage.
func InitOpenField(world *core.World) (t1, a1, t2, a2 *core.Tank) {
	sw := world.ScreenWidth()
	sh := world.ScreenHeight()

	// set player cash
	//-----------------
	world.SetCash(50, 50)

	// build base 1
	//------------------------
	build(world, core.RedBase, 150, sh/2+200)

	// build base 2 (mirror)
	//------------------------
	build(world, core.BlueBase, sw-150, sh-(sh/2+200))

	// build neutral
	//------------------------
	build(world, core.NeutralRock, 400, 750)
	build(world, core.NeutralRock, 325, 785)
	build(world, core.NeutralRock, 300, 840)

	build(world, core.NeutralRock, sw-400, sh-750)
	build(world, core.NeutralRock, sw-325, sh-785)
	build(world, core.NeutralRock, sw-300, sh-840)

	// Tanks
	//------------------------

	// battle tank 1
	var err error
	t1, err = core.NewTank(world, core.RedTank, 55, 20, core.WeaponCannon)
	if err != nil {
		panic(err)
	}
	t1.SetPosition(core.NewPosition(430, 340), core.Southeast)
	world.AddTank(t1)

	// battle tank 2
	t2, err = core.NewTank(world, core.BlueTank, 55, 20, core.WeaponCannon)
	if err != nil {
		panic(err)
	}
	t2.SetPosition(core.NewPosition(sw-430, sh-340), core.Northwest)
	world.AddTank(t2)

	// artillery 1
	a1, err = core.NewTank(world, core.RedTank, 15, 60, core.WeaponRockets)
	if err != nil {
		panic(err)
	}
	a1.SetPosition(core.NewPosition(260, 500), core.East)
	world.AddTank(a1)

	// artillery 2
	a2, err = core.NewTank(world, core.BlueTank, 15, 60, core.WeaponRockets)
	if err != nil {
		panic(err)
	}
	a2.SetPosition(core.NewPosition(sw-260, sh-500), core.West)
	world.AddTank(a2)

	// return
	return
}
