package main

import "fmt"

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
	if lightningRolled == false && w.Weathers.RainLeft >= 4 &&
		w.Weathers.RainIntensity > 0.4 {
		if roll := w.Rng.Intn(10); roll > 5 {
			lightningDays = 3
			lightningRolled = true
		}
	}

	if lightningRolled == true {
		if lightningDays == 0 {
			lightningRolled = false
			lightningDays = 0
			return
		}
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
