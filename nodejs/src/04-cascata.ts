// Example 4: canceling the parent cancels the child.
//
//   node src/04-cascata.ts
//
// In Go every derived context is linked to its parent automatically. In Node
// the link has to be declared: AbortSignal.any combines the parent signal
// with the child's own 5s timeout, and whichever fires first wins.

import { cronometro } from "./apoio.ts";

console.log("== pai cancelado, filho com timeout de 5s ==");
const fim = cronometro();

const pai = new AbortController();
const filho = AbortSignal.any([pai.signal, AbortSignal.timeout(5000)]);
const neto = AbortSignal.any([filho]);

filho.addEventListener("abort", () => console.log("  evento abort no filho"));
neto.addEventListener("abort", () => console.log("  evento abort no neto"));

pai.abort();

console.log("  filho.aborted:", filho.aborted, `(${(filho.reason as Error).name})`);
console.log("  neto.aborted: ", neto.aborted, `(${(neto.reason as Error).name})`);
fim("  tempo (o timeout era 5s)");
