package resources

import (
	"bytes"
	"embed"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"io"
	"log"
	"time"
)

// Mp3SampleRate is the sample rate of all sound files
const Mp3SampleRate = 44100

// Sounds is the global variable that holds all sound resources
var Sounds *SoundResources

// SoundResources is a collection of all sounds
type SoundResources struct {
	Fire      []byte
	Explosion []byte
}

var otoContext *oto.Context

func init() {
	// files
	Sounds = &SoundResources{
		Fire:      loadGameSound("sound/fire.mp3"),
		Explosion: loadGameSound("sound/explosion.mp3"),
	}
	// init oto context
	c, ready, err := oto.NewContext(Mp3SampleRate, 2, 2)
	if err != nil {
		panic(err)
	}
	<-ready
	otoContext = c
}

//go:embed sound
var sFS embed.FS

func loadGameSound(name string) []byte {
	// open reader
	r, err := sFS.Open(name)
	if err != nil {
		log.Fatalf("err: loadGameSound: %v\n", err)
	}
	// read all
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("err: loadGameSound: %v\n", err)
	}
	// return
	return b
}

//--------------------------------------------------------------------------------------------------------------------//

// MuteSound can disable all sounds (used for tests and simulations)
var MuteSound = false

// PlaySound play the given sound.
// Use MuteSound to disable this behavior.
func PlaySound(b []byte) {
	if MuteSound {
		return
	}

	go func(b []byte) {
		// decode
		d, err := mp3.NewDecoder(bytes.NewReader(b))
		if err != nil {
			panic(err)
		}

		// create player
		p := otoContext.NewPlayer(d)
		defer func(p oto.Player) {
			_ = p.Close()
		}(p)
		p.Play()

		// wait
		for {
			time.Sleep(time.Second)
			if !p.IsPlaying() {
				break
			}
		}
	}(b)
}
