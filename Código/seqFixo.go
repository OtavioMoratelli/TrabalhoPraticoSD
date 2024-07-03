package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func cliente(channel chan string) {
	value := rand.Intn(100)

	time.Sleep(time.Second * 5)

	fmt.Println("Valor gerado pelo cliente e enviado para o sequenciador fixo =", value)
	channel <- strconv.Itoa(value)

	select {
	case ack := <-channel:
		if ack == "Ok" {
			fmt.Println("ACK recebido com sucesso")
		} else {
			fmt.Println("Erro de confirmação de recebimento da mensagem")
			return
		}
	case <-time.After(time.Second * 3):
		fmt.Println("Timeout de recebimento do ACK")
		return
	}

}

func sequenciadorFixo(inchan1 chan string, inchan2 chan string, inchan3 chan string, outchan1 chan string, outchan2 chan string) {
	var data [3]string

	data[0] = <-inchan1
	data[1] = <-inchan2
	data[2] = <-inchan3

	fmt.Println()
	fmt.Println("Valores recebidos pelo sequenciador e ordenados =", data)
	fmt.Println()

	for i := 0; i < 3; i++ {
		outchan1 <- data[i]
		outchan2 <- data[i]
	}

	time.Sleep(time.Second * 1)
	ack1 := <-outchan1
	ack2 := <-outchan2

	if ack1 != "Ok" || ack2 != "Ok" {
		fmt.Println("Erro de confirmação de recebimento da mensagem")
		return
	}

	inchan1 <- "Ok"
	inchan2 <- "Ok"
	inchan3 <- "Ok"

}

func servico(channel chan string) {
	var valorFinal string

	for i := 0; i < 3; i++ {
		temp := <-channel
		valorFinal += temp
	}

	if len(valorFinal) <= 0 {
		fmt.Println("Erro no recebimento da mensagem no serviço")
		return
	}

	channel <- "Ok"
	fmt.Println("Valor final concatenado pelo processo servidor =", valorFinal)
}

func main() {

	chan1 := make(chan string)
	chan2 := make(chan string)
	chan3 := make(chan string)

	chan4 := make(chan string, 3)
	chan5 := make(chan string, 3)

	go cliente(chan1)
	go cliente(chan2)
	go cliente(chan3)

	go sequenciadorFixo(chan1, chan2, chan3, chan4, chan5)

	go servico(chan4)
	go servico(chan5)

	fmt.Scanln()

}
