package main

import (
	"fmt"
	"math"
)

func (w *World) decorateGrid(g [][]Cell) {
	w.decoRain(g, w.Weathers.RainIntensity)
	w.decoLightning(g, w.Weathers.RainIntensity)
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

// lightning decoration

var (
	lightningDays   int
	lightningRolled bool
)

func (w *World) decoLightning(g [][]Cell, Intensity float64) {
	// if lightning is not going on, roll one.
	if lightningRolled == false && w.Weathers.RainLeft >= 4 &&
		w.Weathers.RainIntensity > 0.4 {
		if roll := w.Rng.Intn(10); roll > 5 {
			lightningDays = 3 // lasts for 3 days
			lightningRolled = true
		}
	}

	if lightningRolled == true {
		if lightningDays == 0 {
			lightningRolled = false
			lightningDays = 0
			return
		}
		// gradiently modify the colors
		v := 255 - int(float64(3-lightningDays)/2*(255-64))
		c := fmt.Sprintf("\033[38;2;%d;%d;%dm", v, v, v)
		for y := range g {
			for x := range g[y] {
				g[y][x].Color = c
			}
		}
		lightningDays--
	}

}

func (w *World) drawLakeFlow() {
	const (
		wavelength = 8.0 // cells per ripple
		period     = 6.0 // ticks per ripple
		acc        = 1.0 // faster at the mouth
	)
	for i, row := range w.Lakes {
		for j, _ := range row {
			l := w.Lakes[i][j]
			if l == nil || l.Height == 0.0 {
				continue
			}
			speed := 1.0
			phase := float64(l.Position)/wavelength + float64(w.Day)*speed/period
			l.Phase = 2 * math.Pi * phase
		}
	}
}
