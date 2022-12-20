package resources

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/jpeg" // needed for ebitenutil.NewImageFromReader()
	_ "image/png"  // needed for ebitenutil.NewImageFromReader()
	"log"
)

// Imgs is the global variable that holds all image resources
var Imgs *ImgResources

// ImgResources is a collection of all images
type ImgResources struct {
	VictoryRed  *ebiten.Image // 285 x 120
	VictoryBlue *ebiten.Image // 285 x 120
	VictoryDraw *ebiten.Image // 285 x 120

	Ball      *ebiten.Image // 20 x 19
	Error     *ebiten.Image // 40 x 40
	Explosion *ebiten.Image // 64 x 64
	Grass     *ebiten.Image // 64 x 64
	Rock      *ebiten.Image // 64 x 64
	Logo      *ebiten.Image // 64 x 64

	Base1          *ebiten.Image // 64 x 64
	Tank1North     *ebiten.Image // 64 x 64
	Tank1Northeast *ebiten.Image // 64 x 64
	Tank1East      *ebiten.Image // 64 x 64
	Tank1Southeast *ebiten.Image // 64 x 64
	Tank1South     *ebiten.Image // 64 x 64
	Tank1Southwest *ebiten.Image // 64 x 64
	Tank1West      *ebiten.Image // 64 x 64
	Tank1Northwest *ebiten.Image // 64 x 64

	Base2          *ebiten.Image // 64 x 64
	Tank2North     *ebiten.Image // 64 x 64
	Tank2Northeast *ebiten.Image // 64 x 64
	Tank2East      *ebiten.Image // 64 x 64
	Tank2Southeast *ebiten.Image // 64 x 64
	Tank2South     *ebiten.Image // 64 x 64
	Tank2Southwest *ebiten.Image // 64 x 64
	Tank2West      *ebiten.Image // 64 x 64
	Tank2Northwest *ebiten.Image // 64 x 64
}

func init() {
	Imgs = &ImgResources{
		VictoryRed:  loadGameImg("img/victory_red.png"),
		VictoryBlue: loadGameImg("img/victory_blue.png"),
		VictoryDraw: loadGameImg("img/victory_draw.png"),

		Ball:      loadGameImg("img/ball.png"),
		Error:     loadGameImg("img/error.png"),
		Explosion: loadGameImg("img/explosion.png"),
		Grass:     loadGameImg("img/grass.png"),
		Rock:      loadGameImg("img/rock.png"),
		Logo:      loadGameImg("img/logo.png"),

		Base1:          loadGameImg("img/base1.png"),
		Tank1North:     loadGameImg("img/tank1_north.png"),
		Tank1Northeast: loadGameImg("img/tank1_northeast.png"),
		Tank1East:      loadGameImg("img/tank1_east.png"),
		Tank1Southeast: loadGameImg("img/tank1_southeast.png"),
		Tank1South:     loadGameImg("img/tank1_south.png"),
		Tank1Southwest: loadGameImg("img/tank1_southwest.png"),
		Tank1West:      loadGameImg("img/tank1_west.png"),
		Tank1Northwest: loadGameImg("img/tank1_northwest.png"),

		Base2:          loadGameImg("img/base2.png"),
		Tank2North:     loadGameImg("img/tank2_north.png"),
		Tank2Northeast: loadGameImg("img/tank2_northeast.png"),
		Tank2East:      loadGameImg("img/tank2_east.png"),
		Tank2Southeast: loadGameImg("img/tank2_southeast.png"),
		Tank2South:     loadGameImg("img/tank2_south.png"),
		Tank2Southwest: loadGameImg("img/tank2_southwest.png"),
		Tank2West:      loadGameImg("img/tank2_west.png"),
		Tank2Northwest: loadGameImg("img/tank2_northwest.png"),
	}
}

//go:embed img
var gFS embed.FS

func loadGameImg(name string) *ebiten.Image {
	// open reader
	r, err := gFS.Open(name)
	if err != nil {
		log.Fatalf("err: loadGameImg: %v\n", err)
	}
	// get image
	eim, _, err := ebitenutil.NewImageFromReader(r)
	if err != nil {
		log.Fatalf("err: loadGameImg: %v\n", err)
	}
	// return
	return eim
}
