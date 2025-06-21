package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	bl "tetris/block"
	fig "tetris/figure"
	move "tetris/movement"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	// оступ слева и справа, 1300 - ширина окна
	Margin = (1300 - blockSize*20) / 2
)

// для работы текста, выполняется раньше всех функций
func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

}

const (
	blockSize = 51
)

type Game struct {
	frame         int
	gameOver      bool
	board         [18][11]int
	index         int
	figures       [7]string
	colors        []int
	currentFigure fig.Figure
	currentColor  int
	frameCount    int
	score         int
}

// Функция для инициализации сессии
func NewGame(frame int) *Game {
	// создаём тут, чтобы обратиться к их содержимому по индексу ниже
	figures := fig.RandomBag()
	colors := fig.GetRandomNums()
	return &Game{
		frame:         frame,
		gameOver:      false,
		board:         [18][11]int{},
		index:         1,
		figures:       figures,
		colors:        colors,
		currentFigure: fig.MapGetFigure[figures[0]],
		currentColor:  colors[0],
		frameCount:    0,
		score:         0,
	}
}

// функция для проверки заполненности строк
func (g *Game) checkLines() {
	var linesCount int
	for y := 17; y >= 0; y-- {
		flag := true
		for x := 1; x < 11; x++ {
			if g.board[y][x] == 0 {
				flag = false
				break
			}
		}

		// если строка одного цвета
		if flag {
			//считаем количество строк
			linesCount++
			for newY := y; newY > 0; newY-- {
				for x := 0; x < 11; x++ {
					g.board[newY][x] = g.board[newY-1][x]
				}
			}
			// верхняя строка
			for x := 0; x < 11; x++ {
				g.board[0][x] = 0
			}
			// чтобы убирать несколько строк подряд
			y++
		}
	}
	// Считаем очки за определенное количество строк
	switch linesCount {
	case 1:
		g.score += 100
	case 2:
		g.score += 300
	case 3:
		g.score += 700
	case 4:
		g.score += 1500
	}
}

// перезапуск при проигрыше
func restartGame(g *Game) *Game {
	g.gameOver = false
	g.board = [18][11]int{}
	g.index = 1
	g.figures = fig.RandomBag()
	g.colors = fig.GetRandomNums()
	g.currentFigure = fig.MapGetFigure[g.figures[0]]
	g.currentColor = g.colors[0]
	g.frameCount = 0
	return g
}

// логика игры
func (g *Game) Update() error {
	// для перазуска - нажать R
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			restartGame(g)
		}
		return nil
	}

	g.frameCount++
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.currentFigure = move.Spin(g.currentFigure, g.board)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.currentFigure = move.Right(g.currentFigure, g.board)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.currentFigure = move.Left(g.currentFigure, g.board)
	}
	// если нажата S ИЛИ количество прошедших кадров кратно выставленному обновлению (сложности), то опустить фигуру
	if g.frameCount%g.frame == 0 || ebiten.IsKeyPressed(ebiten.KeySpace) {
		if move.CanGoDown(g.currentFigure, g.board) {
			g.currentFigure = move.Down(g.currentFigure, g.board)
		} else {
			for _, block := range g.currentFigure.Blocks {
				x, y := block[0], block[1]
				g.board[y][x] = g.currentColor
			}
			g.checkLines()
			if g.index >= 7 {
				g.figures = fig.RandomBag()
				g.colors = fig.GetRandomNums()
				g.index = 0
			}
			g.currentFigure = fig.MapGetFigure[g.figures[g.index]]
			g.currentColor = g.colors[g.index]
			g.index++
		}
	}
	return nil
}

// окантовка из серых блоков
func (g *Game) DrawEdging(screen *ebiten.Image) {
	// серый квадрат со стороной BlockSize
	color := bl.MapGetColor[8]
	block := bl.GetBlock(blockSize, color.Basic, color.Darker, color.Lighter)
	for y := 0; y < 20; y++ {
		// рисуем слева
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(Margin)+0, float64(blockSize*y))
		screen.DrawImage(block, op)
		// рисуем справа
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(float64(Margin)+blockSize*13-float64(blockSize), float64(blockSize*y))
		screen.DrawImage(block, op2)
	}
	for x := 0; x < 13; x++ {
		// рисуем снизу
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(Margin)+float64(blockSize*x), float64(screen.Bounds().Dy())-float64(blockSize))
		screen.DrawImage(block, op)
		// рисуем сверху
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(float64(Margin)+float64(blockSize*x), 0)
		screen.DrawImage(block, op2)
	}
}

// вывод очков
func (g *Game) drawText(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(Margin)+720, 0)
	strScore := fmt.Sprintf("Current score: %d", g.score)
	text.Draw(screen, strScore, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)
}

// главная функция, отвечающая за отрисовку
func (g *Game) Draw(screen *ebiten.Image) {
	// проверка на заполненность поля
	if g.board[0][5] != 0 {
		g.DrawEdging(screen)
		g.gameOver = true
	}

	screen.Fill(color.Transparent)
	g.DrawEdging(screen)
	g.drawText(screen)

	// рисование с матрицы
	for x := 0; x < 11; x++ {
		for y := 0; y < 18; y++ {
			if color := g.board[y][x]; color != 0 {
				colorPal := bl.MapGetColor[color]
				block := bl.GetBlock(blockSize, colorPal.Basic, colorPal.Lighter, colorPal.Darker)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(Margin)+float64(blockSize+x*blockSize), float64(blockSize+y*blockSize))
				screen.DrawImage(block, op)
			}
		}
	}
	// рисование с временной фигуры в воздухе, она НЕ принадлежит матрице, становится её частью лишь после приземления
	blocks := g.currentFigure.Blocks
	for _, block := range blocks {
		x, y := block[0], block[1]
		colorPal := bl.MapGetColor[g.currentColor]
		block := bl.GetBlock(blockSize, colorPal.Basic, colorPal.Lighter, colorPal.Darker)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(Margin)+float64(blockSize+x*blockSize), float64(blockSize+y*blockSize))
		screen.DrawImage(block, op)
	}
}

// получение координатной системы внутри окошка, беру на всё окно
func (g *Game) Layout(
	outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
