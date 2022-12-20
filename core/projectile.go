package core

import (
	"github.com/SchnorcherSepp/TankWars/gui/resources"
)

// Projectile is created by a weapon.
// It moves in the world. It can collide with other objects and can explode.
type Projectile struct {
	world  *World
	parent *Tank // don't explode on your parent tank

	pos       Position // current projectile position
	startPos  Position // tank (start) position
	endPos    Position // target (end) position
	angle     int      // 0 is north; 45 is northeast; 90 is east; ...
	distance  int      // maximum flight distance
	speed     int      // object speed
	damage    int      // explosion damage
	aoeRadius int      // aoe radius
	collision bool     // true - only explodes on contact; false - only explodes at the end of the distance

	exploded uint // if >0 this projectile is exploded already; count the updates since then
}

// NewProjectile create a new projectile.
// When created in a world, it interacts with other objects.
// Forward the parents, otherwise the projectile will explode at the start.
//
// For parameter description see Projectile.Pos(), Projectile.Angle(), Projectile.Distance(), Projectile.Speed(),
// Projectile.Damage(), Projectile.AoERadius() and Projectile.Collision().
func NewProjectile(world *World, parent *Tank, pos Position, angle, distance, speed, damage, aoeRadius int, collision bool) *Projectile {
	return &Projectile{
		world:     world,
		parent:    parent,
		pos:       pos,
		startPos:  pos,
		endPos:    CalcPosFromAngle(pos, angle, distance),
		angle:     angle,
		distance:  distance,
		speed:     speed,
		damage:    damage,
		aoeRadius: aoeRadius,
		collision: collision,
	}
}

//---------------- GETTER --------------------------------------------------------------------------------------------//

// Parent returns the tank from which the projectile was fired.
func (p *Projectile) Parent() *Tank {
	return p.parent
}

// Pos is the current position of this projectile.
// Is changed by Update().
func (p *Projectile) Pos() Position {
	return p.pos
}

// StartPos is the initial position of this projectile.
func (p *Projectile) StartPos() Position {
	return p.startPos
}

// EndPos returns the planned target position,
// but a collisions can happen earlier.
func (p *Projectile) EndPos() Position {
	return p.endPos
}

// Angle is the start angle of the projectile.
// see North, South, East, ...
func (p *Projectile) Angle() int {
	return p.angle
}

// Distance between start position and current position.
func (p *Projectile) Distance() int {
	return p.distance
}

// Speed at which the bullet moves per tick.
func (p *Projectile) Speed() int {
	return p.speed
}

// Damage that is included in the calculation on a hit.
// see Tank.Hit()
func (p *Projectile) Damage() int {
	return p.damage
}

// AoERadius determines how many objects are hit in an explosion.
// see Explode()
func (p *Projectile) AoERadius() int {
	return p.aoeRadius
}

// Collision defines two types of projectiles:
//
//	true: the projectile flies until it hits a target or disappears after the maximum distance (WeaponCannon).
//	false: the projectile flies without any interaction and explode after the maximum distance (WeaponArtillery).
func (p *Projectile) Collision() bool {
	return p.collision
}

// Exploded return true if the projectile is exploded.
func (p *Projectile) Exploded() bool {
	return p.exploded > 0
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// Explode destroys the projectile at the current position.
// All objects in range (see AoERadius) are hit (see Tank.Hit).
func (p *Projectile) Explode() {
	if p.world != nil {

		// aoe radius
		aoeRadius := p.aoeRadius
		if aoeRadius < BallRadius { // min aoe to hit the target (ball radius)
			aoeRadius = BallRadius

		}

		// hit all around
		for _, tank := range p.world.tanks {
			if tank != nil && IsCollided(p.pos, aoeRadius, tank.pos, BlockRadius) {
				tank.Hit(p.damage)
			}
		}

		// projectile removed by Update()
		p.exploded = 1

		// play sound
		resources.PlaySound(resources.Sounds.Explosion)
	}
}

// Remove this projectile from the world list World.Projectiles.
func (p *Projectile) Remove() {
	if p.world != nil {
		newProjectiles := make([]*Projectile, 0, len(p.world.projectiles))
		for _, proj := range p.world.projectiles {
			if proj != p {
				newProjectiles = append(newProjectiles, proj)
			}
		}
		p.world.projectiles = newProjectiles
	}
}

//---------------- UPDATE --------------------------------------------------------------------------------------------//

// Update move the projectile, remove exploded projectile, calculate collisions and enforce the max distance.
func (p *Projectile) Update() {

	// exploded already -> do nothing
	// see ShowExplosionIterations
	if p.Exploded() {
		p.exploded++
		if p.exploded > ShowExplosionIterations {
			p.Remove() // remove after n updates
		}
		return // EXIT
	}

	// move
	p.pos.Update(p.angle, p.speed)

	// check collisions
	// only if projectile can collide (see collision)
	if p.world != nil && p.collision {
		for _, tank := range p.world.tanks {
			if tank != nil && tank != p.parent && IsCollided(p.pos, BallRadius, tank.pos, BlockRadius) {
				p.Explode() // bum & remove
				return      // EXIT
			}
		}
	}

	// max distance
	if float64(p.distance) < Distance(p.startPos, p.pos) {
		if p.collision {
			// Cannon: explode on hit only
			p.Remove() // no hit; remove
		} else {
			// Artillery: projectile without collision explode at the end
			p.Explode() // bum & remove
		}
	}
}
