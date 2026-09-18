// Example 2: cancellation that vanishes halfway down the call chain.
//
//   node src/02-propagacao.ts
//
// caller -> servico -> repositorio -> fetch. The caller gives up at 300ms.
// In the first run every layer hands the signal along and the fetch is really
// aborted. In the second run servico simply does not pass it on, and the
// request runs for the full 2s even though nobody is waiting anymore.
// Nothing warns you at runtime: the only clue is the server log.

import { BASE_URL, cronometro, exigirServidor } from "./apoio.ts";

await exigirServidor();
await rodar("signal repassado em todas as camadas", servicoCorreto);
console.log();
await rodar("servico nao repassa o signal", servicoQuebrado);

async function rodar(titulo: string, servico: (signal: AbortSignal) => Promise<number>) {
  console.log(`== ${titulo} ==`);
  const fim = cronometro();

  const signal = AbortSignal.timeout(300);
  try {
    const status = await servico(signal);
    console.log("  resultado: status", status);
  } catch (err) {
    console.log("  resultado:", (err as Error).name);
  }
  console.log("  signal.aborted de quem chamou:", signal.aborted);
  fim("  tempo ate a chamada voltar");
}

function servicoCorreto(signal: AbortSignal) {
  return repositorio(signal);
}

// signal is optional in repositorio, so forgetting it compiles fine.
function servicoQuebrado(_signal: AbortSignal) {
  return repositorio();
}

async function repositorio(signal?: AbortSignal) {
  const res = await fetch(`${BASE_URL}/delay/2000`, { signal });
  await res.text();
  return res.status;
}
