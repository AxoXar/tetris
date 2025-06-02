package newBag

import (
	"math/rand/v2"
)

func FigureInNewBag(newBag [7]string, str string) bool {
	for _, current := range newBag {
		if current == str {
			return true
		}
	}

	return false
}

func randomBag() [7]string {
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
