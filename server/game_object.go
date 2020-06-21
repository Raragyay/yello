package main

import (
	"bufio"
	"encoding/json"
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
wall tile = 0b10000000
p1 tile = 0b01000000
p2 tile = 0b00100000
p3 tile = 0b00010000
p4 tile = 0b00001000
p5 tile = 0b00000100
pellet tile = 0b00000010
super-pellet tile = 0b00000001
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
	//TODO parse into enum
}

// func readInput() (string, error) { //make so that
// 	buffer := make([]byte, 100)

// 	cnt, err := //WHATEVER THE SERVER SENDS (UP DOWN LEFT RIGHT)
// 	if err != nil {
// 		return "", err
// 	}

// 	if cnt == 1 && buffer[0] == 0x1b {
// 		return "ESC", nil
// 	} else if cnt >= 3 {
// 		if buffer[0] == 0x1b && buffer[1] == '[' {
// 			switch buffer[2] {
// 			case 'A':
// 				return "UP", nil
// 			case 'B':
// 				return "DOWN", nil
// 			case 'C':
// 				return "RIGHT", nil
// 			case 'D':
// 				return "LEFT", nil
// 			}
// 		}
// 	}
// 	return "", nil
// }
//
//func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
//	newRow, newCol = oldRow, oldCol
//
//	switch dir {
//	case "UP":
//		newRow = newRow - 1
//		if newRow < 0 {
//			newRow = len(maze) - 1
//		}
//	case "DOWN":
//		newRow = newRow + 1
//		if newRow == len(maze)-1 {
//			newRow = 0
//		}
//	case "RIGHT":
//		newCol = newCol + 1
//		if newCol == len(maze[0]) {
//			newCol = 0
//		}
//	case "LEFT":
//		newCol = newCol - 1
//		if newCol < 0 {
//			newCol = len(maze[0]) - 1
//		}
//	}
//
//	if maze[newRow][newCol] == wall {
//		newRow = oldRow
//		newCol = oldCol
//	}
//	return
//}
//
//func movePlayer(playerindex string, dir string) {
//	player.row, player.col = makeMove(players[playerindex].player.row, players[playerindex].player.col, dir)
//	if players[playerindex].player.ghost == false {
//		removeDot := func(row, col int) {
//			maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
//		}
//
//		switch maze[player.row][player.col] {
//		case '.':
//			numDots--
//			score++
//			removeDot(player.row, player.col)
//		}
//	} else {
//
//	}
//
//}
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

func yes() {
	//flag.Parse()
	//
	//// initialize game
	//initialise()
	//defer cleanup()

	// load resources
	// maze := loadAndParseMazeFile(*mazeFile)
	// println(len(maze))
	//// process input (async)
	//input := make(chan string)
	//go func(ch chan<- string) {
	//	for {
	//		input, err := readInput()
	//		if err != nil {
	//			log.Print("error reading input:", err)
	//			ch <- "ESC"
	//		}
	//		ch <- input
	//	}
	//}(input)
	//
	//// game loop
	//for {
	//	// process movement
	//	select {
	//	case inp := <-input:
	//		if inp == "ESC" {
	//			lives = 0
	//		}
	//		movePlayer(inp)
	//		moveGhosts(inp)
	//		// check game over
	//		if numDots == 0 || lives <= 0 {
	//			if lives == 0 {
	//				fmt.Print("GAME OVER")
	//			}
	//			break
	//		}
	//	}
	//}
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
