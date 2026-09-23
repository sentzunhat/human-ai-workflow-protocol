<script lang="ts">
  import { downloadUrl, githubUrl, loadLatestRelease, platformDownloads, releaseLoadingVersion, type LatestRelease } from '$lib/content/site';
  import ReleaseVersion from '$lib/components/ReleaseVersion.svelte';
  import { onMount } from 'svelte';

  let release = $state<LatestRelease | null>(null);
  let releaseLoading = $state(true);

  onMount(async () => {
    try {
      release = await loadLatestRelease();
    } finally {
      releaseLoading = false;
    }
  });

  function assetUrl(asset: string): string {
    return release?.assets.find((candidate) => candidate.name === asset)?.browser_download_url ?? downloadUrl(asset);
  }

  let releaseVersion = $derived(release?.tag_name.replace(/^v/, '') ?? releaseLoadingVersion);
</script>

<section id="downloads" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28">
  <div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">Latest release</p>
        <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Get HAWP for your platform.</h2>
        <p class="mt-5 max-w-2xl text-lg leading-8 text-[var(--hawp-muted)]">Download the HAWP CLI from the latest GitHub release. The version and asset links update from the release feed when a new build is published.</p>
      </div>
      <a class="shrink-0 text-sm font-semibold text-[var(--hawp-clay)] underline decoration-[var(--hawp-line-strong)] underline-offset-4 hover:text-[var(--hawp-ink)]" href={release?.html_url ?? `${githubUrl}/releases`}>View release notes ↗</a>
    </div>

    <p class="sr-only" aria-live="polite">
      {releaseLoading ? 'Checking GitHub for the latest HAWP release.' : `Showing HAWP version ${releaseVersion}.`}
    </p>

    <div class="mt-10 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {#each platformDownloads as platform}
        <a class="flex items-center justify-between rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-5 transition hover:border-[var(--hawp-red)] hover:bg-[var(--hawp-panel)]" href={assetUrl(platform.asset)}>
          <span class="flex items-center gap-3">
            <span class="grid size-9 place-items-center rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] text-sm text-[var(--hawp-clay)]" aria-hidden="true">{platform.icon}</span>
            <span>
              <strong class="block text-sm text-[var(--hawp-ink)]">{platform.label}</strong>
              <span class="mt-1 block text-xs text-[var(--hawp-muted-2)]"><ReleaseVersion version={releaseVersion} loading={releaseLoading} /></span>
            </span>
          </span>
          <span class="text-lg text-[var(--hawp-clay)]" aria-hidden="true">↓</span>
        </a>
      {/each}
    </div>
  </div>
</section>
