package game

import (
	"image/color"
	bl "tetris/block"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) DrawEdging(screen *ebiten.Image) {
	// серый квадрат со стороной n
	block := bl.GetBlock(blockSize, color.RGBA{119, 119, 119, 255}, color.RGBA{153, 153, 153, 255}, color.RGBA{49, 49, 49, 255})
	for y := 0; y < 20; y++ {
		// рисуем слева
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(blockSize*y))
		screen.DrawImage(block, op)
		// рисуем справа
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(float64(screen.Bounds().Dx())-float64(blockSize), float64(blockSize*y))
		screen.DrawImage(block, op2)
	}
	for x := 0; x < 12; x++ {
		// рисуем снизу
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(blockSize*x), float64(screen.Bounds().Dy())-float64(blockSize))
		screen.DrawImage(block, op)
		// рисуем сверху
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(float64(blockSize*x), 0)
		screen.DrawImage(block, op2)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawEdging(screen)
}

func (g *Game) Layout(
	outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 612, 1020
}

const (
	blockSize = 51
)

var board [18][10]int
