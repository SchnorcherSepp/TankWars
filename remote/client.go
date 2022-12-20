package remote

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"sync"
)

//---------------- TcpClient -----------------------------------------------------------------------------------------//

// TcpClient is an API to access a server
// and to remote control player tanks.
type TcpClient struct {
	conn *net.TCPConn
	tp   *textproto.Reader
	mux  *sync.Mutex
}

// NewTcpClient init a TcpClient
func NewTcpClient(host, port string) *TcpClient {

	// address
	tcpAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("NewTcpClient: %v\n", err)
	}

	// connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("NewTcpClient: %v\n", err)
	}

	// config client
	tc := &TcpClient{
		conn: conn,
		tp:   textproto.NewReader(bufio.NewReader(conn)),
		mux:  new(sync.Mutex),
	}

	// read first line after connect
	first, err := tc.tp.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	println(first)

	// return client
	return tc
}

//---------------- GETTER --------------------------------------------------------------------------------------------//

// MyName returns the active player of this connection (RedTank or BlueTank)
func (tc *TcpClient) MyName() string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, "MyName")
}

// GameStatus returns a json with all world data.
func (tc *TcpClient) GameStatus() string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, "GameStatus")
}

// TankStatus returns a json with all data of a requested tank.
func (tc *TcpClient) TankStatus(tankID string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("TankStatus %s", tankID))
}

// CloseTargets returns all objects in the world that are theoretical in weapon range.
// The weapon type is irrelevant (WeaponCannon or WeaponArtillery) and the angle of the tank is ignored.
// The list is sorted by distance (from the closest to the farthest).
func (tc *TcpClient) CloseTargets(tankID, f1, f2, f3, f4, f5 string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("CloseTargets %s %s %s %s %s %s", tankID, f1, f2, f3, f4, f5))
}

// PossibleTargets extends CloseTargets.
// It only returns objects that can actually be attacked,depending on the weapon type.
// However, it may be necessary for the battle tank to change its angle.
// The list is sorted by the rotation required to reach the target.
func (tc *TcpClient) PossibleTargets(tankID, f1, f2, f3, f4, f5 string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("PossibleTargets %s %s %s %s %s %s", tankID, f1, f2, f3, f4, f5))
}

//---------------- SETTER --------------------------------------------------------------------------------------------//

// Exit kills the server (for tests only).
func (tc *TcpClient) Exit() string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, "Exit")
}

// BuyTank buy a new tank and place it near the home base.
func (tc *TcpClient) BuyTank(armor, damage int, weapon string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("BuyTank %d %d %s", armor, damage, weapon))
}

// Fire creates a new projectile.
// The attributes fireAngle and distance determine the direction and distance of the shot.
// Cannons can fire in vehicle angle only.
// The distance is limited by the weapon range.
func (tc *TcpClient) Fire(tankID string, angle, distance int) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("Fire %s %d %d", tankID, angle, distance))
}

// FireAt is a wrapper for Fire() and convert the position to fireAngle and distance.
func (tc *TcpClient) FireAt(tankID string, x, y int) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("FireAt %s %d %d", tankID, x, y))
}

// Forward send the tank forward.
func (tc *TcpClient) Forward(tankID string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("Forward %s", tankID))
}

// Backward send the tank back.
func (tc *TcpClient) Backward(tankID string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("Backward %s", tankID))
}

// Stop the movement.
// Weapons can only build up when the tank is stationary
func (tc *TcpClient) Stop(tankID string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("Stop %s", tankID))
}

// Left turn the tank direction 45° left.
func (tc *TcpClient) Left(tankID string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("Left %s", tankID))
}

// Right turn the tank direction 45° right.
func (tc *TcpClient) Right(tankID string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("Right %s", tankID))
}

// SetMacroMoveTo sets a special macro with a position that is called with every update.
func (tc *TcpClient) SetMacroMoveTo(tankID string, x, y int) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("SetMacroMoveTo %s %d %d", tankID, x, y))
}

// SetMacro sets a macro that is called with every update.
func (tc *TcpClient) SetMacro(tankID, macro string) string {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	return command(tc, fmt.Sprintf("SetMacro %s %s", tankID, macro))
}

//---------------- HELPER --------------------------------------------------------------------------------------------//

// command send the cmd to the server and return the response
func command(tc *TcpClient, cmd string) string {
	if tc == nil || tc.conn == nil || tc.tp == nil {
		return "err: TcpClient connection closed."
	}

	// remove protocol break
	cmd = strings.ReplaceAll(cmd, "\n", "")
	cmd = strings.ReplaceAll(cmd, "\r", "")
	cmd = strings.ReplaceAll(cmd, "  ", " ")

	// send command
	_, err := tc.conn.Write([]byte(fmt.Sprintf("%s\r\n", cmd)))
	if err != nil {
		return fmt.Sprintf("err: TcpClient write: %v", err)
	}

	// read response
	resp, err := tc.tp.ReadLine()
	if err != nil {
		return fmt.Sprintf("err: TcpClient read: %v", err)
	}

	// return server response
	return resp
}
