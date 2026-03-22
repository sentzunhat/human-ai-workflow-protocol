import assert from "node:assert/strict";
import test from "node:test";

import { Mode, OutputFormat } from "../models/index.js";
import { parseArgs } from "../cli.js";

test("parseArgs parses default dry-run mode (parse-only)", () => {
  const parsed = parseArgs(["backlog", "upgrade"]);

  assert.ok(parsed);
  assert.equal(parsed.command, "backlog");
  assert.equal(parsed.subcommand, "upgrade");
  assert.equal(parsed.mode, Mode.DryRun);
  assert.equal(parsed.validate, false);
  assert.equal(parsed.format, OutputFormat.Text);
  assert.equal(parsed.forceDirty, false);
  assert.equal(parsed.verbose, false);
});

test("parseArgs parses apply and optional flags (no execution)", () => {
  const parsed = parseArgs([
    "backlog",
    "upgrade",
    "--apply",
    "--validate",
    "--format",
    "json",
    "--output",
    "result.txt",
    "--export-plan",
    "plan.json",
    "--export-research-queue",
    "research.json",
    "--force-dirty",
    "--verbose",
  ]);

  assert.ok(parsed);
  assert.equal(parsed.mode, Mode.Apply);
  assert.equal(parsed.validate, true);
  assert.equal(parsed.format, OutputFormat.Json);
  assert.equal(parsed.output, "result.txt");
  assert.equal(parsed.exportPlan, "plan.json");
  assert.equal(parsed.exportResearchQueue, "research.json");
  assert.equal(parsed.forceDirty, true);
  assert.equal(parsed.verbose, true);
});

test("parseArgs parses help and version shortcuts", () => {
  const helpParsed = parseArgs(["--help"]);
  assert.ok(helpParsed);
  assert.equal(helpParsed.help, true);

  const versionParsed = parseArgs(["-v"]);
  assert.ok(versionParsed);
  assert.equal(versionParsed.version, true);
});
