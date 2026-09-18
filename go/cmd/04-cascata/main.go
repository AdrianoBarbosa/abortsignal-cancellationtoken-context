// Example 4: canceling the parent cancels every child.
//
//	go run ./cmd/04-cascata
//
// Any context derived from another one is linked to it automatically. The
// child below has its own 5s timeout, but it ends the moment the parent is
// canceled. .NET needs CreateLinkedTokenSource and Node needs AbortSignal.any
// to get the same thing.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AdrianoBarbosa/abortsignal-cancellationtoken-context/go/internal/apoio"
)

func main() {
	fmt.Println("== pai cancelado, filho com timeout de 5s ==")
	c := apoio.NovoCronometro()

	parent, cancelParent := context.WithCancel(context.Background())
	child, cancelChild := context.WithTimeout(parent, 5*time.Second)
	defer cancelChild()

	neto, cancelNeto := context.WithCancel(child)
	defer cancelNeto()

	cancelParent()

	<-child.Done()
	<-neto.Done()
	fmt.Println("  child.Err():", child.Err())
	fmt.Println("  neto.Err(): ", neto.Err())
	c.Fim("  tempo (o timeout era 5s)")
}
