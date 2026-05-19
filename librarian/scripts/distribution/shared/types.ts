export interface DistributionVariant {
  outputFile: string;
  sectionFiles: string[];
  scriptSourceFile: string;
  ref: string;
}

export interface CompositionPlan {
  variants: DistributionVariant[];
}

export interface BuildResult {
  outputPath: string;
  content: string;
}
