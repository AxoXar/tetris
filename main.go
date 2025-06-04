package main

import (
	g "tetris/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 20 блоков по 51 в высоту
	ebiten.SetWindowSize(1300, 1020)
	ebiten.SetWindowTitle("Tetris")
	// фпс
	ebiten.SetTPS(15)
	// первый параметр вызова NewGame - то, как часто фигура опускается сама по себе, т.е. сложность
	// по заданию нужно добавить сложности, как по мне, 10 - очень сложно уже, 7 - средне, 5 - легко
	game := g.NewGame(10)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

}
