// Example 1: an HTTP call canceled by a timeout.
//
//   node src/01-timeout.ts
//
// The server takes 2s to answer and the client gives up at 300ms. Watch the
// server log: the request shows up as "cliente desistiu", because aborting
// the fetch closes the connection.

import { BASE_URL, cronometro, exigirServidor } from "./apoio.ts";

await exigirServidor();
await comAbortController();
console.log();
await comAbortSignalTimeout();

// The general tool: a controller that cancels on demand. Here the "demand"
// is a setTimeout, which you have to clear yourself afterwards.
async function comAbortController() {
  console.log("== AbortController + setTimeout ==");
  const fim = cronometro();

  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 300);

  try {
    const res = await fetch(`${BASE_URL}/delay/2000`, { signal: controller.signal });
    console.log("  status:", res.status);
  } catch (err) {
    console.log("  desistiu:", (err as Error).name);
  } finally {
    clearTimeout(timeout);
  }
  fim("  tempo ate desistir");
}

// The shortcut for the common case. Note the different error name:
// TimeoutError, not AbortError.
async function comAbortSignalTimeout() {
  console.log("== AbortSignal.timeout ==");
  const fim = cronometro();

  try {
    const res = await fetch(`${BASE_URL}/delay/2000`, { signal: AbortSignal.timeout(300) });
    console.log("  status:", res.status);
  } catch (err) {
    console.log("  desistiu:", (err as Error).name);
  }
  fim("  tempo ate desistir");
}
