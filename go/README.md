# Exemplos em Go

Requer Go 1.22 ou mais novo (os exemplos usam `context.AfterFunc` e rotas com método no `http.ServeMux`).

```bash
# servidor de apoio, em outro terminal (01 e 02 precisam dele)
go run ./cmd/00-mock-server

go run ./cmd/01-timeout           # context.WithTimeout numa chamada HTTP
go run ./cmd/02-propagacao        # ctx repassado x context.Background() no meio do caminho
go run ./cmd/03-avisos            # select em ctx.Done() e context.AfterFunc
go run ./cmd/04-cascata           # pai cancelado derruba filho e neto
go run ./cmd/05-goroutine-leak    # goroutine presa sem ctx.Done()

go vet ./...                      # inclui o check lostcancel
```

Para apontar para outro servidor: `BASE_URL=http://localhost:9000 go run ./cmd/01-timeout`.

O módulo é `github.com/AdrianoBarbosa/abortsignal-cancellationtoken-context/go`. Se o
repositório for publicado com outro nome, ajuste com `go mod edit -module ...` e
atualize os imports de `internal/apoio`.
