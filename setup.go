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
	generateTerrain(w, 4)
	// Sheeps
	generateSheep(w, 9)

	return w
}
