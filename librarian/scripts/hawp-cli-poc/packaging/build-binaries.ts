import { chmodSync, mkdirSync, writeFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { spawnSync } from "node:child_process";

import { binaryTargets, currentBinaryTarget } from "./targets";

interface BuildRecord {
  target: string;
  status: "built" | "skipped";
  output?: string;
  reason?: string;
}

interface BuildOptions {
  target: string;
}

const parseArgs = (argv: string[]): BuildOptions => {
  const targetIndex = argv.indexOf("--target");
  if (targetIndex === -1) {
    return { target: "current" };
  }

  const target = argv[targetIndex + 1];
  if (target === undefined || target.startsWith("--")) {
    throw new Error("--target requires a target id or 'all'");
  }

  return { target };
};

const repoRelative = (...parts: string[]): string => resolve(process.cwd(), ...parts);

const buildDir = repoRelative("build-bin");
const bundlePath = repoRelative("build-bin", "hawp-poc.cjs");
const releaseDir = repoRelative("build-bin", "releases");
const seaConfigDir = repoRelative("build-bin", "sea");
const manifestPath = repoRelative("build-bin", "binary-manifest.json");

const ensureBundle = (): void => {
  const result = spawnSync(
    "npx",
    [
      "esbuild",
      "scripts/hawp-cli-poc/bin/hawp-poc.ts",
      "--bundle",
      "--platform=node",
      "--format=cjs",
      "--outfile=build-bin/hawp-poc.cjs",
    ],
    {
      cwd: process.cwd(),
      encoding: "utf8",
      stdio: "inherit",
    },
  );

  if (result.status !== 0) {
    throw new Error(`esbuild failed with exit code ${result.status ?? "unknown"}`);
  }
};

const supportsBuildSea = (nodeExecutable: string): boolean => {
  const majorVersion = Number.parseInt(process.versions.node.split(".")[0] ?? "0", 10);
  if (nodeExecutable === process.execPath && majorVersion < 26) {
    return false;
  }

  const result = spawnSync(nodeExecutable, ["--help"], {
    encoding: "utf8",
  });

  return result.status === 0 && result.stdout.includes("--build-sea");
};

const buildTarget = (targetId: string): BuildRecord => {
  const target =
    targetId === "current"
      ? currentBinaryTarget()
      : binaryTargets.find((candidate) => candidate.id === targetId);

  if (target === undefined) {
    return {
      target: targetId,
      status: "skipped",
      reason: "target is not in the PoC binary matrix",
    };
  }

  const isCurrentHost = target.os === process.platform && target.arch === process.arch;
  const nodeExecutable = isCurrentHost ? process.execPath : process.env[target.nodeExecutableEnv];

  if (nodeExecutable === undefined) {
    return {
      target: target.id,
      status: "skipped",
      reason: `set ${target.nodeExecutableEnv} to a Node 26 executable for cross-target SEA builds`,
    };
  }

  if (!supportsBuildSea(nodeExecutable)) {
    return {
      target: target.id,
      status: "skipped",
      reason: `${nodeExecutable} does not expose Node 26 --build-sea`,
    };
  }

  const outputPath = resolve(releaseDir, target.fileName);
  const seaConfigPath = resolve(seaConfigDir, `${target.id}.sea.json`);
  mkdirSync(dirname(outputPath), { recursive: true });
  mkdirSync(dirname(seaConfigPath), { recursive: true });

  writeFileSync(
    seaConfigPath,
    `${JSON.stringify(
      {
        main: bundlePath,
        mainFormat: "commonjs",
        executable: nodeExecutable,
        output: outputPath,
        disableExperimentalSEAWarning: true,
        useSnapshot: false,
        useCodeCache: false,
        execArgvExtension: "env",
      },
      null,
      2,
    )}\n`,
  );

  const result = spawnSync(nodeExecutable, ["--build-sea", seaConfigPath], {
    encoding: "utf8",
    stdio: "inherit",
  });

  if (result.status !== 0) {
    return {
      target: target.id,
      status: "skipped",
      reason: `node --build-sea failed with exit code ${result.status ?? "unknown"}`,
    };
  }

  if (target.os !== "win32") {
    chmodSync(outputPath, 0o755);
  }

  if (target.os === "darwin") {
    const signResult = spawnSync("codesign", ["--force", "--sign", "-", "--timestamp=none", outputPath], {
      encoding: "utf8",
      stdio: "inherit",
    });

    if (signResult.status !== 0) {
      return {
        target: target.id,
        status: "skipped",
        reason: `codesign failed with exit code ${signResult.status ?? "unknown"}`,
      };
    }
  }

  return {
    target: target.id,
    status: "built",
    output: outputPath,
  };
};

const selectTargets = (target: string): string[] => {
  if (target === "all") {
    return binaryTargets.map((candidate) => candidate.id);
  }

  return [target];
};

export const runBuildBinaries = (argv: string[]): number => {
  const options = parseArgs(argv);
  mkdirSync(buildDir, { recursive: true });
  buildBundleAndBinaries(options);
  return 0;
};

const buildBundleAndBinaries = (options: BuildOptions): void => {
  buildBundleAndManifestPrelude();
  const records = selectTargets(options.target).map(buildTarget);
  writeManifest(records);
};

const buildBundleAndManifestPrelude = (): void => {
  buildBundle();
};

const buildBundle = (): void => {
  ensureBundle();
};

const writeManifest = (records: BuildRecord[]): void => {
  const manifest = {
    generatedAt: new Date().toISOString(),
    node: process.version,
    platform: process.platform,
    arch: process.arch,
    bundle: bundlePath,
    targets: binaryTargets,
    records,
  };

  writeFileSync(manifestPath, `${JSON.stringify(manifest, null, 2)}\n`);
  process.stdout.write(`${JSON.stringify(manifest, null, 2)}\n`);
};

process.exitCode = runBuildBinaries(process.argv.slice(2));
