package main

import (
	"fmt"
	//"sync"

	"os"
	"strconv"
	"strings"
	"time"

	BEB "../BestEffortBroadcast"
)

var addresses []string
var chAction chan string = make(chan string)

const bombSize = 3
const maxPlayers = 3
const boardSize = 10

func findPlayer(address string) int {
	for i := 0; i < len(addresses); i++ {
		if addresses[i] == address {
			return i
		}
	}
	return -1
}

func networkinit() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses = os.Args[1:]
	fmt.Println(addresses)

	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	beb.Init(addresses[0])

	// enviador de broadcasts
	go func() {

		for {
			msg := <-chAction
			req := BEB.BestEffortBroadcast_Req_Message{
				Addresses: addresses[1:],
				Message:   msg}
			beb.Req <- req
		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			fmt.Println("Received: " + in.Message)
			chAction <- treatMessage(in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}

func treatMessage(msg string) string {
	aux := strings.Split(msg, " ")
	var playerId int
	if aux[0] == "move" {
		playerId = findPlayer(aux[1])
		return aux[0] + " " + strconv.Itoa(playerId) + " " + aux[2]
	} else if aux[0] == "bomb" {
		playerId = findPlayer(aux[1])
		return aux[0] + " " + strconv.Itoa(playerId) + " " + strconv.Itoa(bombSize)
	}
	return "The server received an invalid message!"
}

func main() {
	go networkinit()
	time.Sleep(1000000000)
	var block chan int = make(chan int)
	<-block
}
