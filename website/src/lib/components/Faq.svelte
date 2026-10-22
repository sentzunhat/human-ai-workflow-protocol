<script lang="ts">
  import { faqs } from '$lib/content/site';
</script>

<section id="faq" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28">
  <div class="mx-auto max-w-4xl px-5 sm:px-8">
    <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">FAQ</p>
    <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Small protocol, deliberate boundaries.</h2>

    <div class="mt-10 divide-y divide-[var(--hawp-line)] border-y border-[var(--hawp-line)]">
      {#each faqs as faq}
        <article class="py-6">
          <h3 class="text-lg font-semibold text-[var(--hawp-ink)]">{faq[0]}</h3>
          <p class="mt-2 text-sm leading-6 text-[var(--hawp-muted-2)]">{faq[1]}</p>
        </article>
      {/each}
    </div>
  </div>
</section>
