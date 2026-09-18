// Package apoio holds what every example shares: the address of the
// support server, a check that it is running and a simple stopwatch.
package apoio

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// BaseURL points to the repository's support server.
// Start it before running the network examples:
//
//	go run ./cmd/00-mock-server
//	node mock-server/server.js
func BaseURL() string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "http://localhost:8080"
}

// ExigirServidor stops the example with a clear message when the support
// server is not running, instead of a raw "connection refused". It also
// catches another program answering on the same port.
func ExigirServidor() {
	c := http.Client{Timeout: time.Second}
	resp, err := c.Get(BaseURL() + "/")
	if err != nil {
		pararSemServidor(fmt.Sprintf("O servidor de apoio nao respondeu em %s.", BaseURL()))
	}
	defer resp.Body.Close()

	inicio := make([]byte, 64)
	n, _ := io.ReadFull(resp.Body, inicio)
	if !strings.HasPrefix(string(inicio[:n]), "Servidor de apoio") {
		pararSemServidor(fmt.Sprintf("Outro programa esta respondendo em %s, nao o servidor de apoio.", BaseURL()))
	}
}

func pararSemServidor(motivo string) {
	fmt.Fprintln(os.Stderr, motivo)
	fmt.Fprintln(os.Stderr, "Suba o servidor de apoio antes, em outro terminal:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  go run ./cmd/00-mock-server")
	fmt.Fprintln(os.Stderr, "  (ou, na raiz do repositorio: node mock-server/server.js)")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Para usar outra porta: PORT=9000 no servidor e BASE_URL=http://localhost:9000 no exemplo.")
	os.Exit(1)
}

// Cronometro prints how much time has passed since it was created.
type Cronometro struct{ inicio time.Time }

func NovoCronometro() Cronometro { return Cronometro{inicio: time.Now()} }

func (c Cronometro) Fim(rotulo string) {
	fmt.Printf("%-32s %dms\n", rotulo, time.Since(c.inicio).Milliseconds())
}
