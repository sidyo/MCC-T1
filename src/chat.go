package main

import "fmt"
import "os"
import "bufio"
import BEB "../BestEffortBroadcast"

func main() {

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

/*

go run servidor.go 127.0.0.1:1001 127.0.0.1:2001 127.0.0.1:3001

go run main.go 127.0.0.1:2001  127.0.0.1:1001

go run main.go 127.0.0.1:2001  127.0.0.1:1001

*/
