<script lang="ts">
  import { benchmarkMetrics } from '$lib/content/site';
  import ReleaseVersion from '$lib/components/ReleaseVersion.svelte';

  let { releaseVersion, releaseLoading }: { releaseVersion: string; releaseLoading: boolean } = $props();
</script>

<section id="benchmarks" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28">
  <div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
    <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase"><ReleaseVersion version={releaseVersion} loading={releaseLoading} /> repository evidence</p>
    <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Less context. Faster retrieval. Explicit uncertainty.</h2>
    <p class="mt-5 max-w-3xl text-lg leading-8 text-[var(--hawp-muted)]">Current repository evidence covers search retrieval, context shaping, and structured intake shaping. These are HAWP-local measurements intended to make workflow claims inspectable.</p>

    <div class="mt-12 grid gap-3 md:grid-cols-3">
      {#each benchmarkMetrics as metric}
        <article class="rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-7">
          <strong class="block text-4xl font-bold tracking-[-0.04em] text-[var(--hawp-peach)]">{metric[0]}</strong>
          <span class="mt-4 block text-sm leading-6 text-[var(--hawp-muted-2)]">{metric[1]}</span>
        </article>
      {/each}
    </div>
    <p class="mt-5 text-xs text-[var(--hawp-muted-2)]">Benchmarks are project-local measurements, not universal model performance claims.</p>
  </div>
</section>
