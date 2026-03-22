import { describe, it } from "node:test";
import assert from "node:assert/strict";
import { mkdirSync, writeFileSync, rmSync } from "node:fs";
import { join, dirname } from "node:path";
import { tmpdir } from "node:os";

import { checkFileNaming } from "../validations/file-naming";
import { checkRequiredFiles } from "../validations/required-files";
import { checkInternalLinks } from "../validations/internal-links";

const makeKit = (files: Record<string, string>): string => {
  const dir = join(tmpdir(), `kit-validate-test-${Date.now()}`);
  for (const [rel, content] of Object.entries(files)) {
    const full = join(dir, rel);
    mkdirSync(dirname(full), { recursive: true });
    writeFileSync(full, content);
  }
  return dir;
};

describe("checkFileNaming", () => {
  it("passes for lowercase-hyphen filenames", () => {
    const dir = makeKit({ "start-here.md": "", "README.md": "" });
    assert.deepEqual(checkFileNaming(dir), []);
    rmSync(dir, { recursive: true });
  });

  it("flags uppercase filenames", () => {
    const dir = makeKit({ "GUARDRAIL_ADR.md": "" });
    const issues = checkFileNaming(dir);
    assert.equal(issues.length, 1);
    assert.ok(issues[0] !== undefined && issues[0].message.includes("lowercase-hyphen"));
    rmSync(dir, { recursive: true });
  });
});

describe("checkRequiredFiles", () => {
  it("flags missing required files", () => {
    const dir = makeKit({ "start-here.md": "" });
    const issues = checkRequiredFiles(dir);
    assert.ok(issues.length > 0);
    assert.ok(issues.every((i) => i.message === "required kit file is missing"));
    rmSync(dir, { recursive: true });
  });

  it("passes when all required files exist", () => {
    const dir = makeKit({
      "start-here.md": "",
      "usage/status-report.md": "",
      "usage/intake-workflow.md": "",
      "usage/init.md": "",
      "references/spec.md": "",
      "references/backlog-alignment.md": "",
    });
    assert.deepEqual(checkRequiredFiles(dir), []);
    rmSync(dir, { recursive: true });
  });
});

describe("checkInternalLinks", () => {
  it("flags a broken relative link", () => {
    const dir = makeKit({ "start-here.md": "[missing](missing-file.md)" });
    const issues = checkInternalLinks(dir);
    assert.equal(issues.length, 1);
    assert.ok(issues[0] !== undefined && issues[0].message.includes("broken link"));
    rmSync(dir, { recursive: true });
  });

  it("passes for a valid relative link", () => {
    const dir = makeKit({
      "start-here.md": "[other](usage/init.md)",
      "usage/init.md": "",
    });
    assert.deepEqual(checkInternalLinks(dir), []);
    rmSync(dir, { recursive: true });
  });

  it("ignores external links", () => {
    const dir = makeKit({ "start-here.md": "[ext](https://example.com)" });
    assert.deepEqual(checkInternalLinks(dir), []);
    rmSync(dir, { recursive: true });
  });
});
