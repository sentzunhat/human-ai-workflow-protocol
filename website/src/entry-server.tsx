import { renderToString } from 'react-dom/server';
import { App } from './app';

export function renderPage() {
  return renderToString(<App />);
}
