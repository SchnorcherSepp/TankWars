package maps

import "github.com/SchnorcherSepp/TankWars/core"

// InitFortress is a map with many defensive barriers and heavy artillery.
func InitFortress(world *core.World) (t1, a1, t2, a2 *core.Tank) {

	// build base 1
	//------------------------
	build(world, core.RedBase, 150, 150)

	build(world, core.RedRock, 470, 170)
	build(world, core.RedRock, 410, 210)
	build(world, core.RedRock, 370, 270)
	build(world, core.RedRock, 350, 330)
	build(world, core.RedRock, 310, 390)
	build(world, core.RedRock, 230, 420)

	// build base 2 (mirror)
	//------------------------
	sw := world.ScreenWidth()
	sh := world.ScreenHeight()

	build(world, core.BlueBase, sw-150, sh-150)

	build(world, core.BlueRock, sw-470, sh-170)
	build(world, core.BlueRock, sw-410, sh-210)
	build(world, core.BlueRock, sw-370, sh-270)
	build(world, core.BlueRock, sw-350, sh-330)
	build(world, core.BlueRock, sw-310, sh-390)
	build(world, core.BlueRock, sw-230, sh-420)

	// build neutral
	//------------------------
	build(world, core.NeutralRock, 400, 750)
	build(world, core.NeutralRock, 325, 785)
	build(world, core.NeutralRock, 300, 840)

	build(world, core.NeutralRock, sw/2, sh/2)

	build(world, core.NeutralRock, sw-400, sh-750)
	build(world, core.NeutralRock, sw-325, sh-785)
	build(world, core.NeutralRock, sw-300, sh-840)

	// Tanks
	//------------------------

	// battle tank 1
	var err error
	t1, err = core.NewTank(world, core.RedTank, 30, 32, core.WeaponCannon)
	if err != nil {
		panic(err)
	}
	t1.SetPosition(core.NewPosition(500, 100), core.Southeast)
	world.AddTank(t1)

	// battle tank 2
	t2, err = core.NewTank(world, core.BlueTank, 30, 32, core.WeaponCannon)
	if err != nil {
		panic(err)
	}
	t2.SetPosition(core.NewPosition(sw-500, sh-100), core.Northwest)
	world.AddTank(t2)

	// artillery 1
	a1, err = core.NewTank(world, core.RedTank, 15, 50, core.WeaponArtillery)
	if err != nil {
		panic(err)
	}
	a1.SetPosition(core.NewPosition(260, 300), core.East)
	world.AddTank(a1)

	// artillery 2
	a2, err = core.NewTank(world, core.BlueTank, 15, 50, core.WeaponArtillery)
	if err != nil {
		panic(err)
	}
	a2.SetPosition(core.NewPosition(sw-260, sh-300), core.West)
	world.AddTank(a2)

	// return
	return
}
