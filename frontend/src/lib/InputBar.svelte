<script>
  import { t } from './i18n.js';
  export let running = false;
  export let onStart;
  export let onStop;

  let address = '';
  let touched = false;

  const RE_DOMAIN = /^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)+$/;
  const RE_IPV4 = /^(\d{1,3}\.){3}\d{1,3}$/;

  function validate(raw) {
    const v = raw.trim();
    if (!v) return '';
    let host = v, port = null;
    const c = v.lastIndexOf(':');
    if (c > 0) { const p = v.substring(c+1); if (/^\d+$/.test(p)) { port = parseInt(p,10); host = v.substring(0,c); } }
    if (port !== null && (port < 1 || port > 65535)) return $t('val.portInvalid', { port });
    if (RE_IPV4.test(host)) { if (host.split('.').map(Number).some(o => o > 255)) return $t('val.octet'); return ''; }
    if (RE_DOMAIN.test(host)) return host.length > 253 ? $t('val.tooLong') : '';
    if (/\s/.test(host)) return $t('val.spaces');
    if (/[^a-zA-Z0-9._:-]/.test(host)) return $t('val.badChars');
    if (!host.includes('.')) return $t('val.incomplete');
    return $t('val.invalid');
  }

  $: error = validate(address);
  $: canStart = address.trim() && !error && !running;

  function start() { touched = true; if (!canStart) return; onStart(address.trim()); }
  function onKeydown(e) { if (e.key === 'Enter') start(); }
</script>

<div class="input-bar">
  <div class="mc-input-wrap" class:input-error={touched && error}>
    <input type="text" bind:value={address} placeholder={$t('input.placeholder')}
      on:keydown={onKeydown} on:input={() => { touched = true; }}
      disabled={running} spellcheck="false" autocomplete="off" />
  </div>
  {#if running}
    <button class="mc-btn mc-btn-red" on:click={onStop}>{$t('btn.stop')}</button>
  {:else}
    <button class="mc-btn mc-btn-green" on:click={start} disabled={!canStart}>{$t('btn.start')}</button>
  {/if}
</div>
{#if touched && error}<div class="validation-msg">{error}</div>{/if}

<style>
  .input-bar { display: flex; gap: 8px; margin-bottom: 18px; }
</style>
