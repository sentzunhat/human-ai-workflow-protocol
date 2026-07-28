import assert from "node:assert/strict";
import { mkdtempSync, mkdirSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import { checkBacklogConsistency } from "../validations/backlog-consistency";
import { checkClosedTaskCompleteness } from "../validations/closed-task-completeness";
import {
  checkEvidenceIntegrity,
  collectClosedPlanFiles,
} from "../validations/evidence-integrity";
import { checkVerificationClarity } from "../validations/verification-clarity";

const withTempWorkDir = (fn: (workDir: string) => void): void => {
  const root = mkdtempSync(join(tmpdir(), "hawp-validations-"));
  try {
    const workDir = join(root, ".hawp", "work");
    mkdirSync(join(workDir, "active"), { recursive: true });
    mkdirSync(join(workDir, "parked"), { recursive: true });
    mkdirSync(join(workDir, "closed"), { recursive: true });
    fn(workDir);
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
};

const writeClosedPlan = (
  workDir: string,
  date: { year: string; month: string; day: string },
  fileName: string,
  content: string,
): string => {
  const dir = join(workDir, "closed", date.year, date.month, date.day);
  mkdirSync(dir, { recursive: true });
  const filePath = join(dir, fileName);
  writeFileSync(filePath, content, "utf-8");
  return filePath;
};

test("checkVerificationClarity accepts plain and annotated Verification headings", () => {
  withTempWorkDir((workDir) => {
    const plain = writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "01" },
      "TASK-900.md",
      "# Task\n\n## Verification\n\n- [x] checked. Evidence: output captured\n",
    );
    const annotated = writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "01" },
      "TASK-901.md",
      "# Task\n\n## Verification (filled at close)\n\n- [ ] NOT YET VERIFIED claim\n",
    );

    const result = checkVerificationClarity([plain, annotated]);
    assert.equal(result.total, 2);
    assert.equal(result.proven, 1);
    assert.equal(result.unproven.length, 1);
    assert.equal(result.ambiguous.length, 0);
    assert.equal(typeof result.unproven[0]?.lineNumber, "number");
    assert.match(result.unproven[0]?.filePath ?? "", /TASK-901\.md$/);
    assert.equal(result.status, "WARN");
  });
});

test("checkVerificationClarity warns on checklist claims without evidence labels", () => {
  withTempWorkDir((workDir) => {
    const ambiguous = writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "02" },
      "TASK-902.md",
      "# Task\n\n## Verification\n\n- [x] Smoke test completed\n",
    );

    const result = checkVerificationClarity([ambiguous]);
    assert.equal(result.total, 1);
    assert.equal(result.proven, 0);
    assert.equal(result.unproven.length, 0);
    assert.equal(result.ambiguous.length, 1);
    assert.equal(result.ambiguous[0]?.id, "TASK-902");
    assert.equal(result.ambiguous[0]?.lineNumber, 2);
    assert.match(result.ambiguous[0]?.filePath ?? "", /TASK-902\.md$/);
    assert.equal(result.status, "WARN");
  });
});

test("checkVerificationClarity ignores evidence follow-up scaffolding entries", () => {
  withTempWorkDir((workDir) => {
    const plan = writeClosedPlan(
      workDir,
      { year: "2026", month: "07", day: "09" },
      "TASK-903.md",
      [
        "# Task",
        "",
        "## Verification",
        "",
        "- [x] Smoke test completed. Evidence: command output recorded",
        "",
        "### Evidence Follow-Up",
        "",
        "- [ ] Research evidence for: Old claim",
        "- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.",
        "",
      ].join("\n"),
    );

    const result = checkVerificationClarity([plan]);
    assert.equal(result.total, 1);
    assert.equal(result.proven, 1);
    assert.equal(result.unproven.length, 0);
    assert.equal(result.ambiguous.length, 0);
    assert.equal(result.status, "PASS");
  });
});

test("checkEvidenceIntegrity validates links and rejects evidence-folder escapes", () => {
  withTempWorkDir((workDir) => {
    const evidenceDir = join(workDir, "evidence", "2026", "06", "01");
    mkdirSync(evidenceDir, { recursive: true });
    writeFileSync(join(evidenceDir, "TASK-900-proof.md"), "proof", "utf-8");

    const plan = writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "01" },
      "TASK-900.md",
      [
        "# Task",
        "",
        "## Verification",
        "",
        "- Evidence: link to ../evidence/2026/06/01/TASK-900-proof.md",
        "- Evidence: link to ../evidence/2026/06/01/missing.md",
        "- Evidence: link to ../evidence/../../escape.md",
        "",
      ].join("\n"),
    );

    const result = checkEvidenceIntegrity(workDir, [plan]);
    assert.equal(result.total, 2, "escaping link must not be counted");
    assert.equal(result.valid, 1);
    assert.equal(result.broken.length, 1);
    assert.equal(result.status, "WARN");
  });
});

test("collectClosedPlanFiles walks only YYYY/MM/DD folders", () => {
  withTempWorkDir((workDir) => {
    writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "01" },
      "TASK-900.md",
      "# Task\n",
    );
    writeFileSync(join(workDir, "closed", "README.md"), "readme", "utf-8");

    const files = collectClosedPlanFiles(join(workDir, "closed"));
    assert.equal(files.length, 1);
    assert.ok(files[0]?.endsWith("TASK-900.md"));
  });
});

test("checkBacklogConsistency passes when rows match files and fails on missing plans", () => {
  withTempWorkDir((workDir) => {
    writeFileSync(
      join(workDir, "active", "TASK-100.md"),
      "# plan\n",
      "utf-8",
    );

    const okResult = checkBacklogConsistency(workDir, {
      active: [
        {
          id: "TASK-100",
          type: "task",
          title: "sample",
          status: "in-progress",
          detail: "[plan](active/TASK-100.md)",
        },
      ],
      closed: [],
      parked: [],
    });
    assert.equal(okResult.status, "PASS");
    assert.equal(okResult.activeWork.found, 1);

    const missingResult = checkBacklogConsistency(workDir, {
      active: [
        {
          id: "TASK-200",
          type: "task",
          title: "missing plan",
          status: "in-progress",
          detail: undefined,
        },
      ],
      closed: [],
      parked: [],
    });
    assert.equal(missingResult.status, "FAIL");
    assert.deepEqual(missingResult.activeWork.missing, ["TASK-200"]);
  });
});

test("checkBacklogConsistency matches short-UUID rows to full-UUID files without orphans", () => {
  withTempWorkDir((workDir) => {
    const uuid = "0e1c4afa-9668-4d61-b5b6-1e27be42ca23";
    writeFileSync(join(workDir, "active", `${uuid}.md`), "# plan\n", "utf-8");

    const result = checkBacklogConsistency(workDir, {
      active: [
        {
          id: "0e1c4afa",
          type: "improvement",
          title: "short-uuid display row",
          status: "in-progress",
          detail: `[plan](active/${uuid}.md)`,
        },
      ],
      closed: [],
      parked: [],
    });
    assert.equal(result.status, "PASS");
    assert.equal(result.activeWork.found, 1);
    assert.deepEqual(result.orphanedFiles, []);
  });
});

test("checkBacklogConsistency resolves UUID-named active plan files without orphans", () => {
  withTempWorkDir((workDir) => {
    const uuid = "361fb08e-6457-4ed5-80bd-76337b6f0e89";
    writeFileSync(join(workDir, "active", `${uuid}.md`), "# plan\n", "utf-8");

    const result = checkBacklogConsistency(workDir, {
      active: [
        {
          id: uuid,
          type: "task",
          title: "uuid-native item",
          status: "in-progress",
          detail: `[plan](active/${uuid}.md)`,
        },
      ],
      closed: [],
      parked: [],
    });
    assert.equal(result.status, "PASS");
    assert.equal(result.activeWork.found, 1);
    assert.deepEqual(result.orphanedFiles, []);
  });
});

test("checkClosedTaskCompleteness recognizes a shared closed record via its Closes: line", () => {
  withTempWorkDir((workDir) => {
    writeClosedPlan(
      workDir,
      { year: "2026", month: "07", day: "26" },
      "shared-audit.md",
      [
        "# Shared audit",
        "",
        "**Closes:** `abc12345`, `def67890`",
        "",
        "## Outcome",
        "",
        "Both fixed.",
        "",
        "## Verification",
        "",
        "- [x] checked. Evidence: output captured",
        "",
        "## Close Checklist",
        "",
        "- [x] done",
        "",
      ].join("\n"),
    );

    const result = checkClosedTaskCompleteness(workDir, []);
    assert.equal(result.total, 1);
    assert.equal(result.untypedCurrent.length, 0);
    assert.equal(result.failing.length, 0);
    assert.equal(result.status, "PASS");
  });
});

test("checkBacklogConsistency matches multiple rows against one shared closed record", () => {
  withTempWorkDir((workDir) => {
    writeClosedPlan(
      workDir,
      { year: "2026", month: "07", day: "26" },
      "shared-audit.md",
      "# Shared audit\n\n**Closes:** `abc12345`, `def67890`\n",
    );

    const result = checkBacklogConsistency(workDir, {
      active: [],
      closed: [
        {
          id: "abc12345",
          type: "fix",
          title: "first fix",
          status: "done",
          detail: undefined,
        },
        {
          id: "def67890",
          type: "fix",
          title: "second fix",
          status: "done",
          detail: undefined,
        },
      ],
      parked: [],
    });
    assert.equal(result.status, "PASS");
    assert.equal(result.recentlyClosed.found, 2);
  });
});

test("checkClosedTaskCompleteness fails post-cutoff plans missing required sections", () => {
  withTempWorkDir((workDir) => {
    writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "01" },
      "TASK-910.md",
      "# Task\n\n## Outcome\n\nDone.\n\n## Verification\n\n- [x] ok\n\n## Close Checklist\n\n- [x] done\n",
    );
    writeClosedPlan(
      workDir,
      { year: "2026", month: "06", day: "01" },
      "TASK-911.md",
      "# Task\n\nNo close sections here.\n",
    );

    const result = checkClosedTaskCompleteness(workDir, []);
    assert.equal(result.total, 2);
    assert.equal(result.failing.length, 1);
    assert.equal(result.failing[0]?.id, "TASK-911");
    assert.equal(result.status, "FAIL");
  });
});
