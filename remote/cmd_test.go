package remote

import (
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"strings"
	"testing"
)

func Test_id2Tank(t *testing.T) {
	w := core.NewWorld(100, 200)
	red, _ := core.NewTank(w, core.RedTank, 5, 40, core.WeaponCannon)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)
	blue, _ := core.NewTank(w, core.BlueTank, 5, 40, core.WeaponCannon)
	blue.SetPosition(core.NewPosition(200, 100), core.East)
	w.AddTank(blue)

	if nt, err := id2Tank(w, "", red.ID()); nt != red || err != nil { // ok
		t.Error(err)
	}
	if nt, err := id2Tank(w, "", blue.ID()); nt != blue || err != nil { // ok
		t.Error(err)
	}

	if nt, err := id2Tank(w, "red", red.ID()); nt != red || err != nil { // ok
		t.Error(err)
	}
	if _, err := id2Tank(w, "red", blue.ID()); err == nil { // err
		t.Error("wrong value")
	}

	if _, err := id2Tank(w, "blue", red.ID()); err == nil { // err
		t.Error("wrong value")
	}
	if nt, err := id2Tank(w, "blue", blue.ID()); nt != blue || err != nil { // ok
		t.Error(err)
	}

	if _, err := id2Tank(w, "red_base", red.ID()); err == nil { // err
		t.Error("wrong value")
	}
	if _, err := id2Tank(w, "blue_base", blue.ID()); err == nil { // err
		t.Error("wrong value")
	}
}

func TestCloseTargets(t *testing.T) {

	w := core.NewWorld(100, 200)
	red, _ := core.NewTank(w, core.RedTank, 5, 40, core.WeaponCannon)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)
	blue, _ := core.NewTank(w, core.RedTank, 5, 40, core.WeaponCannon)
	blue.SetPosition(core.NewPosition(200, 100), core.East)
	w.AddTank(blue)

	if ct := CloseTargets(w, "id"); ct != "err: tank not found" {
		t.Error(ct)
	}
	if ct := CloseTargets(w, red.ID()); !strings.Contains(ct, `,"distance":100,"relativeAngle":90}]`) {
		t.Error(ct)
	}
}

func TestSetMacroMoveTo(t *testing.T) {
	w := core.NewWorld(1000, 1000)
	nt, _ := core.NewTank(w, core.RedTank, 5, 15, core.WeaponCannon)
	nt.SetPosition(core.NewPosition(500, 500), core.North)
	w.AddTank(nt)
	rock, _ := core.NewTank(w, core.NeutralRock, 5, 15, core.WeaponNone)
	rock.SetPosition(core.NewPosition(600, 600), core.Northwest)
	w.AddTank(rock)

	// test nil
	if s := SetMacroMoveTo(nil, "", "", "", ""); s != "err: tank not found" {
		t.Error("wrong value", s)
	}
	// wrong id
	if s := SetMacroMoveTo(w, "", "wrong", "", ""); s != "err: tank not found" {
		t.Error("wrong value", s)
	}
	// wrong X
	if s := SetMacroMoveTo(w, "", nt.ID(), "w", ""); s != "err: X: strconv.Atoi: parsing \"w\": invalid syntax" {
		t.Error("wrong value", s)
	}
	// wrong Y
	if s := SetMacroMoveTo(w, "", nt.ID(), "700", "w"); s != "err: Y: strconv.Atoi: parsing \"w\": invalid syntax" {
		t.Error("wrong value", s)
	}

	// success
	if s := SetMacroMoveTo(w, "", nt.ID(), "700", "700"); s != "ok" {
		t.Error("wrong value", s)
	}

	// check position
	if nt.Pos().X != 500 || nt.Pos().Y != 500 {
		t.Error("wrong value", nt.Pos().X, nt.Pos().Y)
	}
	w.UpdateN(30)
	if nt.Pos().X == 500 || nt.Pos().Y == 500 {
		t.Error("wrong value", nt.Pos().X, nt.Pos().Y)
	}
}

func TestPossibleTargets(t *testing.T) {
	w := core.NewWorld(100, 200)
	red, _ := core.NewTank(w, core.RedTank, 5, 40, core.WeaponCannon)
	red.SetPosition(core.NewPosition(100, 100), core.East)
	w.AddTank(red)

	if pt := PossibleTargets(w, "id"); pt != "err: tank not found" {
		t.Error(pt)
	}
	if pt := PossibleTargets(w, red.ID()); pt != "[]" {
		t.Error(pt)
	}
}

func TestGameStatus(t *testing.T) {
	w := core.NewWorld(100, 200)

	if gs := GameStatus(nil); gs != "err: invalid world status" {
		t.Error(gs)
	}
	if gs := GameStatus(w); len(gs) < 100 {
		t.Error(gs)
	}
}

func TestTankStatus(t *testing.T) {
	w := core.NewWorld(100, 200)
	red, _ := core.NewTank(w, core.RedTank, 5, 40, core.WeaponRockets)
	w.AddTank(red)

	if ts := TankStatus(w, "id"); ts != "err: tank not found" {
		t.Error(ts)
	}
	if ts := TankStatus(nil, "id"); ts != "err: tank not found" {
		t.Error(ts)
	}
	if ts := TankStatus(w, red.ID()); len(ts) < 100 {
		t.Error(ts)
	}
}

func TestFire(t *testing.T) {
	resources.MuteSound = true // mute sound for tests

	w := core.NewWorld(100, 200)
	w.UpdateN(100)

	red, _ := core.NewTank(w, core.RedTank, 5, 40, core.WeaponRockets)
	red.SetPosition(core.NewPosition(100, 100), core.North)
	w.AddTank(red)
	blue, _ := core.NewTank(w, core.BlueTank, 5, 15, core.WeaponCannon)
	blue.SetPosition(core.NewPosition(200, 100), core.West)
	w.AddTank(blue)

	// test Fire
	if txt := Fire(w, "", "id", "0", "100"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}
	if txt := Fire(w, "", red.ID(), "0!", "100"); txt != "err: angle: strconv.Atoi: parsing \"0!\": invalid syntax" {
		t.Error("wrong value", txt)
	}
	if txt := Fire(w, "", red.ID(), "0", "100!"); txt != "err: distance: strconv.Atoi: parsing \"100!\": invalid syntax" {
		t.Error("wrong value", txt)
	}
	if txt := Fire(w, "", red.ID(), "0", "100"); txt != "err: Preparing" {
		t.Error("wrong value", txt)
	}

	w.UpdateN(int(red.Weapon().PreparationTime()))

	if txt := Fire(w, "", red.ID(), "90", "100"); txt != "ok" {
		t.Error("wrong value", txt)
	}
	for i := 0; i < 100; i++ {
		w.Update()
	}
	if h := blue.Health(); h != 94 {
		t.Error("wrong value", h)
	}

	// test FireAt
	if txt := FireAt(w, "", "id", "200", "100"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}
	if txt := FireAt(w, "", red.ID(), "200!", "100"); txt != "err: X: strconv.Atoi: parsing \"200!\": invalid syntax" {
		t.Error("wrong value", txt)
	}
	if txt := FireAt(w, "", red.ID(), "200", "100!"); txt != "err: Y: strconv.Atoi: parsing \"100!\": invalid syntax" {
		t.Error("wrong value", txt)
	}
	if txt := FireAt(w, "", red.ID(), "200", "100"); txt != "ok" {
		t.Error("wrong value", txt)
	}
	if txt := FireAt(w, "", red.ID(), "200", "100"); txt != "err: Reloading" {
		t.Error("wrong value", txt)
	}

	w.UpdateN(int(red.Weapon().ReloadTime()))

	if txt := FireAt(w, "", red.ID(), "200", "100"); txt != "ok" {
		t.Error("wrong value", txt)
	}
	for i := 0; i < 100; i++ {
		w.Update()
	}
	if h := blue.Health(); h != 82 {
		t.Error("wrong value", h)
	}
}

func TestSetMacro(t *testing.T) {
	w := core.NewWorld(100, 200)
	nt, _ := core.NewTank(w, core.RedTank, 5, 15, core.WeaponRockets)
	w.AddTank(nt)

	// set
	if txt := SetMacro(w, "", nt.ID(), core.MacroGuardMode); txt != "ok" || nt.ActiveMacro() != true {
		t.Error("wrong value", txt)
	}
	nt.Update()
	nt.SetMacro(nil)
	if txt := SetMacro(w, "", nt.ID(), core.MacroFireWall); txt != "ok" || nt.ActiveMacro() != true {
		t.Error("wrong value", txt)
	}
	nt.Update()
	nt.SetMacro(nil)
	if txt := SetMacro(w, "", nt.ID(), core.MacroFireAndManeuver); txt != "ok" || nt.ActiveMacro() != true {
		t.Error("wrong value", txt)
	}
	nt.Update()
	nt.SetMacro(nil)
	if txt := SetMacro(w, "", nt.ID(), core.MacroAttackMove); txt != "ok" || nt.ActiveMacro() != true {
		t.Error("wrong value", txt)
	}
	nt.Update()

	// remove
	if txt := SetMacro(w, "", nt.ID(), "nothing"); txt != "err: macro not found" || nt.ActiveMacro() != false {
		t.Error("wrong value", txt)
	}
	nt.Update()
	if txt := SetMacro(w, "", nt.ID(), ""); txt != "ok: disable macro" || nt.ActiveMacro() != false {
		t.Error("wrong value", txt)
	}
	nt.Update()
	if txt := SetMacro(w, "", nt.ID(), core.MacroReset); txt != "ok: disable macro" || nt.ActiveMacro() != false {
		t.Error("wrong value", txt)
	}
	nt.Update()

	// wrong id
	if txt := SetMacro(w, "", "id", "nil"); txt != "err: tank not found" || nt.ActiveMacro() != false {
		t.Error("wrong value", txt)
	}
	nt.Update()
}

func TestMyName(t *testing.T) {
	if MyName(core.RedTank) != "red" {
		t.Error("wrong value")
	}
}

func TestBuyTank(t *testing.T) {
	w := core.NewWorld(100, 200)
	hb, _ := core.NewTank(w, core.RedBase, 5, 15, core.WeaponNone)
	w.AddTank(hb)
	w.SetCash(100, 0)

	// invalid input
	if txt := BuyTank(w, core.RedTank, "10!", "20", core.WeaponCannon); txt != "err: armor: strconv.Atoi: parsing \"10!\": invalid syntax" {
		t.Error("wrong value", txt)
	}
	if txt := BuyTank(w, core.RedTank, "10", "20!", core.WeaponCannon); txt != "err: damage: strconv.Atoi: parsing \"20!\": invalid syntax" {
		t.Error("wrong value", txt)
	}
	if txt := BuyTank(w, core.RedTank, "10", "20", "no"); txt != "err: unknown weapon: no" {
		t.Error("wrong value", "'", txt, "'")
	}

	// no cash
	if txt := BuyTank(w, core.RedTank, "10", "20", core.WeaponCannon); txt[:2] != "ok" {
		t.Error("wrong value", txt)
	}
	if txt := BuyTank(w, core.RedTank, "10", "20", core.WeaponCannon); txt != "err: not enough tank budget or unknown owner" {
		t.Error("wrong value", txt)
	}

	// no base
	w.SetCash(100, 0)
	hb.Remove()
	if txt := BuyTank(w, core.RedTank, "10", "20", core.WeaponCannon); txt != "err: home base not found" {
		t.Error("wrong value", txt)
	}
}

func TestMovement(t *testing.T) {
	w := core.NewWorld(100, 200)
	nt, _ := core.NewTank(w, core.RedTank, 5, 15, core.WeaponCannon)
	nt.SetPosition(core.NewPosition(100, 200), core.Northwest)
	w.AddTank(nt)

	// CHECK NIL
	if txt := Left(w, "", "nil"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}
	if txt := Right(w, "", "nil"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}
	if txt := Forward(w, "", "nil"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}
	if txt := Backward(w, "", "nil"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}
	if txt := Stop(w, "", "nil"); txt != "err: tank not found" {
		t.Error("wrong value", txt)
	}

	// init check
	if nt.Angle() != core.Northwest || nt.Pos().X != 100 || nt.Pos().Y != 200 {
		t.Error("wrong value", nt.Angle(), nt.Pos().X, nt.Pos().Y)
	}

	// check left & right
	w.UpdateN(core.TankRotationDelay)
	if txt := Left(w, "", nt.ID()); nt.Angle() != core.West || txt != "ok" {
		t.Error("wrong value", nt.Angle(), txt)
	}
	if txt := Left(w, "", nt.ID()); txt != "err: Preparing" { // timer
		t.Error("wrong value", txt)
	}
	w.UpdateN(core.TankRotationDelay)
	if txt := Right(w, "", nt.ID()); nt.Angle() != core.Northwest || txt != "ok" {
		t.Error("wrong value", nt.Angle(), txt)
	}
	if txt := Right(w, "", nt.ID()); txt != "err: Preparing" { // timer
		t.Error("wrong value", txt)
	}

	// check forward
	if txt := Forward(w, "", nt.ID()); txt != "ok" {
		t.Error("wrong value", nt.Angle(), txt)
	}
	w.Update()
	w.Update()
	if nt.Moving() != true || nt.Pos().X != 97 || nt.Pos().Y != 197 {
		t.Error("wrong value", nt.Moving(), nt.Pos().X, nt.Pos().Y)
	}

	// check stop
	if txt := Stop(w, "", nt.ID()); txt != "ok" {
		t.Error("wrong value", nt.Angle(), txt)
	}
	w.Update()
	w.Update()
	if nt.Moving() != false {
		t.Error("wrong value", nt.Moving())
	}

	// check back
	if txt := Backward(w, "", nt.ID()); txt != "ok" {
		t.Error("wrong value", nt.Angle(), txt)
	}
	w.Update()
	w.Update()
	if nt.Moving() != true || nt.Pos().X != 100 || nt.Pos().Y != 200 {
		t.Error("wrong value", nt.Moving(), nt.Pos().X, nt.Pos().Y)
	}
}
