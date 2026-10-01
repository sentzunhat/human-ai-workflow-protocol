<script lang="ts">
  import { githubUrl, navItems } from '$lib/content/site';
  import ThemeToggle from './ThemeToggle.svelte';

  function closeMobileMenu(event: MouseEvent) {
    (event.currentTarget as HTMLElement).closest('details')?.removeAttribute('open');
  }
</script>

<header class="sticky top-0 z-50 border-b border-[var(--hawp-line)] bg-[var(--hawp-bg)]/92 backdrop-blur-xl">
  <div class="mx-auto flex h-[72px] max-w-[1280px] items-center justify-between px-4 sm:px-8 lg:px-10">
    <a href="#top" class="flex items-center gap-2.5" aria-label="HAWP home">
      <span class="relative grid size-9 shrink-0 place-items-center rounded-[11px] border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel)]">
        <span class="size-4 rotate-45 border border-[var(--hawp-peach)]"></span>
        <span class="absolute size-2 rotate-45 bg-[var(--hawp-red)]"></span>
      </span>
      <span class="text-xl font-bold tracking-tight text-[var(--hawp-ink)]">HAWP</span>
    </a>

    <nav class="hidden items-center gap-8 md:flex" aria-label="Main navigation">
      {#each [...navItems, { label: 'Downloads', href: '#downloads' }] as item}
        <a href={item.href} class="text-sm text-[var(--hawp-muted)] transition hover:text-[var(--hawp-ink)]">{item.label}</a>
      {/each}
    </nav>

    <div class="flex items-center gap-1.5 sm:gap-3">
      <ThemeToggle />

      <details class="group relative md:hidden">
        <summary class="flex min-h-10 cursor-pointer list-none items-center justify-center rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-3 text-sm font-semibold text-[var(--hawp-ink)] marker:content-none">
          Menu
        </summary>
        <nav class="absolute right-0 top-12 min-w-44 overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-2 shadow-2xl shadow-black/10" aria-label="Mobile navigation">
          {#each [...navItems, { label: 'Downloads', href: '#downloads' }] as item}
            <a href={item.href} class="block rounded-xl px-3 py-2.5 text-sm font-semibold text-[var(--hawp-muted)] transition hover:bg-[var(--hawp-panel)] hover:text-[var(--hawp-ink)]" onclick={closeMobileMenu}>{item.label}</a>
          {/each}
          <a href={githubUrl} class="mt-1 block rounded-xl border-t border-[var(--hawp-line)] px-3 py-2.5 text-sm font-semibold text-[var(--hawp-ink)]" rel="noreferrer" onclick={closeMobileMenu}>GitHub ↗</a>
        </nav>
      </details>

      <a href={githubUrl} class="hidden rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-4 py-2.5 text-sm font-semibold text-[var(--hawp-ink)] transition hover:border-[var(--hawp-red)] hover:bg-[var(--hawp-panel)] md:inline-flex" rel="noreferrer">GitHub</a>
    </div>
  </div>
</header>
