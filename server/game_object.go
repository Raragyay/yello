package main

// import (
// 	"bufio"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"sync"
// 	"time"
// )

// var (
// 	configFile = flag.String("config-file", "config.json", "/")
// 	mazeFile   = flag.String("maze-file", "maze01.txt", "/")
// )

// type player struct {
// 	row      int
// 	col      int
// 	startRow int
// 	startCol int
// }

// type ghost struct {
// 	position player
// 	status   GhostStatus
// }

// type GhostStatus string

// const (
// 	GhostStatusNormal GhostStatus = "Normal"
// 	GhostStatusBlue   GhostStatus = "Blue"
// )

// var ghostsStatusMx sync.RWMutex

// type config struct {
// 	Player           string        `json:"player"`
// 	Ghost            string        `json:"ghost"`
// 	GhostBlue        string        `json:"ghost_blue"`
// 	Wall             string        `json:"wall"`
// 	Dot              string        `json:"dot"`
// 	Pill             string        `json:"pill"`
// 	Death            string        `json:"death"`
// 	Space            string        `json:"space"`
// 	UseEmoji         bool          `json:"use_emoji"`
// 	PillDurationSecs time.Duration `json:"pill_duration_secs"`
// }

// func loadConfig(file string) error {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	decoder := json.NewDecoder(f)
// 	err = decoder.Decode(&cfg)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func loadMaze(file string) error {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		maze = append(maze, line)
// 	}

// 	for row, line := range maze {
// 		for col, char := range line {
// 			switch char {
// 			case 'P':
// 				player = player{row, col, row, col}
// 			case 'G':
// 				ghosts = append(ghosts, &ghost{player{row, col, row, col}, GhostStatusNormal})
// 			case '.':
// 				numDots++
// 			}
// 		}
// 	}

// 	return nil
// }

// func readInput() (string, error) {
// 	buffer := make([]byte, 100)

// 	cnt, err := os.Stdin.Read(buffer)
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

// func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
// 	newRow, newCol = oldRow, oldCol

// 	switch dir {
// 	case "UP":
// 		newRow = newRow - 1
// 		if newRow < 0 {
// 			newRow = len(maze) - 1
// 		}
// 	case "DOWN":
// 		newRow = newRow + 1
// 		if newRow == len(maze)-1 {
// 			newRow = 0
// 		}
// 	case "RIGHT":
// 		newCol = newCol + 1
// 		if newCol == len(maze[0]) {
// 			newCol = 0
// 		}
// 	case "LEFT":
// 		newCol = newCol - 1
// 		if newCol < 0 {
// 			newCol = len(maze[0]) - 1
// 		}
// 	}

// 	if maze[newRow][newCol] == '#' {
// 		newRow = oldRow
// 		newCol = oldCol
// 	}

// 	return
// }

// func movePlayer(dir string) {
// 	player.row, player.col = makeMove(player.row, player.col, dir)

// 	removeDot := func(row, col int) {
// 		maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
// 	}

// 	switch maze[player.row][player.col] {
// 	case '.':
// 		numDots--
// 		score++
// 		removeDot(player.row, player.col)
// 	case 'X':
// 		score += 10
// 		removeDot(player.row, player.col)
// 		go processPill()
// 	}
// }

// func drawDirection() string {
// 	dir := rand.Intn(4)
// 	move := map[int]string{
// 		0: "UP",
// 		1: "DOWN",
// 		2: "RIGHT",
// 		3: "LEFT",
// 	}
// 	return move[dir]
// }

// func test_main() {
// 	flag.Parse()

// 	// initialize game
// 	initialise()
// 	defer cleanup()

// 	// load resources
// 	err := loadMaze(*mazeFile)
// 	if err != nil {
// 		log.Println("failed to load maze:", err)
// 		return
// 	}

// 	// process input (async)
// 	input := make(chan string)
// 	go func(ch chan<- string) {
// 		for {
// 			input, err := readInput()
// 			if err != nil {
// 				log.Print("error reading input:", err)
// 				ch <- "ESC"
// 			}
// 			ch <- input
// 		}
// 	}(input)

// 	// game loop
// 	for {
// 		// process movement
// 		select {
// 		case inp := <-input:
// 			if inp == "ESC" {
// 				lives = 0
// 			}
// 			movePlayer(inp)

// 			moveGhosts()
// 			// check game over
// 			if numDots == 0 || lives <= 0 {
// 				if lives == 0 {
// 					fmt.Print("GAME OVER")
// 				}
// 				break
// 			}

// 			// repeat
// 			time.Sleep(200 * time.Millisecond)
// 		}
// 	}
// }
