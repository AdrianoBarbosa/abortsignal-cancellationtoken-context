# Exemplos em Node.js

Requer Node 22.18 ou mais novo. Os arquivos são `.ts` e rodam direto, sem build:
o Node remove os tipos antes de executar. Em versões anteriores aparece o erro
`Unknown file extension ".ts"`.

```bash
# servidor de apoio, em outro terminal (01 e 02 precisam dele)
npm run mock

npm run 01    # AbortController e AbortSignal.timeout numa chamada HTTP
npm run 02    # signal repassado x signal esquecido no meio do caminho
npm run 03    # throwIfAborted (polling) e evento abort (callback)
npm run 04    # AbortSignal.any: pai cancelado derruba filho e neto
npm run 05    # once() sem signal deixa listeners pendurados

npm install && npm run typecheck   # opcional, só para os tipos
```

Para apontar para outro servidor: `BASE_URL=http://localhost:9000 npm run 01`.
