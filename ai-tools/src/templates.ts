import { existsSync, readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

/** Directorio de este archivo (ai-tools/src). */
const HERE = dirname(fileURLToPath(import.meta.url));

/** Raíz de ai-tools (un nivel arriba de src/). */
export const AI_TOOLS_DIR = resolve(HERE, "..");

/** Carpeta con las plantillas de spec. */
export const TEMPLATES_DIR = join(AI_TOOLS_DIR, "templates", "spec");

/** Variantes disponibles. La primera es el default. */
export const VARIANTS = [
  "project",
  "bdd",
  "contract",
  "test-first",
  "full",
] as const;

export type Variant = (typeof VARIANTS)[number];

export const DEFAULT_VARIANT: Variant = "project";

export function isVariant(value: string): value is Variant {
  return (VARIANTS as readonly string[]).includes(value);
}

/** Lee el contenido de una plantilla. */
export function readTemplate(variant: Variant): string {
  const file = join(TEMPLATES_DIR, `${variant}.md`);
  if (!existsSync(file)) {
    throw new Error(`No existe la plantilla "${variant}" en ${file}`);
  }
  return readFileSync(file, "utf8");
}

/**
 * Sube desde `start` hasta encontrar un directorio con `.git`.
 * Devuelve la raíz del repo, o `start` si no encuentra ninguno.
 */
export function findRepoRoot(start: string = AI_TOOLS_DIR): string {
  let dir = resolve(start);
  while (true) {
    if (existsSync(join(dir, ".git"))) return dir;
    const parent = dirname(dir);
    if (parent === dir) return resolve(start); // llegamos a la raíz del FS
    dir = parent;
  }
}

/** Carpeta donde viven los specs generados. */
export function specsDir(): string {
  return join(findRepoRoot(), "docs", "specs");
}

/** Próximo número de spec (padding 3) mirando los `NNN-` ya existentes. */
export function nextSpecNumber(existing: string[]): string {
  const max = existing.reduce((acc, name) => {
    const m = name.match(/^(\d+)-/);
    return m ? Math.max(acc, parseInt(m[1], 10)) : acc;
  }, 0);
  return String(max + 1).padStart(3, "0");
}

/** Convierte un nombre libre en un slug apto para nombre de archivo. */
export function slugify(name: string): string {
  return name
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "") // saca acentos
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}
