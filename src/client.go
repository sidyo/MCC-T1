package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	BEB "../BestEffortBroadcast"
	term "../termbox"
)

type player struct {
	x, y int
}

const maxPlayers = 4
const boardSize = 10
const currentPlayers = 2
const bombExplosionTime = 3 * time.Second
const fireFadeTime = 500 * time.Millisecond

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
		y: 0,
	},
	player{
		x: 0,
		y: boardSize - 2,
	},
}
var chAction chan string = make(chan string)
var chUpdateGame chan string = make(chan string)
var address string
var players [maxPlayers]player
var board [boardSize][boardSize]bool
var bombs [boardSize][boardSize]bool
var fires [boardSize][boardSize]int

func playerInit() {
	for i := 0; i < currentPlayers; i++ {
		players[i] = basePositions[i]
	}
}

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
}

func resetScreen() {
	term.Sync() // cosmestic purpose
	printGame()
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
				break keyPressListenerLoop
			case term.KeyArrowUp:
				//fmt.Println("Arrow Up pressed")
				chAction <- "move " + address + " up"
			case term.KeyArrowDown:
				//fmt.Println("Arrow Down pressed")
				chAction <- "move " + address + " down"
			case term.KeyArrowLeft:
				//fmt.Println("Arrow Left pressed")
				chAction <- "move " + address + " left"
			case term.KeyArrowRight:
				//fmt.Println("Arrow Right pressed")
				chAction <- "move " + address + " right"
			case term.KeySpace:
				//fmt.Println("Space pressed")
				chAction <- "bomb " + address
			default:
				// we only want to read a single character or one key pressed event
				//fmt.Println("ASCII : ", ev.Ch)
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
				for p := 0; p < currentPlayers; p++ {
					if players[p].x == i && players[p].y == j {
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
			chUpdateGame <- in.Message
		}
	}()

	blq := make(chan int)
	<-blq
}

func updateGame() {
	for {
		msg := <-chUpdateGame
		aux := strings.Split(msg, " ")
		if aux[0] == "bomb" {
			player := aux[1]
			size := aux[2]
			bomb(player, size)
		} else if aux[0] == "move" {
			player := aux[1]
			direction := aux[2]
			move(player, direction)
		} else {
			fmt.Println("MESSAGE ERROR!")
		}
	}
}

func bomb(playerId string, size string) {
	p, _ := strconv.Atoi(playerId)
	p = p - 1
	s, _ := strconv.Atoi(size)
	x := players[p].x
	y := players[p].y
	go bombAux(x, y, s)
}

func bombAux(x int, y int, size int) {
	bombs[x][y] = true
	resetScreen()
	time.Sleep(2 * time.Second)
	bombs[x][y] = false
	resetScreen()
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
	for i := 0; i < currentPlayers; i++ {
		x := players[i].x
		y := players[i].y
		if fires[x][y] > 0 {
			fmt.Printf("PLAYER %v TOOK DAMAGE!\n", i+1)
		}
	}
}

func move(playerId string, direction string) {
	p, _ := strconv.Atoi(playerId)
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
	for p := 0; p < currentPlayers; p++ {
		if players[p].x == x && players[p].y == y {
			playerPos = true
		}
	}
	return playerPos || board[x][y] || bombs[x][y]
}
func main() {
	go networkinit()
	go updateGame()
	buildBoard()
	playerInit()
	go keyListener()
	time.Sleep(1000000000)
	resetScreen()

	var block chan int = make(chan int)
	<-block
}

//go run client.go 127.0.0.1:2001 127.0.0.1:1001
//go run client.go 127.0.0.1:3001 127.0.0.1:1001
