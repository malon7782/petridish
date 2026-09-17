package main

import "math"

type Liquid interface {
	Level() float64
	SetLevel(v float64)

	Frozen() bool
	SetFrozen(v bool)

	// properties
	FreezePoint() float64
	EvapRate() float64 // how fast the liquid dries out

	Source() bool
}

func handleLiquidEvaporation[L Liquid](grid [][]L, w *World) {
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

func handleLiquidFlow[L interface{Liquid; comparable}](grid [][]L, w *World) {
	var q []pair
	var NIL L

	for y := range w.Height {
		for x := range w.Width {
			l := grid[y][x]
			if l != NIL && l.Level() > 0.0 {
				q = append(q, pair{y: y, x: x})
			}
		}
	}

	type Type struct {
		diff float64
		coor pair
	}

	delta := make([][]float64, w.Height)
	for i := 0; i < w.Height; i++ {
		delta[i] = make([]float64, w.Width)
	}
	for _, cur := range q {
		l := grid[cur.y][cur.x]
		var nb []Type
		for _, d := range dirs4 {
			nxt := pair{y: cur.y + d[0], x: cur.x + d[1]}
			if !w.inMap(nxt.y, nxt.x) {
				continue
			}

			l2 := grid[nxt.y][nxt.x]
			curLvl := w.Map[cur.y][cur.x].Height + l.Level()
			nxtLvl := w.Map[nxt.y][nxt.x].Height + l2.Level()
			if curLvl > nxtLvl {
				nb = append(nb, Type{diff: curLvl - nxtLvl, coor: nxt})
			}
		}

		m := len(nb)
		if m == 0 {
			continue
		}
		mn := math.Inf(1)
		s := 0.0
		for _, n := range nb {
			mn = min(mn, n.diff)
			s += n.diff
		}
		out := min(l.Level(), mn/2)
		for _, n := range nb {
			d := out * n.diff / s
			delta[cur.y][cur.x] -= d
			delta[n.coor.y][n.coor.x] += d
		}
	}

	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			grid[i][j].SetLevel(grid[i][j].Level() + delta[i][j])
		}
	}

	// finally, to avoid flooding, water should go out of the map, into the occean
	for i := 0; i < w.Height; i++ {
		d := 1
		if 0 < i && i+1 < w.Height {
			d = w.Width - 1
		}
		for j := 0; j < w.Width; j += d {
			l := grid[i][j]
			if !l.Source() {
				l.SetLevel(l.Level() * 0.2)
			}
		}
	}

}
