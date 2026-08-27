package main

import "math/rand"

func randomWorld(seed int64, width, height int) *World {
	rng := rand.New(rand.NewSource(seed))
	w := &World{
		Width:  width,
		Height: height,
		Day:    0,
		rng:    rng,
	}
	// Sheeps
	for i := 0; i < rng.Intn(9); i++ {
		w.Entities = append(w.Entities, &Sheep{X: rng.Intn(width),
			Y: rng.Intn(height)})
	}
	return w
}
