// Example 3: how running code finds out it should stop.
//
//   node src/03-avisos.ts
//
// Two ways: checking the signal between steps (polling), or reacting to the
// "abort" event (callback). The callback is how you stop something that
// knows nothing about AbortSignal, like the legacy task below.

import { cronometro, dormir } from "./apoio.ts";

await comPolling();
console.log();
await comEvento();

async function comPolling() {
  console.log("== polling: throwIfAborted entre um lote e outro ==");
  const fim = cronometro();

  const signal = AbortSignal.timeout(250);
  let lotes = 0;
  try {
    while (true) {
      signal.throwIfAborted();
      await dormir(50); // one batch of work, deliberately without the signal
      lotes++;
    }
  } catch (err) {
    console.log(`  parou depois de ${lotes} lotes: ${(err as Error).name}`);
  }
  fim("  tempo");
}

async function comEvento() {
  console.log("== callback: evento abort parando uma tarefa legada ==");
  const fim = cronometro();

  const signal = AbortSignal.timeout(350);

  // A legacy task: it only knows how to be stopped, it takes no signal.
  let ticks = 0;
  const legado = setInterval(() => console.log(`  tarefa legada: tick ${++ticks}`), 100);

  await new Promise<void>((resolve) => {
    signal.addEventListener(
      "abort",
      () => {
        clearInterval(legado);
        console.log("  evento abort: parando a tarefa legada");
        resolve();
      },
      { once: true },
    );
  });
  fim("  tempo");

  const ticksAoParar = ticks;
  await dormir(300);
  console.log(`  ticks depois de parar: ${ticks - ticksAoParar}`);
}
