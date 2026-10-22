<script lang="ts">
  import { githubUrl } from '$lib/content/site';
</script>

<section class="py-24 sm:py-28">
  <div class="mx-auto max-w-[1180px] px-5 sm:px-8">
    <div class="site-grid overflow-hidden rounded-[2rem] border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-6 py-12 text-center shadow-2xl shadow-black/20 sm:px-12 sm:py-16">
      <h2 class="mx-auto max-w-3xl text-4xl font-bold tracking-[-0.04em] text-[var(--hawp-ink)] sm:text-5xl">Give AI agent work a shape that survives the session.</h2>
      <p class="mx-auto mt-5 max-w-2xl text-lg leading-8 text-[var(--hawp-muted)]">HAWP is open source, Apache-2.0 licensed, and designed to stay readable even when the agent changes.</p>
      <div class="mt-8 flex flex-wrap justify-center gap-3">
        <a href={githubUrl} class="rounded-xl border border-[var(--hawp-red)] bg-[var(--hawp-red-deep)] px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-[var(--hawp-red)]">Read the repository ↗</a>
        <a href={`${githubUrl}#install--pick-your-agent`} class="rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] px-6 py-3.5 text-sm font-semibold text-[var(--hawp-ink)]">Get started</a>
      </div>
    </div>
  </div>
</section>
