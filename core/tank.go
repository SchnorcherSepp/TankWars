package core

import (
	"fmt"
	"math"
	"sync/atomic"
)

// globalIdPool is used to generate unique IDs
var globalIdPool uint64 = 1234

// Tank is an object in World. It can be a tank, a building, a rock, ...
type Tank struct {
	// system (set by New())
	world  *World  // world ref to interact with other objects
	id     string  // unique id
	owner  string  // tank of the owner
	weapon *Weapon // weapon struct

	// base attribute
	health int // health points of the tank
	armor  int // armor reduces damage by this number; heavy armor reduces speed
	speed  int // speed depends on the weight of armor and weapon

	// move
	pos        Position // tank position
	command    int      // 1 is forward; 0 is stop; -1 is back
	angle      int      // 0 is north; 45 is northeast; 90 is east; ...
	isBlocked  bool     // movement has ended because the path was blocked
	lastRotate uint64   // iteration of the last rotate command

	// macro function
	macro func(t *Tank) // is called by update
}

// NewTank return a new tank.
// The tank must be added to the world manually (see World.AddTank()).
//
// The attributes armor, damage and speed are related.
// Every unused budget point (see TankBudget) is converted into speed.
//
//	TankMinArmor < armor < TankMaxArmor
//	TankMinDamage < damage < TankMaxDamage
//	TankMinSpeed < speed
func NewTank(world *World, owner string, armor, damage int, weapon string) (*Tank, error) {

	// calc cost
	speed := TankBudget - armor - damage

	// check speed
	if speed < TankMinSpeed {
		return nil, fmt.Errorf("with armor and weapons the speed would be %d but min. is %d", speed, TankMinSpeed)
	}
	// check armor
	if armor < TankMinArmor || armor > TankMaxArmor {
		return nil, fmt.Errorf("the armor must be between %d and %d", TankMinArmor, TankMaxArmor)
	}
	// check damage
	if damage < TankMinDamage || damage > TankMaxDamage {
		return nil, fmt.Errorf("the damage must be between %d and %d", TankMinDamage, TankMaxDamage)
	}

	// add speed factor for more fun!
	speedFactor := int(math.Pow(float64(speed-TankMinSpeed), 6) * 0.0000000007)
	speed += speedFactor

	// build tank
	t := &Tank{
		// system
		world:  world,
		id:     fmt.Sprintf("%d", atomic.AddUint64(&globalIdPool, 1)),
		owner:  owner,
		weapon: nil, // is set below

		// base attribute
		health: 100,
		armor:  armor,
		speed:  speed,

		// move
		command: 0,
		angle:   South,
	}

	// set weapon
	switch weapon {
	case WeaponNone:
		t.weapon = nil // no weapon
	case WeaponCannon:
		t.weapon = NewWeaponCannon(t.world, t, damage)
	case WeaponArtillery:
		t.weapon = NewWeaponArtillery(t.world, t, damage)
	case WeaponRockets:
		t.weapon = NewWeaponRocketLauncher(t.world, t, damage)
	default:
		return nil, fmt.Errorf("unknown weapon: %s", weapon)
	}

	// return
	return t, nil
}

//---------------- GETTER --------------------------------------------------------------------------------------------//

// ID is a unique id of this tank.
func (t *Tank) ID() string {
	return t.id
}

// Owner returns who control this object.
func (t *Tank) Owner() string {
	return t.owner
}

// Weapon of the tank. see NewWeaponCannon, NewWeaponArtillery.
func (t *Tank) Weapon() *Weapon {
	return t.weapon
}

// Health points of the tank. Is reduced by Hit(). Is used by Alive().
func (t *Tank) Health() int {
	return t.health
}

// Armor reduces damage with each hit.
// Armor can be so high that no damage is dealt.
func (t *Tank) Armor() int {
	return t.armor
}

// Speed of the tank.
func (t *Tank) Speed() int {
	return t.speed
}

// Pos returns the current position.
func (t *Tank) Pos() Position {
	return t.pos
}

// Command returns the current move command of the tank.
// 1 is forward; 0 is stop; -1 is backward
func (t *Tank) Command() int {
	return t.command
}

// Angle of the tank. see North, South, East, ...
func (t *Tank) Angle() int {
	return t.angle
}

// Alive returns true if the Health is not 0.
func (t *Tank) Alive() bool {
	return t.health > 0
}

// Moving returns the status true if the tank is moving.
// see Command()
func (t *Tank) Moving() bool {
	return t.command != 0
}

// Blocked return true if movement has ended because the path was blocked.
// Set by Update() and reset by Forward() and Backward().
func (t *Tank) Blocked() bool {
	return t.isBlocked
}

// LastRotate returns at which iteration the last rotation was.
func (t *Tank) LastRotate() uint64 {
	return t.lastRotate
}

// ActiveMacro returns if this tank is controlled by a macro.
// see SetMacro().
func (t *Tank) ActiveMacro() bool {
	return t.macro != nil
}

// Status returns the weapon status: (StatusMoving, StatusPreparing, StatusReloading, StatusReady or StatusNoWeapon).
// see Weapon.Status
func (t *Tank) Status() (rdy bool, status string) {
	if t.weapon != nil && t.Alive() {
		return t.weapon.Status()
	} else {
		return false, StatusNoWeapon
	}
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// Fire is a simple call for Weapon.Fire()
func (t *Tank) Fire(fireAngle, distance int) (success bool, txt string) {
	if t.weapon != nil && t.Alive() {
		return t.weapon.Fire(t.pos, t.angle, fireAngle, distance)
	} else {
		return false, StatusNoWeapon
	}
}

// FireAt is a wrapper for Fire() and convert the position to fireAngle and distance.
func (t *Tank) FireAt(pos Position) (success bool, txt string) {
	fireAngle := RelativeAngle(t.pos, pos)
	distance := Distance(t.pos, pos)
	return t.Fire(fireAngle, int(distance))
}

// Hit calculate the damage on a direct hit.
// Reduce the damage by Armor().
// Call Remove() for death tanks.
func (t *Tank) Hit(damage int) {

	// calc damage with armor
	damage -= t.armor
	if damage < 1 {
		damage = 1 // minimum damage
	}

	// remove HP
	t.health -= damage

	// check death
	if !t.Alive() {
		t.Remove()
	}
}

// Remove the tank from World.
func (t *Tank) Remove() {
	if t != nil && t.world != nil {
		// remove from list
		newTanks := make([]*Tank, 0, len(t.world.tanks))
		for _, tank := range t.world.tanks {
			if tank != t {
				newTanks = append(newTanks, tank)
			}
		}
		t.world.tanks = newTanks
		// kill
		t.health = -1
	}
}

// SetPosition set a new tank position without collision check.
func (t *Tank) SetPosition(pos Position, angle int) {
	t.pos = pos
	t.angle = angle
}

// SetMacro sets a macro that is called with every update.
// Remove it with 'nil'.
func (t *Tank) SetMacro(macro func(t *Tank)) {
	t.macro = macro
}

//---------------- MOVE (Setter) -------------------------------------------------------------------------------------//

// Forward send the tank forward.
// see Command().
func (t *Tank) Forward() {
	if t.weapon != nil {
		t.weapon.Update(true)
	}
	t.command = 1
	t.isBlocked = false // reset
}

// Backward send the tank back.
// see Command().
func (t *Tank) Backward() {
	if t.weapon != nil {
		t.weapon.Update(true)
	}
	t.command = -1
	t.isBlocked = false // reset
}

// Stop the movement.
// Weapons can only build up when the tank is stationary
// see Command().
// see Weapon.PreparationTime()
func (t *Tank) Stop() {
	t.command = 0
}

// Left turn the tank direction 45° left.
// see Angle()
func (t *Tank) Left() (success bool, status string) {
	return t.rotate(-45)
}

// Right turn the tank direction 45° right.
// see Angle()
func (t *Tank) Right() (success bool, status string) {
	return t.rotate(45)
}

// rotate is a helper methode for Left() and Right().
// Manipulate the angle and set valid values (North, South, East, ...)
func (t *Tank) rotate(a int) (success bool, status string) {
	// check last rotation
	if t.world != nil && t.LastRotate()+TankRotationDelay > t.world.iteration {
		return false, StatusPreparing
	}

	// new angle
	na := t.angle
	na += a
	if na < 0 {
		na += 360
	}
	na %= 360 // modulo (Non-negative only)

	// set new angle and lastRotate
	t.angle = na
	if t.world != nil {
		t.lastRotate = t.world.iteration
	}

	// rotation is move!
	if t.weapon != nil {
		t.weapon.Update(true)
	}

	// success
	return true, StatusReady
}

//---------------- UPDATE --------------------------------------------------------------------------------------------//

// Update calculate movement, check collisions, check world borders
// and set Weapon.LastMove().
func (t *Tank) Update() {

	// update weapon (move == lock)
	if t.weapon != nil {
		t.weapon.Update(t.command != 0) // update with every tick
	}

	// update position
	if t.command != 0 {
		// move (update)
		oldPos := t.pos
		t.pos.Update(t.angle, t.speed*t.command)

		// check borders
		if t.world != nil && CheckBorders(t.pos, BlockRadius, t.world.ScreenWidth(), t.world.ScreenHeight()) {
			t.pos = oldPos     // reset position
			t.Stop()           // stop movement
			t.isBlocked = true // set status
		}

		// check other tanks
		if t.world != nil {
			for _, ot := range t.world.tanks {
				if t != ot && IsCollided(t.pos, BlockRadius, ot.pos, BlockRadius) {
					t.pos = oldPos     // reset position
					t.Stop()           // stop movement
					t.isBlocked = true // set status
					break
				}
			}
		}
	}

	// call macro function
	if t.macro != nil {
		t.macro(t)
	}
}
