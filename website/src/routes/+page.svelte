<script lang="ts">
  import Agents from '$lib/components/Agents.svelte';
  import Benchmarks from '$lib/components/Benchmarks.svelte';
  import Faq from '$lib/components/Faq.svelte';
  import FinalCta from '$lib/components/FinalCta.svelte';
  import Footer from '$lib/components/Footer.svelte';
  import Header from '$lib/components/Header.svelte';
  import Hero from '$lib/components/Hero.svelte';
  import Problem from '$lib/components/Problem.svelte';
  import Protocol from '$lib/components/Protocol.svelte';
  import ReleaseDownloads from '$lib/components/ReleaseDownloads.svelte';
  import Workflow from '$lib/components/Workflow.svelte';
  import { loadLatestRelease, releaseLoadingVersion } from '$lib/content/site';
  import { onMount } from 'svelte';

  let releaseVersion = $state(releaseLoadingVersion);
  let releaseLoading = $state(true);

  onMount(async () => {
    try {
      const release = await loadLatestRelease();
      if (release) releaseVersion = release.tag_name.replace(/^v/, '');
    } finally {
      releaseLoading = false;
    }
  });
</script>

<svelte:head>
  <title>HAWP: Human-AI Workflow Protocol for AI Agent Work</title>
  <meta name="description" content="HAWP is an open-source Human-AI Workflow Protocol for shaping AI agent work, preserving context, reducing drift, and improving handoffs across tools." />
  <link rel="canonical" href="https://hawp.online/" />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content="HAWP" />
  <meta property="og:title" content="HAWP: Shape AI agent work before execution" />
  <meta property="og:description" content="An open-source workflow protocol for clearer human-AI execution, evidence, handoffs, and less drift." />
  <meta property="og:url" content="https://hawp.online/" />
  <meta property="og:image" content="https://hawp.online/social-card.png" />
  <meta property="og:image:type" content="image/png" />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="HAWP: Human-AI Workflow Protocol" />
  <meta name="twitter:description" content="Shape AI agent work before execution. Preserve intent, constraints, evidence, and handoffs." />
  <meta name="twitter:image" content="https://hawp.online/social-card.png" />
  <meta name="application-name" content="HAWP" />
  <meta name="author" content="Sentzunhat" />
  <meta name="subject" content="Human-AI workflow protocol and AI agent work" />
</svelte:head>

<Header />
<main>
  <Hero {releaseVersion} {releaseLoading} />
  <Problem />
  <Protocol />
  <Agents {releaseVersion} {releaseLoading} />
  <Benchmarks {releaseVersion} {releaseLoading} />
  <Workflow />
  <ReleaseDownloads />
  <Faq />
  <FinalCta />
</main>
<Footer />
