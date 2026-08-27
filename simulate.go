package main

// we don't wanna create a complete world structure every tick,
// so let's modify the world in-place.

func (w *World) simulateWorld() {
	w.handleSheep()
	w.Day += 1
}
