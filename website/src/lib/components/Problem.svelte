<script lang="ts">
  import { problemItems } from '$lib/content/site';
</script>

<section class="border-b border-[var(--hawp-line)] py-24 sm:py-28">
  <div class="mx-auto grid max-w-[1280px] gap-14 px-5 sm:px-8 lg:grid-cols-[.92fr_1.08fr] lg:px-10">
    <div>
      <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">The failure mode</p>
      <h2 class="mt-5 max-w-xl text-4xl font-bold leading-[1.08] tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">AI work often loses shape between the idea and execution.</h2>
      <p class="mt-6 max-w-xl text-lg leading-8 text-[var(--hawp-muted)]">The expensive part is not always generation. It is re-explaining intent, rediscovering constraints, recovering decisions, and correcting drift across sessions and tools.</p>
    </div>

    <div class="divide-y divide-[var(--hawp-line)] overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)]">
      {#each problemItems as item}
        <article class="grid grid-cols-[46px_1fr] gap-4 p-6 sm:p-7">
          <span class="font-mono text-sm text-[var(--hawp-clay)]">{item[0]}</span>
          <div><h3 class="text-xl font-semibold text-[var(--hawp-ink)]">{item[1]}</h3><p class="mt-2 text-sm leading-6 text-[var(--hawp-muted-2)]">{item[2]}</p></div>
        </article>
      {/each}
    </div>
  </div>
</section>
