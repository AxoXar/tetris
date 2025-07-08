package game

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	bl "tetris/block"
	fig "tetris/figure"
	move "tetris/movement"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	// оступ слева и справа
	Margin = (1920 - blockSize*13) / 2
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
	Difficulty    int
	GameOver      bool
	Board         [18][11]int
	Index         int
	Figures       [7]string
	CurrentFigure fig.Figure
	CurrentColor  int
	FrameCount    int
	Score         int
	Records       []string
	Started       bool
	Paused        bool
	HoldKeys      map[string]int
}

// Функция для инициализации сессии
func NewGame() *Game {
	// создаём тут, чтобы обратиться к их содержимому по индексу ниже
	figures := fig.RandomBag()
	holdKeys := map[string]int{
		"W":     0,
		"A":     0,
		"D":     0,
		"M1":    0,
		"Space": 0,
		"Esc":   0,
	}
	return &Game{
		Difficulty:    0,
		GameOver:      false,
		Board:         [18][11]int{},
		Index:         1,
		Figures:       figures,
		CurrentFigure: fig.MapGetFigure[figures[0]],
		CurrentColor:  fig.MapGetColor[figures[0]],
		FrameCount:    0,
		Score:         0,
		Records:       nil,
		HoldKeys:      holdKeys,
	}

}

// перезапуск при проигрыше
func restartGame(g *Game) *Game {
	g.WriteRecord()
	g.GameOver = false
	g.Board = [18][11]int{}
	g.Index = 1
	g.Figures = fig.RandomBag()
	g.CurrentFigure = fig.MapGetFigure[g.Figures[0]]
	g.CurrentColor = fig.MapGetColor[g.Figures[0]]
	g.FrameCount = 0
	g.Score = 0
	g.Records = ReadRecord(g.Difficulty)
	return g
}

// логика игры
func (g *Game) Update() error {
	// пауза
	pressedEsc := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if pressedEsc {
		g.HoldKeys["Esc"]++
		if g.HoldKeys["Esc"] == 1 {
			g.Paused = !(g.Paused)
		}
	} else if !pressedEsc {
		g.HoldKeys["Esc"] = 0
	}

	// не входить в логику игры при паузе
	if g.Paused {
		return nil
	}

	pressedMouse := ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	if pressedMouse {
		g.HoldKeys["M1"]++
		if g.HoldKeys["M1"] == 1 {
			x, y := ebiten.CursorPosition()
			for i := 0; i < 3; i++ {
				if (!g.Started || g.GameOver) && (x >= 110 && x <= 510) && (y >= (75+i*200) && y <= (75+i*200+100)) {
					g.Difficulty = 3 - i
					g.Records = ReadRecord(g.Difficulty)
					if !g.Started {
						g.Started = true
					} else if g.Started {
						restartGame(g)
					}
				}
			}

			if (g.Started && !g.GameOver) && (x >= 110 && x <= 510) && (y >= (75+3*200) && y <= (75+3*200+100)) {
				err := g.SaveProgress("save.json")
				if err != nil {
					fmt.Println("Ошибка сохранения:", err)
				}
			} else if (!g.Started || g.GameOver) && (x >= 110 && x <= 510) && (y >= (75+4*200) && y <= (75+4*200+100)) {
				err := g.LoadProgress("save.json")
				if err != nil {
					fmt.Println("Ошибка загрузки:", err)
				}
			}
		}
	} else if !pressedMouse {
		g.HoldKeys["M1"] = 0
	}

	// для перезапуска после поражения - нажать R
	if g.GameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			restartGame(g)
		}
		return nil
	}

	// не входить в логику, если игра не запущена
	if !g.Started {
		return nil
	}
	g.FrameCount++

	pressedW := ebiten.IsKeyPressed(ebiten.KeyW)
	if pressedW {
		g.HoldKeys["W"]++
		if g.HoldKeys["W"] == 1 {
			g.CurrentFigure = move.Spin(g.CurrentFigure, g.Board)
		}
	} else if !pressedW {
		g.HoldKeys["W"] = 0
	}

	pressedD := ebiten.IsKeyPressed(ebiten.KeyD)
	if pressedD {
		g.HoldKeys["D"]++
		if g.HoldKeys["D"] == 1 {
			g.CurrentFigure = move.Right(g.CurrentFigure, g.Board)
		} else if g.HoldKeys["D"] > 5 && g.HoldKeys["D"]%2 == 0 {
			g.CurrentFigure = move.Right(g.CurrentFigure, g.Board)
		}
	} else if !pressedD {
		g.HoldKeys["D"] = 0
	}

	pressedA := ebiten.IsKeyPressed(ebiten.KeyA)
	if pressedA {
		g.HoldKeys["A"]++
		if g.HoldKeys["A"] == 1 {
			g.CurrentFigure = move.Left(g.CurrentFigure, g.Board)
		} else if g.HoldKeys["A"] > 5 && g.HoldKeys["A"]%2 == 0 {
			g.CurrentFigure = move.Left(g.CurrentFigure, g.Board)
		}
	} else if !pressedA {
		g.HoldKeys["A"] = 0
	}

	// если нажата S ИЛИ количество прошедших кадров кратно выставленному обновлению (сложности), то опустить фигуру
	pressedSpace := ebiten.IsKeyPressed(ebiten.KeySpace)
	if pressedSpace {
		g.HoldKeys["Space"]++
		if g.HoldKeys["Space"] == 1 {
			g.loweringFigure()
		} else if g.HoldKeys["Space"] > 5 && g.HoldKeys["Space"]%2 == 0 {
			g.loweringFigure()
		}
	} else if (g.FrameCount)%(g.Difficulty*3) == 0 {
		g.loweringFigure()
	} else if !pressedSpace {
		g.HoldKeys["Space"] = 0
	}

	return nil
}

func (g *Game) SaveProgress(filename string) error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (g *Game) LoadProgress(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, g)
}

func (g *Game) loweringFigure() {
	if move.CanGoDown(g.CurrentFigure, g.Board) {
		g.CurrentFigure = move.Down(g.CurrentFigure, g.Board)
	} else {
		for _, block := range g.CurrentFigure.Blocks {
			x, y := block[0], block[1]
			g.Board[y][x] = g.CurrentColor
		}
		g.checkLines()
		g.CurrentFigure = fig.MapGetFigure[g.Figures[g.Index]]
		g.CurrentColor = fig.MapGetColor[g.Figures[g.Index]]
		g.Index++

		if g.Index >= 7 {
			g.Figures = fig.RandomBag()
			g.Index = 0
		}
	}
}

// функция для проверки заполненности строк
func (g *Game) checkLines() {
	var linesCount int
	for y := 17; y >= 0; y-- {
		flag := true
		for x := 0; x < 11; x++ {
			if g.Board[y][x] == 0 {
				flag = false
				break
			}
		}

		// если строка заполнена
		if flag {
			//считаем количество строк
			linesCount++
			for newY := y; newY > 0; newY-- {
				for x := 0; x < 11; x++ {
					g.Board[newY][x] = g.Board[newY-1][x]
				}
			}
			// верхняя строка
			for x := 0; x < 11; x++ {
				g.Board[0][x] = 0
			}
			// чтобы убирать несколько строк подряд
			y++
		}
	}
	// Считаем очки за определенное количество строк
	switch linesCount {
	case 1:
		g.Score += 100
	case 2:
		g.Score += 300
	case 3:
		g.Score += 700
	case 4:
		g.Score += 1500
	}
}

// главная функция, отвечающая за отрисовку
func (g *Game) Draw(screen *ebiten.Image) {
	// заполним внутренность чёрным, а наружу - серым
	screen.Fill(color.RGBA{45, 45, 45, 255})
	block := ebiten.NewImage(11*blockSize, 18*blockSize)
	block.Fill(color.Black)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(Margin)+blockSize, blockSize)
	screen.DrawImage(block, op)

	g.drawButtons(screen)
	g.DrawEdging(screen)

	// Если игра ещё не начата
	if !g.Started {
		g.drawRules(screen)
		return
	}
	// Если игра на паузе
	if g.Paused {
		g.drawPauseInfo(screen)
		return
	}

	g.drawRecords(screen)

	// проверка на заполненность поля
	if g.Board[0][5] != 0 {
		g.GameOver = true
		g.drawLose(screen)
		return
	}

	g.DrawBoard(screen)
	g.drawNextFig(screen)
	g.drawScore(screen)
	g.drawDifficulty(screen)
	// рисование с матрицы
	for x := 0; x < 11; x++ {
		for y := 0; y < 18; y++ {
			if color := g.Board[y][x]; color != 0 {
				colorPal := bl.MapGetColor[color]
				block := bl.GetBlock(blockSize, blockSize, 2, colorPal.Basic, colorPal.Lighter, colorPal.Darker)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(Margin)+float64(blockSize+x*blockSize), float64(blockSize+y*blockSize))
				screen.DrawImage(block, op)
			}
		}
	}
	// рисование с временной фигуры в воздухе, она НЕ принадлежит матрице, становится её частью лишь после приземления
	blocks := g.CurrentFigure.Blocks
	for _, block := range blocks {
		x, y := block[0], block[1]
		colorPal := bl.MapGetColor[g.CurrentColor]
		block := bl.GetBlock(blockSize, blockSize, 2, colorPal.Basic, colorPal.Lighter, colorPal.Darker)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(Margin)+float64(blockSize+x*blockSize), float64(blockSize+y*blockSize))
		screen.DrawImage(block, op)
	}
}

func (g *Game) drawButtons(screen *ebiten.Image) {
	buttonText := make(map[int]string)
	buttonText[0] = "Лёгкий"
	buttonText[1] = "Средний"
	buttonText[2] = "Сложный"
	buttonText[3] = "Сохранить"
	buttonText[4] = "Продолжить"

	face := text.GoTextFace{
		Source: mplusFaceSource,
		Size:   17,
	}

	for i := 0; i < 5; i++ {
		// кнопки-блоки
		button := bl.GetBlock(400, 100, 2, color.RGBA{30, 144, 255, 255}, color.RGBA{56, 78, 255, 255}, color.RGBA{56, 78, 255, 255})
		opButton := &ebiten.DrawImageOptions{}
		opButton.GeoM.Translate(110, float64(75+i*200))
		screen.DrawImage(button, opButton)
		// Печать текста на кнопках
		opText := &text.DrawOptions{}
		opText.GeoM.Translate(110+(400-4.5*float64(len(buttonText[i])))/2, float64(75+i*200)+34)
		text.Draw(screen, buttonText[i], &face, opText)
	}
}

func (g *Game) drawRules(screen *ebiten.Image) {
	lines := []string{
		"Управление: W - поворот фигуры; A, D, Space - перемещение",
		"Escape - поставить/снять паузу",
		"Выберите сложность из предложенных",
		"В любой момент вы можете сохранить и загрузить игру",
		"Задача: полностью наполнять линии блоками",
	}
	for i := 0; i < 5; i++ {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(Margin)+blockSize*2, float64(400+50*i))
		text.Draw(screen, lines[i], &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   16,
		}, op)
	}
}

func (g *Game) drawPauseInfo(screen *ebiten.Image) {
	lines := []string{
		"Игра приоставновлена",
		"Нажмите Esc, чтобы продолжить",
	}
	for i := 0; i < 2; i++ {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(Margin)+blockSize*4, float64(400+50*i))
		text.Draw(screen, lines[i], &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   16,
		}, op)
	}
}

// окантовка из серых блоков
func (g *Game) DrawEdging(screen *ebiten.Image) {
	// серый квадрат со стороной BlockSize
	color := bl.MapGetColor[8]
	block := bl.GetBlock(blockSize, blockSize, 2, color.Basic, color.Darker, color.Lighter)
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

func (g *Game) DrawBoard(screen *ebiten.Image) {
	block := bl.GetBlock(blockSize, blockSize, 1, color.RGBA{0, 0, 0, 0}, color.RGBA{25, 25, 25, 255}, color.RGBA{25, 25, 25, 255})
	for y := 1; y < 19; y++ {
		for x := 1; x < 12; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(Margin+blockSize*x), float64(blockSize*y))
			screen.DrawImage(block, op)
		}
	}
}

// вывод окошка "следующая фигура"
func (g *Game) drawNextFig(screen *ebiten.Image) {
	// фон для блока следующей фигуры
	block1 := ebiten.NewImage(350, 350)
	block1.Fill(color.RGBA{25, 25, 25, 255})
	opBlock1 := &ebiten.DrawImageOptions{}
	opBlock1.GeoM.Translate(float64(Margin)+775, 575)

	block2 := ebiten.NewImage(300, 300)
	block2.Fill(color.Black)
	opBlock2 := &ebiten.DrawImageOptions{}
	opBlock2.GeoM.Translate(float64(Margin)+800, 600)

	screen.DrawImage(block1, opBlock1)
	screen.DrawImage(block2, opBlock2)
	// вывод текста
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(Margin)+879, 602)

	text.Draw(screen, "Next Figure", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)

	blocks := fig.MapGetFigure[g.Figures[g.Index]].Blocks
	for _, block := range blocks {
		x, y := block[0], block[1]
		colorPal := bl.MapGetColor[fig.MapGetColor[g.Figures[g.Index]]]
		block := bl.GetBlock(blockSize, blockSize, 2, colorPal.Basic, colorPal.Lighter, colorPal.Darker)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(Margin)+float64(blockSize+x*blockSize)+600, float64(blockSize+y*blockSize)+620)
		screen.DrawImage(block, op)
	}
}

// вывод текущих очков
func (g *Game) drawScore(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(Margin)+817, 20)
	strScore := fmt.Sprintf("Current score: %d", g.Score)
	text.Draw(screen, strScore, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)
}

// вывод текущих очков
func (g *Game) drawDifficulty(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	difficulties := []string{
		"Сложный",
		"Средний",
		"Лёгкий",
	}
	op.GeoM.Translate(167, 20)
	strScore := fmt.Sprintf("Difficulty level: %s", difficulties[g.Difficulty-1])
	text.Draw(screen, strScore, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)
}

// вывод при поражении
func (g *Game) drawLose(screen *ebiten.Image) {
	scoreText := "Вы проиграли, набрав " + strconv.Itoa(g.Score) + " очков"
	lines := []string{
		scoreText,
		"Для перезапуска нажмите R",
	}
	for i := 0; i < 2; i++ {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(Margin)+blockSize*4, float64(400+50*i))
		text.Draw(screen, lines[i], &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   16,
		}, op)
	}
}

// вывод рекордов
func (g *Game) drawRecords(screen *ebiten.Image) {
	// фон для блока рекордов
	block1 := ebiten.NewImage(350, 350)
	block1.Fill(color.RGBA{56, 78, 255, 255})
	opBlock1 := &ebiten.DrawImageOptions{}
	opBlock1.GeoM.Translate(float64(Margin)+790, 75)
	screen.DrawImage(block1, opBlock1)

	block2 := ebiten.NewImage(300, 300)
	block2.Fill(color.RGBA{30, 144, 255, 255})
	opBlock2 := &ebiten.DrawImageOptions{}
	opBlock2.GeoM.Translate(float64(Margin)+815, 100)
	screen.DrawImage(block2, opBlock2)

	// вывод текста
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(Margin)+915, 100)

	text.Draw(screen, "Records", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)

	line := g.Records
	for i := 0; i < 5; i++ {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(Margin)+820, float64(150+i*50))

		str := strconv.Itoa(i+1) + ". " + line[i] + "\n"
		text.Draw(screen, str, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   20,
		}, op)
	}

}

// чтение файла с рекордами
func ReadRecord(num int) []string {
	readFile, err := os.Open("records.txt")
	if os.IsNotExist(err) {
		f, _ := os.Create("records.txt")
		f.WriteString("0 0 0 0 0\n")
		f.WriteString("0 0 0 0 0\n")
		f.WriteString("0 0 0 0 0\n")

		f.Close()

		readFile, err = os.Open("records.txt")
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return strings.Fields(fileLines[num-1])
}

// перезапись файла с рекордами
func (g *Game) WriteRecord() {
	readFile, err := os.Open("records.txt")
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	splittedLine := strings.Fields(fileLines[g.Difficulty-1])
	for counter, field := range splittedLine {
		num, _ := strconv.Atoi(field)
		if g.Score > num {
			copiedLine := make([]string, len(splittedLine))
			copy(copiedLine, splittedLine)
			splittedLine[counter] = strconv.Itoa(g.Score)

			if counter != 4 {
				for i := counter + 1; i < 5; i++ {
					splittedLine[i] = copiedLine[i-1]
				}
			}
			break
		}
	}

	fileLines[g.Difficulty-1] = strings.Join(splittedLine, " ")

	os.WriteFile("records.txt", []byte(strings.Join(fileLines, "\n")), 0666)
}

// получение координатной системы внутри окошка, беру на всё окно
func (g *Game) Layout(
	outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
