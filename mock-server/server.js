// Support server for the article examples.
// No dependencies: node mock-server/server.js
//
// Routes:
//   GET /delay/:ms  -> answers 200 after :ms milliseconds
//   GET /           -> help (also used by the examples to check the server is up)
//
// The log shows when each request arrived, when it was answered and when the
// client gave up before the answer. That last line is the point of the
// article: cancellation reaching the other side of the connection.

import http from "node:http";

const PORT = Number(process.env.PORT ?? 8080);
const inicio = Date.now();
let sequencia = 0;

const agora = () => String(Date.now() - inicio).padStart(6, " ");

const servidor = http.createServer((req, res) => {
  const url = new URL(req.url ?? "/", "http://localhost");
  const partes = url.pathname.split("/").filter(Boolean);

  if (partes[0] !== "delay") {
    res.writeHead(200, { "content-type": "text/plain; charset=utf-8" });
    res.end(
      [
        "Servidor de apoio dos exemplos.",
        "",
        "GET /delay/:ms  responde 200 depois de :ms",
        "",
      ].join("\n"),
    );
    return;
  }

  const id = ++sequencia;
  const chegada = Date.now();
  let espera = Number(partes[1] ?? 0);
  if (!Number.isFinite(espera) || espera < 0) espera = 0;

  console.log(`[${agora()}ms] #${id} chegou     ${req.method} ${req.url}`);

  const timer = setTimeout(() => {
    console.log(`[${agora()}ms] #${id} respondeu  200`);
    res.writeHead(200, { "content-type": "application/json; charset=utf-8" });
    res.end(JSON.stringify({ id, esperaMs: espera }));
  }, espera);

  // The connection closed before the answer went out: the client canceled.
  // Stop the pending work instead of answering nobody.
  res.on("close", () => {
    if (res.writableFinished) return;
    clearTimeout(timer);
    console.log(`[${agora()}ms] #${id} cliente desistiu depois de ${Date.now() - chegada}ms`);
  });
});

servidor.on("error", (erro) => {
  if (erro.code !== "EADDRINUSE") throw erro;
  console.error(`A porta ${PORT} ja esta em uso. Suba o servidor em outra porta:`);
  console.error();
  console.error(`  PORT=9000 node mock-server/server.js`);
  console.error();
  console.error("e aponte os exemplos para ela com BASE_URL=http://localhost:9000");
  process.exit(1);
});

servidor.listen(PORT, () => {
  console.log(`Servidor de apoio ouvindo em http://localhost:${PORT}`);
});
