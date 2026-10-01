<script lang="ts">
  import { providers } from '$lib/content/site';
  import ReleaseVersion from '$lib/components/ReleaseVersion.svelte';

  let { releaseVersion, releaseLoading }: { releaseVersion: string; releaseLoading: boolean } = $props();
</script>

<section id="agents" class="border-b border-[var(--hawp-line)] py-24 sm:py-28">
  <div class="mx-auto grid max-w-[1280px] items-center gap-14 px-5 sm:px-8 lg:grid-cols-2 lg:px-10">
    <div class="overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-alt)] shadow-2xl shadow-black/20">
      <div class="flex gap-2 border-b border-[var(--hawp-line)] px-5 py-4"><i class="size-2.5 rounded-full bg-[var(--hawp-red)]"></i><i class="size-2.5 rounded-full bg-[var(--hawp-line-strong)]"></i><i class="size-2.5 rounded-full bg-[var(--hawp-muted-2)]"></i></div>
      <div class="p-6 font-mono text-sm leading-7 text-[var(--hawp-muted)] sm:p-8">
        <div><span class="text-[var(--hawp-clay)]">$</span> hawp work new</div>
        <div class="text-[var(--hawp-clay)]">✓ work item scaffolded</div>
        <div class="text-[var(--hawp-muted-2)]">&nbsp;&nbsp;input&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; → captured</div>
        <div class="text-[var(--hawp-muted-2)]">&nbsp;&nbsp;context&nbsp;&nbsp;&nbsp;&nbsp; → bounded</div>
        <div class="text-[var(--hawp-muted-2)]">&nbsp;&nbsp;mission&nbsp;&nbsp;&nbsp;&nbsp; → explicit</div>
        <div class="text-[var(--hawp-muted-2)]">&nbsp;&nbsp;constraints → preserved</div>
        <div class="text-[var(--hawp-muted-2)]">&nbsp;&nbsp;output&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; → verifiable</div>
        <div class="mt-5"><span class="text-[var(--hawp-clay)]">$</span> hawp mcp</div>
        <div class="text-[var(--hawp-clay)]">✓ hawp_work_intake · compound shaping</div>
        <div class="text-[var(--hawp-clay)]">✓ hawp_work_doc · canonical work paths</div>
      </div>
    </div>

    <div>
      <p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">Tool-independent by design</p>
      <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">The work survives the agent.</h2>
      <p class="mt-5 text-lg leading-8 text-[var(--hawp-muted)]">HAWP uses Markdown artifacts and exposes search and workflows through MCP, so the same project shape can move across coding agents without getting trapped in one vendor's memory or chat history.</p>
      <p class="mt-4 text-sm leading-6 text-[var(--hawp-muted-2)]">In <ReleaseVersion version={releaseVersion} loading={releaseLoading} />, compound intake can retrieve relevant context, reshape the request through a local LLM, and return structured mission, constraints, output specification, and done signal states without creating a work item automatically.</p>
      <div class="mt-7 flex flex-wrap gap-2">
        {#each providers as provider}<span class="rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-3 py-1.5 text-xs text-[var(--hawp-muted)]">{provider}</span>{/each}
      </div>
    </div>
  </div>
</section>
