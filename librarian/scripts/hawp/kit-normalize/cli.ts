export interface KitNormalizeArgs {
  help: boolean;
  apply: boolean;
  kitPath?: string;
}

export class CLIUsageError extends Error {
  public constructor(message: string) {
    super(message);
    this.name = "CLIUsageError";
  }
}

export const parseArgs = (argv: string[]): KitNormalizeArgs => {
  const opts: KitNormalizeArgs = {
    help: false,
    apply: false,
  };

  let i = 0;
  while (i < argv.length) {
    const arg = argv[i];

    if (arg === "--help" || arg === "-h") {
      opts.help = true;
      return opts;
    }

    if (arg === "--apply") {
      opts.apply = true;
      i += 1;
      continue;
    }

    if (arg === "--dry-run") {
      opts.apply = false;
      i += 1;
      continue;
    }

    if (arg === "--kit-path") {
      const kitPath = argv[i + 1];
      if (!kitPath || kitPath.startsWith("-")) {
        throw new CLIUsageError("Error: --kit-path requires a directory path");
      }
      opts.kitPath = kitPath;
      i += 2;
      continue;
    }

    if (arg?.startsWith("-")) {
      throw new CLIUsageError(`Error: Unknown flag '${arg}'`);
    }

    i += 1;
  }

  return opts;
};

export const getHelpText = (): string =>
  [
    "kit:normalize — normalize .hawp/kit structure and links",
    "",
    "Usage:",
    "  npm --prefix librarian run kit:normalize [--apply] [--kit-path <path>]",
    "",
    "Behavior:",
    "  - Defaults to dry-run",
    "  - --apply renames nonconforming kit files and updates internal links",
    "  - Fails closed if apply mode is requested on a dirty working tree",
    "",
    "Options:",
    "  --apply              Apply file renames and link updates in place",
    "  --dry-run            Show planned changes without mutating files (default)",
    "  --kit-path <path>    Override the default .hawp/kit path",
    "  --help, -h           Show this help message",
    "",
  ].join("\n");
