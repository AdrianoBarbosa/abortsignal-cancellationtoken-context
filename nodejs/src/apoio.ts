// What every example shares: the support server address, a check that it
// is running and a simple stopwatch.
//
// Start the server before running the network examples:
//   node mock-server/server.js
//   (or npm run mock, from inside nodejs/)

export const BASE_URL = process.env.BASE_URL ?? "http://localhost:8080";

// Stops the example with a clear message when the support server is not
// running, instead of a raw "fetch failed". It also catches another program
// answering on the same port.
export async function exigirServidor(): Promise<void> {
  let texto: string;
  try {
    const res = await fetch(`${BASE_URL}/`, { signal: AbortSignal.timeout(1000) });
    texto = await res.text();
  } catch {
    return pararSemServidor(`O servidor de apoio nao respondeu em ${BASE_URL}.`);
  }
  if (!texto.startsWith("Servidor de apoio")) {
    pararSemServidor(`Outro programa esta respondendo em ${BASE_URL}, nao o servidor de apoio.`);
  }
}

function pararSemServidor(motivo: string): never {
  console.error(motivo);
  console.error("Suba o servidor de apoio antes, em outro terminal:");
  console.error();
  console.error("  npm run mock");
  console.error("  (ou, na raiz do repositorio: node mock-server/server.js)");
  console.error();
  console.error("Para usar outra porta: PORT=9000 no servidor e BASE_URL=http://localhost:9000 no exemplo.");
  process.exit(1);
}

export function cronometro() {
  const inicio = performance.now();
  return (rotulo: string) => {
    const ms = Math.round(performance.now() - inicio);
    console.log(`${rotulo.padEnd(32)} ${ms}ms`);
  };
}

export const dormir = (ms: number) => new Promise((r) => setTimeout(r, ms));
