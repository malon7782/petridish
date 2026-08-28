package main

func (w *World) simulateWorld() {
	for _, e := range w.Entities {
		e.Simulate(w)
	}
	cumLake(w, 20)
	w.Day += 1
	// note that we don't simulate topographical elements. this part will be added
	// in future updates!
}
