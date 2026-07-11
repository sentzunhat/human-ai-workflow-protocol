import { Command, Flags } from "@oclif/core";

import { validateKit } from "../../../domain/kit/validate-kit";
import { resolveKitPath } from "../../../infrastructure/filesystem/hawp-repo";

export default class KitValidateCommand extends Command {
  static override readonly description =
    "Validate .hawp/kit structure using the Zacatl-aligned CLI PoC.";

  static override readonly examples = [
    "<%= config.bin %> kit validate",
    "<%= config.bin %> kit validate --kit-path .hawp/kit",
  ];

  static override readonly flags = {
    "kit-path": Flags.string({
      description: "Path to the kit directory. Defaults to the nearest repo .hawp/kit.",
      required: false,
    }),
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(KitValidateCommand);
    const kitPathInput =
      flags["kit-path"] === undefined
        ? { cwd: process.cwd() }
        : { cwd: process.cwd(), kitPath: flags["kit-path"] };
    const kitPath = resolveKitPath(kitPathInput);
    const result = validateKit({ kitPath });

    this.log("kit:validate");
    this.log("============");
    this.log(`kit: ${result.kitPath}`);
    this.log("");

    if (result.issueCount === 0) {
      this.log(`OK ${result.checks} checks passed, 0 issues`);
      return;
    }

    for (const issue of result.issues) {
      this.error(`${issue.file}: ${issue.message}`, { exit: false });
    }

    this.exit(1);
  }
}
