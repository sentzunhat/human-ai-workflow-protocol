import { runKitValidateScript } from "../../../hawp/kit-validate/script";

export interface ValidateKitInput {
  kitPath: string;
}

export interface ValidateKitOutput {
  checks: number;
  issueCount: number;
  issues: Array<{ file: string; message: string }>;
  kitPath: string;
}

export const validateKit = (input: ValidateKitInput): ValidateKitOutput => {
  const result = runKitValidateScript(input.kitPath);

  return {
    checks: result.checks,
    issueCount: result.issues.length,
    issues: result.issues,
    kitPath: result.kitPath,
  };
};
