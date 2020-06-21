package main

import (
	"strconv"
	"time"
)

type game struct {
}

//game_maker queues players up and puts them under a certain game_manager. It queues them up and delegates them into different games depending on whether a game is
//full etc.

const (
	queueMessageSendCooldown time.Duration = time.Second * 2
	playersInEachGame        int           = 5
)

var queuedPlayersChannel = make(chan *clientPlayer, 5) //channel to handle players joining queue
var queuedPlayers = make([]*clientPlayer, 0, 20)       //slice to store all those who are looking for games.

func queueSystem() {
	for serverActive {
		select {
		case newPlayer := <-queuedPlayersChannel:
			if newPlayer.activeGame != nil {
				newPlayer.writeChannel <- &writeRequest{
					message: "PONG INVALID", //possibly do more than just this.
				}
				break
			}
			go handleQueuedPlayer(newPlayer)
			break
		}
	}
}

func handleQueuedPlayer(newPlayer *clientPlayer) {
	queuedPlayers = append(queuedPlayers, newPlayer)
	breakFlag := false //to handle player having quit while waiting and other terrible misconduct
	if len(queuedPlayers) >= playersInEachGame {
		//if inside here, the last player added completes the lobby. but we must check whether other players are still valid..
		for {
			breakFlag = false

			if len(queuedPlayers) > 1 && !queuedPlayers[0].valid {
				queuedPlayers = queuedPlayers[1:]
				breakFlag = true
			}

			//and now other indices

			if !breakFlag {
				break
			}
		}
	}
	if len(queuedPlayers) >= playersInEachGame {
		//go initializeGameServer(queuedPlayers[0], queuedPlayers[1])

		//remove from array

		// queuedPlayers[0] = queuedPlayers[len(queuedPlayers)-1]
		// queuedPlayers[len(queuedPlayers)-1] = nil
		// queuedPlayers[1] = queuedPlayers[len(queuedPlayers)-2]
		// queuedPlayers[len(queuedPlayers)-2] = nil
		// queuedPlayers = queuedPlayers[:len(queuedPlayers)-2]

		queuedPlayers = queuedPlayers[5:]

	} else {
		defer recover()
		for newPlayer.activeGame == nil && newPlayer.valid {
			tendedPlayersMutex.RLock()
			newPlayer.writeChannel <- &writeRequest{
				message: "MD QUEUE " + strconv.Itoa(tendedPlayers.Len()) + "\n",
			}
			tendedPlayersMutex.RUnlock()
			time.Sleep(queueMessageSendCooldown)
		}
	}
}
