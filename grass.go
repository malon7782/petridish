package main

const (
	GrassColor            = "\033[38;5;70m"
	GrassIcon             = '#'
	MinLakeHeightForGrass = 0.2
)

func generateGrass(w *World) {
	w.Grass = make([][]bool, w.Height)
	for y := 0; y < w.Height; y++ {
		w.Grass[y] = make([]bool, w.Width)
	}
	// Needswork: how are we going to put grass blocks
	// on the map in the setup phase?
}

func (w *World) simulateGrass() {
	for y := range w.Height {
		for x := range w.Width {
			if w.Lakes[y][x].Height < MinLakeHeightForGrass {
				h := w.Map[y][x].Height
				if h < 0 {
					h = 0
				}

				if w.Moisture[y][x] < 10.0 || w.Weathers.Temperature > 35 || w.Weathers.Temperature < 10 {
					// Needswork: the model now is rather simple. as long as the humidity is too low
					// or the temperature is not habitable, the grass start to die, but by a hardcoded chance
					if roll := w.Rng.Intn(100); roll > 90 {
						w.Grass[y][x] = false
					}
				} else {
					roll := float64(w.Rng.Intn(9)) + 1

					// influenced by both humidity and height
					p := (0.3 * w.Moisture[y][x]) / float64((h+1)*(h+1))

					if p > roll {
						w.Grass[y][x] = true
					}
				}
			}
		}
	}
}
