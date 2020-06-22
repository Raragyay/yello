package main

import (
	"fmt"
	"strconv"
	"strings"
)

type direction string

const (
	startingPacLives int = 3

	left  direction = "L"
	right direction = "R"
	up    direction = "U"
	down  direction = "D"
)

type playerGameData struct { //TODO LOAD GAME DATA
	p               *clientPlayer
	latestDirection direction
}

func initializeGameServer(p1, p2 *clientPlayer) {
	gameInstance := &game{
		p1: &playerGameData{p: p1},
		p2: &playerGameData{p: p2},
		//p3:     &playerGameData{p: p3},
		//p4:     &playerGameData{p: p4},
		//p5:     &playerGameData{p: p5},
		active: true,
	}

	//handle maze
	maze, numRows := loadAndParseMazeFile(*mazeFile)

	if maze == nil {
		writeToAllPlayers(gameInstance, "PONG GAME-INVALID")
		disconnectAllPlayers(gameInstance)
		return //terminate game early!
	}

	var mazeMsg strings.Builder
	fmt.Fprintf(&mazeMsg, "%s-%s", strconv.Itoa(len(maze[0])), strconv.Itoa(numRows))
	for _, val := range maze {
		for _, el := range val {
			fmt.Fprintf(&mazeMsg, "-%s", string(el))
		}
	}

	//now construct bit maze for us yes
	gameInstance.maze = constructBitMaze(maze)

	gameInstance.pacLives = startingPacLives

	//handle game initialization
	p1.writeChanneledMessage("PONG GAME-INIT " + "P1-" + p2.name) //Who needs JSON when you got -?
	p2.writeChanneledMessage("PONG GAME-INIT " + "P2-" + p1.name)

	//p1.writeChanneledMessage("PONG GAME-INIT " + "P1-" + p2.name + "-" + p3.name + "-" + p4.name + "-" + p5.name) //Who needs JSON when you got -?
	//p2.writeChanneledMessage("PONG GAME-INIT " + "P2-" + p1.name + "-" + p3.name + "-" + p4.name + "-" + p5.name)
	//p3.writeChanneledMessage("PONG GAME-INIT " + "P3-" + p1.name + "-" + p2.name + "-" + p4.name + "-" + p5.name)
	//p4.writeChanneledMessage("PONG GAME-INIT " + "P4-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p5.name)
	//p5.writeChanneledMessage("PONG GAME-INIT " + "P5-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name)

	writeToAllPlayers(gameInstance, "PONG SET-LEVEL "+mazeMsg.String())

	p1.activeGame = gameInstance
	p2.activeGame = gameInstance
	//p3.activeGame = &gameInstance
	//p4.activeGame = &gameInstance
	//p5.activeGame = &gameInstance

	tendGame(gameInstance)
}

func tendGame(g *game) {

	//WRITE CODE HERE

	for g.active {
		//ze game loop!
	}
}

//INNER PROCESSES

func checkUpdateObjectStates() {

}

//TO ZE CLIENTSSSS

func updateObjectPosition(g *game, id string, v *posVector) {
	writeToAllPlayers(g, "PONG GAME-ENTITY-POS "+id+"-"+v.toString())
}

func sendPelletConsumed(g *game, v *posVector) {
	writeToAllPlayers(g, "PONG GAME-PELLET-HOM "+v.toString()) //client should get rid of pellet AND increase score for pac-man (if it is keeping track for visual purposes)
}

func sendScared(g *game) {
	writeToAllPlayers(g, "PONG GAME-SCARED 1")
}

func sendCeaseScared(g *game) {
	writeToAllPlayers(g, "PONG GAME-SCARED 0")
}

//func updateObjectState(g *game, o *gameObject, state objectState) {
//	writeToAllPlayers(g, "PONG GAME-OBJECT-STATE " + o.string_ID + "-" + string(objectState))
//}

//FROM ZE CLIENT

func playerUpdateDesiredDirection(req *playerRequest, argument string) {
	switch argument {
	case "R":
		break
	case "L":
		break
	case "U":
		break
	case "D":
		break
	}
}

//FUNCTIONALITY

func moveEntity()

//INNER UTILS

func writeToAllPlayers(g *game, msg string) {
	g.p1.p.writeChanneledMessage(msg)
	g.p2.p.writeChanneledMessage(msg)
}

func disconnectAllPlayers(g *game) {
	g.p1.p.disconnectChannel <- struct{}{}
	g.p2.p.disconnectChannel <- struct{}{}
}

func getNumberOfPellets(g *game) int {
	return 0
}

func constructBitMaze(sMaze [][]string) [][]tile {
	for (int row = 0; row < sMaze.length; row++) {
		for (int col = 0; col < sMaze[row].length; col++) {
			tileMaze[row][col] := (tileToBit(sMaze[row][col]))
		}
	}
	return(tileMaze)
}
