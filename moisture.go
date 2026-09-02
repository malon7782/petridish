package main

func generateMoisture(w *World) {
	if w.Moisture == nil {
		w.Moisture = make([][]float64, w.Height)
		for y := 0; y < w.Height; y++ {
			w.Moisture[y] = make([]float64, w.Width)
		}
	}
	w.updateMoisture()
}

// this function should only modify w.Moisture. never let it modify w.Lakes.

func (w *World) updateMoisture() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Lakes[y][x].Height > 0.0 {
				w.Moisture[y][x] = 30.0
			}

		}
	}

	// the diffusion algo is ai-generated.

	nextMoisture := make([][]float64, w.Height)
	for y := range nextMoisture {
		nextMoisture[y] = make([]float64, w.Width)
	}

	diffusionRate := 0.4

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			current := w.Moisture[y][x]
			neighborSum := 0.0
			neighborCount := 0.0

			dirs := []struct{ dx, dy int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
			for _, d := range dirs {
				ny, nx := y+d.dy, x+d.dx
				if ny >= 0 && ny < w.Height && nx >= 0 && nx < w.Width {
					neighborSum += w.Moisture[ny][nx]
					neighborCount++
				}
			}

			avgNeighbor := neighborSum / neighborCount
			nextMoisture[y][x] = current*(1.0-diffusionRate) + avgNeighbor*diffusionRate
		}
	}

	w.Moisture = nextMoisture
}
