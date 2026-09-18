// Support server, the Go version of mock-server/server.js.
//
//	go run ./cmd/00-mock-server
//
// Routes:
//
//	GET /delay/{ms}  answers 200 after ms milliseconds
//	GET /            help (also used by the examples to check the server is up)
//
// The handler waits on r.Context().Done(): when the client gives up, the
// request context is canceled and the server stops working for nobody.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	inicio    = time.Now()
	sequencia atomic.Int64
)

func registrar(formato string, args ...any) {
	fmt.Printf("[%6dms] %s\n", time.Since(inicio).Milliseconds(), fmt.Sprintf(formato, args...))
}

func main() {
	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Servidor de apoio dos exemplos.")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "GET /delay/:ms  responde 200 depois de :ms")
	})

	mux.HandleFunc("GET /delay/{ms}", func(w http.ResponseWriter, r *http.Request) {
		espera, err := strconv.Atoi(r.PathValue("ms"))
		if err != nil || espera < 0 {
			espera = 0
		}
		id := sequencia.Add(1)
		chegada := time.Now()
		registrar("#%d chegou     %s %s", id, r.Method, r.URL)

		t := time.NewTimer(time.Duration(espera) * time.Millisecond)
		defer t.Stop()

		select {
		case <-t.C:
			registrar("#%d respondeu  200", id)
			w.Header().Set("content-type", "application/json; charset=utf-8")
			fmt.Fprintf(w, `{"id":%d,"esperaMs":%d}`, id, espera)
		case <-r.Context().Done():
			registrar("#%d cliente desistiu depois de %dms", id, time.Since(chegada).Milliseconds())
		}
	})

	ln, err := net.Listen("tcp", ":"+porta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Nao foi possivel ouvir na porta %s: %v\n", porta, err)
		fmt.Fprintln(os.Stderr, "Se ela estiver em uso, suba o servidor em outra porta:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  PORT=9000 go run ./cmd/00-mock-server")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "e aponte os exemplos para ela com BASE_URL=http://localhost:9000")
		os.Exit(1)
	}

	fmt.Printf("Servidor de apoio ouvindo em http://localhost:%s\n", porta)
	log.Fatal(http.Serve(ln, mux))
}
