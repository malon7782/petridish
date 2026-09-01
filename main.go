package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// ----------------------
// 双鬼拍门，petridish最重要的?个数据结构

type Entity interface {
	Simulate(w *World)
	Pos() (x, y int)
	Icon() byte
	Layer() int
	Color() string
}

// Needswork:
// type Terrain interface...

type World struct {
	Width  int
	Height int
	Day    int

	// for living spieces
	Entities []Entity
	// for (semi-)stationary elements of the world
	Map   [][]*Mountain
	Lakes [][]*Lake
	Grass [][]bool

	Logger *Logger
	Rng    *rand.Rand
}

// -----------------------

func main() {
	fmt.Print("\033[2J")
	fmt.Print("\033[?25l")

	seed := flag.Int64("seed", time.Now().UnixNano(), "Seed of the world")
	maxDays := flag.Int("days", 100, "Number of days")
	tickMs := flag.Int("tick", 200, "Rate of time")
	flag.Parse()

	fmt.Print("\033[2J")
	w := randomWorld(*seed, 40, 20)

	for w.Day < *maxDays {
		w.renderWorld()
		time.Sleep(time.Duration(*tickMs) * time.Millisecond)
		w.simulateWorld()
	}

	w.renderWorld()
	fmt.Println("Done.")
}
