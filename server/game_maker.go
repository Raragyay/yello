package main

import (
	"fmt"
	"strconv"
	"time"
)

//game_maker queues players up and puts them under a certain game_manager. It queues them up and delegates them into different games depending on whether a game is
//full etc.

const (
	queueMessageSendCooldown time.Duration = time.Second * 2
	playersInEachGame        int           = 4
)

var queuedPlayersChannel = make(chan *clientPlayer, 5) //channel to handle players joining queue
var queuedPlayers = make([]*clientPlayer, 0, 20)       //slice to store all those who are looking for games.

func queuePlayer(req *playerRequest, argument string) {
	if req.p.queued {
		fmt.Println("player has tried to queue up but he is already queued up: " + req.p.name)
		return
	}
	fmt.Println("player has queued up: " + req.p.name)
	queuedPlayersChannel <- req.p
}

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
	newPlayer.queued = true
	queuedPlayers = append(queuedPlayers, newPlayer)
	if len(queuedPlayers) >= playersInEachGame {
		//if inside here, the last player added completes the lobby. but we must check whether other players are still valid..
		for i := 0; i < playersInEachGame; i++ {
			if !queuedPlayers[i].valid {
				for j := i + 1; j < playersInEachGame; j++ {
					queuedPlayers[j-1] = queuedPlayers[j]
				}
				queuedPlayers = queuedPlayers[:len(queuedPlayers)-1] //discard last element that is now a duplicate
				break                                                //TODO FIX THIS. I NEED EVERYONE FOR THIS OR AT LEAST 2 PEOPLE. QUEUED PLAYERS SLICE DOES NOT GET FIXED PROPERLY
			}
		}
	}
	if len(queuedPlayers) >= playersInEachGame {

		go initializeGameServer(queuedPlayers[0], queuedPlayers[1], queuedPlayers[2], queuedPlayers[3]) //add others later
		for i := 0; i < playersInEachGame; i++ {
			queuedPlayers[i].queued = false
		}
		queuedPlayers = queuedPlayers[playersInEachGame:] //TODO SLICE LENGTH ISSUE

	} else {
		fmt.Println("player", newPlayer.name, "joined queue and the queue now has at most", len(queuedPlayers), "players.")
		defer recover()
		for newPlayer.activeGame == nil && newPlayer.valid {
			tendedPlayersMutex.RLock()
			newPlayer.writeChannel <- &writeRequest{
				message: "PONG QUEUE " + strconv.Itoa(tendedPlayers.Len()) + "\n",
			}
			tendedPlayersMutex.RUnlock()
			time.Sleep(queueMessageSendCooldown)
		}
	}
}
