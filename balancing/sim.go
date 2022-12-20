package main

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"github.com/SchnorcherSepp/TankWars/macro"
	"runtime"
	"sync"
	"time"
)

// MassSimulation starts hundreds of single duels.
// All variants (armor, damage, speed) are tested against each other.
func MassSimulation(w1, w2 string) []Result {
	start := time.Now()
	list := make([]Result, 0)
	listMux := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	// simulate all
	const steps = 3
	fmt.Printf("BUILD SIMULATION\n")

	// left options
	for a1 := core.TankMinArmor; a1 <= core.TankMaxArmor; a1 += steps {
		for d1 := core.TankMinDamage; d1 <= core.TankMaxDamage; d1 += steps {
			if a1+d1+core.TankMinSpeed > core.TankBudget {
				continue // invalid config for left
			}
			// right options
			for a2 := core.TankMinArmor; a2 <= core.TankMaxArmor; a2 += steps {
				for d2 := core.TankMinDamage; d2 <= core.TankMaxDamage; d2 += steps {
					if a2+d2+core.TankMinSpeed > core.TankBudget {
						continue // invalid config for right
					}
					// battle tests
					wg.Add(1)
					go func(a1, d1 int, w1 string, a2, d2 int, w2 string) {
						defer wg.Done()
						//------------------------------------------------------
						// sim
						world := core.NewWorld(core.WorldXWidth, core.WorldYHeight)
						result, err := BattleTest(world, true, a1, d1, w1, a2, d2, w2)
						if err != nil {
							return
						}
						// add to list
						listMux.Lock()
						list = append(list, result)
						listMux.Unlock()
						//------------------------------------------------------
					}(a1, d1, w1, a2, d2, w2)
				}
			}
		}

		fmt.Printf("%d/%d (%d)\n", a1-core.TankMinArmor+1, core.TankMaxArmor-core.TankMinArmor+1, runtime.NumGoroutine())
	}

	fmt.Printf("WAIT SIMULATION\n")
	wg.Wait()
	fmt.Printf("FIN SIMULATION after %s\n", time.Since(start))

	return list
}

// BattleTest simulates a duel. The side with the shorter range always attacks.
// The duel ends after when one participant is destroyed, or after 3600 iteration.
func BattleTest(world *core.World, update bool, a1, d1 int, w1 string, a2, d2 int, w2 string) (Result, error) {
	resources.MuteSound = true

	// left
	left, err := core.NewTank(world, core.BlueTank, a1, d1, w1)
	if err != nil {
		return Result{}, err
	}
	world.AddTank(left)
	left.SetPosition(core.NewPosition(100, 100), core.East)

	// right
	right, err := core.NewTank(world, core.RedTank, a2, d2, w2)
	if err != nil {
		return Result{}, err
	}
	world.AddTank(right)
	right.SetPosition(core.NewPosition(100+left.Weapon().Range()+right.Weapon().Range(), 100), core.West)

	// run test
	for {
		// macros
		if left.Weapon().Range() >= right.Weapon().Range() {
			macro.GuardMode(left, "")
			macro.AttackMove(right, "")
		} else {
			macro.AttackMove(left, "")
			macro.GuardMode(right, "")
		}

		// exit
		var winner string
		if !left.Alive() {
			winner = "right"
			for len(world.Projectiles()) > 0 {
				world.Update()
			}
		} else if !right.Alive() {
			winner = "left"
			for len(world.Projectiles()) > 0 {
				world.Update()
			}
		} else if world.Iteration() > 5000 {
			winner = "none"
		}

		if winner != "" {
			return Result{
				A1: a1, D1: d1, S1: left.Speed(), W1: w1, A2: a2, D2: d2, S2: right.Speed(), W2: w2,
				Winner: winner, LeftHP: left.Health(), RightHP: right.Health(), Iteration: world.Iteration(),
			}, nil
		}

		// update
		if update {
			world.Update()
		}
	}
}
