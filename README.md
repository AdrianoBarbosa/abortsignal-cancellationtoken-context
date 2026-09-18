**Português** · [English](README.en.md)

# `context.Context` vs `CancellationToken`: quem avisa que já pode parar

<!-- TODO: trocar pelo link do artigo quando for publicado -->
Código de apoio do artigo publicado no LinkedIn, continuação de
[Três jeitos de esperar](https://github.com/AdrianoBarbosa/event-loop-async-goroutines).

O mesmo conjunto de exemplos escrito em **Node.js**, **.NET** e **Go**, para comparar
`AbortSignal`, `CancellationToken` e `context.Context` com código que roda de verdade
e mostra, no log do servidor, se o cancelamento chegou do outro lado da conexão.

## Pré-requisitos

| Plataforma | Versão | Por quê |
|---|---|---|
| Go | 1.22+ | `context.AfterFunc` (1.21), rotas com método no `http.ServeMux` e variável de loop por iteração (1.22) |
| Node.js | 22.18+ | executa `.ts` direto, sem passo de build |
| .NET | SDK 8+ | `CancelAsync` (8); o projeto usa `RollForward=Major` e roda nos runtimes mais novos |

Você não precisa das três para começar: cada pasta funciona sozinha.

## Servidor de apoio

Os exemplos 01 e 02 batem em um servidor local que responde depois de um tempo
controlado. Os outros três não usam rede.

```bash
node mock-server/server.js        # ou: cd go && go run ./cmd/00-mock-server
```

| Rota | O que faz |
|---|---|
| `GET /delay/:ms` | responde 200 depois de `:ms` milissegundos |

O log mostra quando cada requisição chegou, quando foi respondida e quando o cliente
desistiu antes da resposta. Nesse último caso o servidor também para o trabalho, em
vez de responder para ninguém. É onde dá para ver o cancelamento atravessando a rede.

Se a porta 8080 estiver ocupada, suba o servidor com `PORT=9000` e rode os exemplos com
`BASE_URL=http://localhost:9000`. Se o servidor não estiver no ar, ou se outro programa
estiver respondendo na porta, os exemplos 01 e 02 avisam e dizem como resolver.

## Mapa dos exemplos

| Tema do artigo | Node.js | .NET | Go |
|---|---|---|---|
| Chamada HTTP cancelada por timeout | [01-timeout.ts](nodejs/src/01-timeout.ts) | [TimeoutHttp.cs](dotnet/Cancelamento/Demos/TimeoutHttp.cs) | [01-timeout](go/cmd/01-timeout/main.go) |
| Cancelamento que some no meio do caminho | [02-propagacao.ts](nodejs/src/02-propagacao.ts) | [Propagacao.cs](dotnet/Cancelamento/Demos/Propagacao.cs) | [02-propagacao](go/cmd/02-propagacao/main.go) |
| Como cada camada é avisada | [03-avisos.ts](nodejs/src/03-avisos.ts) | [Avisos.cs](dotnet/Cancelamento/Demos/Avisos.cs) | [03-avisos](go/cmd/03-avisos/main.go) |
| Cancelamento em cadeia | [04-cascata.ts](nodejs/src/04-cascata.ts) | [Cascata.cs](dotnet/Cancelamento/Demos/Cascata.cs) | [04-cascata](go/cmd/04-cascata/main.go) |
| Espera sem saída | [05-espera-sem-saida.ts](nodejs/src/05-espera-sem-saida.ts) | [EsperaSemSaida.cs](dotnet/Cancelamento/Demos/EsperaSemSaida.cs) | [05-goroutine-leak](go/cmd/05-goroutine-leak/main.go) |
| `defer cancel()` esquecido | não se aplica | não se aplica | `go vet`, veja [abaixo](#o-que-observar) |

Os tempos são menores que os do artigo (300ms em vez de 3s) para os exemplos rodarem rápido.

## Como rodar

Cada bloco parte da raiz do repositório.

```bash
cd nodejs
npm run 01
```

```bash
cd dotnet
dotnet run --project Cancelamento -- 01
```

```bash
cd go
go run ./cmd/01-timeout
```

Os detalhes de cada pasta estão em [nodejs/README.md](nodejs/README.md),
[dotnet/README.md](dotnet/README.md) e [go/README.md](go/README.md).

## O que observar

**Timeout.** Os três desistem em cerca de 300ms, e o servidor registra
`cliente desistiu depois de 300ms`. Muda só a forma do erro: `AbortError` ou
`TimeoutError` no Node (depende de ter usado `AbortController` ou `AbortSignal.timeout`),
`TaskCanceledException` no .NET, e em Go um `*url.Error` que embrulha
`context.DeadlineExceeded`, reconhecido por `errors.Is` através do embrulho.

**Propagação.** O exemplo mais revelador. Três camadas, e a do meio "esquece" de
repassar o cancelamento:

```
Go    ctx repassado                       volta em 301ms
Go    context.Background() no meio        volta em 2011ms, erro nil
.NET  token repassado                     volta em 305ms
.NET  CancellationToken.None no meio      volta em 2011ms, status 200
Node  signal repassado                    volta em 307ms
Node  signal esquecido no meio            volta em 2019ms, status 200
```

Quem chamou já tinha desistido em 300ms, mas ficou preso esperando 2 segundos e ainda
recebeu um "sucesso". Nenhuma das três plataformas avisa em tempo de execução. A única
pista está no log do servidor:

```
[   568ms] #1 chegou     GET /delay/2000
[   871ms] #1 cliente desistiu depois de 303ms
[   872ms] #2 chegou     GET /delay/2000
[  2882ms] #2 respondeu  200
```

**Como cada camada é avisada.** Por polling, o código só percebe no próximo ponto de
checagem: 5 lotes de 50ms, parada em cerca de 250ms. Por callback (`addEventListener`,
`Register`, `context.AfterFunc`), a reação é imediata e serve para parar o que não
conhece o mecanismo de cancelamento: um timer legado no Node e no .NET, e em Go uma
leitura bloqueada numa conexão, que só volta quando o callback fecha a conexão.

**Cascata.** O filho tem timeout próprio de 5s, mas termina em 0 a 3ms quando o pai é
cancelado, e o neto vai junto. Em Go o vínculo é automático para qualquer contexto
derivado; no .NET e no Node precisa ser declarado com `CreateLinkedTokenSource` e
`AbortSignal.any`.

**Espera sem saída.** Três leitores esperando algo que nunca vai chegar:

```
Go    sem ctx.Done()          goroutines vivas: 4 (eram 1)
Go    com ctx.Done()          goroutines vivas: 4 (eram 4)
.NET  ReadAsync sem token     tarefas ainda esperando: 3 de 3
.NET  ReadAsync com token     tarefas ainda esperando: 0 de 3
Node  once() sem signal       listeners ainda registrados: 3
Node  once() com signal       listeners ainda registrados: 0
```

Na segunda linha de Go, o "eram 4" são as três goroutines da primeira parte, que
continuam presas até o processo terminar. No Node o vazamento é mais discreto: uma
Promise pendente não segura o processo, mas os listeners ficam pendurados no emitter.
E uma pegadinha que apareceu escrevendo o exemplo: o timer do `AbortSignal.timeout`
não mantém o processo vivo. Sem mais nada pendente, o Node encerra antes de ele disparar.

**`defer cancel()` esquecido.** Em [go/cmd/01-timeout/main.go](go/cmd/01-timeout/main.go),
troque `ctx, cancel :=` por `ctx, _ :=`, apague a linha `defer cancel()` e rode
`go vet ./...`:

```
the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
```

Apagar só o `defer` não basta para ver isso: o compilador reclama antes, com
`declared and not used`.

## Estrutura

```
.
├── mock-server/server.js           servidor de apoio, sem dependencias
├── nodejs/src/                     exemplos em TypeScript, rodam sem build
├── dotnet/Cancelamento/Demos/      um exemplo por arquivo, dotnet run -- NN
└── go/cmd/                         um exemplo por pasta, go run ./cmd/NN-nome
```
