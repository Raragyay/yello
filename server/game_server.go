package main

type direction string

const (
	left  direction = "L"
	right direction = "R"
	up    direction = "U"
	down  direction = "D"
)

type playerGameData struct { //TODO LOAD GAME DATA
	p               *clientPlayer
	latestDirection direction
	mazeP           *player
}

func initializeGameServer(p1, p2, p3, p4, p5 *clientPlayer) {
	gameInstance := game{
		p1: &playerGameData{p: p1},
		p2: &playerGameData{p: p2},
		//p3:     &playerGameData{p: p3},
		//p4:     &playerGameData{p: p4},
		//p5:     &playerGameData{p: p5},
		active: true,
	}

	p1.writeChanneledMessage("PONG GAME-INIT " + "P1-" + p2.name + "-" + p3.name + "-" + p4.name + "-" + p5.name) //Who needs JSON when you got -?
	p2.writeChanneledMessage("PONG GAME-INIT " + "P2-" + p1.name + "-" + p3.name + "-" + p4.name + "-" + p5.name)
	//p3.writeChanneledMessage("PONG GAME-INIT " + "P3-" + p1.name + "-" + p2.name + "-" + p4.name + "-" + p5.name)
	//p4.writeChanneledMessage("PONG GAME-INIT " + "P4-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p5.name)
	//p5.writeChanneledMessage("PONG GAME-INIT " + "P5-" + p1.name + "-" + p2.name + "-" + p3.name + "-" + p4.name)

	p1.activeGame = &gameInstance
	p2.activeGame = &gameInstance
	//p3.activeGame = &gameInstance
	//p4.activeGame = &gameInstance
	//p5.activeGame = &gameInstance

	tendGame(&gameInstance)
}

func tendGame(g *game) {
	for g.active {
		//ze game loop!
	}
}

//INNER PROCESSES

func checkUpdateObjectStates() {

}

//TO ZE CLIENTSSSS

func updateObjectPosition(g *game, o *gameObject, v *posVector) {
	p1.writeChanneledMessage("PONG GAME-OBJECT-POS " + o.string_ID + "-" + v.toString())
	//...
}

func updateObjectState(g *game, o *gameObject, state objectState) {
	p1.writeChanneledMessage("PONG GAME-OBJECT-STATE " + o.string_ID + "-" + string(objectState))
}

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
