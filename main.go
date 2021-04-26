package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Задание №1: Напишите программу, которая запускает n потоков и дожидается завершения их всех:

// func main() {
// 	var wg = sync.WaitGroup{}
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			fmt.Printf("This is the %d routine\n", i)
// 		}(i)
// 	}
// 	wg.Wait()
// 	fmt.Println("main done")
// }

// Задание №3: Протестируйте производительность операций чтения и записи на множестве действительных чисел,
//  безопасность которого обеспечивается sync.Mutex и sync.RWMutex для разных вариантов использования:
//  10% запись, 90% чтение; 50% запись, 50% чтение; 90% запись, 10% чтение

var (
	globalMap   = map[int]int{}
	globalMapMu = sync.Mutex{}

	globalRWMap   = map[int]int{}
	globalMapRWMu = sync.RWMutex{}
)

func classicMutex(arg float32) {
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			globalMapMu.Lock()
			defer globalMapMu.Unlock()

			if rand.Float32() < arg/100 {
				globalMap[i] = i * 100
				return
			}
			fmt.Println(globalMap[i])
		}(i)
	}

	wg.Wait()

	fmt.Println("main done", len(globalMap))
}

func rwMutex(arg float32) {
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if rand.Float32() < arg/100 {
				globalMapMu.Lock()
				globalMap[i] = i * 100
				globalMapMu.Unlock()
				return
			}
			globalMapRWMu.RLock()
			fmt.Println(globalRWMap[i])
			globalMapRWMu.RUnlock()
		}(i)
	}

	wg.Wait()

	fmt.Println("main done")
}
