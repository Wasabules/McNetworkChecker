<script>
  import { onMount } from 'svelte';
  import { t } from './i18n.js';
  import { UploadPaste, GetPasteServices } from '../../wailsjs/go/main/App';

  export let report = '';

  let copied = false;
  let services = [];
  let menuOpen = false;
  let uploading = false;
  let pasteURL = '';
  let pasteError = '';

  onMount(async () => {
    try { services = await GetPasteServices(); } catch(_) {}
  });

  // Reset paste state when a new report arrives
  let prevReport = '';
  $: if (report !== prevReport) {
    prevReport = report;
    pasteURL = '';
    pasteError = '';
    copied = false;
    menuOpen = false;
  }

  function copy() {
    navigator.clipboard.writeText(report).then(() => {
      copied = true; setTimeout(() => { copied = false; }, 2000);
    });
  }

  function toggleMenu() {
    if (uploading) return;
    menuOpen = !menuOpen;
  }

  function closeMenu() {
    menuOpen = false;
  }

  async function upload(serviceId) {
    menuOpen = false;
    uploading = true;
    pasteURL = '';
    pasteError = '';
    try {
      const result = await UploadPaste(serviceId, report);
      if (result.error) {
        pasteError = result.error;
      } else {
        pasteURL = result.url;
        // Auto-copy the URL
        navigator.clipboard.writeText(result.url).catch(() => {});
      }
    } catch (e) {
      pasteError = e.toString();
    }
    uploading = false;
  }

  function openURL() {
    if (pasteURL) window.open(pasteURL, '_blank');
  }

  // Close menu on outside click
  function onWindowClick(e) {
    if (menuOpen) menuOpen = false;
  }
</script>

<svelte:window on:click={onWindowClick} />

{#if report}
  <div class="mc-panel report-section">
    <div class="report-top">
      <span class="report-title">{$t('report.title')}</span>
      <div class="report-actions">
        <!-- Copy button -->
        <button class="mc-btn mc-btn-blue" on:click={copy}>
          {copied ? $t('report.copied') : $t('btn.copy')}
        </button>

        <!-- Upload split button -->
        <div class="paste-split" on:click|stopPropagation>
          {#if pasteURL}
            <button class="mc-btn mc-btn-green paste-result" on:click={openURL}>
              {$t('paste.open')}
            </button>
          {:else}
            <button class="mc-btn mc-btn-gold paste-main" on:click={toggleMenu} disabled={uploading}>
              {uploading ? $t('paste.uploading') : $t('paste.upload')} &#9662;
            </button>
          {/if}

          {#if menuOpen}
            <div class="paste-menu">
              {#each services as svc}
                <button class="paste-menu-item" on:click={() => upload(svc.id)}>
                  <span class="paste-svc-name">{svc.name}</span>
                  <span class="paste-svc-desc">{svc.desc}</span>
                </button>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Paste result feedback -->
    {#if pasteURL}
      <div class="paste-feedback paste-ok">
        {$t('paste.success')} <a href={pasteURL} target="_blank" rel="noopener">{pasteURL}</a>
      </div>
    {/if}
    {#if pasteError}
      <div class="paste-feedback paste-err">
        {$t('paste.error', { err: pasteError })}
      </div>
    {/if}

    <div class="mc-console report-console">
      <pre>{report}</pre>
    </div>
  </div>
{/if}

<style>
  .report-section { margin-top: 4px; }
  .report-console pre { margin: 0; white-space: pre-wrap; word-break: break-all; line-height: 1.55; color: #D0D0D0; }

  .report-actions { display: flex; gap: 6px; align-items: center; }

  /* Split button container */
  .paste-split { position: relative; }

  .paste-main, .paste-result {
    font-size: 10px;
  }

  /* Dropdown menu */
  .paste-menu {
    position: absolute;
    right: 0;
    top: calc(100% + 4px);
    background: #3A3A3A;
    border: 2px solid;
    border-color: #555 #222 #222 #555;
    min-width: 200px;
    z-index: 100;
    display: flex;
    flex-direction: column;
  }

  .paste-menu-item {
    display: flex;
    flex-direction: column;
    gap: 1px;
    padding: 7px 12px;
    background: none;
    border: none;
    border-bottom: 1px solid #2A2A2A;
    color: #DDD;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    font-size: 12px;
  }
  .paste-menu-item:last-child { border-bottom: none; }
  .paste-menu-item:hover { background: #4A4A4A; }

  .paste-svc-name { font-weight: bold; color: #FFAA00; font-size: 11px; }
  .paste-svc-desc { color: #888; font-size: 10px; }

  /* Feedback bar */
  .paste-feedback {
    padding: 6px 12px;
    font-family: 'Cascadia Mono', 'Consolas', monospace;
    font-size: 11px;
    border-top: 1px solid #2A2A2A;
  }
  .paste-ok { color: #55FF55; }
  .paste-ok a { color: #88CCFF; text-decoration: underline; }
  .paste-ok a:hover { color: #BBDDFF; }
  .paste-err { color: #FF5555; }
</style>
