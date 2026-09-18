[Português](README.md) · **English**

# `context.Context` vs `CancellationToken`: who says it is time to stop

<!-- TODO: replace with the article link once it is published -->
Companion code for the article published on LinkedIn (in Portuguese), a follow-up to
[Three ways of waiting](https://github.com/AdrianoBarbosa/event-loop-async-goroutines).

The same set of examples written in **Node.js**, **.NET** and **Go**, to compare
`AbortSignal`, `CancellationToken` and `context.Context` with code that actually runs
and shows, in the server log, whether the cancellation reached the other side of the
connection.

## Requirements

| Platform | Version | Why |
|---|---|---|
| Go | 1.22+ | `context.AfterFunc` (1.21), method-based routes in `http.ServeMux` and per-iteration loop variables (1.22) |
| Node.js | 22.18+ | runs `.ts` directly, with no build step |
| .NET | SDK 8+ | `CancelAsync` (8); the project uses `RollForward=Major` and runs on newer runtimes |

You do not need all three to get started: each folder works on its own.

## Support server

Examples 01 and 02 hit a local server that answers after a controlled delay. The
other three do not use the network.

```bash
node mock-server/server.js        # or: cd go && go run ./cmd/00-mock-server
```

| Route | What it does |
|---|---|
| `GET /delay/:ms` | answers 200 after `:ms` milliseconds |

The log shows when each request arrived, when it was answered and when the client
gave up before the answer. In that last case the server also stops working, instead
of answering nobody. That is where you can watch cancellation crossing the network.

If port 8080 is taken, start the server with `PORT=9000` and run the examples with
`BASE_URL=http://localhost:9000`. If the server is not up, or another program is
answering on the port, examples 01 and 02 say so and tell you how to fix it.

## Map of the examples

| Topic from the article | Node.js | .NET | Go |
|---|---|---|---|
| HTTP call canceled by a timeout | [01-timeout.ts](nodejs/src/01-timeout.ts) | [TimeoutHttp.cs](dotnet/Cancelamento/Demos/TimeoutHttp.cs) | [01-timeout](go/cmd/01-timeout/main.go) |
| Cancellation that vanishes halfway down | [02-propagacao.ts](nodejs/src/02-propagacao.ts) | [Propagacao.cs](dotnet/Cancelamento/Demos/Propagacao.cs) | [02-propagacao](go/cmd/02-propagacao/main.go) |
| How each layer gets notified | [03-avisos.ts](nodejs/src/03-avisos.ts) | [Avisos.cs](dotnet/Cancelamento/Demos/Avisos.cs) | [03-avisos](go/cmd/03-avisos/main.go) |
| Chained cancellation | [04-cascata.ts](nodejs/src/04-cascata.ts) | [Cascata.cs](dotnet/Cancelamento/Demos/Cascata.cs) | [04-cascata](go/cmd/04-cascata/main.go) |
| A wait with no way out | [05-espera-sem-saida.ts](nodejs/src/05-espera-sem-saida.ts) | [EsperaSemSaida.cs](dotnet/Cancelamento/Demos/EsperaSemSaida.cs) | [05-goroutine-leak](go/cmd/05-goroutine-leak/main.go) |
| Forgotten `defer cancel()` | not applicable | not applicable | `go vet`, see [below](#what-to-look-for) |

Timings are shorter than in the article (300ms instead of 3s) so the examples run fast.

## Running them

Each block starts from the repository root.

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

The details of each folder live in [nodejs/README.md](nodejs/README.md),
[dotnet/README.md](dotnet/README.md) and [go/README.md](go/README.md), in Portuguese.

## What to look for

**Timeout.** All three give up at around 300ms, and the server logs
`cliente desistiu depois de 300ms` (client gave up). Only the shape of the error
changes: `AbortError` or `TimeoutError` in Node (depending on `AbortController` or
`AbortSignal.timeout`), `TaskCanceledException` in .NET, and in Go a `*url.Error`
wrapping `context.DeadlineExceeded`, which `errors.Is` sees through the wrapping.

**Propagation.** The most revealing example. Three layers, and the middle one
"forgets" to pass the cancellation on:

```
Go    ctx passed along                      returns in 301ms
Go    context.Background() in the middle    returns in 2011ms, nil error
.NET  token passed along                    returns in 305ms
.NET  CancellationToken.None in the middle  returns in 2011ms, status 200
Node  signal passed along                   returns in 307ms
Node  signal forgotten in the middle        returns in 2019ms, status 200
```

The caller had already given up at 300ms, but stayed stuck for 2 seconds and still
got a "success". None of the three platforms warns you at runtime. The only clue is
the server log:

```
[   568ms] #1 chegou     GET /delay/2000
[   871ms] #1 cliente desistiu depois de 303ms
[   872ms] #2 chegou     GET /delay/2000
[  2882ms] #2 respondeu  200
```

**How each layer gets notified.** With polling, the code only notices at the next
check: 5 batches of 50ms, stopping at around 250ms. With a callback
(`addEventListener`, `Register`, `context.AfterFunc`), the reaction is immediate and
works for stopping things that know nothing about cancellation: a legacy timer in
Node and .NET, and in Go a read blocked on a connection, which only returns when the
callback closes it.

**Cascade.** The child has its own 5s timeout, but it ends in 0 to 3ms when the
parent is canceled, and the grandchild goes with it. In Go the link is automatic for
any derived context; in .NET and Node it has to be declared with
`CreateLinkedTokenSource` and `AbortSignal.any`.

**A wait with no way out.** Three readers waiting for something that will never
arrive:

```
Go    without ctx.Done()      live goroutines: 4 (were 1)
Go    with ctx.Done()         live goroutines: 4 (were 4)
.NET  ReadAsync, no token     tasks still waiting: 3 of 3
.NET  ReadAsync with token    tasks still waiting: 0 of 3
Node  once() without signal   listeners still attached: 3
Node  once() with signal      listeners still attached: 0
```

In the second Go line, the "were 4" are the three goroutines from the first part,
which stay stuck until the process ends. In Node the leak is quieter: a pending
Promise does not keep the process alive, but the listeners stay attached to the
emitter. And a gotcha found while writing the example: the `AbortSignal.timeout`
timer does not keep the process alive. With nothing else pending, Node exits before
it fires.

**Forgotten `defer cancel()`.** In [go/cmd/01-timeout/main.go](go/cmd/01-timeout/main.go),
replace `ctx, cancel :=` with `ctx, _ :=`, delete the `defer cancel()` line and run
`go vet ./...`:

```
the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak
```

Deleting only the `defer` is not enough to see this: the compiler complains first,
with `declared and not used`.

## Layout

```
.
├── mock-server/server.js           support server, no dependencies
├── nodejs/src/                     TypeScript examples, run without a build
├── dotnet/Cancelamento/Demos/      one example per file, dotnet run -- NN
└── go/cmd/                         one example per folder, go run ./cmd/NN-name
```

---

Translation of [README.md](README.md), which is the source of truth. Code comments
are in English; the folder READMEs and the console output are in Portuguese.
