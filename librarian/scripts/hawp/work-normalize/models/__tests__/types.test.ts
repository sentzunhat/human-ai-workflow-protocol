/**
 * Tests for backlog upgrade data models: factory functions, type guards,
 * and JSON round-trip serialization.
 */

import assert from "node:assert/strict";
import test from "node:test";

import {
  Mode,
  OutputFormat,
  ExitCode,
  createBlockedItem,
  isBlockedItem,
  createBacklogFixOperation,
  createBacklogFixPlan,
  isBacklogFixPlan,
  createDetectionReport,
  isDetectionReport,
} from "../index.js";

test("enums expose expected members", () => {
  assert.ok(Mode.DryRun);
  assert.ok(OutputFormat.Json);
  assert.equal(ExitCode.Success, 0);
});

test("BlockedItem factory output passes its type guard and survives JSON round-trip", () => {
  const blocked = createBlockedItem(
    "BLOCKED-001",
    "B1",
    "TASK-001",
    0.68,
    ["task", "bug", "improvement"],
    "Ambiguous type",
    { data: "test" },
    "Manual review",
  );

  assert.equal(isBlockedItem(blocked), true);
  assert.equal(isBlockedItem(JSON.parse(JSON.stringify(blocked))), true);
  assert.equal(isBlockedItem({ not: "a blocked item" }), false);
});

test("BacklogFixPlan factory output passes its type guard", () => {
  const operation = createBacklogFixOperation(
    "OP-001",
    "add-field",
    "TASK-001",
    ".hawp/work/BACKLOG.md",
    [1, 5],
    "Add field",
    "safe",
    0.95,
  );

  const plan = createBacklogFixPlan(
    "PLAN-abc123",
    "sha256_hash",
    new Date().toISOString(),
    ".hawp/work/BACKLOG.md",
    5,
    10,
    [operation],
    { version: "1.0.0" },
  );

  assert.equal(isBacklogFixPlan(plan), true);
  assert.equal(plan.operations.length, 1);
});

test("DetectionReport factory output passes its type guard", () => {
  const plan = createBacklogFixPlan(
    "PLAN-abc123",
    "sha256_hash",
    new Date().toISOString(),
    ".hawp/work/BACKLOG.md",
    5,
    10,
    [],
    { version: "1.0.0" },
  );

  const report = createDetectionReport(
    "REPORT-001",
    new Date().toISOString(),
    plan,
    [".hawp/work/BACKLOG.md"],
    ["TASK-001"],
  );

  assert.equal(isDetectionReport(report), true);
});
