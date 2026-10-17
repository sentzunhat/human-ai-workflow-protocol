<script lang="ts">
  import { githubUrl } from '$lib/content/site';
</script>

<footer class="border-t border-[var(--hawp-line)] py-10">
  <div class="mx-auto flex max-w-[1280px] flex-col gap-5 px-5 text-sm text-[var(--hawp-muted-2)] sm:px-8 md:flex-row md:items-center md:justify-between lg:px-10">
    <span>HAWP · Human-AI Workflow Protocol · Sentzunhat</span>
    <div class="flex gap-5"><a class="hover:text-[var(--hawp-ink)]" href={githubUrl}>GitHub</a><a class="hover:text-[var(--hawp-ink)]" href={`${githubUrl}/blob/main/LICENSE`}>Apache-2.0</a><a class="hover:text-[var(--hawp-ink)]" href={`${githubUrl}/discussions`}>Discussions</a></div>
  </div>
</footer>
