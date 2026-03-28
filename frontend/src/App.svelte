<script>
  import { onMount } from 'svelte';
  import { fly } from 'svelte/transition';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import { RunDiagnostic, StopDiagnostic, GetSystemLocale, SetLocale } from '../wailsjs/go/main/App';
  import { locale, t } from './lib/i18n.js';

  import LangPicker from './lib/LangPicker.svelte';
  import InputBar from './lib/InputBar.svelte';
  import ProgressBar from './lib/ProgressBar.svelte';
  import StepCard from './lib/StepCard.svelte';
  import Report from './lib/Report.svelte';

  let running = false;
  let started = false;
  let steps = {};
  let logs = {};
  let report = '';

  const STEP_IDS = ['dns','ping','traceroute','tcp','tcpTraceroute','minecraft'];
  const STEP_KEYS = {
    dns: 'step.dns', ping: 'step.ping', traceroute: 'step.traceroute',
    tcp: 'step.tcp', tcpTraceroute: 'step.tcpTraceroute', minecraft: 'step.minecraft',
  };

  // Detect OS language at startup
  onMount(async () => {
    try {
      const lang = await GetSystemLocale();
      if (lang) {
        locale.set(lang);
        await SetLocale(lang);
      }
    } catch(_) {}
  });

  $: stepStates = STEP_IDS.map((id, i) => {
    const data = steps[id] || { status: 'pending', result: null, reason: null };
    const stepLogs = logs[id] || [];
    const reason = data.reason || null;
    let cls = 'pending', badgeText = $t('badge.pending');

    if (data.status === 'skipped') {
      if (reason === 'stopped') { cls = 'stopped'; badgeText = $t('badge.stopped'); }
      else { cls = 'skipped'; badgeText = $t('badge.skip'); }
    } else if (data.status === 'running') {
      cls = 'running'; badgeText = $t('badge.running');
    } else if (data.status === 'done') {
      if (data.result && data.result.success) { cls = 'success'; badgeText = $t('badge.ok'); }
      else { cls = 'error'; badgeText = $t('badge.fail'); }
    }

    const showConsole = data.status === 'running' || data.status === 'done'
        || (data.status === 'skipped' && reason === 'user' && stepLogs.length > 0);

    return {
      id, num: i + 1, label: $t(STEP_KEYS[id]),
      cls, badgeText, stepLogs, reason, showConsole,
      summary: computeSummary(id, data),
      result: data.result, status: data.status,
    };
  });

  $: completedCount = stepStates.filter(s => s.cls !== 'pending').length;

  function computeSummary(id, data) {
    if (!data || data.status !== 'done' || !data.result) return '';
    const r = data.result;
    switch (id) {
      case 'dns': return r.success ? `${r.duration}ms` : '';
      case 'ping': return r.success
        ? `${r.received}/${r.sent} - ${$t('sum.avg')} ${(r.avgMs||0).toFixed(0)}ms`
        : `${r.received||0}/${r.sent||0} ${$t('sum.packets')}`;
      case 'traceroute':
      case 'tcpTraceroute': return r.success ? `${r.hops} ${$t('sum.hops')}` : '';
      case 'tcp': return r.success ? `${r.latencyMs}ms` : '';
      case 'minecraft': return r.success ? `${r.version} - ${r.playersOnline}/${r.playersMax}` : '';
    }
    return '';
  }

  function onStart(address) {
    running = true; started = true; report = '';
    steps = {}; logs = {};
    for (const id of STEP_IDS) {
      steps[id] = { status: 'pending', result: null, reason: null };
      logs[id] = [];
    }
    steps = steps; logs = logs;
    RunDiagnostic(address);
  }

  function onStop() { StopDiagnostic(); }

  EventsOn('check:update', (data) => {
    if (data && data.step) {
      const prev = steps[data.step];
      steps[data.step] = {
        status: data.status,
        result: (data.result != null) ? data.result : (prev ? prev.result : null),
        reason: data.reason || null,
      };
      steps = steps;
    }
  });

  EventsOn('check:log', (data) => {
    if (data && data.step && data.line != null) {
      if (!logs[data.step]) logs[data.step] = [];
      logs[data.step] = [...logs[data.step], data.line];
      logs = logs;
    }
  });

  EventsOn('check:report', (data) => { report = data; });

  EventsOn('check:finished', () => {
    running = false;
    let changed = false;
    for (const id of STEP_IDS) {
      const st = steps[id];
      if (st && (st.status === 'running' || st.status === 'pending')) {
        steps[id] = { status: 'skipped', result: st.result, reason: 'stopped' };
        changed = true;
      }
    }
    if (changed) steps = steps;
  });
</script>

<main>
  <LangPicker />
  <header>
    <h1 class="mc-title">{$t('app.title')}</h1>
    <p class="mc-sub">{$t('app.subtitle')}</p>
  </header>

  <InputBar {running} {onStart} {onStop} />

  {#if started}
    <div in:fly={{ y: -10, duration: 200 }}>
      <ProgressBar completed={completedCount} total={STEP_IDS.length} />
    </div>
    <div class="steps">
      {#each stepStates as step (step.id)}
        <StepCard {step} />
      {/each}
    </div>
  {/if}

  <Report {report} />

  <footer class="app-footer">
    Made with <span class="heart">&#9829;</span> by
    <a href="https://github.com/Wasabules/McNetworkChecker" target="_blank" rel="noopener">Wasabules</a>
  </footer>
</main>

<style>
  main { max-width: 860px; margin: 0 auto; padding: 20px 16px 10px; position: relative; }
  header { text-align: center; margin-bottom: 24px; }
  .steps { display: flex; flex-direction: column; gap: 6px; margin-bottom: 20px; }
  .app-footer {
    text-align: center; padding: 18px 0 4px; font-size: 11px; color: #555;
    font-family: 'Press Start 2P', monospace;
  }
  .app-footer .heart { color: #FF5555; }
  .app-footer a { color: #FFAA00; text-decoration: none; }
  .app-footer a:hover { text-decoration: underline; }
</style>
