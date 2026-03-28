<script>
  export let lines = [];
  export let isRunning = false;

  function autoScroll(node, _dep) {
    function scroll() {
      requestAnimationFrame(() => { node.scrollTop = node.scrollHeight; });
    }
    scroll();
    return { update() { scroll(); } };
  }
</script>

{#if lines.length > 0}
  <div class="mc-console" use:autoScroll={lines}>
    {#each lines as line}
      <div class="log-line"><span class="log-gt">&gt;</span> {line}</div>
    {/each}
    {#if isRunning}
      <div class="log-cursor">_</div>
    {/if}
  </div>
{/if}
