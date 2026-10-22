<script lang="ts">
  import { onMount } from 'svelte';

  const measurementId = 'G-1FHPSYNPBS';
  const preferenceKey = 'hawp-analytics-consent-v1';

  let choice = $state<'unknown' | 'accepted' | 'declined'>('unknown');
  let visible = $state(false);

  function loadAnalytics() {
    if (typeof window === 'undefined' || document.querySelector(`script[data-hawp-analytics="${measurementId}"]`)) return;

    const analyticsWindow = window as Window & {
      dataLayer?: unknown[][];
      gtag?: (...args: unknown[]) => void;
    };
    analyticsWindow.dataLayer = analyticsWindow.dataLayer ?? [];
    analyticsWindow.gtag = (...args: unknown[]) => analyticsWindow.dataLayer?.push(args);
    analyticsWindow.gtag('js', new Date());
    analyticsWindow.gtag('config', measurementId);

    const script = document.createElement('script');
    script.async = true;
    script.dataset.hawpAnalytics = measurementId;
    script.src = `https://www.googletagmanager.com/gtag/js?id=${measurementId}`;
    document.head.append(script);
  }

  function choose(value: 'accepted' | 'declined') {
    if (value === 'declined') {
      const analyticsWindow = window as Window & { gtag?: (...args: unknown[]) => void };
      analyticsWindow.gtag?.('consent', 'update', {
        analytics_storage: 'denied',
        ad_storage: 'denied',
        ad_user_data: 'denied',
        ad_personalization: 'denied'
      });
    }
    choice = value;
    localStorage.setItem(preferenceKey, value);
    visible = false;
    if (value === 'accepted') loadAnalytics();
  }

  function openSettings() {
    visible = true;
  }

  onMount(() => {
    const saved = localStorage.getItem(preferenceKey);
    if (saved === 'accepted' || saved === 'declined') {
      choice = saved;
      visible = false;
      if (saved === 'accepted') loadAnalytics();
    } else {
      visible = true;
    }
    window.addEventListener('hawp:analytics-settings', openSettings);
    return () => window.removeEventListener('hawp:analytics-settings', openSettings);
  });
</script>

{#if visible}
  <section class="analytics-consent" aria-label="Analytics preference">
    <div class="analytics-consent__copy">
      <strong>Help improve HAWP</strong>
      <p>Optional Google Analytics helps us understand which pages are useful. It stays off unless you allow it.</p>
    </div>
    <div class="analytics-consent__actions">
      <button type="button" class="analytics-consent__allow" onclick={() => choose('accepted')}>Allow analytics</button>
      <button type="button" class="analytics-consent__decline" onclick={() => choose('declined')}>No thanks</button>
    </div>
  </section>
{/if}

<button class="analytics-settings" type="button" aria-label="Analytics settings" onclick={openSettings}>
  {choice === 'accepted' ? 'Analytics on · settings' : choice === 'declined' ? 'Analytics off · settings' : 'Analytics settings'}
</button>

<style>
  .analytics-consent {
    position: fixed;
    z-index: 50;
    inset: auto 1rem 4.5rem 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.25rem;
    max-width: 760px;
    margin-inline: auto;
    padding: 1rem 1.1rem;
    border: 1px solid var(--hawp-line-strong);
    border-radius: 1rem;
    background: var(--hawp-raised);
    color: var(--hawp-ink);
    box-shadow: 0 18px 50px var(--hawp-shadow);
  }

  .analytics-consent__copy strong { font-size: .9rem; }
  .analytics-consent__copy p { max-width: 32rem; margin: .25rem 0 0; color: var(--hawp-muted); font-size: .78rem; line-height: 1.45; }
  .analytics-consent__actions { display: flex; flex: 0 0 auto; align-items: center; gap: .5rem; }
  .analytics-consent button, .analytics-settings { cursor: pointer; font: inherit; }
  .analytics-consent__allow, .analytics-consent__decline {
    padding: .6rem .8rem;
    border-radius: .65rem;
    font-size: .76rem;
    font-weight: 650;
  }
  .analytics-consent__allow { border: 1px solid var(--hawp-red); background: var(--hawp-red-deep); color: white; }
  .analytics-consent__decline { border: 1px solid var(--hawp-line-strong); background: transparent; color: var(--hawp-ink); }
  .analytics-settings {
    position: fixed;
    z-index: 49;
    right: 1rem;
    bottom: 1rem;
    padding: .45rem .65rem;
    border: 1px solid var(--hawp-line-strong);
    border-radius: 999px;
    background: var(--hawp-raised);
    color: var(--hawp-muted);
    font-size: .7rem;
  }
  .analytics-consent button:focus-visible, .analytics-settings:focus-visible { outline: 2px solid var(--hawp-clay); outline-offset: 3px; }

  @media (max-width: 560px) {
    .analytics-consent { align-items: stretch; flex-direction: column; gap: .85rem; bottom: 4rem; }
    .analytics-consent__actions { justify-content: flex-start; }
  }
</style>
