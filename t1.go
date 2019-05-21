package main

import (
	"fmt"
	//"sync"
	//"time"
	"bufio"
    	"os"
    	"log"
		
        term "github.com/nsf/termbox-go"
)

var reader = bufio.NewReader(os.Stdin)
const maxPlayers = 3
const tabTamanho = 5

type celula struct {
	top, bottom, left, right bool
}

type tabuleiro struct {
	tab [tabTamanho][tabTamanho]celula
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

func reset() {
        term.Sync() // cosmestic purpose
}
func keyListener(){
	err := term.Init()
	if err != nil {
		 panic(err)
	}

	defer term.Close()

	fmt.Println("Enter any key to see their ASCII code or press ESC button to quit")

	keyPressListenerLoop:
		for {
			switch ev := term.PollEvent(); ev.Type {
			case term.EventKey:
				switch ev.Key {
					case term.KeyEsc:
							break keyPressListenerLoop
					case term.KeyArrowUp:
							reset()
							fmt.Println("Arrow Up pressed")
					case term.KeyArrowDown:
							reset()
							fmt.Println("Arrow Down pressed")
					case term.KeyArrowLeft:
							reset()
							fmt.Println("Arrow Left pressed")
					case term.KeyArrowRight:
							reset()
							fmt.Println("Arrow Right pressed")
					case term.KeySpace:
							reset()
							fmt.Println("Space pressed")
					default:
							// we only want to read a single character or one key pressed event
							reset()
							fmt.Println("ASCII : ", ev.Ch)
				}
			case term.EventError:
				 panic(ev.Err)
			}
		}
}


func main() {
	test := montarTabuleiro()

	for i := 0; i < tabTamanho; i++ {
		for j := 0; j < tabTamanho; j++ {
			fmt.Println(" ", test.tab)
		}
	}

	keyListener()


}
