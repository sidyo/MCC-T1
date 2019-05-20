package main

import (
	"fmt"
	//"sync"
	//"time"
	"bufio"
    	"os"
    	"log"
    	"time"
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

func readKey(input chan rune) {
    char, _, err := reader.ReadRune()
    if err != nil {
        log.Fatal(err)
    }
    input <- char
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

func main() {
	test := montarTabuleiro()

	for i := 0; i < tabTamanho; i++ {
		for j := 0; j < tabTamanho; j++ {
			fmt.Println(" ", test.tab)
		}
	}

	input := make(chan rune, 1)
   	fmt.Println("Checking keyboard input...")
    	go readKey(input)

    	select {
        	case i := <-input:
           		fmt.Printf("Input : %v\n", i)
        	case <-time.After(500000 * time.Millisecond):
            		fmt.Println("Time out!")
    	}


}

