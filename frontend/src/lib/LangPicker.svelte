<script>
  import { locale, LOCALES } from './i18n.js';
  import { SetLocale } from '../../wailsjs/go/main/App';

  function pick(code) {
    locale.set(code);
    SetLocale(code);
  }
</script>

<div class="lang-picker">
  {#each LOCALES as l}
    <button
      class="lang-btn"
      class:active={$locale === l.code}
      on:click={() => pick(l.code)}
      title={l.label}
    >{l.flag}</button>
  {/each}
</div>

<style>
  .lang-picker {
    position: absolute;
    top: 12px;
    left: 12px;
    display: flex;
    gap: 4px;
  }
  .lang-btn {
    background: #333;
    border: 1px solid #444;
    padding: 2px 6px;
    cursor: pointer;
    font-size: 16px;
    line-height: 1;
    opacity: 0.5;
    transition: opacity 0.15s;
  }
  .lang-btn:hover { opacity: 0.8; }
  .lang-btn.active {
    opacity: 1;
    border-color: #FFAA00;
    background: #3A3000;
  }
</style>
