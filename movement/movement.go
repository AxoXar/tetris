package movement

import (
	fig "tetris/figure"
)

// Нужно добавить функцию для поворота на W

func Right(piece fig.Figure, board [18][11]int) fig.Figure {
	for _, block := range piece.Blocks {
		x, y := block[0], block[1]
		if x > 9 || board[y][x+1] != 0 {
			return piece
		}
	}

	var newBlocks [4][2]int
	for i, block := range piece.Blocks {
		newBlocks[i][0] = block[0] + 1 // x + 1 - сдвигаем направо
		newBlocks[i][1] = block[1]     // y - оставляем
	}

	return fig.Figure{Blocks: newBlocks}
}

func Left(piece fig.Figure, board [18][11]int) fig.Figure {
	for _, block := range piece.Blocks {
		x, y := block[0], block[1]
		if x < 1 || board[y][x-1] != 0 {
			return piece
		}
	}

	var newBlocks [4][2]int
	for i, block := range piece.Blocks {
		newBlocks[i][0] = block[0] - 1 // x - 1 - сдвигаем налево
		newBlocks[i][1] = block[1]     //  y - оставляем
	}

	return fig.Figure{Blocks: newBlocks}
}

func CanGoDown(piece fig.Figure, board [18][11]int) bool {
	for _, block := range piece.Blocks {
		x, y := block[0], block[1]
		if y > 16 || board[y+1][x] != 0 {
			return false
		}
	}
	return true
}

func Down(piece fig.Figure, board [18][11]int) fig.Figure {
	var newBlocks [4][2]int
	for i, block := range piece.Blocks {
		newBlocks[i][0] = block[0]     // x — без изменений
		newBlocks[i][1] = block[1] + 1 // y +1 — вниз
	}
	return fig.Figure{Blocks: newBlocks}
}
