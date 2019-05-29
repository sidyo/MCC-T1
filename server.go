package main

import (
    "fmt"
    //"sync"
    "time"
    "bufio"
    "os"
)

var reader = bufio.NewReader(os.Stdin)

const maxPlayers = 3
const boardSize = 10

var board [boardSize][boardSize]bool
var pl [2]player = playerInit()

type player struct {
    x, y int
}

func buildBoard() {
    for j := 0; j < boardSize; j++ {
        for i := 0; i < boardSize; i++ {
            if j == 0 || j == boardSize - 1 || i == 0 || i == boardSize - 1 {
                board[j][i] = true
            } else {
                board[j][i] = false
            }
        }
    }
}

func playerInit() [2]player {
    p := [2]player{player{}, player{}}

    p[0].x = 1
    p[0].y = 1

    p[1].x = boardSize - 2
    p[1].y = boardSize - 2

    return p
}

func bomb() {
    //fmt.Printf("Bomb planted at: %d, %d\n",pl.x,pl.y)
    time.Sleep(2000000000)
    fmt.Print("Booom!\n")
}

/*
func move(player int, action int) {
	switch action {
		case 1:
			if(pl.x < boardSize - 1) {
				pl.x = pl.x + 1
			}
		case 2:
			if(pl.y > 0) {
				pl.y = pl.y - 1
			}
		case 3:
			if(pl.x > 0) {
				pl.x = pl.x - 1
			}
		case 4:
			if(pl.y < boardSize - 1) {
				pl.y = pl.y + 1
			}
		default:
			fmt.Print("move default\n")
	}
}*/

func printGame () {
    for j := 0; j < boardSize; j++ {
        for i := 0; i < boardSize; i++ {
            if board[j][i] {
                fmt.Printf("# ")
            } else if pl[0].x == i && pl[0].y == j {
                fmt.Printf("1 ")
            } else if pl[1].x == i && pl[1].y == j {
                fmt.Printf("2 ")
            } else {
                fmt.Printf("  ")
            }
        }
        fmt.Printf("\n")
    }
}

/*
func networkinit() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	beb.Init(addresses[0])

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
			}
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
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)

		}
	}()

	blq := make(chan int)
	<-blq
}
*/

func main() {
    //go keyListener()
    buildBoard()
    printGame()
    //var block chan int = make(chan int)
    //<-block
}
