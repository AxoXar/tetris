package main

import (
	"fmt"
	g "tetris/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// 20 блоков по 51 в высоту
	ebiten.SetWindowSize(1920, 1020)
	ebiten.SetWindowTitle("Tetris")
	// фпс
	ebiten.SetTPS(30)
	game := g.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
	defer saveAtexit(game)
}

func saveAtexit(game *g.Game) {
	err := game.SaveProgress("save.json")
	if err != nil {
		fmt.Println("Ошибка сохранения:", err)
	}
}
