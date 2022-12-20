package remote

import (
	"bufio"
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"
)

// RunServer runs a server (BLOCKING!).
// The server receives commands from the clients and implements them in "World".
// The first connecting client controls player red.
// The second connecting client controls player blue.
func RunServer(host, port string, world *core.World) {
	world.Freeze(true) // wait for all player

	// Listen for incoming connections.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("RunServer: %v\n", err)
	}

	// Close the listener when the application closes.
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)

	fmt.Println("START SERVER [" + host + ":" + port + "]")
	for i := uint64(1); true; i++ {
		// wait for incoming connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		// Handle connections in a new goroutine.
		if i == 1 {
			// player 1: red
			owner := core.RedTank
			go handleRequest(conn, world, owner)
			fmt.Printf("player %d (%s) from %v\n", i, owner, conn.RemoteAddr())

		} else if i == 2 {
			// player 2: blue
			owner := core.BlueTank
			go handleRequest(conn, world, owner)
			fmt.Printf("player %d (%s) from %v\n", i, owner, conn.RemoteAddr())

			// START GAME with player 2!!
			world.Freeze(false)
			fmt.Printf("START GAME\n")

		} else {
			// server full
			owner := fmt.Sprintf("observer-%d", i-2)
			go handleRequest(conn, world, owner)
			fmt.Printf("%s from %v\n", owner, conn.RemoteAddr())
		}
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, w *core.World, owner string) {

	// prepare line reader
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// close at end
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// welcome
	_, _ = conn.Write([]byte("welcome player " + owner + "\n"))

	// loop
	for {
		// read one line (ended with \n or \r\n)
		line, err := tp.ReadLine()
		if err != nil {
			break // EXIT
		}

		// trim line and split args
		args := strings.Split(strings.TrimSpace(line), " ")

		// extract com
		var com string
		if len(args) > 0 {
			com = args[0]
		}

		// CHECK COMMANDS
		switch com {
		case "Exit":
			println("EXIT by player", owner)
			os.Exit(0)
		case "MyName":
			comResponse(conn, MyName(owner))
		case "GameStatus":
			comResponse(conn, GameStatus(w))
		case "TankStatus":
			tankID, _, _, _, _, _ := saveArgs(args)
			comResponse(conn, TankStatus(w, tankID))
		case "CloseTargets":
			tankID, filter1, filter2, filter3, filter4, filter5 := saveArgs(args)
			comResponse(conn, CloseTargets(w, tankID, filter1, filter2, filter3, filter4, filter5))
		case "PossibleTargets":
			tankID, filter1, filter2, filter3, filter4, filter5 := saveArgs(args)
			comResponse(conn, PossibleTargets(w, tankID, filter1, filter2, filter3, filter4, filter5))
		case "BuyTank":
			armor, damage, weapon, _, _, _ := saveArgs(args)
			comResponse(conn, BuyTank(w, owner, armor, damage, weapon))
		case "Fire":
			tankID, angle, distance, _, _, _ := saveArgs(args)
			comResponse(conn, Fire(w, owner, tankID, angle, distance))
		case "FireAt":
			tankID, x, y, _, _, _ := saveArgs(args)
			comResponse(conn, FireAt(w, owner, tankID, x, y))
		case "Forward":
			tankID, _, _, _, _, _ := saveArgs(args)
			comResponse(conn, Forward(w, owner, tankID))
		case "Backward":
			tankID, _, _, _, _, _ := saveArgs(args)
			comResponse(conn, Backward(w, owner, tankID))
		case "Stop":
			tankID, _, _, _, _, _ := saveArgs(args)
			comResponse(conn, Stop(w, owner, tankID))
		case "Left":
			tankID, _, _, _, _, _ := saveArgs(args)
			comResponse(conn, Left(w, owner, tankID))
		case "Right":
			tankID, _, _, _, _, _ := saveArgs(args)
			comResponse(conn, Right(w, owner, tankID))
		case "SetMacroMoveTo":
			tankID, x, y, _, _, _ := saveArgs(args)
			comResponse(conn, SetMacroMoveTo(w, owner, tankID, x, y))
		case "SetMacro":
			tankID, macro, _, _, _, _ := saveArgs(args)
			comResponse(conn, SetMacro(w, owner, tankID, macro))
		default:
			comResponse(conn, "err: invalid command")
		}
	}

	// exit
	fmt.Printf("player %s has left\n", owner)
}

// comResponse is a helper function and send messages back to the clients.
func comResponse(conn net.Conn, s string) {
	_, err := conn.Write([]byte(fmt.Sprintf("%s\r\n", s)))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

// saveArgs is a helper function and return 6 string arguments from the client commands
func saveArgs(args []string) (a1, a2, a3, a4, a5, a6 string) {
	sArgs := make([]string, 7)
	copy(sArgs, args)
	return sArgs[1], sArgs[2], sArgs[3], sArgs[4], sArgs[5], sArgs[6]
}
