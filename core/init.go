package core

// TestInitialization allows setting non-exported variables outside the core packet.
func (p *Projectile) TestInitialization(world *World, parent *Tank, pos, startPos, endPos Position, angle, distance, speed, damage, aoeRadius int, collision bool, exploded uint) {
	p.world = world
	p.parent = parent
	p.pos = pos
	p.startPos = startPos
	p.endPos = endPos
	p.angle = angle
	p.distance = distance
	p.speed = speed
	p.damage = damage
	p.aoeRadius = aoeRadius
	p.collision = collision
	p.exploded = exploded
}

// TestInitialization allows setting non-exported variables outside the core packet.
func (t *Tank) TestInitialization(world *World, id, owner string, weapon *Weapon, health, armor, speed int, pos Position, command, angle int, isBlocked bool, lastRotate uint64, macro func(t *Tank)) {
	t.world = world
	t.id = id
	t.owner = owner
	t.weapon = weapon
	t.health = health
	t.armor = armor
	t.speed = speed
	t.pos = pos
	t.command = command
	t.angle = angle
	t.isBlocked = isBlocked
	t.lastRotate = lastRotate
	t.macro = macro
}

// TestInitialization allows setting non-exported variables outside the core packet.
func (w *Weapon) TestInitialization(world *World, parent *Tank, typ string, rng int, prepTime, reloadTime uint64, projSpeed, damage, aoeRadius int, projCollision, anyFireAngle bool, lastMove, lastFire uint64) {
	w.world = world
	w.parent = parent
	w.typ = typ
	w.rng = rng
	w.prepTime = prepTime
	w.reloadTime = reloadTime
	w.projSpeed = projSpeed
	w.damage = damage
	w.aoeRadius = aoeRadius
	w.projCollision = projCollision
	w.anyFireAngle = anyFireAngle
	w.lastMove = lastMove
	w.lastFire = lastFire
}

// TestInitialization allows setting non-exported variables outside the core packet.
func (w *World) TestInitialization(xWidth, yHeight int, iteration uint64, tanks []*Tank, projectiles []*Projectile, freeze bool, cashRed, cashBlue float64) {
	w.xWidth = xWidth
	w.yHeight = yHeight
	w.iteration = iteration
	w.tanks = tanks
	w.projectiles = projectiles
	w.freeze = freeze
	w.cashRed = cashRed
	w.cashBlue = cashBlue
}
