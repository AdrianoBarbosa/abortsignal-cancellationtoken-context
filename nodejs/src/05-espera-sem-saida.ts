// Example 5: a wait with no way out, the Node version of a goroutine leak.
//
//   node src/05-espera-sem-saida.ts
//
// Three readers wait for an event that nobody will ever emit. A pending
// Promise does not keep the process alive, so the leak here is quieter than
// in Go: the listeners stay attached to the emitter for as long as it lives.
// With a signal, once() removes them when the cancellation arrives.

import { EventEmitter, once } from "node:events";
import { cronometro, dormir } from "./apoio.ts";

await semSignal();
console.log();
await comSignal();

async function semSignal() {
  console.log("== once() sem signal: os listeners ficam pendurados ==");
  const fila = new EventEmitter();

  for (let i = 0; i < 3; i++) {
    void once(fila, "item"); // nobody will ever emit "item"
  }

  await dormir(300);
  console.log("  listeners ainda registrados:", fila.listenerCount("item"));
}

async function comSignal() {
  console.log("== once() com signal: todos saem ==");
  const fim = cronometro();
  const fila = new EventEmitter();

  // Not AbortSignal.timeout(300): its timer does not keep the process alive,
  // and with nothing else pending Node would exit before it ever fired.
  const controller = new AbortController();
  setTimeout(() => controller.abort(), 300);
  const signal = controller.signal;

  await Promise.all(
    [1, 2, 3].map(async (id) => {
      try {
        await once(fila, "item", { signal });
      } catch {
        console.log(`  leitor ${id}: saindo (${(signal.reason as Error).name})`);
      }
    }),
  );

  console.log("  listeners ainda registrados:", fila.listenerCount("item"));
  fim("  tempo");
}
