export interface BinaryTarget {
  id: string;
  os: NodeJS.Platform;
  arch: NodeJS.Architecture;
  fileName: string;
  nodeExecutableEnv: string;
}

export const binaryTargets: BinaryTarget[] = [
  {
    id: "linux-x64",
    os: "linux",
    arch: "x64",
    fileName: "hawp-poc-linux-x64",
    nodeExecutableEnv: "HAWP_POC_NODE_LINUX_X64",
  },
  {
    id: "linux-arm64",
    os: "linux",
    arch: "arm64",
    fileName: "hawp-poc-linux-arm64",
    nodeExecutableEnv: "HAWP_POC_NODE_LINUX_ARM64",
  },
  {
    id: "darwin-x64",
    os: "darwin",
    arch: "x64",
    fileName: "hawp-poc-darwin-x64",
    nodeExecutableEnv: "HAWP_POC_NODE_DARWIN_X64",
  },
  {
    id: "darwin-arm64",
    os: "darwin",
    arch: "arm64",
    fileName: "hawp-poc-darwin-arm64",
    nodeExecutableEnv: "HAWP_POC_NODE_DARWIN_ARM64",
  },
  {
    id: "win32-x64",
    os: "win32",
    arch: "x64",
    fileName: "hawp-poc-win32-x64.exe",
    nodeExecutableEnv: "HAWP_POC_NODE_WIN32_X64",
  },
  {
    id: "win32-arm64",
    os: "win32",
    arch: "arm64",
    fileName: "hawp-poc-win32-arm64.exe",
    nodeExecutableEnv: "HAWP_POC_NODE_WIN32_ARM64",
  },
];

export const currentBinaryTarget = (): BinaryTarget | undefined =>
  binaryTargets.find((target) => target.os === process.platform && target.arch === process.arch);
