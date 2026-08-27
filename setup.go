package main

import "math/rand"

func randomWorld(seed int64, width, height int) *World {
	rng := rand.New(rand.NewSource(seed))
	w := &World{
		Width:  width,
		Height: height,
		Day:    0,
		Rng:    rng,
	}
	// Mountains
	generateMountain(w, 4)
	// Sheeps
	generateSheep(w, 9)

	return w
}

// -----------
// helper functions

func (w *World) getHeight(x, y int) int {
	if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return 0
	}
	return w.Map[x][y].Height
}
