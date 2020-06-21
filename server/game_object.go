package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	configFile = flag.String("config-file", "config.json", "/")
	mazeFile   = flag.String("maze-file", "maze01.txt", "/")
)

type player struct {
	name	 string
	row      int
	col      int
	startRow int
	startCol int
	ghost    bool
}

type ghost struct {
	position player
	status   GhostStatus
}

type GhostStatus string

const (
	GhostStatusNormal GhostStatus = "Normal"
	GhostStatusBlue   GhostStatus = "Blue"
)

var players [5]player

func createplayer(name string){
	if name == "pacman"{
		players[0] = player(PLAYERNAME, STARTROW, STARTCOL, STARTROW, STARTCOL,false)
	}
	else if player[1] == nil{
		player[1] = player(PLAYERNAME, STARTROW, STARTCOL, STARTROW, STARTCOL,true)
	}
	else if player[2] == nil{
		player[2] = player(PLAYERNAME, STARTROW, STARTCOL, STARTROW, STARTCOL,true)
	}
	else if player[3] == nil{
		player[3] = player(PLAYERNAME, STARTROW, STARTCOL, STARTROW, STARTCOL,true)
	}
	else if player[4] == nil{
		player[4] = player(PLAYERNAME, STARTROW, STARTCOL, STARTROW, STARTCOL,true)
	}
	else if player[5] == nil{
		player[5] = player(PLAYERNAME, STARTROW, STARTCOL, STARTROW, STARTCOL,true)
	}
}

var ghostsStatusMx sync.RWMutex

type config struct {
	Player    string `json:"player"`
	Ghost     string `json:"ghost"`
	GhostBlue string `json:"ghost_blue"`
	Wall      string `json:"wall"`
	Dot       string `json:"dot"`
	Death     string `json:"death"`
	Space     string `json:"space"`
}

func loadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}
	return nil
}

func loadMaze(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}

	for row, line := range maze {
		for col, char := range line {
			switch char {
			case 'P':
				player = player{row, col, row, col}
			case 'G':
				ghosts = append(ghosts, &ghost{player{row, col, row, col}, GhostStatusNormal})
			case '.':
				numDots++
			}
		}
	}

	return nil
}

func readInput() (string, error) {
	buffer := make([]byte, 100)

	cnt, err := //WHATEVER THE SERVER SENDS (UP DOWN LEFT RIGHT)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}
	return "", nil
}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze)-1 {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	}

	if maze[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}
	return
}

func movePlayer(playerindex string, dir string) {
	player.row, player.col = makeMove(players[playerindex].player.row, players[playerindex].player.col, dir)
	if players[playerindex].player.ghost == false{
		removeDot := func(row, col int) {
			maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
		}

		switch maze[player.row][player.col] {
		case '.':
			numDots--
			score++
			removeDot(player.row, player.col)
	}
	}
	else{
		
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

func main() {
	flag.Parse()

	// initialize game
	initialise()
	defer cleanup()

	// load resources
	err := loadMaze(*mazeFile)
	if err != nil {
		log.Println("failed to load maze:", err)
		return
	}

	// process input (async)
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := readInput()
			if err != nil {
				log.Print("error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)

	// game loop
	for {
		// process movement
		select {
		case inp := <-input:
			if inp == "ESC" {
				lives = 0
			}
			movePlayer(inp)
			moveGhosts(inp)
			// check game over
			if numDots == 0 || lives <= 0 {
				if lives == 0 {
					fmt.Print("GAME OVER")
				}
				break
			}
		}
	}
}
