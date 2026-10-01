<script lang="ts">
  import { githubUrl } from '$lib/content/site';
  import ReleaseVersion from '$lib/components/ReleaseVersion.svelte';
  import { onMount } from 'svelte';

  let { releaseVersion, releaseLoading }: { releaseVersion: string; releaseLoading: boolean } = $props();

  type NodePosition = { x: number; y: number; label: string; detail: string };

  const nodePositions = $state<NodePosition[]>([
    { x: 50, y: 10, label: 'Input', detail: 'request' },
    { x: 86, y: 39, label: 'Context', detail: 'background' },
    { x: 75, y: 88, label: 'Mission', detail: 'objective' },
    { x: 25, y: 88, label: 'Constraints', detail: 'limits' },
    { x: 14, y: 39, label: 'Output', detail: 'done state' }
  ]);

  let visual: HTMLDivElement;
  let connectionPaths = $state<string[]>([]);
  let draggingIndex = $state(-1);

  function connectionPath(node: NodePosition, index: number): string {
    const dx = node.x - 50;
    const dy = node.y - 50;
    const distance = Math.max(Math.hypot(dx, dy), 1);
    const perpendicularX = -dy / distance;
    const perpendicularY = dx / distance;
    const sway = [2.4, -2.8, 2.1, -2.1, 2.8][index];
    const c1x = 50 + dx * 0.34 + perpendicularX * sway;
    const c1y = 50 + dy * 0.34 + perpendicularY * sway;
    const c2x = node.x - dx * 0.34 + perpendicularX * sway;
    const c2y = node.y - dy * 0.34 + perpendicularY * sway;
    return `M 50 50 C ${c1x.toFixed(2)} ${c1y.toFixed(2)}, ${c2x.toFixed(2)} ${c2y.toFixed(2)}, ${node.x.toFixed(2)} ${node.y.toFixed(2)}`;
  }

  function updateConnections() {
    connectionPaths = nodePositions.map(connectionPath);
  }

  function startDrag(event: PointerEvent, index: number) {
    if (event.pointerType === 'mouse' && event.button !== 0) return;
    draggingIndex = index;
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
    event.preventDefault();
  }

  function drag(event: PointerEvent) {
    if (draggingIndex < 0 || !visual) return;
    const bounds = visual.getBoundingClientRect();
    const x = Math.max(8, Math.min(92, ((event.clientX - bounds.left) / bounds.width) * 100));
    const y = Math.max(8, Math.min(92, ((event.clientY - bounds.top) / bounds.height) * 100));
    nodePositions[draggingIndex].x = x;
    nodePositions[draggingIndex].y = y;
    updateConnections();
  }

  function stopDrag(event: PointerEvent) {
    if (draggingIndex < 0) return;
    const target = event.currentTarget as HTMLElement;
    if (target.hasPointerCapture(event.pointerId)) target.releasePointerCapture(event.pointerId);
    draggingIndex = -1;
  }

  onMount(() => {
    updateConnections();
    const resizeObserver = new ResizeObserver(updateConnections);
    resizeObserver.observe(visual);
    return () => resizeObserver.disconnect();
  });
</script>

<section id="top" class="site-grid relative overflow-hidden border-b border-[var(--hawp-line)]">
  <div class="pointer-events-none absolute -right-40 -top-52 size-[38rem] rounded-full bg-[var(--hawp-red-deep)]/16 blur-3xl"></div>

  <div class="mx-auto grid min-h-[710px] max-w-[1280px] items-center gap-14 px-5 py-20 sm:px-8 lg:grid-cols-[1.08fr_.92fr] lg:px-10 lg:py-24">
    <div class="relative z-10">
      <div class="mb-8 flex items-center gap-3 text-[11px] font-semibold tracking-[0.08em] text-[var(--hawp-peach)] uppercase">
        <span class="size-3 rounded-full bg-[var(--hawp-red)]"></span>
        <ReleaseVersion version={releaseVersion} loading={releaseLoading} variant="badge" /> · Open protocol · Markdown-first · MCP-ready
      </div>

      <h1 class="max-w-[760px] text-5xl font-bold leading-[1.02] tracking-[-0.045em] text-[var(--hawp-ink)] sm:text-6xl lg:text-[5rem]">
        Shape the work before the model <span class="text-[var(--hawp-peach)]">touches it.</span>
      </h1>

      <p class="mt-7 max-w-[730px] text-lg leading-7 text-[var(--hawp-muted)] sm:text-xl sm:leading-8">
        HAWP gives AI agent work a stable shape for intent, context, mission, constraints, and output, so agents can execute with less drift and hand off without starting from zero.
      </p>

      <div class="mt-8 flex flex-wrap gap-3">
        <a href={githubUrl} class="rounded-xl border border-[var(--hawp-red)] bg-[var(--hawp-red-deep)] px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-[var(--hawp-red)]">Explore on GitHub ↗</a>
        <a href="#protocol" class="rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] px-6 py-3.5 text-sm font-semibold text-[var(--hawp-ink)] transition hover:border-[var(--hawp-line-strong)]">See the protocol ↓</a>
      </div>

      <p class="mt-6 text-xs text-[var(--hawp-muted-2)]">Apache-2.0 · Go CLI · Static Markdown artifacts · No lock-in</p>
    </div>

    <div
      class="protocol-visual"
      bind:this={visual}
      role="group"
      aria-label="Five HAWP fields converging into a stable work packet"
      onpointermove={drag}
      onpointerup={stopDrag}
      onpointercancel={stopDrag}
    >
      <div class="orbit one"></div><div class="orbit two"></div><div class="orbit three"></div>
      <svg class="connections" viewBox="0 0 100 100" aria-hidden="true" focusable="false">
        {#each connectionPaths as path, index}
          <path class="connector-shadow" d={path} pathLength="1" style={`--connector-delay: ${index * -0.35}s`} />
          <path class="connector" d={path} pathLength="1" style={`--connector-delay: ${index * -0.35}s`} />
        {/each}
      </svg>
      <div class="core"><span>WORK</span></div>
      {#each nodePositions as node, index}
        <button
          type="button"
          class:dragging={draggingIndex === index}
          class="node"
          style={`--node-x: ${node.x}%; --node-y: ${node.y}%; --float-delay: ${index * -0.7}s`}
          aria-describedby="drag-instruction"
          onpointerdown={(event) => startDrag(event, index)}
        >
          <span class="node-content"><strong>{node.label}</strong>{node.detail}</span>
        </button>
      {/each}
      <p id="drag-instruction" class="drag-hint">Drag a field to reshape the work packet</p>
    </div>
  </div>
</section>
