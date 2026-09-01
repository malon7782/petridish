package main

const (
	GrassColor = "\033[38;5;70m"
	GrassIcon  = '#'
)

func generateGrass(w *World) {
	w.Grass = make([][]bool, w.Height)
	for y := 0; y < w.Height; y++ {
		w.Grass[y] = make([]bool, w.Width)
	}
	// Needswork: how are we going to put grass blocks
	// on the map in the setup phase?
}
