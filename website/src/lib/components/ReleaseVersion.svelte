<script lang="ts">
  let {
    version,
    loading = false,
    variant = 'inline'
  }: {
    version: string;
    loading?: boolean;
    variant?: 'inline' | 'badge';
  } = $props();
</script>

<span class:release-version-loading={loading} class:release-version-badge={variant === 'badge'} class="release-version">
  <span class="release-version-text">v{version}</span>
  {#if loading}
    <span class="release-skeleton-sheen" aria-hidden="true"></span>
  {/if}
</span>

<style>
  .release-version {
    position: relative;
    display: inline-block;
    width: fit-content;
    min-width: 3.7rem;
    overflow: hidden;
    border-radius: 999px;
    isolation: isolate;
  }

  .release-version-loading {
    color: var(--hawp-ink);
    background: color-mix(in srgb, var(--hawp-panel-soft) 72%, transparent);
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--hawp-peach) 22%, transparent);
  }

  .release-version-badge {
    padding: 0.18rem 0.45rem;
  }

  .release-version-text {
    position: relative;
    z-index: 1;
  }

  .release-skeleton-sheen {
    position: absolute;
    inset: 0;
    z-index: 0;
    pointer-events: none;
    background: linear-gradient(105deg, transparent 22%, color-mix(in srgb, var(--hawp-peach) 32%, transparent) 46%, color-mix(in srgb, white 28%, transparent) 51%, transparent 72%);
    background-size: 230% 100%;
    animation: release-skeleton-shimmer 1.65s ease-in-out infinite;
  }

  @keyframes release-skeleton-shimmer {
    from { background-position: 100% 0; }
    to { background-position: -100% 0; }
  }

  @media (prefers-reduced-motion: reduce) {
    .release-skeleton-sheen { animation: none; }
  }
</style>
