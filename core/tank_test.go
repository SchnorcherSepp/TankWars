package core

import (
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"testing"
)

func TestTank_Getter(t *testing.T) {
	w := NewWorld(333, 444)
	w.UpdateN(500)

	// tank
	tank, err := NewTank(w, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	tank.SetPosition(NewPosition(13, 17), West)

	// test tank 2
	tank2, _ := NewTank(w, "TestOwner2", 24, 18, WeaponArtillery)

	if tank.ID() == tank2.ID() {
		t.Error("wrong value")
	}
	if tank.Owner() != "TestOwner" {
		t.Error("wrong value")
	}
	if tank.Weapon() == nil {
		t.Error("wrong value")
	}
	if tank.Health() != 100 {
		t.Error("wrong value")
	}
	if tank.Armor() != 11 {
		t.Error("wrong value")
	}
	if tank.Speed() != 75 { // speed 70 + bonus
		t.Errorf("wrong value: %d", tank.Speed())
	}
	if tank.Pos().X != 13 || tank.Pos().Y != 17 {
		t.Error("wrong value")
	}
	if tank.Command() != 0 {
		t.Error("wrong value")
	}
	if tank.Angle() != West {
		t.Error("wrong value")
	}
	if tank.Alive() != true {
		t.Error("wrong value")
	}
	if tank.Moving() != false {
		t.Error("wrong value")
	}
	if rdy, txt := tank.Status(); rdy != false || txt != StatusPreparing {
		t.Error("wrong value", txt)
	}

	// no weapon status
	tank.weapon = nil
	if rdy, txt := tank.Status(); rdy != false || txt != StatusNoWeapon {
		t.Error("wrong value")
	}
}

func TestNewTank(t *testing.T) {
	// ok
	if _, err := NewTank(nil, "", TankMinArmor, TankMinDamage, WeaponArtillery); err != nil {
		t.Error("wrong value")
	}
	// wrong weapon
	if _, err := NewTank(nil, "", TankMinArmor, TankMinDamage, "none"); err == nil {
		t.Error("wrong value")
	}
	// min damage
	if _, err := NewTank(nil, "", TankMinArmor, TankMinDamage-1, WeaponArtillery); err == nil {
		t.Error("wrong value")
	}
	// max damage
	if _, err := NewTank(nil, "", TankMinArmor, TankMaxDamage+1, WeaponArtillery); err == nil {
		t.Error("wrong value")
	}
	// min armor
	if _, err := NewTank(nil, "", TankMinArmor-1, TankMinDamage, WeaponArtillery); err == nil {
		t.Error("wrong value")
	}
	// max armor
	if _, err := NewTank(nil, "", TankMaxArmor+1, TankMinDamage, WeaponArtillery); err == nil {
		t.Error("wrong value")
	}
	// min speed
	if _, err := NewTank(nil, "", TankMaxArmor, TankMaxDamage, WeaponArtillery); err == nil {
		t.Error("wrong value")
	}
}

func TestTank_Hit(t *testing.T) {
	// tank
	tank, err := NewTank(nil, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}

	// hit tank
	tank.Hit(0)
	if tank.Health() != 99 { // minimum damage 1
		t.Error("wrong value")
	}
	tank.Hit(11)
	if tank.Health() != 98 { // minimum damage 1
		t.Error("wrong value")
	}
	tank.Hit(20)
	if tank.Health() != 98-(20-11) {
		t.Error("wrong value")
	}
	tank.Hit(90000)
	if tank.Health() >= 0 || tank.Alive() {
		t.Error("wrong value")
	}
}

func TestTank_Remove(t *testing.T) {
	// tank
	w := NewWorld(9999, 9999)
	w.UpdateN(500)

	tank, err := NewTank(w, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	w.AddTank(tank)

	// test tank 2
	tank2, _ := NewTank(w, "TestOwner2", 22, 33, WeaponArtillery)
	w.AddTank(tank2)

	// check remove
	if len(w.tanks) != 2 {
		t.Error("wrong value")
	}
	tank.Remove()
	if len(w.tanks) != 1 {
		t.Error("wrong value")
	}
}

func TestTank_Fire(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	w := NewWorld(333, 444)
	w.UpdateN(500)

	// tank
	tank, err := NewTank(w, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}

	// fire:  error preparing
	if ss, err := tank.Fire(9999, 0); ss != false || err != StatusPreparing {
		t.Error("wrong value", ss, err)
	}

	// fire:  ok
	w.UpdateN(100)
	if ss, err := tank.Fire(9999, 0); ss != true || err != StatusReady {
		t.Error("wrong value", ss, err)
	}

	// fire:  no weapon
	tank.weapon = nil
	if ss, err := tank.Fire(9999, 0); ss != false || err != StatusNoWeapon {
		t.Error("wrong value", ss, err)
	}

}

func TestTank_FireAt(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	w := NewWorld(333, 444)
	w.UpdateN(500)

	// tank
	tank, _ := NewTank(w, "TestOwner", 5, 70, WeaponArtillery)
	tank.SetPosition(NewPosition(100, 100), North)
	w.AddTank(tank)
	// other
	other, _ := NewTank(w, "TestOwner", 5, 70, WeaponCannon)
	other.SetPosition(NewPosition(177, 290), North)
	w.AddTank(other)

	//-------------------------------------------

	if !other.Alive() {
		t.Error("other not alive")
	}

	for i := 0; i < 500; i++ {
		tank.FireAt(NewPosition(177, 290))
		w.Update()
		if !other.Alive() {
			break
		}
	}

	if other.Alive() {
		t.Error("other alive")
	}
}

func TestTank_Left(t *testing.T) {
	w := NewWorld(333, 444)
	w.UpdateN(500)

	// tank
	tank, err := NewTank(w, "TestOwner", 11, 19, WeaponNone)
	if err != nil {
		t.Fatal(err)
	}
	tank.SetPosition(NewPosition(0, 0), North)

	// LEFT
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != Northwest {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != West {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != Southwest {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != South {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != Southeast {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != East {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != Northeast {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != North {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Left(); tank.Angle() != Northwest {
		t.Error("wrong value", tank.Angle())
	}
}

func TestTank_Right(t *testing.T) {
	w := NewWorld(333, 444)
	w.UpdateN(500)

	// tank
	tank, err := NewTank(w, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	tank.SetPosition(NewPosition(100, 100), North)

	// Right
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != Northeast {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != East {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != Southeast {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != South {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != Southwest {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != West {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != Northwest {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != North {
		t.Error("wrong value", tank.Angle())
	}
	w.UpdateN(TankRotationDelay)
	if tank.Right(); tank.Angle() != Northeast {
		t.Error("wrong value", tank.Angle())
	}
}

func TestTank_Rotation(t *testing.T) {
	w := NewWorld(222, 333)
	w.UpdateN(500)

	// tank
	tank, err := NewTank(w, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}
	tank.SetPosition(NewPosition(100, 100), North)

	// a rotation is a move!
	if tank.weapon.LastMove() != 500 {
		t.Error("wrong value", tank.weapon.LastMove())
	}
	w.UpdateN(100)
	if tank.weapon.LastMove() != 500 {
		t.Error("wrong value", tank.weapon.LastMove())
	}
	tank.Left()
	if tank.weapon.LastMove() != 600 {
		t.Error("wrong value", tank.weapon.LastMove())
	}

	// check TankRotationDelay
	w.UpdateN(TankRotationDelay)
	if suc, txt := tank.Left(); suc != true || txt != StatusReady {
		t.Error("wrong value")
	}
	w.UpdateN(TankRotationDelay - 1)
	if suc, txt := tank.Left(); suc != false || txt != StatusPreparing {
		t.Error("wrong value")
	}
	w.UpdateN(1)
	if suc, txt := tank.Left(); suc != true || txt != StatusReady {
		t.Error("wrong value")
	}
}

func TestTank_Command(t *testing.T) {
	tank, err := NewTank(nil, "TestOwner", 11, 19, WeaponCannon)
	if err != nil {
		t.Fatal(err)
	}

	if tank.Command() != 0 {
		t.Error("wrong value", tank.Command())
	}
	tank.Forward()
	if tank.Command() != 1 {
		t.Error("wrong value", tank.Command())
	}
	tank.Stop()
	if tank.Command() != 0 {
		t.Error("wrong value", tank.Command())
	}
	tank.Backward()
	if tank.Command() != -1 {
		t.Error("wrong value", tank.Command())
	}
	tank.Forward()
	if tank.Command() != 1 {
		t.Error("wrong value", tank.Command())
	}
}

func TestTank_CheckBorders(t *testing.T) {
	world := NewWorld(2000, 1000)
	tank, err := NewTank(world, "tester", 20, 40, WeaponArtillery)
	if err != nil {
		t.Error(err)
	}

	// check top
	tank.SetPosition(NewPosition(500, 0+BlockRadius+1), North)
	tank.Forward()                                                                                 // start
	if tank.pos.X != 500 || tank.pos.Y != 0+BlockRadius+1 || tank.command != 1 || tank.Blocked() { // no move
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != 500 || tank.pos.Y != 0+BlockRadius+0 || tank.command != 1 || tank.Blocked() { // next step: no collision
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != 500 || tank.pos.Y != 0+BlockRadius+0 || tank.command != 0 || !tank.Blocked() { // next step: collision! & same position
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}

	// check bottom
	tank.SetPosition(NewPosition(500, world.ScreenHeight()-BlockRadius-1), South)
	tank.Forward()                                                                                                    // start
	if tank.pos.X != 500 || tank.pos.Y != world.ScreenHeight()-BlockRadius-1 || tank.command != 1 || tank.Blocked() { // no move
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != 500 || tank.pos.Y != world.ScreenHeight()-BlockRadius-0 || tank.command != 1 || tank.Blocked() { // next step: no collision
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != 500 || tank.pos.Y != world.ScreenHeight()-BlockRadius-0 || tank.command != 0 || !tank.Blocked() { // next step: collision! & same position
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}

	// check left
	tank.SetPosition(NewPosition(0+BlockRadius+1, 500), West)
	tank.Forward()                                                                                 // start
	if tank.pos.X != 0+BlockRadius+1 || tank.pos.Y != 500 || tank.command != 1 || tank.Blocked() { // no move
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != 0+BlockRadius+0 || tank.pos.Y != 500 || tank.command != 1 || tank.Blocked() { // next step: no collision
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != 0+BlockRadius+0 || tank.pos.Y != 500 || tank.command != 0 || !tank.Blocked() { // next step: collision! & same position
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}

	// check right
	tank.SetPosition(NewPosition(world.ScreenWidth()-BlockRadius-1, 500), East)
	tank.Forward()                                                                                                   // start
	if tank.pos.X != world.ScreenWidth()-BlockRadius-1 || tank.pos.Y != 500 || tank.command != 1 || tank.Blocked() { // no move
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != world.ScreenWidth()-BlockRadius-0 || tank.pos.Y != 500 || tank.command != 1 || tank.Blocked() { // next step: no collision
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
	tank.Update()
	if tank.pos.X != world.ScreenWidth()-BlockRadius+0 || tank.pos.Y != 500 || tank.command != 0 || !tank.Blocked() { // next step: collision! & same position
		t.Error("wrong value", tank.pos.X, tank.pos.Y, tank.command)
	}
}

func TestTank_CheckBorders_OtherTank(t *testing.T) {
	world := NewWorld(2000, 1000)

	tank1, err := NewTank(world, "tester", 20, 40, WeaponRockets)
	if err != nil {
		t.Error(err)
	}
	tank1.SetPosition(NewPosition(100, 100), East)
	tank1.Forward()
	world.AddTank(tank1)

	tank2, err := NewTank(world, "tester", 20, 40, WeaponArtillery)
	if err != nil {
		t.Error(err)
	}
	tank2.SetPosition(NewPosition(300, 100), West)
	tank2.Forward()
	world.AddTank(tank2)

	// check collision
	if tank1.isBlocked || tank2.isBlocked {
		t.Error("wrong value")
	}
	world.Update()
	if tank1.isBlocked || tank2.isBlocked {
		t.Error("wrong value")
	}
	for i := 0; i < 100; i++ {
		world.Update()
	}
	if !tank1.isBlocked || !tank2.isBlocked {
		t.Error("wrong value")
	}
}

func TestTank_SetMacro(t *testing.T) {
	nt, _ := NewTank(nil, "ss", 11, 22, WeaponNone)

	if nt.ActiveMacro() != false {
		t.Error("wrong value")
	}

	nt.SetMacro(func(t *Tank) {
		t.id = "hallo"
	})

	if nt.ActiveMacro() != true {
		t.Error("wrong value")
	}

	if nt.ID() == "hallo" {
		t.Error("wrong value")
	}
	nt.Update()
	if nt.ID() != "hallo" {
		t.Error("wrong value")
	}
}
