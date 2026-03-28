<script>
  import { t } from './i18n.js';
  import MotdRender from './MotdRender.svelte';
  export let result;
</script>

<div class="mc-server-info">
  {#if result.favicon}
    <img src={result.favicon} alt="Favicon" class="mc-favicon" />
  {/if}
  <div class="mc-server-text">
    <MotdRender spans={result.motdSpans} fallback={result.motd} />
    <div class="mc-meta">
      <span>{result.version}</span>
      <span class="mc-sep">|</span>
      <span>{result.playersOnline}/{result.playersMax} {$t('sum.players')}</span>
      <span class="mc-sep">|</span>
      <span>{$t('mc.protocol')} {result.protocolVersion}</span>
    </div>
    {#if result.playersSample && result.playersSample.length > 0}
      <div class="mc-players">
        {#each result.playersSample as p}
          <span class="mc-player-tag">{p.name}</span>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .mc-server-text { flex: 1; min-width: 0; }
</style>
