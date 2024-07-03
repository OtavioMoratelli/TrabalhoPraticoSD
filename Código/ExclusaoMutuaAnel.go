package main

import (
	"fmt"
	"time"
)

// Define o número de processos no anel
const numProcesses = 5

// Estrutura para representar um processo no anel
type Process struct {
	id       int        // Identificador do processo
	token    chan bool  // Canal para receber o token
	nextProc *Process   // Ponteiro para o próximo processo no anel
}

// Função para simular a entrada na região crítica
func (p *Process) enterCriticalSection() {
	fmt.Printf("Process %d entrando na região crítica\n", p.id)
	time.Sleep(1 * time.Second) // Simula o tempo gasto na região crítica
	fmt.Printf("Process %d saindo da região crítica\n", p.id)
}

// Função para passar o token para o próximo processo
func (p *Process) passToken() {
	fmt.Printf("Process %d passando o token para o Processo %d\n", p.id, p.nextProc.id)
	p.nextProc.token <- true // Envia o token para o próximo processo
}

// Função que representa o ciclo de vida do processo
func (p *Process) start() {
	for {
		select {
		case <-p.token: // Aguarda receber o token
			fmt.Printf("Process %d recebeu o token\n", p.id)
			p.enterCriticalSection() // Entra na região crítica
			p.passToken()            // Passa o token para o próximo processo
		}
	}
}

func main() {
	// Cria uma lista de processos
	processes := make([]*Process, numProcesses)
	for i := 0; i < numProcesses; i++ {
		processes[i] = &Process{
			id:    i,
			token: make(chan bool, 1), // Canal bufferizado para o token
		}
	}

	// Configura o próximo processo para cada processo, formando o anel
	for i := 0; i < numProcesses; i++ {
		processes[i].nextProc = processes[(i+1)%numProcesses]
	}

	// Inicia todas as goroutines dos processos
	for i := 0; i < numProcesses; i++ {
		go processes[i].start()
	}

	// Gera e envia o token inicial para o primeiro processo
	fmt.Println("Enviando o token inicial para o Processo 0")
	processes[0].token <- true

	// Deixa a simulação rodar por um tempo determinado
	time.Sleep(10 * time.Second)
}
