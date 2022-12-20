package goai

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/remote"
	"time"
)

//####################################################################################################################//
//###############  MY AI  ############################################################################################//
//####################################################################################################################//

// MyAI run demo ai
func MyAI() {
	// connect to server
	client := remote.NewTcpClient("localhost", "3333")

	// get my player
	me := client.MyName()

	// auto request world status
	var world remote.JsonWorld
	go func() {
		for {
			resp := client.GameStatus() // request game status
			world.Set(resp)             // update world
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// wait for world data
	for len(world.Tanks) < 5 {
		time.Sleep(50 * time.Millisecond)
	}

	// ===== START AI =====
	println("===== START AI =====")

	// find my tanks
	//---------------
	var tankID string
	var rocketID string

	for _, tank := range world.Tanks {
		if tank.Owner == me {
			// I am the owner of this tank
			if tank.Weapon.Typ == core.WeaponCannon {
				// found battle tank
				tankID = tank.ID
			}
			if tank.Weapon.Typ == core.WeaponRockets {
				// found rocket launcher
				rocketID = tank.ID
			}
		}
	}

	// get center position
	// (use world.ScreenWidth and world.ScreenHeight)
	//------------------------------------------------
	centerX := world.ScreenWidth / 2
	centerY := world.ScreenHeight / 2
	if me == core.RedTank {
		centerX -= 150
	} else {
		centerX += 150
	}

	// use MoveTo macro to move tank to the center
	//---------------------------------------------
	resp := client.SetMacroMoveTo(tankID, centerX, centerY)
	println("move tank", resp)

	resp = client.SetMacroMoveTo(rocketID, centerX, centerY)
	println("move rocket", resp)

	// wait for tanks
	//----------------
	for world.Iteration < 800 {
		time.Sleep(50 * time.Millisecond)
	}

	// use GuardMode to defend the center
	//------------------------------------
	resp = client.SetMacro(tankID, core.MacroGuardMode)
	println("macro tank", resp)

	resp = client.SetMacro(rocketID, core.MacroGuardMode)
	println("macro rocket", resp)

	// wait for more money
	//------------------------
	for world.CashRed < 101 && world.CashBlue < 101 {
		time.Sleep(50 * time.Millisecond)
	}

	// buy Artillery
	//---------------
	resp = client.BuyTank(5, 70, core.WeaponArtillery)
	println("buy", resp)
	time.Sleep(250 * time.Millisecond)

	var artilleryID string

	for _, tank := range world.Tanks {
		if tank.Owner == me {
			// I am the owner of this tank
			if tank.Weapon.Typ == core.WeaponArtillery {
				// found battle tank
				artilleryID = tank.ID
			}
		}
	}

	// use MacroAttackMove on Artillery to find a way to the enemy
	//-------------------------------------------------------------
	resp = client.SetMacro(artilleryID, core.MacroAttackMove)
	println("macro artillery", resp)

	resp = client.SetMacro(tankID, core.MacroAttackMove)
	println("macro tank", resp)

	resp = client.SetMacro(rocketID, core.MacroAttackMove)
	println("macro rocket", resp)

	// buy more tanks (loops)
	//------------------------
	for {
		client.BuyTank(5, 70, core.WeaponArtillery)
		time.Sleep(250 * time.Millisecond)

		for _, tank := range world.Tanks {
			if tank.Owner == me && !tank.ActiveMacro {
				client.SetMacro(tank.ID, core.MacroAttackMove)
			}
		}
	}
}
