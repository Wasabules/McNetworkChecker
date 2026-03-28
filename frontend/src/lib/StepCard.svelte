<script>
  import { slide } from 'svelte/transition';
  import { t } from './i18n.js';
  import Console from './Console.svelte';
  import MinecraftInfo from './MinecraftInfo.svelte';
  import { SkipStep } from '../../wailsjs/go/main/App';

  export let step;

  // Auto-expand when running, collapsed by default otherwise
  let manualToggle = null; // null = no manual override
  $: expanded = manualToggle !== null
    ? manualToggle
    : step.status === 'running'; // auto-open when running

  // Reset manual toggle when step changes status (so it auto-opens on run)
  let prevStatus = step.status;
  $: if (step.status !== prevStatus) {
    prevStatus = step.status;
    manualToggle = null;
  }

  function toggle() {
    manualToggle = !expanded;
  }

  let copied = false;
  function copyLogs() {
    const text = (step.stepLogs || []).join('\n');
    if (!text) return;
    navigator.clipboard.writeText(text).then(() => {
      copied = true; setTimeout(() => { copied = false; }, 1500);
    });
  }

  $: hasContent = step.showConsole || (step.status === 'skipped' && (step.reason === 'auto' || step.reason === 'user'))
      || (step.id === 'minecraft' && step.result && step.result.success);
</script>

<div class="mc-panel step-card {step.cls}" class:active={step.showConsole || step.status === 'done'}>
  <!-- Header — clickable to toggle -->
  <button class="step-header" on:click={toggle}>
    {#if hasContent}
      <span class="step-chevron" class:open={expanded}>&#9656;</span>
    {:else}
      <span class="step-chevron-spacer"></span>
    {/if}
    <span class="step-num">[{step.num}]</span>
    <span class="step-label">{step.label}</span>
    <span class="step-sum">{step.summary}</span>
    {#if step.status === 'running'}
      <button class="mc-btn-skip" on:click|stopPropagation={() => SkipStep()}>{$t('btn.skip')}</button>
    {/if}
    {#if step.stepLogs && step.stepLogs.length > 0 && step.status !== 'running'}
      <button class="mc-btn-copy" on:click|stopPropagation={copyLogs}>
        {copied ? $t('btn.copied') : $t('btn.copy')}
      </button>
    {/if}
    <span class="step-badge badge-{step.cls}">{step.badgeText}</span>
  </button>

  <!-- Collapsible body -->
  {#if expanded}
    <div transition:slide={{ duration: 150 }}>
      {#if step.status === 'skipped' && step.reason === 'auto'}
        <div class="skip-msg"><span class="log-gt">&gt;</span> {$t('skip.auto')}</div>
      {:else if step.status === 'skipped' && step.reason === 'user'}
        <div class="skip-msg skip-user"><span class="log-gt">&gt;</span> {$t('skip.user')}</div>
      {/if}

      {#if step.showConsole}
        <Console lines={step.stepLogs} isRunning={step.status === 'running'} />
      {/if}

      {#if step.id === 'minecraft' && step.result && step.result.success}
        <MinecraftInfo result={step.result} />
      {/if}
    </div>
  {/if}
</div>
