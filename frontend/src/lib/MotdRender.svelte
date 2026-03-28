<script>
  import { afterUpdate, onDestroy } from 'svelte';

  export let spans = [];
  export let fallback = '';

  let container;
  let obfIntervals = [];

  function spanStyle(span) {
    let s = `color:${span.color || '#fff'};`;
    if (span.bold) s += 'font-weight:bold;';
    if (span.italic) s += 'font-style:italic;';
    let deco = [];
    if (span.underlined) deco.push('underline');
    if (span.strikethrough) deco.push('line-through');
    if (deco.length) s += `text-decoration:${deco.join(' ')};`;
    return s;
  }

  function esc(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function buildHTML(spans) {
    if (!spans || !spans.length) {
      return `<span style="color:#fff">${esc(fallback)}</span>`;
    }
    let h = '';
    for (const span of spans) {
      const style = spanStyle(span);
      if (span.obfuscated) {
        h += `<span style="${style}" data-obf="${span.text.length}">${esc(span.text)}</span>`;
        continue;
      }
      // Handle newlines: split and insert <br> with no surrounding whitespace
      const parts = span.text.split('\n');
      for (let i = 0; i < parts.length; i++) {
        if (i > 0) h += '<br>';
        if (parts[i]) h += `<span style="${style}">${esc(parts[i])}</span>`;
      }
    }
    return h;
  }

  $: html = buildHTML(spans);

  // Set up obfuscated text effect after each render
  afterUpdate(() => {
    obfIntervals.forEach(clearInterval);
    obfIntervals = [];
    if (!container) return;
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%&';
    container.querySelectorAll('[data-obf]').forEach(el => {
      const len = parseInt(el.dataset.obf) || 1;
      const id = setInterval(() => {
        let r = '';
        for (let i = 0; i < len; i++) r += chars[Math.floor(Math.random() * chars.length)];
        el.textContent = r;
      }, 50);
      obfIntervals.push(id);
    });
  });

  onDestroy(() => obfIntervals.forEach(clearInterval));
</script>

<div class="mc-motd-render" bind:this={container}>{@html html}</div>
