import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import { resolve, dirname, extname } from 'node:path';

const repoRoot = resolve(import.meta.dirname, '../..');
const roots = ['.hawp', 'docs', 'README.md']
  .map((path) => resolve(repoRoot, path))
  .filter(existsSync);
const markdownFiles = [];

function collect(path) {
  const stats = statSync(path);
  if (stats.isDirectory()) {
    for (const entry of readdirSync(path)) collect(resolve(path, entry));
    return;
  }
  if (extname(path).toLowerCase() === '.md') markdownFiles.push(path);
}

for (const root of roots) collect(root);

const linkPattern = /(?<!!?)\[[^\]]*\]\((?<target>[^)\s]+)(?:\s+['"][^)]*['"])?\)/g;
const failures = [];

for (const file of markdownFiles) {
  const content = readFileSync(file, 'utf8');
  for (const match of content.matchAll(linkPattern)) {
    const target = match.groups.target.replace(/^<|>$/g, '');
    if (!target || target.startsWith('#') || /^(?:[a-z][a-z+.-]*:|\/\/)/i.test(target)) continue;

    const pathTarget = decodeURIComponent(target.split('#', 1)[0]);
    const resolved = pathTarget.startsWith('/')
      ? resolve(repoRoot, `.${pathTarget}`)
      : resolve(dirname(file), pathTarget);

    if (!existsSync(resolved)) {
      failures.push(`${file.slice(repoRoot.length + 1)} -> ${target}`);
    }
  }
}

if (failures.length > 0) {
  console.error('Broken local Markdown links:');
  for (const failure of failures) console.error(`- ${failure}`);
  process.exitCode = 1;
} else {
  console.log(`Checked ${markdownFiles.length} Markdown file(s): local links are valid.`);
}
