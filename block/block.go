package block

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func GetBlock(n int, basic color.Color, darker color.Color, lighter color.Color) *ebiten.Image {
	block := ebiten.NewImage(n, n)
	block.Fill(basic)
	for i := 0; i < n; i++ {
		block.Set(i, 1, darker)
		block.Set(i, 2, darker)

		block.Set(1, i, darker)
		block.Set(2, i, darker)

		block.Set(i, n, lighter)
		block.Set(i, n-1, lighter)

		block.Set(n, i, lighter)
		block.Set(n-1, i, lighter)
	}
	return block
}
