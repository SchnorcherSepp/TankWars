package core

// game settings
const (
	GameSpeed    = 30            // iterations per second
	MovePerTick  = 0.02          // percent of movement per tick
	WorldXWidth  = 28            // world dimension X (28 * 64 = 1792pxl)
	WorldYHeight = 15            // world dimension Y (15 * 64 = 960pxl)
	BlockSize    = 64            // image size of tanks, barriers, buildings, ...
	BlockRadius  = BlockSize / 2 // radius of blocks (tanks)
	BallSize     = 20            // size of projectiles
	BallRadius   = BallSize / 2  // radius of projectiles
)

// macros
const (
	MacroAttackMove      = "AttackMove"
	MacroFireAndManeuver = "FireAndManeuver"
	MacroFireWall        = "FireWall"
	MacroGuardMode       = "GuardMode"
	MacroReset           = "nil"
)

// weapons
const (
	ShowExplosionIterations = 10               // duration of the explosion animation
	WeaponNone              = "None"           // no weapon for neutral objects (rock)
	WeaponCannon            = "Tank"           // weapon of a battle tank
	WeaponArtillery         = "Artillery"      // weapon of an artillery
	WeaponRockets           = "RocketLauncher" // weapon of a rocket launcher
)

// tank attr
const (
	TankRotationDelay = 467 * GameSpeed / 1000 // rotation delay in iterations (~467 ms)
	TankBudget        = 100                    // max. points = armor + damage + speed
	TankMinSpeed      = 25                     // min. Speed (calc budget-armor-damage)
	TankMinArmor      = 5                      // min armor
	TankMaxArmor      = 55                     // max armor
	TankMinDamage     = 15                     // min damage
	TankMaxDamage     = 70                     // max damage
)

// tank angle (movement)
const (
	North     = 0
	Northeast = 45
	East      = 90
	Southeast = 135
	South     = 180
	Southwest = 225
	West      = 270
	Northwest = 315
)

// tank status
const (
	StatusMoving    = "Moving"    // tank is moving and can't fire
	StatusPreparing = "Preparing" // tank prepare for fire after moving (see prepTime)
	StatusReady     = "Ready"     // tank can fire
	StatusReloading = "Reloading" // tank reload weapon after fire
	StatusNoWeapon  = "NoWeapon"  // error: no weapon or not alive
)

// player tanks
const (
	RedTank  = "red"  // player 1 (red)
	BlueTank = "blue" // player 2 (blue)
)

// buildings
const (
	NeutralRock = "neutral_rock" // rock, assigned to no player
	RedRock     = "red_rock"     // rock of gamer 1 (red)
	BlueRock    = "blue_rock"    // rock of gamer 2 (blue)
	RedBase     = "red_base"     // base of gamer 1 (red)
	BlueBase    = "blue_base"    // base of gamer 2 (blue)
)
