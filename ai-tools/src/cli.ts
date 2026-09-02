import { VARIANTS } from "./templates.js";
import { runNew } from "./commands/new.js";
import { runList } from "./commands/list.js";

function usage(): void {
  console.log(`spec — scaffolding de specs (spec-driven development)

Uso:
  spec new <nombre> [--variant=${VARIANTS.join("|")}] [--force]
  spec list

Ejemplos:
  npm run spec -- new checkout
  npm run spec -- new pricing --variant=bdd
  npm run spec -- list

Default variant: ${VARIANTS[0]}`);
}

function main(): void {
  const [command, ...rest] = process.argv.slice(2);

  switch (command) {
    case "new":
      runNew(rest);
      break;
    case "list":
      runList();
      break;
    case undefined:
    case "help":
    case "--help":
    case "-h":
      usage();
      break;
    default:
      console.error(`Comando desconocido: ${command}\n`);
      usage();
      process.exitCode = 1;
  }
}

try {
  main();
} catch (err) {
  console.error(`✗ ${err instanceof Error ? err.message : String(err)}`);
  process.exitCode = 1;
}
