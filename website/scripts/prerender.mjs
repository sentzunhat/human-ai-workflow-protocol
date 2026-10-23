import { readFile, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

const outputDir = resolve('build');
const [template, { renderPage }] = await Promise.all([
  readFile(resolve(outputDir, 'index.html'), 'utf8'),
  import(pathToFileURL(resolve('server-dist/entry-server.js')).href)
]);

const markup = renderPage();
if (!template.includes('<div id="root"></div>')) {
  throw new Error('Static page template is missing the React root');
}
await writeFile(resolve(outputDir, 'index.html'), template.replace('<div id="root"></div>', `<div id="root">${markup}</div>`));
