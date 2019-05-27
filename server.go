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
const tabTamanho = 5

var tab = montarTabuleiro()
var pl []player = playerInit()

type player struct {
	x, y int
}

type celula struct {
	top, bottom, left, right bool
}

type tabuleiro struct {
	tab [tabTamanho][tabTamanho]celula
}

func playerInit() []player {
	p := [player{}, player{}]
	
	p[0].x = 0
	p[0].y = 0
	
	p[1].x = tabTamanho - 1
	p[1].y = tabTamanho - 1
	
	return p
}

func montarTabuleiro() tabuleiro {
	tab := tabuleiro{}
	for i := 0; i < tabTamanho; i++ {
		for j := 0; j < tabTamanho; j++ {
			cel := celula{}

			cel.top = true
			cel.bottom = true
			cel.left = true
			cel.right = true
			
			if j == 0 {
				cel.top = false	
			} else if j == tabTamanho - 1 {
				cel.bottom = false
			}

			if i == 0 {
				cel.left = false	
			} else if i == tabTamanho - 1 {
				cel.right = false
			}
			
			tab.tab[i][j] = cel
		}	
	}
	return tab
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
			if(pl.x < tabTamanho - 1) {
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
			if(pl.y < tabTamanho - 1) {
				pl.y = pl.y + 1
			} 
		default:
			fmt.Print("move default\n")
	}
}*/

func printGame () {
	//fmt.Printf("Player position: %d, %d\n", pl.x, pl.y)
}

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

func main() {
	printGame()
	go keyListener()
	
	var block chan int = make(chan int)
	<-block
}
