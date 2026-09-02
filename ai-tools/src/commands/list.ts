import { existsSync, readFileSync, readdirSync } from "node:fs";
import { join, relative } from "node:path";
import { specsDir } from "../templates.js";

/** Extrae la línea "**Estado:** ..." del contenido de un spec, si existe. */
function readEstado(content: string): string {
  const match = content.match(/^\*\*Estado:\*\*\s*(.+)$/m);
  return match ? match[1].trim() : "—";
}

export function runList(): void {
  const dir = specsDir();
  if (!existsSync(dir)) {
    console.log("No hay specs todavía (no existe docs/specs/).");
    return;
  }

  const files = readdirSync(dir)
    .filter((f) => f.endsWith(".md"))
    .sort();

  if (files.length === 0) {
    console.log("No hay specs todavía en docs/specs/.");
    return;
  }

  console.log(`Specs en ${relative(process.cwd(), dir)}:\n`);
  for (const f of files) {
    const content = readFileSync(join(dir, f), "utf8");
    console.log(`  ${f.padEnd(32)} ${readEstado(content)}`);
  }
}
