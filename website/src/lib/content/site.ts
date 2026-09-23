export const githubUrl = 'https://github.com/sentzunhat/human-ai-workflow-protocol';
export const currentVersion = '0.0.24';
export const releaseLoadingVersion = '0.0.0';
export const releaseApiUrl = 'https://api.github.com/repos/sentzunhat/human-ai-workflow-protocol/releases?per_page=1';

export type ReleaseAsset = {
  name: string;
  browser_download_url: string;
};

export type LatestRelease = {
  tag_name: string;
  html_url: string;
  assets: ReleaseAsset[];
};

let latestReleasePromise: Promise<LatestRelease | null> | null = null;

export const platformDownloads = [
  { label: 'macOS Apple Silicon', asset: 'hawp-darwin-arm64', icon: '⌘' },
  { label: 'macOS Intel', asset: 'hawp-darwin-amd64', icon: '⌘' },
  { label: 'Linux x64', asset: 'hawp-linux-amd64', icon: '◈' },
  { label: 'Linux ARM64', asset: 'hawp-linux-arm64', icon: '◈' },
  { label: 'Windows x64', asset: 'hawp-windows-amd64.exe', icon: '⊞' },
  { label: 'Windows ARM64', asset: 'hawp-windows-arm64.exe', icon: '⊞' }
] as const;

export async function loadLatestRelease(): Promise<LatestRelease | null> {
  if (latestReleasePromise) return latestReleasePromise;

  latestReleasePromise = fetchLatestRelease();
  return latestReleasePromise;
}

async function fetchLatestRelease(): Promise<LatestRelease | null> {
  try {
    const response = await fetch(releaseApiUrl, {
      headers: { Accept: 'application/vnd.github+json' }
    });
    if (!response.ok) return null;
    const releases = (await response.json()) as LatestRelease[];
    return releases[0] ?? null;
  } catch {
    return null;
  }
}

export function downloadUrl(asset: string): string {
  return `${githubUrl}/releases/latest/download/${asset}`;
}

export const navItems = [
  { label: 'Protocol', href: '#protocol' },
  { label: 'Agents', href: '#agents' },
  { label: 'Benchmarks', href: '#benchmarks' },
  { label: 'FAQ', href: '#faq' }
] as const;

export const problemItems = [
  ['01', 'Every session starts cold.', 'The goal and constraints get reconstructed from conversation instead of carried with the work.'],
  ['02', 'Execution drifts.', 'Agents optimize locally while the original mission, boundaries, or expected output fade from context.'],
  ['03', 'Handoffs lose decisions.', 'The next model, tool, or person re-derives what was already known and repeats avoidable work.']
] as const;

export const protocolFields = [
  ['01', 'Input', 'The request as received, before interpretation changes it.'],
  ['02', 'Context', 'The smallest useful background the agent needs to reason correctly.'],
  ['03', 'Mission', 'The concrete objective in one clear, execution-ready statement.'],
  ['04', 'Constraints', 'Hard boundaries, quality bars, dependencies, and things that must not change.'],
  ['05', 'Output', 'The observable end state that defines what “done” actually means.']
] as const;

export const providers = ['Claude Code', 'GitHub Copilot', 'Codex', 'Cursor', 'Continue'] as const;

export const benchmarkMetrics = [
  ['19%', 'aggregate context-shaping token reduction across the 10-query benchmark, with dense cases saving up to 38%.'],
  ['95%', 'downstream token reduction across the 9 successful v0.0.24 Ollama intake-shaping benchmark runs.'],
  ['<1 ms', 'lexical search latency in the current kit benchmark; hybrid search measured 72 ms with 10/10 reported quality.']
] as const;

export const workflowSteps = [
  ['01 / SHAPE', 'Lock intent before execution.', 'Capture the request, relevant context, mission, constraints, and expected output as a work artifact.'],
  ['02 / EXECUTE', 'Give the agent bounded context.', 'Search the kit and work index, then pack only the useful context instead of repeatedly loading everything.'],
  ['03 / HAND OFF', 'Carry the shape forward.', 'The next agent or session reads the same portable state instead of reconstructing the project from conversation history.']
] as const;

export const faqs = [
  ['What is HAWP?', 'HAWP is an open-source Human-AI Workflow Protocol for shaping AI agent work before execution. It carries intent, context, mission, constraints, output, evidence, and handoffs in portable work artifacts.'],
  ['Is HAWP an AI agent framework?', 'No. HAWP is a workflow protocol and CLI for shaping and carrying work. It does not replace your agent runtime, orchestrator, or model.'],
  ['Does HAWP require a database or cloud service?', 'No. The work model is Markdown-first and repository-local. Search indexing can be generated locally, so the protocol remains a set of portable files rather than a hosted dependency.'],
  ['Why does HAWP use MCP?', 'MCP gives connected agents a standard way to search the HAWP kit and work index without forcing the protocol into a specific editor or vendor.'],
  ['Can I use HAWP with more than one coding agent?', 'Yes. Provider guides in the repository cover Claude Code, GitHub Copilot, Codex, Cursor, and Continue.']
] as const;
