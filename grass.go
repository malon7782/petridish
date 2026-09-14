package main

type Plant interface {
	IsAlive() bool
	MoistureContent() float64 // should be less than 100
}

type Grass struct {
	Alive        bool
	WaterContent float64
}

func (g *Grass) MoistureContent() float64 { return 60.0 }
func (g *Grass) IsAlive() bool            { return g.Alive }

const (
	GrassColor            = "\033[38;5;70m"
	GrassIcon             = '#'
	MinLakeHeightForGrass = 0.2
)

func generateGrass(w *World) {
	w.GrassMap = make([][]*Grass, w.Height)
	for y := 0; y < w.Height; y++ {
		w.GrassMap[y] = make([]*Grass, w.Width)
		for x := 0; x < w.Width; x++ {
			w.GrassMap[y][x] = &Grass{Alive: false, WaterContent: 60.0}
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
			if w.Lakes[y][x].Height < MinLakeHeightForGrass {
				h := w.Map[y][x].Height
				if h < 0 {
					h = 0
				}

				if w.Lakes[y][x].Height >= 0.05 {
					w.GrassMap[y][x].Alive = false
					continue
				}

				if w.Moisture[y][x] < 10.0 || w.Weathers.Temperature > 35 ||
					w.Weathers.Temperature < 10 {
					// Needswork: the model now is rather simple. as long as the humidity is too low
					// or the temperature is not habitable, the grass start to die, but by a hardcoded chance
					if roll := w.Rng.Intn(100); roll > 90 {
						w.GrassMap[y][x].Alive = false
					}
				} else {
					roll := w.Rng.Float64()

					// influenced by both humidity and height
					p := (0.002 * w.Moisture[y][x]) / float64((h+1)*(h+1))

					if p > roll {
						w.GrassMap[y][x].Alive = true
					}
				}
			}
		}
	}
}
