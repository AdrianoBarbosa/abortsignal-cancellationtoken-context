# Exemplos em .NET

Requer o SDK do .NET 8 ou mais novo. O projeto mira o .NET 8 com `RollForward=Major`,
então também roda em quem só tem runtimes mais novos.

```bash
# servidor de apoio, em outro terminal (01 e 02 precisam dele)
node ../mock-server/server.js

dotnet run --project Cancelamento            # lista os exemplos
dotnet run --project Cancelamento -- 01      # CancellationTokenSource com timeout numa chamada HTTP
dotnet run --project Cancelamento -- 02      # token repassado x CancellationToken.None no meio do caminho
dotnet run --project Cancelamento -- 03      # ThrowIfCancellationRequested (polling) e Register (callback)
dotnet run --project Cancelamento -- 04      # CreateLinkedTokenSource: pai cancelado derruba filho e neto
dotnet run --project Cancelamento -- 05      # ReadAsync de Channel sem token fica esperando para sempre
```

Para apontar para outro servidor: `BASE_URL=http://localhost:9000 dotnet run --project Cancelamento -- 01`.
