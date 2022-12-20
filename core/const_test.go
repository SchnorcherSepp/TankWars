package core

import "testing"

func Test_Tank_Attr(t *testing.T) {

	// budget limit
	if float64(TankBudget-TankMaxArmor-TankMaxDamage-TankMinSpeed) > -30 {
		t.Error("TankBudget too height")
	}
	if float64(TankBudget-TankMinArmor-TankMinDamage-TankMinSpeed) < 30 {
		t.Error("TankBudget too low")
	}

	// max armor
	if float64(TankBudget-TankMinSpeed-TankMinDamage) < TankMaxArmor {
		t.Error("max armor not possible")
	}

	// max damage
	if float64(TankBudget-TankMinArmor-TankMinSpeed) < TankMaxDamage {
		t.Error("max damage not possible")
	}

	// max speed
	if float64(TankBudget-TankMinArmor-TankMinDamage) < 80 {
		t.Error("max speed not possible")
	}

	// check MovePerTick AND TankMinSpeed
	tank, err := NewTank(nil, "test", 50, 50-TankMinSpeed, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	tank.SetPosition(NewPosition(1000, 1000), Southeast)
	tank.Forward()
	tank.Update()
	tank.Update()
	tank.Update()
	if tank.pos.X == 100 {
		t.Errorf("tank can't move with TankMinSpeed=%d and MovePerTick=%f", TankMinSpeed, MovePerTick)
	}
}
