package main

/*
 * the functions in simulateWorld() are called every day. they are supposed to involve the real-time
 * distribution of game elements, such as the flow of rivers, the weathering of rocks, and animal fights.
 * but in principle, they should not determine how these elements are displayed on the screen - at least
 * not directly.
 */

func (w *World) simulateWorld() {
	for _, e := range w.Entities {
		e.Simulate(w)
	}
	// those entities whose 'IsAlive == false' will be kicked out from the list!!
	for i, e := range w.Entities {
		if e.IsAlive() == false {
			w.Entities = append(w.Entities[:i], w.Entities[i+1:]...)
		}
	}

	// w(t) -> w(t + 1)
	w.simulateWeather()
	w.updateMoisture()
	w.simulateGrass()

	w.updateMaxHeightMap()

	w.simulateLake()

	w.Day += 1
}
