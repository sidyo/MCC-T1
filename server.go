package main

import (
    "fmt"
    //"sync"
    "time"
    "bufio"
    "os"
)


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


// -----------------------------------------
//
// -----------------------------------------

type Semaphore struct {
	inc, dec chan struct{}
	val      int
}

func NewSemaphore(v int) *Semaphore {
	s := &Semaphore{
		inc: make(chan struct{}),
		dec: make(chan struct{}),
		val: v}

	go func() {
		for {
			if s.val == 0 {
				<-s.inc
				s.val++
			}
			if s.val > 0 {
				select {
				case <-s.inc:
					s.val++
				case <-s.dec:
					s.val--
				}
			}
		}
	}()
	return s
}

func (s *Semaphore) Signal() {
	s.inc <- struct{}{}
}

func (s *Semaphore) Wait() {
	s.dec <- struct{}{}
}

type Monitor struct {
	mutex      *Semaphore // garante exclusão mutua do monitor
	next       *Semaphore // bloqueia thread que sinaliza em favor de outra - vide signal de condition
	next_count int        // conta threads em next, que podem prossegir
}

func initMonitor() *Monitor {
	m := &Monitor{
		mutex:      NewSemaphore(1),
		next:       NewSemaphore(0),
		next_count: 0}

	return m
}

//  procedimentos do monitor

func (m *Monitor) monitorEntry() {
	m.mutex.Wait() // entrada no monitor ee so passar pelo mutex
}

func (m *Monitor) monitorExit() {
	if m.next_count > 0 { // libera uma thread que ja esteve no monitor, senao libera mutex
		m.next.Signal()
	} else {
		m.mutex.Signal()
	}
}

//  ------------------------------------------------------
//  ------------------------------------------------------
//  estruturas genericas de  variaveis condicao

type Condition struct {
	s     *Semaphore // semaforo para bloquear na condicao
	count int        // contador de bloqueados
	m     *Monitor   // monitor associado aa condicao - quando bloqueia na condicao libera o monitor (next ou mutex)
	name  string
}

func initCondition(n string, m1 *Monitor) *Condition {
	c := &Condition{
		s:     NewSemaphore(0), // 0 inicia bloqueando
		count: 0,               // contadores de bloqueados nesta condicao
		m:     m1,              // o monitor associado
		name:  n}

	return c
}

//  procedimentos de variaveis condicao

func (c *Condition) condWait() {
	// fmt.Println("                                           wait  ", c.name)
	c.count++         // mais uma thread vai bloquear aqui nesta condition
	c.m.monitorExit() // libera o monitor associado aa condition
	c.s.Wait()        // bloqueia !!     fica aqui ate alguem dar signal  !!
	c.count--         // esta linha é executada depois de alguem ter dado signal, entao um bloqueado a menos
}

func (c *Condition) condSignal() {
	if c.count > 0 { // tem alguem para sinalizar ?    se nao tem entao nao faz nada, signal nao tem efeito!
		c.m.next_count++ // opa, tem alguem para sinalizar, entao esta thread se bloqueia em favor da sinalizada
		c.s.Signal()
		c.m.next.Wait()
		c.m.next_count-- // foi desbloqueada (veja monitorExit), aqui desbloqueou, decrementa
	}
}


// -----------------------------------------
//
// -----------------------------------------


type printMonitor struct {
	// sincronizacao
	m     *Monitor
	cond *Condition
}

func printMonitorInit() *printMonitor {
	mon := initMonitor()
	mbc := &printMonitor{
		m:     mon,
		cond: initCondition("test", mon),
	}
	return mbc
}


var reader = bufio.NewReader(os.Stdin)

const maxPlayers = 3
const boardSize = 10

var board [boardSize][boardSize]bool
var bombs [boardSize][boardSize]bool

var pl [2]player = playerInit()

type player struct {
    x, y int
}


func printGame (mbc *printMonitor) {
    mbc.m.monitorEntry()
    for j := 0; j < boardSize; j++ {
        for i := 0; i < boardSize; i++ {
            if board[j][i] {
                fmt.Printf("# ")
            } else if bombs[j][i] {
                fmt.Printf("@ ")
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
    mbc.m.monitorExit()
}


func bomb(mbc *printMonitor, x int, y int) {
    bombs[x][y] = true
    printGame(mbc)
    //fmt.Printf("Bomb planted at: %d, %d\n",pl.x,pl.y)
    time.Sleep(2000000000)
    //fmt.Print("Booom!\n")
    bombs[x][y] = false
    printGame(mbc)
}

func main() {


    mbc := printMonitorInit()

    //go keyListener()
    buildBoard()
    printGame(mbc)
    time.Sleep(2000000000)
    go bomb(mbc, 3,4)
    go bomb(mbc, 4,4)
    var block chan int = make(chan int)
    <-block
}
