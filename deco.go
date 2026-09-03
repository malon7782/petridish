package main

import "fmt"

func (w *World) decorateGrid(g [][]Cell) {
	w.decoRain(g, w.Weathers.RainIntensity)
}

func (w *World) decoRain(g [][]Cell, Intensity float64) {
	for i := 0; i < int(Intensity*300.0); i++ {
		randomy := w.Rng.Intn(w.Height)
		randomx := w.Rng.Intn(w.Width)
		g[randomy][randomx] = Cell{
			Char:  '`',
			Color: fmt.Sprintf("\033[38;5;%dm", 195),
		}
	}
}
