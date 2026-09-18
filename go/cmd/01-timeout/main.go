// Example 1: an HTTP call canceled by a timeout.
//
//	go run ./cmd/01-timeout
//
// The server takes 2s to answer and the client gives up at 300ms. Watch the
// server log: the request shows up as "cliente desistiu", because canceling
// the context closes the connection.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AdrianoBarbosa/abortsignal-cancellationtoken-context/go/internal/apoio"
)

func main() {
	apoio.ExigirServidor()

	fmt.Println("== context.WithTimeout ==")
	c := apoio.NovoCronometro()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel() // go vet's lostcancel check flags a cancel that is discarded and never called

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apoio.BaseURL()+"/delay/2000", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Do wraps the cause in a *url.Error; errors.Is sees through the wrapping.
		fmt.Println("  erro devolvido:", err)
		fmt.Println("  errors.Is(err, context.DeadlineExceeded):", errors.Is(err, context.DeadlineExceeded))
		c.Fim("  tempo ate desistir")
		return
	}
	defer resp.Body.Close()

	fmt.Println("  status:", resp.StatusCode)
	c.Fim("  tempo")
}
