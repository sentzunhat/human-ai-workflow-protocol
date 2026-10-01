<script lang="ts">
  import { workflowSteps } from '$lib/content/site';
</script>

<section class="border-b border-[var(--hawp-line)] py-24 sm:py-28">
  <div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
    <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">How it fits</p>
    <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Shape, execute, and hand off.</h2>

    <div class="mt-12 grid gap-8 md:grid-cols-3">
      {#each workflowSteps as step}
        <article class="border-l border-[var(--hawp-line-strong)] pl-5">
          <span class="font-mono text-xs text-[var(--hawp-clay)]">{step[0]}</span>
          <h3 class="mt-5 text-xl font-semibold text-[var(--hawp-ink)]">{step[1]}</h3>
          <p class="mt-3 text-sm leading-6 text-[var(--hawp-muted-2)]">{step[2]}</p>
        </article>
      {/each}
    </div>
  </div>
</section>
