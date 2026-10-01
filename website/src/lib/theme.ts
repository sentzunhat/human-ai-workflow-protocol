export type ThemePreference = 'light' | 'dark' | 'system';

const storageKey = 'hawp-theme';
const mediaQuery = '(prefers-color-scheme: dark)';

function resolvedTheme(preference: ThemePreference): 'light' | 'dark' {
  if (preference !== 'system') return preference;
  return window.matchMedia(mediaQuery).matches ? 'dark' : 'light';
}

export function getThemePreference(): ThemePreference {
  if (typeof window === 'undefined') return 'system';
  const stored = window.localStorage.getItem(storageKey);
  return stored === 'light' || stored === 'dark' || stored === 'system' ? stored : 'system';
}

export function applyTheme(preference: ThemePreference): void {
  if (typeof window === 'undefined') return;

  const theme = resolvedTheme(preference);
  const root = document.documentElement;
  root.dataset.preference = preference;
  root.dataset.theme = theme;
  window.localStorage.setItem(storageKey, preference);

  const themeColor = document.querySelector<HTMLMetaElement>('meta[name="theme-color"]');
  if (themeColor) themeColor.content = theme === 'dark' ? '#151112' : '#f5eee8';
}

export function watchSystemTheme(onChange: () => void): () => void {
  if (typeof window === 'undefined') return () => undefined;
  const query = window.matchMedia(mediaQuery);
  query.addEventListener('change', onChange);
  return () => query.removeEventListener('change', onChange);
}
