package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	p                  *clientPlayer
	position           *posVector
	currentDirection   direction
	desiredDirection   direction
	tileRepresentation tile
	startingPosition   *posVector
}

type game struct {
	p1, p2, p3, p4, p5 *playerGameData
	active             bool
	maze               [][]tile
	pacLives           int
	pelletsLeft        int
}

func initializeGameServer(p1, p2, p3, p4 *clientPlayer) {

	gameInstance := &game{
		p1: &playerGameData{p: p1},
		p2: &playerGameData{p: p2},
		p3: &playerGameData{p: p3},
		p4: &playerGameData{p: p4},
		//p5:     &playerGameData{p: p5},
		active:      true,
		pelletsLeft: 0,
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

	//more initialization
	gameInstance.updatePlayerPositions()
	gameInstance.updateTileReferences()
	gameInstance.updatePelletCounts()

	gameInstance.setPlayerStartingPositionsToCurrent()

	//handle game initialization
	p1.writeChanneledMessage("PONG GAME-INIT " + "P1-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name) //Who need
	p2.writeChanneledMessage("PONG GAME-INIT " + "P2-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name) //Who needs JSON
	p3.writeChanneledMessage("PONG GAME-INIT " + "P3-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name) //Who needs JSON
	p4.writeChanneledMessage("PONG GAME-INIT " + "P4-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name) //Who needs JSON

	//p1.writeChanneledMessage("PONG GAME-INIT " + "P1-" + p2.name + "-" + p3.name + "-" + p4.name + "-" + p5.name) //Who needs JSON when you got -?
	//p2.writeChanneledMessage("PONG GAME-INIT " + "P2-" + p1.name + "-" + p3.name + "-" + p4.name + "-" + p5.name)
	//p3.writeChanneledMessage("PONG GAME-INIT " + "P3-" + p1.name + "-" + p2.name + "-" + p4.name + "-" + p5.name)
	//p4.writeChanneledMessage("PONG GAME-INIT " + "P4-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p5.name)
	//p5.writeChanneledMessage("PONG GAME-INIT " + "P5-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name)

	writeToAllPlayers(gameInstance, "PONG SET-LEVEL "+mazeMsg.String())

	p1.activeGame = gameInstance
	p2.activeGame = gameInstance
	p3.activeGame = gameInstance
	p4.activeGame = gameInstance
	//p5.activeGame = &gameInstance

	tendGame(gameInstance)
}

func tendGame(g *game) {

	//WRITE CODE HERE
	scaredCountdown := 0
	for g.active {
		//move players
		for idx, player := range []*playerGameData{
			g.p1, g.p2, g.p3, g.p4} { //TODO add more players to iterate over
			projX, projY := getProjected(player, player.desiredDirection)
			if canMoveInDirection(g, projX, projY) {
				moveToTile(g, player, projX, projY)
				player.currentDirection = player.desiredDirection
			} else {
				projX, projY = getProjected(player, player.currentDirection)
				if canMoveInDirection(g, projX, projY) {
					moveToTile(g, player, projX, projY)
				} else {
					stopPlayer(player)
				}
			}
			updateObjectPosition(g, "P"+strconv.Itoa(idx+1), player.position)
		}
		if (*g).maze[g.p1.position.y][g.p1.position.x].isPacHom() {
			if (*g).maze[g.p1.position.y][g.p1.position.x]&superPellet != 0 {
				sendScared(g)
				scaredCountdown = 50
			}
			consumePellet(g, g.p1.position)
			sendPelletConsumed(g, g.p1.position)
		}
		if (*g).maze[g.p1.position.y][g.p1.position.x].isPacGhostCollide() {
			fmt.Println("Ghost collide")
			if scaredCountdown != 0 {
				fmt.Println("Kill ghosts")
				killGhosts(g, getPlayersOnTile(g, (*g).maze[g.p1.position.y][g.p1.position.x]))
			} else {
				fmt.Println("player died")
				resetPlayerPositions(g)
			}
		}
		if scaredCountdown >= 1 {
			scaredCountdown--
			if scaredCountdown == 0 {
				sendCeaseScared(g)
			}
		}

		//ze game loop!
		time.Sleep(100 * time.Millisecond)
	}
}

func consumePellet(g *game, v *posVector) {
	notMask := ^pellet &^ superPellet
	(*g).maze[v.y][v.x] &= notMask
}

func killGhosts(g *game, entityList []*playerGameData) {
	fmt.Println(entityList)
	for _, entity := range entityList {
		fmt.Printf("Entity at %s is killed", entity.position.toString())
		if entity.tileRepresentation != p1 {
			moveToTile(g, entity, entity.startingPosition.x, entity.startingPosition.y)
		}
	}
}

func resetPlayerPositions(g *game) {
	moveToTile(g, g.p1, g.p1.startingPosition.x, g.p1.startingPosition.y)
	moveToTile(g, g.p2, g.p2.startingPosition.x, g.p2.startingPosition.y)
	moveToTile(g, g.p3, g.p3.startingPosition.x, g.p3.startingPosition.y)
	moveToTile(g, g.p4, g.p4.startingPosition.x, g.p4.startingPosition.y)
}

func getProjected(player *playerGameData, direction direction) (int, int) {
	projX, projY := player.position.x, player.position.y
	if direction == right {
		projX += 1
	}
	if direction == left {
		projX -= 1
	}
	if direction == up {
		projY -= 1
	}
	if direction == down {
		projY += 1
	}
	return projX, projY
}

func canMoveInDirection(g *game, projX int, projY int) bool {
	return !(projX < 0 || projX >= len((*g).maze[0]) || projY < 0 || projY >= len(g.maze) || (*g).maze[projY][projX] == wall)
}

func moveToTile(g *game, player *playerGameData, x int, y int) {
	notMask := ^player.tileRepresentation
	(*g).maze[player.position.y][player.position.x] &= notMask
	player.position.x = x
	player.position.y = y
	(*g).maze[y][x] |= player.tileRepresentation

}

func stopPlayer(player *playerGameData) {
	(*player).desiredDirection = ""
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
	p := pickPlayerGameData(req.p.activeGame, req.p)
	switch argument {
	case "R":
		p.desiredDirection = right
		break
	case "L":
		p.desiredDirection = left
		break
	case "U":
		p.desiredDirection = up
		break
	case "D":
		p.desiredDirection = down
		break
	}
}

//FUNCTIONALITY

func moveEntity() {

}

//INNER UTILS

func writeToAllPlayers(g *game, msg string) {
	g.p1.p.writeChanneledMessage(msg)
	g.p2.p.writeChanneledMessage(msg)
	g.p3.p.writeChanneledMessage(msg)
	g.p4.p.writeChanneledMessage(msg)
}

func disconnectAllPlayers(g *game) {
	g.p1.p.disconnectChannel <- struct{}{}
	g.p2.p.disconnectChannel <- struct{}{}
	g.p3.p.disconnectChannel <- struct{}{}
	g.p4.p.disconnectChannel <- struct{}{}
}

func getNumberOfPellets(g *game) int {
	return 0
}

func constructBitMaze(sMaze [][]string) [][]tile {
	tileMaze := make([][]tile, len(sMaze)) //TODO BETTER LEN HANDLING
	for row := 0; row < len(sMaze); row++ {
		tileMaze[row] = make([]tile, len(sMaze[row]))
		for col := 0; col < len(sMaze[row]); col++ {
			tileMaze[row][col] = tileToBit(sMaze[row][col])
		}
	}
	return tileMaze
}

func pickPlayerGameData(g *game, p *clientPlayer) *playerGameData {
	switch p {
	case g.p1.p:
		return g.p1
	case g.p2.p:
		return g.p2
	case g.p3.p:
		return g.p3
	case g.p4.p:
		return g.p4
	}
	return nil
}

func (g *game) updatePlayerPositions() {
	for y := 0; y < len(g.maze); y++ {
		for x := 0; x < len(g.maze[y]); x++ {
			switch g.maze[y][x] {
			case p1:
				g.p1.position = &posVector{x: x, y: y}
				break
			case p2:
				g.p2.position = &posVector{x: x, y: y}
				break
			case p3:
				g.p3.position = &posVector{x: x, y: y}
				break
			case p4:
				g.p4.position = &posVector{x: x, y: y}
				break
				//case p5:
				//	g.p5.position = &posVector{x: x, y: y}
				//	break
			}
		}
	}
}

func (g *game) setPlayerStartingPositionsToCurrent() {
	g.p1.startingPosition = &posVector{}
	g.p2.startingPosition = &posVector{}
	g.p3.startingPosition = &posVector{}
	g.p4.startingPosition = &posVector{}
	// g.p5.startingPosition = &posVector{}

	*g.p1.startingPosition = *g.p1.position
	*g.p2.startingPosition = *g.p2.position
	*g.p3.startingPosition = *g.p3.position
	*g.p4.startingPosition = *g.p4.position
	//*g.p5.startingPosition = *g.p5.position
}

func (g *game) updatePelletCounts() {
	g.pelletsLeft = 0
	for x := 0; x < len(g.maze); x++ {
		for y := 0; y < len(g.maze[x]); y++ {
			if g.maze[x][y] == pellet {
				g.pelletsLeft++
			}
		}
	}
}

func (g *game) updateTileReferences() {
	g.p1.tileRepresentation = p1
	g.p2.tileRepresentation = p2
	g.p3.tileRepresentation = p3
	g.p4.tileRepresentation = p4
	//g.p5.tileRepresentation = p5
}

func (g *game) getPlayerByID(id string) *playerGameData {
	switch id {
	case "P1":
		return g.p1
	case "P2":
		return g.p2
	case "P3":
		return g.p3
	case "P4":
		return g.p4
	case "P5":
		return g.p5
	}
	return nil
}

func getPlayersOnTile(g *game, t tile) []*playerGameData {
	players := t.getPlayerIDs()
	fmt.Println(players)
	returned := make([]*playerGameData, 0)
	for _, val := range players {
		returned = append(returned, g.getPlayerByID(val))
	}
	return returned
}
