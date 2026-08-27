package main

import "math/rand"

func randomWorld(seed int64, width, height, sheepCount int) *World {
	rng := rand.New(rand.NewSource(seed))
	w := &World{
		Width:  width,
		Height: height,
		Day:    0,
		rng:    rng,
	}
	for i := 0; i < sheepCount; i++ {
		w.Sheep = append(w.Sheep, Sheep{X: rng.Intn(width),
			Y: rng.Intn(height)})
	}
	return w
}
