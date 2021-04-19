package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 1. С помощью пула воркеров написать программу, которая запускает 1000 горутин, каждая из которых увеличивает число на 1.
// Дождаться завершения всех горутин и убедиться, что при каждом запуске программы итоговое число равно 1000.

// 2. Написать программу, которая при получении в канал сигнала SIGTERM останавливается не позднее, чем за одну секунду (установить таймаут).

func main() {
	firstExc()
	secondExc()
}

func firstExc() {
	cntChan := make(chan int, 1000)
	cnt := 0
	for i := 0; i < 1000; i++ {
		go func() {
			cntChan <- 1
		}()
		cnt += <-cntChan
	}
	fmt.Println(cnt)
}

func secondExc() {

	chn := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(chn, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-chn
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Ожидание сигнала")
	<-done
	time.Sleep(1 * time.Second)
}
