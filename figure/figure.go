package figure

import (
	"math/rand/v2"
)

type Figure struct {
	Blocks [4][2]int
}

// Фигуры и их координаты спавна, поблочно
var MapGetFigure = map[string]Figure{
	"I": {
		Blocks: [4][2]int{{5, 2}, {5, 0}, {5, 1}, {5, 3}},
	},
	"O": {
		Blocks: [4][2]int{{5, 1}, {5, 0}, {6, 0}, {6, 1}},
	},
	"S": {
		Blocks: [4][2]int{{5, 1}, {5, 0}, {6, 0}, {4, 1}},
	},
	"J": {
		Blocks: [4][2]int{{5, 2}, {5, 0}, {5, 1}, {4, 2}},
	},
	"L": {
		Blocks: [4][2]int{{5, 1}, {5, 0}, {5, 2}, {6, 2}},
	},
	"Z": {
		Blocks: [4][2]int{{6, 1}, {5, 0}, {6, 0}, {7, 1}},
	},
	"T": {
		Blocks: [4][2]int{{5, 0}, {4, 0}, {6, 0}, {5, 1}},
	},
}

func FigureInNewBag(newBag [7]string, str string) bool {
	for _, current := range newBag {
		if current == str {
			return true
		}
	}

	return false
}

func RandomBag() [7]string {
	var (
		randInd int
		bag     [7]string = [7]string{"I", "O", "S", "J", "T", "Z", "L"}
		newBag  [7]string
	)

	newBag[0] = bag[rand.IntN(7)]
	for i := 1; i < 7; i++ {
		randInd = rand.IntN(7)
		for FigureInNewBag(newBag, bag[randInd]) {
			randInd = (randInd + 1) % 7
		}
		newBag[i] = bag[randInd]
	}

	return newBag
}

// 7 рандомных цветов без повторения
func GetRandomNums() []int {
	set := make(map[int]struct{})
	for len(set) < 7 {
		n := rand.IntN(7) + 1
		set[n] = struct{}{}
	}

	var result []int

	for k := range set {
		result = append(result, k)
	}

	return result
}
