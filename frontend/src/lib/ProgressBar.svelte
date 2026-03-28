<script>
  export let completed = 0;
  export let total = 1;
  $: pct = (completed / total * 100);
</script>

<div class="mc-progress-wrap">
  <div class="mc-progress-bar">
    <div class="mc-progress-fill" style="width: {pct}%"></div>
  </div>
  <span class="mc-progress-label">{completed}/{total}</span>
</div>
