package main

import (
	"flag"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/SchnorcherSepp/TankWars/examples/goai"
	"github.com/SchnorcherSepp/TankWars/gui"
	"github.com/SchnorcherSepp/TankWars/macro"
	"github.com/SchnorcherSepp/TankWars/maps"
	"github.com/SchnorcherSepp/TankWars/remote"
	"log"
	"os"
)

const version = "1.0"

func main() {
	mapName := flag.String("map", "field", "maps: 'field', 'fortress' or 'random'")
	speed := flag.Int("speed", 1, "speed multiplier")
	mute := flag.Bool("mute", false, "disable sound")
	aiMode := flag.Bool("ai", false, "start ai client")
	srvMode := flag.Bool("server", false, "start server and wait for two player")
	srvAddr := flag.String("host", "127.0.0.1", "server ip")
	srvPort := flag.String("port", "3333", "server port")

	flag.Parse()

	// print defaults
	if len(os.Args) <= 1 {
		println("TankWars", version)
		println("---------------")
		flag.PrintDefaults()
		println()
		os.Exit(0)
	}

	//------------------------------------------------------------------------------

	// AI MODE
	if *aiMode {
		goai.MyAI(*srvAddr, *srvPort) // blocking
		os.Exit(0)
	}

	// create world
	w := core.NewWorld(core.WorldXWidth, core.WorldYHeight)

	// run server?
	if *srvMode {
		go remote.RunServer(*srvAddr, *srvPort, w)
	}

	// map 'field', 'fortress' or 'random'
	switch *mapName {
	case "field":
		maps.InitOpenField(w)
	case "fortress":
		maps.InitFortress(w)
	case "random":
		maps.InitRandomWorld(w, 1337, func(t *core.Tank) { macro.AttackMove(t) })
	default:
		log.Fatal("unknown map")
	}

	// run gui (blocking)
	if err := gui.RunGame("Tank Wars "+version, w, *speed, *mute); err != nil {
		panic(err)
	}
}
