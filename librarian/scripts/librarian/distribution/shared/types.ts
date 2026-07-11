export interface DistributionVariant {
  outputFile: string;
  sectionFiles: string[];
  scriptPartFiles: string[];
  ref: string;
  provider: string;
  operation: "install" | "update";
}

export interface CompositionPlan {
  variants: DistributionVariant[];
}

export interface BuildResult {
  outputPath: string;
  content: string;
}
