<script lang="ts">
  import { onMount } from 'svelte';
  import { applyTheme, getThemePreference, watchSystemTheme, type ThemePreference } from '$lib/theme';

  const options: { value: ThemePreference; label: string }[] = [
    { value: 'light', label: 'Light' },
    { value: 'dark', label: 'Dark' },
    { value: 'system', label: 'System' }
  ];

  let preference: ThemePreference = 'system';

  function choose(value: ThemePreference) {
    preference = value;
    applyTheme(value);
  }

  onMount(() => {
    preference = getThemePreference();
    applyTheme(preference);

    return watchSystemTheme(() => {
      if (preference === 'system') applyTheme('system');
    });
  });
</script>

<div class="flex items-center rounded-full border border-[var(--hawp-line)] bg-[var(--hawp-raised)]/90 p-1 shadow-sm backdrop-blur" aria-label="Theme preference">
  {#each options as option}
    <button
      type="button"
      class={`flex min-h-8 items-center gap-1.5 rounded-full px-2.5 text-xs font-semibold transition sm:px-3 ${preference === option.value ? 'bg-[var(--hawp-panel-soft)] text-[var(--hawp-ink)] shadow-sm' : 'text-[var(--hawp-muted-2)] hover:text-[var(--hawp-ink)]'}`}
      aria-pressed={preference === option.value}
      aria-label={`Use ${option.label.toLowerCase()} theme`}
      title={`${option.label} theme`}
      onclick={() => choose(option.value)}
    >
      {#if option.value === 'light'}
        <svg aria-hidden="true" viewBox="0 0 24 24" class="size-3.5" fill="none" stroke="currentColor" stroke-width="1.8">
          <circle cx="12" cy="12" r="3.25"></circle>
          <path d="M12 2.5v2.2M12 19.3v2.2M4.7 4.7l1.55 1.55M17.75 17.75l1.55 1.55M2.5 12h2.2M19.3 12h2.2M4.7 19.3l1.55-1.55M17.75 6.25l1.55-1.55"></path>
        </svg>
      {:else if option.value === 'dark'}
        <svg aria-hidden="true" viewBox="0 0 24 24" class="size-3.5" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M20.3 15.4A8.4 8.4 0 0 1 8.6 3.7 8.4 8.4 0 1 0 20.3 15.4Z"></path>
        </svg>
      {:else}
        <svg aria-hidden="true" viewBox="0 0 24 24" class="size-3.5" fill="none" stroke="currentColor" stroke-width="1.8">
          <rect x="3.5" y="4.5" width="17" height="12" rx="2"></rect>
          <path d="M9 20h6M12 16.5V20"></path>
        </svg>
      {/if}
      <span class="hidden lg:inline">{option.label}</span>
    </button>
  {/each}
</div>
