export interface KitValidateArgs {
  kitPath?: string; // override default .hawp/kit path
}

export const parseArgs = (args: string[]): KitValidateArgs => {
  const opts: KitValidateArgs = {};
  let i = 0;
  while (i < args.length) {
    if (args[i] === "--kit-path" && args[i + 1] !== undefined) {
      opts.kitPath = args[i + 1] as string;
      i += 2;
    } else {
      i++;
    }
  }
  return opts;
};
