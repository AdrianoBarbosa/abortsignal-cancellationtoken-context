// Example 5: a goroutine with no way out.
//
//	go run ./cmd/05-goroutine-leak
//
// Three workers read from a channel that nobody writes to or closes. Without
// ctx.Done() in the select they stay blocked for as long as the process
// lives, and runtime.NumGoroutine shows it. With ctx they all leave when the
// cancellation arrives.
package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	semContexto()
	fmt.Println()
	comContexto()
}

func semContexto() {
	fmt.Println("== sem ctx.Done(): as goroutines ficam presas ==")
	antes := runtime.NumGoroutine()

	ch := make(chan int)
	for range 3 {
		go func() {
			for range ch { // blocks forever: ch is never closed
			}
		}()
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("  goroutines vivas: %d (eram %d)\n", runtime.NumGoroutine(), antes)
}

func comContexto() {
	fmt.Println("== com ctx.Done() no select: todas saem ==")
	antes := runtime.NumGoroutine()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	ch := make(chan int)
	var wg sync.WaitGroup
	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, id, ch)
		}()
	}

	wg.Wait()

	// wg.Done runs right before each goroutine returns, so the count can lag
	// for an instant after Wait. Give the runtime a moment to catch up.
	for i := 0; i < 100 && runtime.NumGoroutine() > antes; i++ {
		time.Sleep(time.Millisecond)
	}
	fmt.Printf("  goroutines vivas: %d (eram %d)\n", runtime.NumGoroutine(), antes)
}

func worker(ctx context.Context, id int, ch <-chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  worker %d: saindo (%v)\n", id, ctx.Err())
			return
		case _, ok := <-ch:
			if !ok {
				return
			}
		}
	}
}
