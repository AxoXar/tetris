package block

import (
	c "image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type BlockColor struct {
	Basic   c.Color
	Darker  c.Color
	Lighter c.Color
}

var MapGetColor = map[int]BlockColor{
	1: {Basic: c.RGBA{204, 102, 0, 255}, Darker: c.RGBA{153, 68, 0, 255}, Lighter: c.RGBA{255, 137, 0, 255}},        // оранжевый
	2: {Basic: c.RGBA{204, 0, 0, 255}, Darker: c.RGBA{153, 0, 0, 255}, Lighter: c.RGBA{255, 0, 0, 255}},             // красный
	3: {Basic: c.RGBA{0, 0, 204, 255}, Darker: c.RGBA{0, 0, 153, 255}, Lighter: c.RGBA{0, 0, 255, 255}},             // синий
	4: {Basic: c.RGBA{0, 204, 204, 255}, Darker: c.RGBA{0, 153, 153, 255}, Lighter: c.RGBA{0, 250, 250, 255}},       // голубой
	5: {Basic: c.RGBA{0, 204, 0, 255}, Darker: c.RGBA{0, 153, 0, 255}, Lighter: c.RGBA{0, 255, 0, 255}},             // зелёный
	6: {Basic: c.RGBA{204, 204, 0, 255}, Darker: c.RGBA{153, 153, 0, 255}, Lighter: c.RGBA{250, 250, 0, 255}},       // жёлтый
	7: {Basic: c.RGBA{255, 105, 180, 255}, Darker: c.RGBA{219, 112, 147, 255}, Lighter: c.RGBA{255, 161, 174, 255}}, // розовый
	8: {Basic: c.RGBA{119, 119, 119, 255}, Darker: c.RGBA{153, 153, 153, 255}, Lighter: c.RGBA{49, 49, 49, 255}},    // серый
}

// Получение блока с окантовкой
func GetBlock(n int, basic c.Color, darker c.Color, lighter c.Color) *ebiten.Image {
	block := ebiten.NewImage(n, n)
	block.Fill(basic)
	for i := 0; i < n; i++ {
		block.Set(i, 0, darker)
		block.Set(i, 1, darker)

		block.Set(0, i, darker)
		block.Set(1, i, darker)
	}

	for i := 0; i < n; i++ {
		block.Set(i, n-1, lighter)
		block.Set(i, n-2, lighter)

		block.Set(n-1, i, lighter)
		block.Set(n-2, i, lighter)
	}

	return block
}
