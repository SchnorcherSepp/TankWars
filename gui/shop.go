package gui

import (
	"fmt"
	"github.com/SchnorcherSepp/TankWars/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"strings"
)

// drawShop draws the black shop to buy tanks
func drawShop(screen *ebiten.Image, screenWidth, screenHeight int, active *core.Tank, world *core.World) {

	// draw background
	frame := float64(150)
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0xdd}
	ebitenutil.DrawRect(screen, frame, frame, float64(screenWidth)-2*frame, float64(screenHeight)-2*frame, bgColor)

	// data
	var owner string
	var cash int
	if active != nil && active.Owner() == core.RedBase {
		owner = core.RedTank
		cash, _ = world.CashStat() // red
	} else {
		owner = core.BlueTank
		_, cash = world.CashStat() // blue
	}

	// write text
	txt := new(strings.Builder)
	txt.WriteString(fmt.Sprintf("Player %s:   $%d\n", owner, cash))
	txt.WriteString("--------------------\n\n")

	if cash < 100 {
		txt.WriteString("You don't have enough cash to buy a new tank!\n\n\n")

	} else {
		if shopState == 0 {
			txt.WriteString("Choose weapon type:\n")
			txt.WriteString(fmt.Sprintf(" [C] %s\n", core.WeaponCannon))
			txt.WriteString(fmt.Sprintf(" [R] %s\n", core.WeaponRockets))
			txt.WriteString(fmt.Sprintf(" [A] %s\n", core.WeaponArtillery))
		}
		if shopState > 0 {
			txt.WriteString(fmt.Sprintf("Weapon type: %s\n\n", shopType))
		}
		if shopState == 1 {
			txt.WriteString("Enter armor:\n")
			txt.WriteString(" [1] 5 armor\n")
			txt.WriteString(" [2] 15 armor\n")
			txt.WriteString(" [3] 25 armor\n")
			txt.WriteString(" [4] 35 armor\n")
			txt.WriteString(" [5] 45 armor\n")
			txt.WriteString(" [6] 55 armor\n")
		}
		if shopState > 1 {
			txt.WriteString(fmt.Sprintf("Armor: %d\n\n", shopArmor))
		}
		if shopState == 2 {
			txt.WriteString("Enter damage:\n")
			if shopArmor+15+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [1] 15 damage\n")
			}
			if shopArmor+25+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [2] 25 damage\n")
			}
			if shopArmor+35+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [3] 35 damage\n")
			}
			if shopArmor+45+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [4] 45 damage\n")
			}
			if shopArmor+55+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [5] 55 damage\n")
			}
			if shopArmor+65+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [6] 65 damage\n")
			}
			if shopArmor+70+core.TankMinSpeed <= core.TankBudget {
				txt.WriteString(" [7] 70 damage\n")
			}
		}
		if shopState > 2 {
			txt.WriteString(fmt.Sprintf("Damage: %d\n\n", shopDamage))

			txt.WriteString(fmt.Sprintf("ERROR: %v\n\n", shopErr))
			txt.WriteString("Close with [S]\n")
		}

	}

	// write text
	ebitenutil.DebugPrintAt(screen, txt.String(), int(frame)+50, int(frame)+50)
}
