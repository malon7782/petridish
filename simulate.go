package main

// we don't wanna create a complete world structure every tick,
// so let's modify the world in-place.

func (w *World) simulateWorld() {
	for _, e := range w.Entities {
		e.Simulate(w)
	}
	w.Day += 1
}
