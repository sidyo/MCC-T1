// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti
// ATENCAO: codigo parcialmente encontrado na internet livremente
// usado aqui com objetivo de exemplificacao de sincronizacao.
// Note que a linguagem Go conta com sua propria biblioteca de
// sincronizacao.   Aqui estamos exemplificando como construir a
// semantica de semaforos a partir do uso de canais como forma de
// prover atomicidade.    Funcoes s.wait() e s.signal() tem o
// mesmo significado da literatura.

// este pacote oferece a abstracao de semaforo da literatura
// atraves da estrutura MCCSemaforo.Semaphore
// semaphore.Wait() e .Signal() sao as operacoes de semaforos
// tipicas
// Instrucoes rapidas:
// coloque o pacote MCCSemaforo (este arquivo) dentro de um diretorio
// chamado MCCSemaforo, no diretorio corrente (onde esta seu codigo).
//
// No seu codigo que usa semaforo, faca:
// import (
//	"./MCCSemaforo"
// )
// exemplo de declaracao de um semaforo:
//      MCCSemaforo.Semaphore s = MCCSemaforo.NewSemaphore(1)
//      s.Wait()
//      s.Signal()

package MCCSemaforo

type Semaphore struct {
	wai, sig chan struct{} // canais para wait e signal
	val      int           // valor do semaforo
}

func NewSemaphore(v int) *Semaphore {
	s := &Semaphore{
		wai: make(chan struct{}),
		sig: make(chan struct{}),
		val: v} // inicia semaforo com um valor, deve ser >= 0

	go func() {
		for {
			if s.val == 0 { // se val == 0
				<-s.sig // pode permitir apenas signal, processos fazendo wait bloqueiam
				s.val++ // se acontecer signal, entao incrementa val
			}
			if s.val > 0 { // senao pode permitir tanto wait como signal, alterando val
				select {
				case <-s.sig:
					s.val++
				case <-s.wai:
					s.val--
				}
			}
		}
	}()
	return s
}

func (s *Semaphore) Wait() { // fazer wait ee ter sucesso na sincronizacao em s.wai
	s.wai <- struct{}{} // se nao sincroniza, fica em espera, implementando a espera do Wait()
}

func (s *Semaphore) Signal() { // fazer signal ee ter sicesso na sincronizacao em s.sig
	s.sig <- struct{}{} // esta sincronizacao sempre ee possivel como visto acima
}
