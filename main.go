package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type Sheep struct {
	X, Y int
}

type World struct {
	Width  int
	Height int
	Day    int
	Sheep  []Sheep
	rng    *rand.Rand
}

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

// we don't wanna create a complete world structure every tick,
// so let's modify the world in-place.

func (w *World) simulateWorld() {
	for i := range w.Sheep {
		dx := w.rng.Intn(3) - 1
		dy := w.rng.Intn(3) - 1
		newX := w.Sheep[i].X + dx
		newY := w.Sheep[i].Y + dy

		if newX >= 0 && newX < w.Width {
			w.Sheep[i].X = newX
		}
		if newY >= 0 && newY < w.Height {
			w.Sheep[i].Y = newY
		}
	}
	w.Day += 1
}

func (w *World) renderWorld() {
	fmt.Print("\033[H")

	grid := make([][]byte, w.Height)
	for i := range grid {
		grid[i] = bytes.Repeat([]byte{'.'}, w.Width)
	}

	for _, s := range w.Sheep {
		grid[s.Y][s.X] = 'S'
	}

	fmt.Printf("Day: %d\n", w.Day)
	for _, row := range grid {
		fmt.Println(string(row))
	}
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
