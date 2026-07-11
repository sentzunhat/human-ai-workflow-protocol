import assert from "node:assert/strict";
import { execSync } from "node:child_process";
import {
  existsSync,
  mkdtempSync,
  mkdirSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import { runKitNormalize } from "../script";

const withTempRepo = async (
  fn: (repoRoot: string) => void | Promise<void>,
): Promise<void> => {
  const repoRoot = mkdtempSync(join(tmpdir(), "hawp-kit-normalize-"));
  try {
    await fn(repoRoot);
  } finally {
    rmSync(repoRoot, { recursive: true, force: true });
  }
};

const seedRepo = (repoRoot: string): void => {
  mkdirSync(join(repoRoot, ".hawp", "work"), { recursive: true });
  mkdirSync(join(repoRoot, ".hawp", "kit", "standards", "public", "templates"), {
    recursive: true,
  });
  mkdirSync(join(repoRoot, ".hawp", "kit", "standards", "public", "standards", "docs"), {
    recursive: true,
  });
  mkdirSync(join(repoRoot, ".hawp", "kit", "standards", "public", "standards", "nodejs"), {
    recursive: true,
  });

  writeFileSync(
    join(repoRoot, ".hawp", "work", "BACKLOG.md"),
    "# Backlog\n",
    "utf-8",
  );
};

test("kit normalize dry-run reports file renames and link updates", async () => {
  await withTempRepo(async (repoRoot) => {
    seedRepo(repoRoot);
    const kitRoot = join(repoRoot, ".hawp", "kit");
    writeFileSync(join(kitRoot, "standards", "public", "templates", "ADR.template.md"), "# ADR\n", "utf-8");
    writeFileSync(
      join(kitRoot, "standards", "public", "standards", "docs", "README.md"),
      "Use [template](../../templates/ADR.template.md).\n",
      "utf-8",
    );

    const originalCwd = process.cwd();
    process.chdir(repoRoot);
    try {
      const exitCode = runKitNormalize(["--kit-path", kitRoot]);
      assert.equal(exitCode, 0);
      assert.equal(
        existsSync(join(kitRoot, "standards", "public", "templates", "ADR.template.md")),
        true,
      );
      assert.equal(
        readFileSync(
          join(kitRoot, "standards", "public", "standards", "docs", "README.md"),
          "utf-8",
        ),
        "Use [template](../../templates/ADR.template.md).\n",
      );
    } finally {
      process.chdir(originalCwd);
    }
  });
});

test("kit normalize apply renames files and updates internal links", async () => {
  await withTempRepo(async (repoRoot) => {
    execSync("git init -q", { cwd: repoRoot });
    seedRepo(repoRoot);
    const kitRoot = join(repoRoot, ".hawp", "kit");
    const templatePath = join(kitRoot, "standards", "public", "templates", "ADR.template.md");
    const docsReadme = join(kitRoot, "standards", "public", "standards", "docs", "README.md");
    const nodeGuide = join(kitRoot, "standards", "public", "standards", "nodejs", "guide.md");

    writeFileSync(templatePath, "# ADR\n", "utf-8");
    writeFileSync(docsReadme, "Use [template](../../templates/ADR.template.md).\n", "utf-8");
    writeFileSync(nodeGuide, "Use [template](../../templates/ADR.template.md).\n", "utf-8");

    // commit so the working tree is clean before apply
    execSync("git config user.email 'test@test.com'", { cwd: repoRoot });
    execSync("git config user.name 'Test'", { cwd: repoRoot });
    execSync("git add -A && git commit -m 'seed'", { cwd: repoRoot });

    const originalCwd = process.cwd();
    process.chdir(repoRoot);
    try {
      const exitCode = runKitNormalize(["--apply", "--kit-path", kitRoot]);
      assert.equal(exitCode, 0);
      assert.equal(existsSync(templatePath), false);
      assert.equal(
        existsSync(join(kitRoot, "standards", "public", "templates", "adr-template.md")),
        true,
      );
      assert.match(
        readFileSync(docsReadme, "utf-8"),
        /adr-template\.md/,
      );
      assert.match(
        readFileSync(nodeGuide, "utf-8"),
        /adr-template\.md/,
      );
    } finally {
      process.chdir(originalCwd);
    }
  });
});
