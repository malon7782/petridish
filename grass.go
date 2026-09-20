package main

type Plant interface {
	IsAlive() bool
	MoistureContent() float64 // should be less than 100
}

type Grass struct {
	Alive        bool
	WaterContent float64
}

func (g *Grass) MoistureContent() float64 { return g.WaterContent }
func (g *Grass) IsAlive() bool            { return g.Alive }

func (g *Grass) Icon() byte {
	switch {
	case g.WaterContent < 10:
		return '\''
	case g.WaterContent < 25:
		return '*'
	case g.WaterContent < 40:
		return '*'
	default:
		return '#'
	}
}

func (g *Grass) Color() string {
	switch {
	case g.WaterContent < 15:
		return "\033[38;5;102m" // 灰绿
	case g.WaterContent < 25:
		return "\033[38;5;184m" // 枯黄
	case g.WaterContent < 40:
		return "\033[38;5;34m" // 草绿
	default:
		return "\033[38;5;76m" // 鲜绿
	}
}

func generateGrass(w *World) {
	w.GrassMap = make([][]*Grass, w.Height)
	for y := 0; y < w.Height; y++ {
		w.GrassMap[y] = make([]*Grass, w.Width)
		for x := 0; x < w.Width; x++ {
			w.GrassMap[y][x] = &Grass{Alive: false, WaterContent: 0.0}
		}
	}
}

func (w *World) simulateGrass() {
	for y := range w.Height {
		if y == 0 || y == w.Height-1 {
			continue
		}
		for x := range w.Width {
			if x == 0 || x == w.Width-1 {
				continue
			}
			g := w.GrassMap[y][x]

			if !g.Alive && w.Lakes[y][x].Height <= 0.1 {
				h := w.Map[y][x].Height
				if h < 0 {
					h = 0
				}
				roll := w.Rng.Float64()

				// influenced by both humidity and height
				p := (0.0005 * w.Moisture[y][x]) / float64((h+1)*(h+1))

				if p > roll {
					g.Alive = true
					g.WaterContent = 30.0
				}
			}

			if g.Alive {
				// target WaterContent
				target := 20.0 + 0.5*w.Moisture[y][x] + w.Lakes[y][x].Height*150.0
				if w.Weathers.Temperature < 10 {
					target -= (10 - w.Weathers.Temperature) * 2.0
				}
				if w.Weathers.Temperature > 35 {
					target += (w.Weathers.Temperature - 35) * 2.0
				}

				target = min(100.0, max(0.0, target))
				g.WaterContent = min(100.0, max(0.0, g.WaterContent+(target-g.WaterContent)*0.5))

				// drown
				if w.Lakes[y][x].Height >= 0.2 {
					g.Alive = false
					continue
				}

				// too wet or too dry also dies
				if g.WaterContent < 10 || g.WaterContent > 80 {
					// Needswork: the model now is still simple. as long as the humidity is too low
					// or the temperature is not habitable (baked into the WaterContent target above),
					// the grass dies once its water content crosses too low or too high
					g.Alive = false
					continue
				}
			}

		}
	}
}
