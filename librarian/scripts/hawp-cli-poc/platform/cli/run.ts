import { registeredCommands } from "./command-registry";

const helpText = (): string => {
  const commands = registeredCommands
    .map((entry) => `  ${entry.commandPath.join(" ")}  ${entry.summary}`)
    .join("\n");

  return [
    "hawp-poc",
    "",
    "Node 26 + oclif proof of concept for the future installable HAWP CLI.",
    "",
    "USAGE",
    "  hawp-poc <command> [options]",
    "",
    "COMMANDS",
    commands,
    "",
    "Run `hawp-poc <command> --help` for command-specific help.",
  ].join("\n");
};

const isHelpRequest = (argv: string[]): boolean =>
  argv.length === 0 || argv.includes("--help") || argv.includes("-h");

export const runCli = async (argv: string[]): Promise<number> => {
  if (isHelpRequest(argv)) {
    process.stdout.write(`${helpText()}\n`);
    return 0;
  }

  const commandEntry = registeredCommands.find((entry) =>
    entry.commandPath.every((part, index) => argv[index] === part),
  );

  if (commandEntry === undefined) {
    process.stderr.write(`Unknown command: ${argv.join(" ")}\n\n${helpText()}\n`);
    return 2;
  }

  const commandArgv = argv.slice(commandEntry.commandPath.length);

  try {
    await commandEntry.command.run(commandArgv);
    return 0;
  } catch (error) {
    const maybeExit = error as { oclif?: { exit?: number }; code?: string };
    if (typeof maybeExit.oclif?.exit === "number") {
      return maybeExit.oclif.exit;
    }

    throw error;
  }
};
