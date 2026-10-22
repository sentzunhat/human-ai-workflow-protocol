<script lang="ts">
  import { protocolFields } from '$lib/content/site';
</script>

<section id="protocol" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28">
  <div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
    <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">The protocol</p>
    <h2 class="mt-5 max-w-3xl text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Five fields. One portable shape for the work.</h2>
    <p class="mt-5 max-w-3xl text-lg leading-8 text-[var(--hawp-muted)]">HAWP keeps the protocol intentionally small. The artifact remains readable by humans, useful to agents, and portable across tools.</p>

    <div class="mt-12 grid gap-3 md:grid-cols-2 lg:grid-cols-5">
      {#each protocolFields as field}
        <article class="min-h-56 rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-6">
          <span class="font-mono text-xs text-[var(--hawp-clay)]">{field[0]}</span>
          <h3 class="mt-8 text-xl font-semibold text-[var(--hawp-ink)]">{field[1]}</h3>
          <p class="mt-3 text-sm leading-6 text-[var(--hawp-muted-2)]">{field[2]}</p>
        </article>
      {/each}
    </div>
  </div>
</section>
