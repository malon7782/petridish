package main

type Liquid interface {
	Level() float64
	SetLevel(v float64)

	Frozen() bool
	SetFrozen(v bool)

	// properties
	FreezePoint() float64
	EvapRate() float64 // how fast the liquid dries out
}

func handleLiquidEvaporation(grid [][]Liquid, w *World) {
	for y := range grid {
		for x := range grid[y] {
			l := grid[y][x]
			if w.Weathers.Temperature > l.FreezePoint() {
				if roll := w.Rng.Intn(10); roll > 7 {
					l.SetFrozen(false)
				}

				lvl := l.Level() - w.Weathers.Temperature*l.EvapRate()
				if lvl < 0 {
					lvl = 0
				}

				l.SetLevel(lvl)
			} else if roll := w.Rng.Intn(10); roll > 7 {
				l.SetFrozen(true)
			}
		}
	}
}
