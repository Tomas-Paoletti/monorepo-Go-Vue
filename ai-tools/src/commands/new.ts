import { existsSync, mkdirSync, writeFileSync } from "node:fs";
import { join, relative } from "node:path";
import {
  DEFAULT_VARIANT,
  VARIANTS,
  isVariant,
  readTemplate,
  slugify,
  specsDir,
  type Variant,
} from "../templates.js";

interface NewOptions {
  name: string;
  variant: Variant;
  force: boolean;
}

function parse(args: string[]): NewOptions {
  let variant: Variant = DEFAULT_VARIANT;
  let force = false;
  const positionals: string[] = [];

  for (const arg of args) {
    if (arg === "--force" || arg === "-f") {
      force = true;
    } else if (arg.startsWith("--variant=")) {
      const value = arg.slice("--variant=".length);
      if (!isVariant(value)) {
        throw new Error(
          `Variante inválida: "${value}". Opciones: ${VARIANTS.join(", ")}`,
        );
      }
      variant = value;
    } else if (arg.startsWith("-")) {
      throw new Error(`Flag desconocido: ${arg}`);
    } else {
      positionals.push(arg);
    }
  }

  const name = positionals.join(" ").trim();
  if (!name) {
    throw new Error(
      'Falta el nombre. Uso: spec new <nombre> [--variant=' +
        VARIANTS.join("|") +
        "] [--force]",
    );
  }

  return { name, variant, force };
}

export function runNew(args: string[]): void {
  const { name, variant, force } = parse(args);
  const slug = slugify(name);
  if (!slug) {
    throw new Error(`El nombre "${name}" no produce un slug válido.`);
  }

  const dir = specsDir();
  mkdirSync(dir, { recursive: true });

  const file = join(dir, `${slug}.md`);
  if (existsSync(file) && !force) {
    throw new Error(
      `Ya existe ${relative(process.cwd(), file)}. Usá --force para sobrescribir.`,
    );
  }

  const content = readTemplate(variant).replace(/<Feature>/g, name);
  writeFileSync(file, content, "utf8");

  console.log(`✓ Spec creado (${variant}): ${relative(process.cwd(), file)}`);
}
