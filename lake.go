package main

import "fmt"

type Lake struct {
	Height   int
	ColorStr string
}

func (l *Lake) Icon() byte {
	if l == nil {
		return '%'
	}
	return '#'
}

func (l *Lake) Color() string {
	if l == nil {
		return ""
	}
	return l.ColorStr
}

// this is stupid. we need more algos e.g. generateCirculerLake(), generateCoast()...
func cumLake(w *World, numLakes int) {
	w.Lakes = make([][]*Lake, w.Height)
	for y := 0; y < w.Height; y++ {
		w.Lakes[y] = make([]*Lake, w.Width)
	}
	for i := 0; i < numLakes; i++ {
		randomy := w.Rng.Intn(w.Height)
		randomx := w.Rng.Intn(w.Width)

		w.Lakes[randomy][randomx] = &Lake{
			Height:   w.Map[randomy][randomx].Height,
			ColorStr: fmt.Sprintf("\033[38;5;%dm", 45),
		}
		// this is important. we need to dig down a unit depth  so that we can
		//* place water. in case different lake gen functions access the same [][]*Mounatin
		// block, make sure that the height of that block is equal or greater than -1,
		w.Map[randomy][randomx].Height -= 1
	}
}
