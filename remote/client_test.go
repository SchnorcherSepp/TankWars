package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/gui/resources"
	"strings"
	"testing"
	"time"
)

func Test_Server_Client(t *testing.T) {
	resources.MuteSound = true

	// init world
	w := core.NewWorld(core.WorldXWidth, core.WorldYHeight)
	w.UpdateN(100)

	nt, _ := core.NewTank(w, core.RedTank, 22, 33, core.WeaponCannon)
	nt.SetPosition(core.NewPosition(111, 222), core.North)

	w.UpdateN(int(nt.Weapon().PreparationTime()))

	nt.Fire(nt.Angle(), 9999)
	w.AddTank(nt)
	nt.SetMacro(func(t *core.Tank) {})

	nt, _ = core.NewTank(w, core.BlueTank, 20, 30, core.WeaponCannon)
	nt.SetPosition(core.NewPosition(50, 444), core.West)
	nt.Fire(nt.Angle(), 9999)
	w.AddTank(nt)

	nt2, _ := core.NewTank(w, core.BlueTank, 20, 30, core.WeaponCannon)
	nt2.SetPosition(core.NewPosition(190, 222), core.West)
	w.AddTank(nt2)

	nt2, _ = core.NewTank(w, core.BlueTank, 20, 30, core.WeaponCannon)
	nt2.SetPosition(core.NewPosition(50, 222), core.West)
	w.AddTank(nt2)

	w.Update()
	w.Update()
	w.Update()

	//---------------------

	// start server and init client
	go RunServer("localhost", "3333", w)
	time.Sleep(400 * time.Millisecond)          // wait for server
	NewTcpClient("localhost", "3333")           // player red
	client := NewTcpClient("localhost", "3333") // player blue
	NewTcpClient("localhost", "3333")           // observer

	//---------------------

	respGS := client.GameStatus()
	if !strings.HasPrefix(respGS, "{\"gameSpeed\":") {
		t.Error(respGS)
	}
	respTS := client.TankStatus("1236")
	if !strings.HasPrefix(respTS, "{\"id\":\"") {
		t.Error(respTS)
	}
	respCT := client.CloseTargets("1236", "", "", "", "", "")
	if !strings.HasPrefix(respCT, "[{\"tankID\":") {
		t.Error(respCT)
	}
	respPT := client.PossibleTargets("1236", "", "", "", "", "")
	if !strings.HasPrefix(respPT, "[{\"tankID\":\"1238\",\"distance\":222,\"relativeAngle\":0}]") {
		t.Error(respPT)
	}
	if resp := client.MyName(); resp != core.BlueTank {
		t.Error(resp)
	}
	if resp := client.Fire("1236", 30, 12); resp != "err: Preparing" {
		t.Error(resp)
	}
	if resp := client.FireAt("1236", 300, 200); resp != "err: Preparing" {
		t.Error(resp)
	}
	if resp := client.BuyTank(5, 70, core.WeaponCannon); resp != "err: not enough tank budget or unknown owner" {
		t.Error(resp)
	}
	if resp := client.BuyTank(5, 70, core.WeaponCannon); resp != "err: not enough tank budget or unknown owner" {
		t.Error(resp)
	}
	if resp := client.Forward("1236"); resp != "ok" {
		t.Error(resp)
	}
	if resp := client.Backward("1236"); resp != "ok" {
		t.Error(resp)
	}
	w.UpdateN(core.TankRotationDelay)
	if resp := client.Left("1236"); resp != "ok" {
		t.Error(resp)
	}
	if resp := client.Right("1236"); resp != "err: Preparing" {
		t.Error(resp)
	}
	if resp := client.Stop("1236"); resp != "ok" {
		t.Error(resp)
	}
	if resp := client.SetMacroMoveTo("1236", 11, 22); resp != "ok" {
		t.Error(resp)
	}
	if resp := client.SetMacro("1236", core.MacroGuardMode); resp != "ok" {
		t.Error(resp)
	}

	// wrong command
	if "err: invalid command" != command(client, "wrong") {
		t.Error("wrong value")
	}

	//--------------------------------

	// print protocol
	s, err := prettyString([]byte(respGS))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\nGameStatus:\n%s\n", s)

	s, err = prettyString([]byte(respTS))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\nTankStatus:\n%s\n", s)

	s, err = prettyString([]byte(respCT))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\nCloseTargets:\n%s\n", s)

	s, err = prettyString([]byte(respPT))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\nPossibleTargets:\n%s\n", s)
}

// prettyString returns JSON as pretty string
func prettyString(str []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, str, "", "  "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
