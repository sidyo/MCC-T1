package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	BEB "../BestEffortBroadcast"
	"../MCCSemaforo"
	term "../termbox"
)

type player struct {
	x, y   int
	active bool
}

const maxPlayers = 4
const boardSize = 10
const bombExplosionTime = 3 * time.Second
const fireFadeTime = 500 * time.Millisecond

var mutex = MCCSemaforo.NewSemaphore(1)
var networkReady = MCCSemaforo.NewSemaphore(0)
var boardReady = MCCSemaforo.NewSemaphore(0)
var chClient = make(chan string)
var chFromServer = make(chan string)
var block = make(chan int)
var address string
var players [maxPlayers]player
var board [boardSize][boardSize]bool
var bombs [boardSize][boardSize]bool
var fires [boardSize][boardSize]int

func buildBoard() {
	for j := 0; j < boardSize; j++ {
		for i := 0; i < boardSize; i++ {
			if j == 0 || j == boardSize-1 || i == 0 || i == boardSize-1 {
				board[j][i] = true
			} else {
				board[j][i] = false
			}
		}
	}
	boardReady.Signal()
}

func resetScreen() {
	mutex.Wait()
	term.Sync() // cosmestic purpose
	printGame()
	mutex.Signal()
}

func keyListener() {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				fmt.Println("Game Closed!")
				chClient <- "close " + address
				break keyPressListenerLoop
			case term.KeyArrowUp:
				//fmt.Println("Arrow Up pressed")
				chClient <- "move " + address + " up"
			case term.KeyArrowDown:
				//fmt.Println("Arrow Down pressed")
				chClient <- "move " + address + " down"
			case term.KeyArrowLeft:
				//fmt.Println("Arrow Left pressed")
				chClient <- "move " + address + " left"
			case term.KeyArrowRight:
				//fmt.Println("Arrow Right pressed")
				chClient <- "move " + address + " right"
			case term.KeySpace:
				//fmt.Println("Space pressed")
				chClient <- "bomb " + address
			default:
				// Ignore any other key press
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

func printGame() {
	for j := 0; j < boardSize; j++ {
		for i := 0; i < boardSize; i++ {
			playerpos := false
			if board[i][j] {
				fmt.Print("# ")
			} else if bombs[i][j] {
				fmt.Print("@ ")
			} else if fires[i][j] > 0 {
				fmt.Print("* ")
			} else {
				for p := 0; p < maxPlayers; p++ {
					if players[p].active && players[p].x == i && players[p].y == j {
						fmt.Printf("%v ", p+1)
						playerpos = true
					}
				}
				if !playerpos {
					fmt.Print("  ")
				}
			}
		}
		fmt.Printf("\n")
	}
}

func networkinit() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	address = os.Args[1]
	addresses := os.Args[1:]
	fmt.Println(addresses)

	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	beb.Init(addresses[0])

	// enviador de broadcasts
	go func() {
		for {
			msg := <-chClient
			req := BEB.BestEffortBroadcast_Req_Message{
				Addresses: addresses[1:],
				Message:   msg}
			beb.Req <- req
			if msg == "close "+address {
				time.Sleep(2 * time.Second)
				block <- 1
			}
		}
	}()
	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			chFromServer <- in.Message
		}
	}()
	networkReady.Signal()
	blq := make(chan int)
	<-blq
}

func messageHandler() {
	for {
		msg := <-chFromServer
		aux := strings.Split(msg, " ")
		if aux[0] == "bomb" {
			player := aux[1]
			size := aux[2]
			bomb(player, size)
		} else if aux[0] == "move" {
			player := aux[1]
			direction := aux[2]
			move(player, direction)
		} else if aux[0] == "new" {
			p, _ := strconv.Atoi(aux[1])
			x, _ := strconv.Atoi(aux[2])
			y, _ := strconv.Atoi(aux[3])
			players[p-1] = player{
				x:      x,
				y:      y,
				active: true,
			}
			resetScreen()
		} else if aux[0] == "info" {
			msg := "info"
			for i := 0; i < maxPlayers; i++ {
				if players[i].active {
					msg += " " + strconv.Itoa(i+1) + " " + strconv.Itoa(players[i].x) + " " + strconv.Itoa(players[i].y)
				}
			}
			chClient <- msg
		} else if aux[0] == "sync" {
			if len(aux) > 1 {
				syncAux := aux[1:]
				for i := 0; i < len(syncAux); i += 3 {
					p, _ := strconv.Atoi(syncAux[i])
					x, _ := strconv.Atoi(syncAux[i+1])
					y, _ := strconv.Atoi(syncAux[i+2])
					players[p-1].x = x
					players[p-1].y = y
					players[p-1].active = true
				}
				resetScreen()
			}
		} else if aux[0] == "close" {
			p, _ := strconv.Atoi(aux[1])
			players[p].active = false
			resetScreen()
		} else {
			fmt.Println("Client received an invalid message: \n" + msg)
		}
	}
}

func bomb(playerNum string, size string) {
	p, _ := strconv.Atoi(playerNum)
	p = p - 1
	s, _ := strconv.Atoi(size)
	x := players[p].x
	y := players[p].y
	go bombAux(x, y, s)
}

func bombAux(x int, y int, size int) {
	bombs[x][y] = true
	resetScreen()
	time.Sleep(bombExplosionTime)
	bombs[x][y] = false
	fires[x][y]++
	for i := 1; i <= size; i++ {
		if x-i > 0 {
			fires[x-i][y]++
		}
		if y-i > 0 {
			fires[x][y-i]++
		}
		if x+i < boardSize {
			fires[x+i][y]++
		}
		if y+i < boardSize {
			fires[x][y+i]++
		}
	}
	resetScreen()
	hitDetection()
	time.Sleep(fireFadeTime)
	fires[x][y]--
	for i := 1; i <= size; i++ {
		if x-i > 0 {
			fires[x-i][y]--
		}
		if y-i > 0 {
			fires[x][y-i]--
		}
		if x+i < boardSize {
			fires[x+i][y]--
		}
		if y+i < boardSize {
			fires[x][y+i]--
		}
	}
	resetScreen()
}

func hitDetection() {
	for i := 0; i < maxPlayers; i++ {
		if players[i].active {
			x := players[i].x
			y := players[i].y
			if fires[x][y] > 0 {
				fmt.Printf("PLAYER %v TOOK DAMAGE!\n", i+1)
			}
		}
	}
}

func move(playerNum string, direction string) {
	p, _ := strconv.Atoi(playerNum)
	p = p - 1
	x := players[p].x
	y := players[p].y
	collision := false
	xUpd := x
	yUpd := y
	switch direction {
	case "up":
		collision = checkCollision(x, y-1)
		yUpd = yUpd - 1
	case "down":
		collision = checkCollision(x, y+1)
		yUpd = yUpd + 1
	case "left":
		collision = checkCollision(x-1, y)
		xUpd = xUpd - 1
	case "right":
		collision = checkCollision(x+1, y)
		xUpd = xUpd + 1
	}
	if !collision {
		players[p].x = xUpd
		players[p].y = yUpd
		resetScreen()
	}
}

func checkCollision(x int, y int) bool {
	if x < 0 || y < 0 {
		return true
	}
	playerPos := false
	for p := 0; p < maxPlayers; p++ {
		if players[p].active {
			if players[p].x == x && players[p].y == y {
				playerPos = true
			}
		}
	}
	return playerPos || board[x][y] || bombs[x][y]
}

func main() {
	buildBoard()
	go keyListener()
	go messageHandler()
	go networkinit()
	networkReady.Wait()
	boardReady.Wait()
	time.Sleep(1 * time.Second) //termbox messes up with the screen sometimes
	chClient <- "new " + address

	<-block
}

//go run client.go 127.0.0.1:3001 127.0.0.1:1001
