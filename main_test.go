package main

import (
	"fmt"
	"testing"
)

func BenchmarkClassic(b *testing.B) {
	fmt.Println("test with", 90)
	classicMutex(float32(90))
}

func BenchmarkRW(b *testing.B) {
	fmt.Println("test with", 90)
	rwMutex(float32(90))
}

//Classic

//10%
//
// 1000000000               0.00250 ns/op
// PASS
// ok      github.com/MelBogdan/go-level2  0.025s

//50%
//
// 1000000000               0.00218 ns/op
// PASS
// ok      github.com/MelBogdan/go-level2  0.022s

//90%
//
// 1000000000               0.000866 ns/op
// PASS
// ok      github.com/MelBogdan/go-level2  0.016s

//--------------------------------------------------------------------------------

//RW

//10%
//
// 1000000000               0.00682 ns/op
// PASS
// ok      github.com/MelBogdan/go-level2  0.025s

//50%
//
// 1000000000               0.00189 ns/op
// PASS
// ok      github.com/MelBogdan/go-level2  0.014s

//90%
//
// 1000000000               0.000346 ns/op
// PASS
// ok      github.com/MelBogdan/go-level2  0.011s
