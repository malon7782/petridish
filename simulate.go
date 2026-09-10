package main

func (w *World) simulateWorld() {
	for _, e := range w.Entities {
		e.Simulate(w)
	}

	// w(t) -> w(t + 1)
	w.simulateWeather()
	w.updateMoisture()
	w.simulateGrass()
	w.simulateLake()

	w.Day += 1
}
