package menu

import (
	"tetris/game"
)

type menu struct {
	startNewG game.Game
	continueG game.Game
	settings  bool
	gameplay  bool
}
