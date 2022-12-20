package core

import (
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"math"
)

// Weapon can be mounted on vehicles or buildings.
// It generates bullets with Fire().
type Weapon struct {
	// world ref to interact with other objects
	world  *World
	parent *Tank

	// attributes
	typ           string
	rng           int
	prepTime      uint64
	reloadTime    uint64
	projSpeed     int
	damage        int
	aoeRadius     int
	projCollision bool
	anyFireAngle  bool

	// update
	lastMove uint64
	lastFire uint64
}

// NewWeaponCannon return the weapon for a battle tank.
// It's fast and collide with other tanks. The damage is very height on single target (+50%).
func NewWeaponCannon(world *World, parent *Tank, damage int) *Weapon {
	nw := &Weapon{
		world:         world,
		parent:        parent,
		typ:           WeaponCannon,
		rng:           338,                     // weapon range
		prepTime:      800 * GameSpeed / 1000,  // ~ 800 ms
		reloadTime:    3000 * GameSpeed / 1000, // ~ 3000 ms
		projSpeed:     600,
		damage:        int(math.Round(float64(damage) * 1.5)), // correction factor
		aoeRadius:     0,
		projCollision: true,
		anyFireAngle:  false,
		lastMove:      0, // set later
	}
	if world != nil {
		nw.lastMove = world.iteration // simulate move for preparation timer
	}
	return nw
}

// NewWeaponArtillery return the weapon for an artillery.
// It's slow and can't collide with other tanks. Explode at the destination with big aoe but less damage (-10%).
func NewWeaponArtillery(world *World, parent *Tank, damage int) *Weapon {
	nw := &Weapon{
		world:         world,
		parent:        parent,
		typ:           WeaponArtillery,
		rng:           736,                     // weapon range
		prepTime:      8000 * GameSpeed / 1000, // ~ 8000 ms
		reloadTime:    6000 * GameSpeed / 1000, // ~ 6000 ms
		projSpeed:     300,
		damage:        int(math.Round(float64(damage) * 0.9)), // correction factor
		aoeRadius:     1.75 * BlockRadius,
		projCollision: false,
		anyFireAngle:  true,
		lastMove:      0, // set later
	}
	if world != nil {
		nw.lastMove = world.iteration // simulate move for preparation timer
	}
	return nw
}

// NewWeaponRocketLauncher return the weapon for a rocket launcher.
// It's like the artillery but the damage is used to reduce the reload time.
func NewWeaponRocketLauncher(world *World, parent *Tank, damage int) *Weapon {
	minRld := 600
	maxRld := 2000
	maxDmg := 20
	nw := &Weapon{
		world:         world,
		parent:        parent,
		typ:           WeaponRockets,
		rng:           387,                                                                                                    // weapon range
		prepTime:      8000 * GameSpeed / 1000,                                                                                // ~ 8000 ms
		reloadTime:    uint64(maxRld-(maxRld-minRld)/(TankMaxDamage-TankMinDamage)*(damage-TankMinDamage)) * GameSpeed / 1000, // mod with damage
		projSpeed:     300,
		damage:        int(float64(damage) / TankMaxDamage * float64(maxDmg)), // RocketLauncher use damage to reduce the reload time
		aoeRadius:     1.00 * BlockRadius,
		projCollision: false,
		anyFireAngle:  true,
		lastMove:      0, // set later
	}
	if world != nil {
		nw.lastMove = world.iteration // simulate move for preparation timer
	}
	return nw
}

//---------------- GETTER --------------------------------------------------------------------------------------------//

// Type is the weapon name.
// see: WeaponCannon or WeaponArtillery.
func (w *Weapon) Type() string {
	return w.typ
}

// Range is the maximum range of the weapon.
// WeaponCannon disappear after the range.
// WeaponArtillery explode after the range.
//
// see Projectile.Collision()
func (w *Weapon) Range() int {
	return w.rng
}

// PreparationTime is the time (=iteration ticks) the weapon has to build up after a movement.
// see status StatusPreparing
func (w *Weapon) PreparationTime() uint64 {
	return w.prepTime
}

// ReloadTime is the time (=iteration ticks) the weapon has to be reloaded after firing.
// see status StatusReloading
func (w *Weapon) ReloadTime() uint64 {
	return w.reloadTime
}

// ProjectileSpeed is the velocity of the bullet.
func (w *Weapon) ProjectileSpeed() int {
	return w.projSpeed
}

// Damage is the base damage on hits.
// Is reduced by armor.
// see Tank.Hit() for calculation.
func (w *Weapon) Damage() int {
	return w.damage
}

// AoERadius is the area damage radius. Objects within the radius will be damaged even without a direct hit.
// see: Projectile.Explode() for calculation.
func (w *Weapon) AoERadius() int {
	return w.aoeRadius
}

// ProjectileCollision if is active, the projectile will be stopped by objects in the fly path and explodes.
// Otherwise, it will explode after the specified range.
//
//	true - only explodes on contact
//	false - only explodes at the end of the distance
func (w *Weapon) ProjectileCollision() bool {
	return w.projCollision
}

// AnyFireAngle enables any launch direction.
// Otherwise, the projectile always fires in the tank direction.
func (w *Weapon) AnyFireAngle() bool {
	return w.anyFireAngle
}

// LastMove returns at witch iteration the last move was.
// see status StatusPreparing
func (w *Weapon) LastMove() uint64 {
	return w.lastMove
}

// LastFire returns at witch iteration the last fire was.
// see status StatusReloading
func (w *Weapon) LastFire() uint64 {
	return w.lastFire
}

// Status returns true if the weapon is ready to fire.
// Otherwise, the reason is returned (StatusMoving, StatusPreparing, StatusReloading or StatusReady).
// see: PreparationTime() and ReloadTime().
func (w *Weapon) Status() (rdy bool, status string) {

	// STATUS: StatusReloading after fire
	if w.world != nil && w.lastFire+w.reloadTime > w.world.iteration {
		return false, StatusReloading
	}

	// STATUS: StatusMoving
	if w.parent != nil && w.parent.Moving() {
		return false, StatusMoving
	}

	// STATUS: StatusPreparing after move
	if w.world != nil && w.lastMove+w.prepTime > w.world.iteration {
		return false, StatusPreparing
	}

	// STATUS: ready
	return true, StatusReady
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// Fire creates a new projectile.
// The attributes fireAngle and distance determine the direction and distance of the shot.
//
//	If AnyFireAngle() is false, attribute fireAngle is overridden by attribute vehicleAngle.
//	If ProjectileCollision() is true, attribute distance is overridden by Range().
//	Attribute distance is limited by Range().
//
// see NewProjectile()
func (w *Weapon) Fire(vehiclePos Position, vehicleAngle, fireAngle, distance int) (success bool, status string) {
	// check status
	success, status = w.Status()
	if !success {
		return // can't fire; return error
	}

	// max range limit
	if w.rng < distance {
		distance = w.rng
	}

	// min range limit (Cannons)
	if w.projCollision {
		distance = w.rng
	}

	// FireAngle limit
	if !w.anyFireAngle {
		fireAngle = vehicleAngle
	}

	// weapon is ready; lunch projectile
	pj := NewProjectile(w.world, w.parent, vehiclePos, fireAngle, distance, w.projSpeed, w.damage, w.aoeRadius, w.projCollision)
	if w.world != nil && w.world.projectiles != nil {
		w.world.projectiles = append(w.world.projectiles, pj)
	}

	// play sound
	resources.PlaySound(resources.Sounds.Fire)

	// set new fire time and return
	if w.world != nil {
		w.lastFire = w.world.iteration
	}
	return
}

//---------------- UPDATE --------------------------------------------------------------------------------------------//

// Update is called with each iteration and updates internal variables.
func (w *Weapon) Update(isMoving bool) {
	if isMoving && w.world != nil {
		w.lastMove = w.world.iteration // set movement
	}
}
