package remote

import (
	"encoding/json"
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
)

//---------------- [1] Position --------------------------------------------------------------------------------------//

// JsonPosition is the protocol struct of core.Position
type JsonPosition struct {
	X  int     `json:"x"`
	Xf float64 `json:"xf"`
	Y  int     `json:"y"`
	Yf float64 `json:"yf"`
}

// NewJsonPosition convert a core object to a json object
func NewJsonPosition(p core.Position) JsonPosition {
	return JsonPosition{
		X:  p.X,
		Xf: p.Xf,
		Y:  p.Y,
		Yf: p.Yf,
	}
}

// Get returns a json representation of this object
func (p *JsonPosition) Get() string {
	b, err := json.Marshal(p)
	if err != nil || p == nil {
		fmt.Printf("err: JsonPosition: %v\n", err)
	}
	return string(b)
}

// Set parse a json string and update the inner variables of this object
func (p *JsonPosition) Set(j string) {
	if err := json.Unmarshal([]byte(j), &p); err != nil {
		fmt.Printf("err: JsonPosition: %v\n", err)
	}
}

//---------------- [2] Target ----------------------------------------------------------------------------------------//

// JsonTarget is the protocol struct of core.Target
type JsonTarget struct {
	TankID        string `json:"tankID"`
	Distance      int    `json:"distance"`
	RelativeAngle int    `json:"relativeAngle"`
}

// NewJsonTarget convert a core object to a json object
func NewJsonTarget(t core.Target) JsonTarget {
	id := ""
	if t.Tank != nil {
		id = t.Tank.ID()
	}

	return JsonTarget{
		TankID:        id,
		Distance:      t.Distance,
		RelativeAngle: t.RelativeAngle,
	}
}

//---------------- [3] Targets (LIST) --------------------------------------------------------------------------------//

// JsonTargets is the protocol list of core.Target
type JsonTargets []JsonTarget

// NewJsonTargets convert a core object list to a json object
func NewJsonTargets(ts []core.Target) JsonTargets {
	ret := make(JsonTargets, 0)
	if ts == nil {
		return ret
	}

	for _, t := range ts {
		ret = append(ret, NewJsonTarget(t))
	}

	return ret
}

// Get returns a json representation of this object
func (ts *JsonTargets) Get() string {
	b, err := json.Marshal(ts)
	if err != nil || ts == nil {
		fmt.Printf("err: JsonTargets: %v\n", err)
	}
	return string(b)
}

// Set parse a json string and update the inner variables of this object
func (ts *JsonTargets) Set(j string) {
	if err := json.Unmarshal([]byte(j), &ts); err != nil {
		fmt.Printf("err: JsonTargets: %v\n", err)
	}
}

//---------------- [4] Projectile ------------------------------------------------------------------------------------//

// JsonProjectile is the protocol struct of core.Projectile
type JsonProjectile struct {
	Parent    string       `json:"parent"`
	Pos       JsonPosition `json:"pos"`
	StartPos  JsonPosition `json:"startPos"`
	Angle     int          `json:"angle"`
	Distance  int          `json:"distance"`
	Speed     int          `json:"speed"`
	Damage    int          `json:"damage"`
	AoeRadius int          `json:"aoeRadius"`
	Collision bool         `json:"collision"`
	Exploded  bool         `json:"exploded"`
}

// NewJsonProjectile convert a core object to a json object
func NewJsonProjectile(p *core.Projectile) JsonProjectile {
	if p == nil {
		return JsonProjectile{}
	}
	jp := JsonProjectile{
		Parent:    "", // set later
		Pos:       NewJsonPosition(p.Pos()),
		StartPos:  NewJsonPosition(p.StartPos()),
		Angle:     p.Angle(),
		Distance:  p.Distance(),
		Speed:     p.Speed(),
		Damage:    p.Damage(),
		AoeRadius: p.AoERadius(),
		Collision: p.Collision(),
		Exploded:  p.Exploded(),
	}
	if p.Parent() != nil {
		jp.Parent = p.Parent().ID()
	}
	return jp
}

// Get returns a json representation of this object
func (p *JsonProjectile) Get() string {
	b, err := json.Marshal(p)
	if err != nil || p == nil {
		fmt.Printf("err: JsonProjectile: %v\n", err)
	}
	return string(b)
}

// Set parse a json string and update the inner variables of this object
func (p *JsonProjectile) Set(j string) {
	if err := json.Unmarshal([]byte(j), &p); err != nil {
		fmt.Printf("err: JsonProjectile: %v\n", err)
	}
}

//---------------- [5] Tank ------------------------------------------------------------------------------------------//

// JsonTank is the protocol struct of core.Tank
type JsonTank struct {
	ID          string       `json:"id"`
	Owner       string       `json:"owner"`
	Health      int          `json:"health"`
	Armor       int          `json:"armor"`
	Speed       int          `json:"speed"`
	Pos         JsonPosition `json:"pos"`
	Command     int          `json:"command"`
	Angle       int          `json:"angle"`
	IsBlocked   bool         `json:"isBlocked"`
	ActiveMacro bool         `json:"activeMacro"`
	Alive       bool         `json:"alive"`
	Moving      bool         `json:"moving"`
	LastRotate  uint64       `json:"lastRotate"`
	Rdy         bool         `json:"rdy"`
	Status      string       `json:"status"`
	Weapon      JsonWeapon   `json:"weapon"`
}

// NewJsonTank convert a core object to a json object
func NewJsonTank(t *core.Tank) JsonTank {
	if t == nil {
		return JsonTank{}
	}

	rdy, stat := t.Status()

	return JsonTank{
		ID:          t.ID(),
		Owner:       t.Owner(),
		Health:      t.Health(),
		Armor:       t.Armor(),
		Speed:       t.Speed(),
		Pos:         NewJsonPosition(t.Pos()),
		Command:     t.Command(),
		Angle:       t.Angle(),
		IsBlocked:   t.Blocked(),
		ActiveMacro: t.ActiveMacro(),
		Alive:       t.Alive(),
		Moving:      t.Moving(),
		LastRotate:  t.LastRotate(),
		Rdy:         rdy,
		Status:      stat,
		Weapon:      NewJsonWeapon(t.Weapon()),
	}
}

// Get returns a json representation of this object
func (t *JsonTank) Get() string {
	b, err := json.Marshal(t)
	if err != nil || t == nil {
		fmt.Printf("err: JsonTank: %v\n", err)
	}
	return string(b)
}

// Set parse a json string and update the inner variables of this object
func (t *JsonTank) Set(j string) {
	if err := json.Unmarshal([]byte(j), &t); err != nil {
		fmt.Printf("err: JsonTank: %v\n", err)
	}
}

//---------------- [6] Weapon ----------------------------------------------------------------------------------------//

// JsonWeapon is the protocol struct of core.Weapon
type JsonWeapon struct {
	Typ           string `json:"typ"`
	Rng           int    `json:"rng"`
	PrepTime      uint64 `json:"prepTime"`
	ReloadTime    uint64 `json:"reloadTime"`
	ProjSpeed     int    `json:"projSpeed"`
	Damage        int    `json:"damage"`
	AoeRadius     int    `json:"aoeRadius"`
	ProjCollision bool   `json:"projCollision"`
	AnyFireAngle  bool   `json:"anyFireAngle"`
	Rdy           bool   `json:"rdy"`
	Status        string `json:"status"`
	LastMove      uint64 `json:"lastMove"`
	LastFire      uint64 `json:"lastFire"`
}

// NewJsonWeapon convert a core object to a json object
func NewJsonWeapon(w *core.Weapon) JsonWeapon {
	if w == nil {
		return JsonWeapon{}
	}

	rdy, stat := w.Status()

	return JsonWeapon{
		Typ:           w.Type(),
		Rng:           w.Range(),
		PrepTime:      w.PreparationTime(),
		ReloadTime:    w.ReloadTime(),
		ProjSpeed:     w.ProjectileSpeed(),
		Damage:        w.Damage(),
		AoeRadius:     w.AoERadius(),
		ProjCollision: w.ProjectileCollision(),
		AnyFireAngle:  w.AnyFireAngle(),
		Rdy:           rdy,
		Status:        stat,
		LastMove:      w.LastMove(),
		LastFire:      w.LastFire(),
	}
}

// Get returns a json representation of this object
func (w *JsonWeapon) Get() string {
	b, err := json.Marshal(w)
	if err != nil || w == nil {
		fmt.Printf("err: JsonWeapon: %v\n", err)
	}
	return string(b)
}

// Set parse a json string and update the inner variables of this object
func (w *JsonWeapon) Set(j string) {
	if err := json.Unmarshal([]byte(j), &w); err != nil {
		fmt.Printf("err: JsonWeapon: %v\n", err)
	}
}

//---------------- [7] World -----------------------------------------------------------------------------------------//

// JsonWorld is the protocol struct of core.World
type JsonWorld struct {
	GameSpeed     int              `json:"gameSpeed"`
	MovePerTick   float64          `json:"movePerTick"`
	TankRadius    int              `json:"tankRadius"`
	BallRadius    int              `json:"ballRadius"`
	RotationDelay int              `json:"rotationDelay"`
	TankBudget    int              `json:"tankBudget"`
	MinSpeed      int              `json:"minSpeed"`
	MinArmor      int              `json:"minArmor"`
	MaxArmor      int              `json:"maxArmor"`
	MinDamage     int              `json:"minDamage"`
	MaxDamage     int              `json:"maxDamage"`
	XWidth        int              `json:"xWidth"`
	YHeight       int              `json:"yHeight"`
	ScreenWidth   int              `json:"screenWidth"`
	ScreenHeight  int              `json:"screenHeight"`
	Iteration     uint64           `json:"iteration"`
	Freeze        bool             `json:"freeze"`
	CashRed       int              `json:"cashRed"`
	CashBlue      int              `json:"cashBlue"`
	Tanks         []JsonTank       `json:"tanks"`
	Projectiles   []JsonProjectile `json:"projectiles"`
	UnitCountRed  int              `json:"unitCountRed"`
	UnitCountBlue int              `json:"unitCountBlue"`
}

// NewJsonWorld convert a core object to a json object
func NewJsonWorld(w *core.World) JsonWorld {
	if w == nil {
		return JsonWorld{}
	}

	cRed, cBlue := w.CashStat()
	uRed, uBlue := w.UnitCount()
	tanks := make([]JsonTank, 0, 1024)
	for _, t := range w.Tanks() {
		tanks = append(tanks, NewJsonTank(t))
	}
	proj := make([]JsonProjectile, 0, 1024)
	for _, p := range w.Projectiles() {
		proj = append(proj, NewJsonProjectile(p))
	}

	return JsonWorld{
		GameSpeed:     core.GameSpeed,
		MovePerTick:   core.MovePerTick,
		TankRadius:    core.BlockRadius,
		BallRadius:    core.BallRadius,
		RotationDelay: core.TankRotationDelay,
		TankBudget:    core.TankBudget,
		MinSpeed:      core.TankMinSpeed,
		MinArmor:      core.TankMinArmor,
		MaxArmor:      core.TankMaxArmor,
		MinDamage:     core.TankMinDamage,
		MaxDamage:     core.TankMaxDamage,
		XWidth:        w.XWidth(),
		YHeight:       w.YHeight(),
		ScreenWidth:   w.ScreenWidth(),
		ScreenHeight:  w.ScreenHeight(),
		Iteration:     w.Iteration(),
		Freeze:        w.IsFrozen(),
		CashRed:       cRed,
		CashBlue:      cBlue,
		Tanks:         tanks,
		Projectiles:   proj,
		UnitCountRed:  uRed,
		UnitCountBlue: uBlue,
	}
}

// Get returns a json representation of this object
func (w *JsonWorld) Get() string {
	b, err := json.Marshal(w)
	if err != nil || w == nil {
		fmt.Printf("err: JsonWorld: %v\n", err)
	}
	return string(b)
}

// Set parse a json string and update the inner variables of this object
func (w *JsonWorld) Set(j string) {
	if err := json.Unmarshal([]byte(j), &w); err != nil {
		fmt.Printf("err: JsonWorld: %v\n", err)
	}
}
