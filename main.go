package main

import (
	g "tetris/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 20 блоков в высоту, 12 в ширину
	ebiten.SetWindowSize(612, 1020)
	ebiten.SetWindowTitle("Tetris")
	game := &g.Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
