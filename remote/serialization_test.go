package remote

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"reflect"
	"regexp"
	"testing"
)

//---------------- Position ------------------------------------------------------------------------------------------//

func TestJsonPosition_Changes(t *testing.T) {
	// detect struct changes
	o := core.NewPosition(67, 89) // NewPosition
	cs := "core.Position{X:67, Xf:67, Y:89, Yf:89}"

	if s := fmt.Sprintf("%#v", o); s != cs {
		println(cs)
		println(s)
		t.Fatal(s)
	}
}

func TestJsonPosition(t *testing.T) {
	org := core.NewPosition(33, 44) // origin
	obj := NewJsonPosition(org)     // JSON Object
	str := obj.Get()                // json string
	newO := JsonPosition{}          // NEW JSON Object
	newO.Set(str)                   // parse

	// check
	if obj.X != newO.X || obj.Xf != newO.Xf || obj.Y != newO.Y || obj.Yf != newO.Yf {
		t.Error("wrong value")
	}
	// test invalid input
	newO.Set("")
	var nilObj *JsonPosition
	nilObj.Get() // test nil
}

//---------------- Target --------------------------------------------------------------------------------------------//

func TestJsonTarget_Changes(t *testing.T) {
	// detect struct changes
	o := &core.Target{Tank: nil, Distance: 3, RelativeAngle: 5} // Target
	cs := "&core.Target{Tank:(*core.Tank)(nil), Distance:3, RelativeAngle:5}"

	if s := fmt.Sprintf("%#v", o); s != cs {
		println(cs)
		println(s)
		t.Fatal(s)
	}
}

func TestJsonTargets(t *testing.T) {
	nt, _ := core.NewTank(nil, "", 11, 22, core.WeaponCannon)

	org := []core.Target{{Tank: nt, Distance: 99, RelativeAngle: 88}}
	obj := NewJsonTargets(org) // JSON Object
	str := obj.Get()           // json string
	newO := JsonTargets{}      // NEW JSON Object
	newO.Set(str)              // parse

	// check
	if obj[0].Distance != newO[0].Distance || obj[0].RelativeAngle != newO[0].RelativeAngle {
		t.Error("wrong value")
	}
	// test invalid input
	newO.Set("")
	NewJsonTargets(nil)
	NewJsonTargets([]core.Target{{Tank: nil}, {Tank: nil}}) // Tank is nil
	var nilObj *JsonTargets
	nilObj.Get() // test nil
}

//---------------- Projectile ----------------------------------------------------------------------------------------//

func TestJsonProjectile_Changes(t *testing.T) {
	// detect struct changes
	pos := core.NewPosition(9, 8)
	o := core.NewProjectile(nil, nil, pos, 11, 22, 33, 44, 55, false) // NewProjectile
	cs := "&core.Projectile{world:(*core.World)(nil), parent:(*core.Tank)(nil), pos:core.Position{X:9, Xf:9, Y:8, Yf:8}, startPos:core.Position{X:9, Xf:9, Y:8, Yf:8}, endPos:core.Position{X:13, Xf:13.197797898283973, Y:-14, Yf:-13.59579803584861}, angle:11, distance:22, speed:33, damage:44, aoeRadius:55, collision:false, exploded:0x0}"

	if s := fmt.Sprintf("%#v", o); s != cs {
		println(cs)
		println(s)
		t.Fatal(s)
	}
}

func TestJsonProjectile(t *testing.T) {
	org := core.NewProjectile(nil, nil, core.NewPosition(1, 2), 3, 4, 5, 6, 7, true) // origin
	obj := NewJsonProjectile(org)                                                    // JSON Object
	str := obj.Get()                                                                 // json string
	newO := JsonProjectile{}                                                         // NEW JSON Object
	newO.Set(str)                                                                    // parse

	// check
	if newO.Pos.X != 1 || newO.Pos.Y != 2 || newO.Angle != 3 || newO.Distance != 4 || newO.Speed != 5 || newO.Damage != 6 || newO.AoeRadius != 7 || newO.Collision != true {
		t.Error("wrong value")
	}
	// test invalid input
	newO.Set("")
	NewJsonProjectile(nil)
	var nilObj *JsonProjectile
	nilObj.Get() // test nil
}

//---------------- Tank ----------------------------------------------------------------------------------------------//

func TestJsonTank_Changes(t *testing.T) {
	// detect struct changes
	o, _ := core.NewTank(nil, core.RedTank, 11, 22, core.WeaponCannon) // NewTank
	cs := "&core.Tank{world:(*core.World)(nil), id:\"9999\", owner:\"red\", weapon:(*core.Weapon)(0x1010101010), health:100, armor:11, speed:70, pos:core.Position{X:0, Xf:0, Y:0, Yf:0}, command:0, angle:180, isBlocked:false, lastRotate:0x0, macro:(func(*core.Tank))(nil)}"

	s := fmt.Sprintf("%#v", o)
	s = fixJsonStrings(s)
	if s != cs {
		println(cs)
		println(s)
		t.Fatal(s)
	}
}

func TestJsonTank(t *testing.T) {
	org, _ := core.NewTank(nil, "owner", 11, 22, core.WeaponArtillery) // origin
	obj := NewJsonTank(org)                                            // JSON Object
	str := obj.Get()                                                   // json string
	newO := JsonTank{}                                                 // NEW JSON Object
	newO.Set(str)                                                      // parse

	// check
	if newO.Armor != 11 || newO.Owner != "owner" {
		t.Error("wrong value")
	}
	// test invalid input
	newO.Set("")
	NewJsonTank(nil)
	var nilObj *JsonTank
	nilObj.Get() // test nil
}

//---------------- Weapon --------------------------------------------------------------------------------------------//

func TestJsonWeapon_Changes(t *testing.T) {
	// detect struct changes
	o := core.NewWeaponRocketLauncher(nil, nil, 11) // NewWeaponRocketLauncher
	cs := "&core.Weapon{world:(*core.World)(nil), parent:(*core.Tank)(nil), typ:\"RocketLauncher\", rng:387, prepTime:0xf0, reloadTime:0x3f, projSpeed:300, damage:3, aoeRadius:32, projCollision:false, anyFireAngle:true, lastMove:0x0, lastFire:0x0}"

	s := fmt.Sprintf("%#v", o)
	s = fixJsonStrings(s)
	if s != cs {
		println(cs)
		println(s)
		t.Fatal(s)
	}
}

func TestJsonWeapon(t *testing.T) {
	org := core.NewWeaponCannon(nil, nil, 10) // origin
	obj := NewJsonWeapon(org)                 // JSON Object
	str := obj.Get()                          // json string
	newO := JsonWeapon{}                      // NEW JSON Object
	newO.Set(str)                             // parse

	// check
	if newO.Typ != core.WeaponCannon || newO.Damage != 15 {
		t.Error("wrong value")
	}
	// test invalid input
	newO.Set("")
	NewJsonWeapon(nil)
	var nilObj *JsonWeapon
	nilObj.Get() // test nil
}

//---------------- World ---------------------------------------------------------------------------------------------//

func TestJsonWorld_Changes(t *testing.T) {
	// detect struct changes
	o := core.NewWorld(33, 44) // NewWorld
	cs := "&core.World{xWidth:33, yHeight:44, iteration:0x0, tanks:[]*core.Tank{}, projectiles:[]*core.Projectile{}, freeze:false, cashRed:0, cashBlue:0}"

	if s := fmt.Sprintf("%#v", o); s != cs {
		println(cs)
		println(s)
		t.Fatal(s)
	}
}

func TestJsonWorld(t *testing.T) {
	w := core.NewWorld(33, 44) // origin
	w.UpdateN(100)

	nt, _ := core.NewTank(w, "red", 11, 22, core.WeaponCannon)
	w.AddTank(nt)

	w.UpdateN(int(nt.Weapon().PreparationTime()))

	nt.Fire(0, 999999)
	w.Update()

	obj := NewJsonWorld(w) // JSON Object
	str := obj.Get()       // json string
	newO := JsonWorld{}    // NEW JSON Object
	newO.Set(str)          // parse

	// check
	if newO.YHeight != 44 || newO.XWidth != 33 || len(newO.Tanks) != 1 || len(newO.Projectiles) != 1 {
		t.Error("wrong value")
	}
	// test invalid input
	newO.Set("")
	NewJsonWorld(nil)
	var nilObj *JsonWorld
	nilObj.Get() // test nil
}

//---------------- World (reverse) -----------------------------------------------------------------------------------//

func TestJsonWorld_CoreWorld(t *testing.T) {
	resources.MuteSound = true

	w := core.NewWorld(111, 222)

	nt1, _ := core.NewTank(w, "owner 1", 22, 33, core.WeaponCannon)
	nt1.SetPosition(core.NewPosition(234, 567), core.Northwest)
	nt1.SetMacro(func(t *core.Tank) {})
	w.AddTank(nt1)

	nt2, _ := core.NewTank(w, "owner 1", 33, 22, core.WeaponArtillery)
	nt2.SetPosition(core.NewPosition(567, 234), core.Southwest)
	w.AddTank(nt2)

	w.UpdateN(300)
	nt1.Fire(0, 100)
	nt2.Fire(0, 100)
	w.UpdateN(17)

	// convert
	jw1 := NewJsonWorld(w)
	jw2 := new(JsonWorld)
	jw2.Set(jw1.Get())
	w2 := jw2.CoreWorld()

	// check macro
	var activeMacro bool
	for _, nt := range w2.Tanks() {
		if nt.ActiveMacro() {
			activeMacro = true
			nt.Update() // call macro
			nt.SetMacro(nil)
		}
	}
	if !activeMacro {
		t.Error("macro not set")
	}
	nt1.SetMacro(nil) // disable macro for next compare

	// compare
	if !reflect.DeepEqual(w, w2) || len(w.Projectiles()) < 2 || len(w.Tanks()) < 2 {
		t.Error("wrong value")
	}

	// simulate ID NOT FOUND
	nt2.Remove()

	// convert
	jw1 = NewJsonWorld(w)
	jw2 = new(JsonWorld)
	jw2.Set(jw1.Get())
	w2 = jw2.CoreWorld()

	if reflect.DeepEqual(w, w2) {
		t.Error("wrong value")
	}
}

//---------------- HELPER --------------------------------------------------------------------------------------------//

// fixTimeStrings set all dates, times, ... to 0s
func fixJsonStrings(s string) string {

	// fix  `weapon:(*core.Weapon)(0xc000632000)`
	reg := regexp.MustCompile(`weapon:\(\*core\.Weapon\)\(0x.+?\)`)
	s = reg.ReplaceAllString(s, "weapon:(*core.Weapon)(0x1010101010)")

	// fix  `id:"1235", owner:`
	reg = regexp.MustCompile(`id:".+?", owner:`)
	s = reg.ReplaceAllString(s, `id:"9999", owner:`)

	return s
}
