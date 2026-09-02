import { test } from "node:test";
import assert from "node:assert/strict";
import { nextSpecNumber } from "./templates.js";

test("nextSpecNumber: siguiente al máximo, padding 3", () => {
  assert.equal(nextSpecNumber(["001-a.md", "007-b.md", "003-c.md"]), "008");
});

test("nextSpecNumber: 001 cuando no hay specs", () => {
  assert.equal(nextSpecNumber([]), "001");
});

test("nextSpecNumber: ignora archivos sin prefijo numérico", () => {
  assert.equal(nextSpecNumber(["README.md", "005-x.md"]), "006");
});
