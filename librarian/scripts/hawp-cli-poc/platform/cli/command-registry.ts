import KitValidateCommand from "../../application/commands/kit/validate";

interface RunnableCommand {
  run(argv?: string[]): Promise<unknown>;
}

export interface RegisteredCommand {
  commandPath: string[];
  command: RunnableCommand;
  summary: string;
}

export const registeredCommands: RegisteredCommand[] = [
  {
    commandPath: ["kit", "validate"],
    command: KitValidateCommand,
    summary: "Validate .hawp/kit structure.",
  },
];
