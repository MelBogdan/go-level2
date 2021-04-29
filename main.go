package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

// 3. Смоделировать ситуацию “гонки”, и проверить программу на наличии “гонки”

func Race() {

	m := make(map[int]int)

	m[2] = 2

	go func() {

		m[1] = 1
	}()

	for i, j := range m {
		fmt.Println(i, "-", j)
	}

}

// ==================
// WARNING: DATA RACE
// Read at 0x00c0000a2150 by main goroutine:
//   runtime.mapiternext()
//       /usr/lib/go-1.13/src/runtime/map.go:851 +0x0
//   main.main()
//       /home/qwerty/go/src/github.com/MelBogdan/go-level2/main.go:22 +0x216

// Previous write at 0x00c0000a2150 by goroutine 7:
//   runtime.mapassign_fast64()
//       /usr/lib/go-1.13/src/runtime/map_fast64.go:92 +0x0
//   main.main.func1()
//       /home/qwerty/go/src/github.com/MelBogdan/go-level2/main.go:19 +0x4d

// Goroutine 7 (running) created at:
//   main.main()
//       /home/qwerty/go/src/github.com/MelBogdan/go-level2/main.go:17 +0xae
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000a2150 by goroutine 7:
//   runtime.mapassign_fast64()
//       /usr/lib/go-1.13/src/runtime/map_fast64.go:92 +0x0
//   main.main.func1()
//       /home/qwerty/go/src/github.com/MelBogdan/go-level2/main.go:19 +0x4d

// Previous read at 0x00c0000a2150 by main goroutine:
//   runtime.mapiterinit()
//       /usr/lib/go-1.13/src/runtime/map.go:802 +0x0
//   main.main()
//       /home/qwerty/go/src/github.com/MelBogdan/go-level2/main.go:22 +0xf7

// Goroutine 7 (running) created at:
//   main.main()
//       /home/qwerty/go/src/github.com/MelBogdan/go-level2/main.go:17 +0xae
// ==================

// 1. Написать программу, которая использует мьютекс для безопасного доступа к данным из нескольких потоков. Выполните трассировку программы

type StoreCache struct {
	sync.Mutex
	cache map[string]string
}

func New() *StoreCache {
	return &StoreCache{
		cache: make(map[string]string),
	}
}

func (sc *StoreCache) Set(key string, value string) {
	sc.cache[key] = value
}

func (sc *StoreCache) Get(key string) string {
	if sc.Count() > 0 {
		item := sc.cache[key]
		return item
	}
	return ""
}

func (sc *StoreCache) Count() int {
	return len(sc.cache)
}

func mutex() {
	tr, _ := os.Create("trace.out")
	defer tr.Close()
	trace.Start(tr)
	defer trace.Stop()

	store := New()
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		store.Set("Go", "Lang")
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		store.Set("Go1", "Lang1")
	}()

	wg.Wait()

	result := store.Get("Go")
	result1 := store.Get("Go1")

	fmt.Println(result, result1)
}

// 2. Написать многопоточную программу, в которой будет использоваться явный вызов планировщика. Выполните трассировку программы

func say(str string) {
	fmt.Println(str)
}

func main() {
	tr, _ := os.Create("trace1.out")
	defer tr.Close()
	trace.Start(tr)
	defer trace.Stop()

	wg := sync.WaitGroup{}
	runtime.GOMAXPROCS(1)

	for i := 0; i < 5; i++ {
		//runtime.Gosched()
		wg.Add(1)

		go func() {
			defer wg.Done()
			say("hello")
		}()
		say("everybody")
	}
	wg.Wait()
}
