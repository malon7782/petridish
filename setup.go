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
	// NOTE: generateOOO() means to initialize, but not to simulate!

	// Mountains
	generateMountain(w, 7)

	// Lakes
	generateLake(w) // actually, to initialize w.Lakes
	generateRiver(w)

	generateBFSLake(w, 9)
	generateBFSLake(w, 9)

	// Weather
	generateWeather(w)

	//Grass
	generateGrass(w)

	// Moisture
	generateMoisture(w)

	// Sheeps
	generateSheep(w, 2)
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
