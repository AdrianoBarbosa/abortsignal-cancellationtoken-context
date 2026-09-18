// Example 3: how running code finds out it should stop.
//
//	go run ./cmd/03-avisos
//
// The idiomatic way is a select on ctx.Done(). Since Go 1.21 there is also
// context.AfterFunc, a callback for when there is no select to put it in,
// such as a blocking read that only returns when the connection is closed.
package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/AdrianoBarbosa/abortsignal-cancellationtoken-context/go/internal/apoio"
)

func main() {
	comSelect()
	fmt.Println()
	comAfterFunc()
}

func comSelect() {
	fmt.Println("== select em ctx.Done() ==")
	c := apoio.NovoCronometro()

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	// Buffered on purpose: if nobody reads the result, the goroutine can still
	// send it and finish. With an unbuffered channel it would leak.
	resultado := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		resultado <- "pronto"
	}()

	select {
	case <-ctx.Done():
		fmt.Println("  desistiu:", ctx.Err())
	case r := <-resultado:
		fmt.Println("  resultado:", r)
	}
	c.Fim("  tempo")
}

func comAfterFunc() {
	fmt.Println("== context.AfterFunc fechando uma conexao presa ==")
	c := apoio.NovoCronometro()

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	cliente, servidor := net.Pipe()
	defer servidor.Close()

	// Read has no ctx parameter: it only returns when data arrives or the
	// connection is closed. The callback turns cancellation into a Close.
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("  AfterFunc: contexto terminou, fechando a conexao")
		cliente.Close()
	})
	defer stop()

	buf := make([]byte, 1)
	_, err := cliente.Read(buf)
	fmt.Println("  Read voltou com:", err)
	c.Fim("  tempo")
}
