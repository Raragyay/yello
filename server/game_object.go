package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type tile uint8

type posVector struct {
	x int
	y int
}

func (v posVector) toString() string {
	return strconv.Itoa(v.x) + "-" + strconv.Itoa(v.y)
}

const (
	// Bitwise enums
	wall        tile = 0b10000000
	p1          tile = 0b01000000
	p2          tile = 0b00100000
	p3          tile = 0b00010000
	p4          tile = 0b00001000
	p5          tile = 0b00000100
	superPellet tile = 0b00000010
	pellet      tile = 0b00000001
	empty       tile = 0b00000000
)

var (
	configFile = flag.String("config-file", "config.json", "/")
	mazeFile   = flag.String("maze-file", "map.txt", "/")
)

var ghostsStatusMx sync.RWMutex

func loadMaze(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	maze := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}
	return maze, nil
}

func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func loadAndParseMazeFile(mazeFileName string) ([][]string, int) {
	rawMaze, err := loadMaze(mazeFileName)
	if err != nil {
		log.Println("failed to load maze:", err)
		return nil, 0
	}
	parsedMaze := make([][]string, len(rawMaze))
	rows := len(rawMaze)
	for i := 0; i < rows; i++ {
		splitRow := strings.Split(rawMaze[i], " ")
		parsedMaze[i] = make([]string, len(splitRow))
		for j := 0; j < len(splitRow); j++ {
			parsedMaze[i][j] = string(splitRow[j])
		}
	}
	return parsedMaze, rows
}

//UTILS

func (t tile) isIllegalCollision() bool {
	return (uint8(t) > uint8(wall))
}

func (t tile) getPlayerIDs() []string {
	returned := make([]string, 0)
	if t&p1 != 0 {
		returned = append(returned, "P1")
	}
	if t&p2 != 0 {
		returned = append(returned, "P2")
	}
	if t&p3 != 0 {
		returned = append(returned, "P3")
	}
	if t&p4 != 0 {
		returned = append(returned, "P4")
	}
	if t&p5 != 0 {
		returned = append(returned, "P5")
	}
	return returned
}

func playerIDToTile(id string) tile {
	switch id {
	case "P1":
		return p1
	case "P2":
		return p2
	case "P3":
		return p3
	case "P4":
		return p4
	case "P5":
		return p5
	}
	return empty
}

func (t tile) isPacHom() bool {
	return pellet+p1 == t || superPellet|p1 == t
}

func (t tile) isPacGhostCollide() bool {
	return (uint8(t) >= uint8(p1+p5)) && (uint8(t) < uint8(wall))
}

func tileToBit(sMaze string) tile {
	switch sMaze {
	case "000":
		return empty
	case "002":
		return pellet
	case "003":
		return superPellet
	case "004":
		return p1
	case "010":
		return p2
	case "011":
		return p3
	case "012":
		return p4
	case "013":
		return p5
	default:
		return wall
	}
}
