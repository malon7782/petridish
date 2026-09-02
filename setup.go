package main

import "math/rand"

func randomWorld(seed int64, width, height int) *World {
	rng := rand.New(rand.NewSource(seed))
	w := &World{
		Width:  width,
		Height: height,
		Logger: &Logger{
			Num: 10},
		Day: 0,
		Rng: rng,
	}
	// Mountains
	generateMountain(w, 7)
	// Lakes
	// cumLake(w, 30)
	generateBFSLake(w, 9)
	generateBFSLake(w, 9)

	// Moisture
	generateMoisture(w)

	// Sheeps
	generateSheep(w, 9)
	return w
}

// -----------
// helper functions
// Needswork: together with main.go, some stuffs associated with
// struct World need to be migrated to a seperate file like world.go.

func (w *World) getHeight(x, y int) int {
	if x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return 0
	}
	return w.Map[x][y].Height
}
