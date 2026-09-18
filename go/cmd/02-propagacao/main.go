// Example 2: cancellation that vanishes halfway down the call chain.
//
//	go run ./cmd/02-propagacao
//
// caller -> servico -> repositorio -> HTTP. The caller gives up at 300ms.
// In the first run every layer hands ctx along and the HTTP call is really
// interrupted. In the second run servico swaps it for context.Background(),
// and the call runs for the full 2s even though nobody is waiting anymore.
// Nothing warns you at runtime: the only clue is the server log.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AdrianoBarbosa/abortsignal-cancellationtoken-context/go/internal/apoio"
)

func main() {
	apoio.ExigirServidor()

	rodar("ctx repassado em todas as camadas", servicoCorreto)
	fmt.Println()
	rodar("servico troca o ctx por context.Background()", servicoQuebrado)
}

func rodar(titulo string, servico func(context.Context) error) {
	fmt.Println("==", titulo, "==")
	c := apoio.NovoCronometro()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err := servico(ctx)
	fmt.Println("  erro devolvido:", err)
	fmt.Println("  ctx.Err() de quem chamou:", ctx.Err())
	c.Fim("  tempo ate a chamada voltar")
}

func servicoCorreto(ctx context.Context) error {
	return repositorio(ctx)
}

// The contextcheck linter flags this: a fresh context instead of the one
// received. The compiler and go vet are fine with it.
func servicoQuebrado(ctx context.Context) error {
	return repositorio(context.Background())
}

func repositorio(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apoio.BaseURL()+"/delay/2000", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}
