package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type World struct {
	Width  int
	Height int
	Day    int
	Sheep  []Sheep
	rng    *rand.Rand
}

func main() {
	seed := flag.Int64("seed", time.Now().UnixNano(), "Seed of the world")
	maxDays := flag.Int("days", 100, "Number of days")
	tickMs := flag.Int("tick", 200, "Rate of time")
	flag.Parse()

	fmt.Print("\033[2J")
	w := randomWorld(*seed, 20, 10, 3)

	for w.Day < *maxDays {
		w.renderWorld()
		time.Sleep(time.Duration(*tickMs) * time.Millisecond)
		w.simulateWorld()
	}

	w.renderWorld()
	fmt.Println("Done.")
}
