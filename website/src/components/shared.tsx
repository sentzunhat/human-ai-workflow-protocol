import { useEffect, useState, type MouseEvent } from "react";
import { githubUrl, navItems } from "../lib/content/site";
import {
  applyTheme,
  getThemePreference,
  watchSystemTheme,
  type ThemePreference,
} from "../lib/theme";

export function ReleaseVersion({
  version,
  loading = false,
  badge = false,
}: {
  version: string;
  loading?: boolean;
  badge?: boolean;
}) {
  return (
    <span
      className={`release-version${loading ? " release-version-loading" : ""}${badge ? " release-version-badge" : ""}`}
    >
      <span className="release-version-text">v{version}</span>
      {loading && (
        <span className="release-skeleton-sheen" aria-hidden="true" />
      )}
    </span>
  );
}

export function ThemeToggle() {
  const [preference, setPreference] = useState<ThemePreference>("system");
  useEffect(() => {
    const initial = getThemePreference();
    setPreference(initial);
    applyTheme(initial);
    return watchSystemTheme(() => {
      if (getThemePreference() === "system") applyTheme("system");
    });
  }, []);
  const options: ThemePreference[] = ["light", "dark", "system"];
  return (
    <div
      className="flex items-center rounded-full border border-[var(--hawp-line)] bg-[var(--hawp-raised)]/90 p-1 shadow-sm backdrop-blur"
      aria-label="Theme preference"
    >
      {options.map((option) => (
        <button
          key={option}
          type="button"
          className={`flex min-h-8 items-center gap-1.5 rounded-full px-2.5 text-xs font-semibold transition sm:px-3 ${preference === option ? "bg-[var(--hawp-panel-soft)] text-[var(--hawp-ink)] shadow-sm" : "text-[var(--hawp-muted-2)] hover:text-[var(--hawp-ink)]"}`}
          aria-pressed={preference === option}
          aria-label={`Use ${option} theme`}
          title={`${option[0].toUpperCase()}${option.slice(1)} theme`}
          onClick={() => {
            setPreference(option);
            applyTheme(option);
          }}
        >
          <span aria-hidden="true">
            {option === "light" ? "☼" : option === "dark" ? "◐" : "▣"}
          </span>
          <span className="hidden lg:inline">
            {option[0].toUpperCase()}
            {option.slice(1)}
          </span>
        </button>
      ))}
    </div>
  );
}

export function Header() {
  const items = [...navItems, { label: "Downloads", href: "#downloads" }];
  const closeMenu = (event: MouseEvent<HTMLAnchorElement>) =>
    event.currentTarget.closest("details")?.removeAttribute("open");
  return (
    <header className="sticky top-0 z-50 border-b border-[var(--hawp-line)] bg-[var(--hawp-bg)]/92 backdrop-blur-xl">
      <div className="mx-auto flex h-[72px] max-w-[1280px] items-center justify-between px-4 sm:px-8 lg:px-10">
        <a
          href="#top"
          className="flex items-center gap-2.5"
          aria-label="HAWP home"
        >
          <svg
            className="size-9 shrink-0"
            viewBox="0 0 100 100"
            aria-hidden="true"
          >
            <path
              d="M50 7 90.9 36.7 75.3 84.8 24.7 84.8 9.1 36.7Z"
              fill="none"
              stroke="var(--hawp-peach)"
              strokeWidth="8"
              strokeLinejoin="round"
            />
          </svg>
          <span className="text-xl font-bold tracking-tight text-[var(--hawp-ink)]">
            HAWP
          </span>
        </a>
        <nav
          className="hidden items-center gap-8 lg:flex"
          aria-label="Main navigation"
        >
          {items.map((item) => (
            <a
              key={item.label}
              href={item.href}
              className="text-sm text-[var(--hawp-muted)] transition hover:text-[var(--hawp-ink)]"
            >
              {item.label}
            </a>
          ))}
        </nav>
        <div className="flex items-center gap-1.5 sm:gap-3">
          <ThemeToggle />
          <details className="group relative lg:hidden">
            <summary className="flex min-h-10 cursor-pointer list-none items-center justify-center rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-3 text-sm font-semibold text-[var(--hawp-ink)] marker:content-none">
              Menu
            </summary>
            <nav
              className="absolute right-0 top-12 min-w-44 overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-2 shadow-2xl shadow-black/10"
              aria-label="Mobile navigation"
            >
              {items.map((item) => (
                <a
                  key={item.label}
                  href={item.href}
                  className="block rounded-xl px-3 py-2.5 text-sm font-semibold text-[var(--hawp-muted)] transition hover:bg-[var(--hawp-panel)] hover:text-[var(--hawp-ink)]"
                  onClick={closeMenu}
                >
                  {item.label}
                </a>
              ))}
              <a
                href={githubUrl}
                className="mt-1 block rounded-xl border-t border-[var(--hawp-line)] px-3 py-2.5 text-sm font-semibold text-[var(--hawp-ink)]"
                onClick={closeMenu}
              >
                GitHub ↗
              </a>
            </nav>
          </details>
          <a
            href={githubUrl}
            className="hidden rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-4 py-2.5 text-sm font-semibold text-[var(--hawp-ink)] transition hover:border-[var(--hawp-red)] hover:bg-[var(--hawp-panel)] lg:inline-flex"
          >
            GitHub
          </a>
        </div>
      </div>
    </header>
  );
}
