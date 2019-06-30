package main

import (
	"fmt"
	//"sync"

	"os"
	"strconv"
	"strings"

	BEB "../BestEffortBroadcast"
)

var addresses []string
var chServer = make(chan string)
var clients [maxPlayers]client
var currentClients = 0

type client struct {
	address string
	pNum    int
	active  bool
}
type player struct {
	x, y int
}

var basePositions = []player{
	player{
		x: 1,
		y: 1,
	},
	player{
		x: boardSize - 2,
		y: boardSize - 2,
	},
	player{
		x: boardSize - 2,
		y: 1,
	},
	player{
		x: 1,
		y: boardSize - 2,
	},
}

const bombSize = 3
const maxPlayers = 4
const boardSize = 10

func findPlayer(address string) int {
	for i := 0; i < maxPlayers; i++ {
		if clients[i].active && clients[i].address == address {
			return clients[i].pNum
		}
	}
	return -1
}

func getAdresses() []string {
	aux := ""
	for i := 0; i < maxPlayers; i++ {
		if clients[i].active {
			aux += " " + clients[i].address
		}
	}
	return strings.Split(strings.Trim(aux, " "), " ")
}

func serverNetworkInit() {

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
			msg := <-chServer
			fmt.Println("--\nMessage: " + msg + "\n Sent to: " + strings.Join(getAdresses(), ", "))

			req := BEB.BestEffortBroadcast_Req_Message{
				Addresses: getAdresses(),
				Message:   msg}
			beb.Req <- req
		}
	}()
	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			fmt.Println("--\nMessage: " + in.Message + "\n From: " + in.From)
			chServer <- treatMessage(in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}

func findOpenSpot() int {
	for i := 0; i < currentClients; i++ {
		if clients[i].active == false {
			return i
		}
	}
	return -1
}

func treatMessage(msg string) string {
	aux := strings.Split(msg, " ")
	var playerNum int
	if aux[0] == "move" {
		playerNum = findPlayer(aux[1])
		return aux[0] + " " + strconv.Itoa(playerNum) + " " + aux[2]
	} else if aux[0] == "bomb" {
		playerNum = findPlayer(aux[1])
		return aux[0] + " " + strconv.Itoa(playerNum) + " " + strconv.Itoa(bombSize)
	} else if aux[0] == "new" {
		if currentClients < maxPlayers {
			if currentClients > 0 {
				chServer <- "info"
			}
			currentClients++
			clients[findOpenSpot()] = client{
				address: aux[1],
				pNum:    currentClients,
				active:  true,
			}
			pNum := findPlayer(aux[1])
			return "new " + strconv.Itoa(pNum) + " " + strconv.Itoa(basePositions[pNum-1].x) + " " + strconv.Itoa(basePositions[pNum-1].y)
		}
		fmt.Println("Connection rejected. Max players reached.")
	} else if aux[0] == "info" {
		msg := "sync"
		for i := 1; i < len(aux); i++ {
			msg += " " + aux[i]
		}
		return msg
	} else if aux[0] == "close" {
		currentClients--
		for i := 0; i < maxPlayers; i++ {
			if clients[i].active && clients[i].address == aux[1] {
				clients[i].active = false
				return "close " + strconv.Itoa(i)
			}
		}
	}
	return "Server received an invalid message: \n" + msg
}

func main() {
	go serverNetworkInit()
	var block = make(chan int)
	<-block
}

//go run server.go 127.0.0.1:1001
