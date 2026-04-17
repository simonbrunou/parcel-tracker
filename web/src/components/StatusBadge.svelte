<script lang="ts">
  import { getStatusLabel } from "../lib/i18n.svelte";

  let { status, size = "sm" }: { status: string; size?: "sm" | "md" } = $props();

  const label = $derived(getStatusLabel(status));
  const statusVar = $derived(`var(--color-status-${status.replace(/_/g, "-")}, var(--color-status-unknown))`);

  const sizeClasses = $derived(
    size === "md" ? "px-3 py-1.5 text-sm gap-2" : "px-2.5 py-1 text-xs gap-1.5"
  );
  const dotSize = $derived(size === "md" ? "w-2 h-2" : "w-1.5 h-1.5");
  const iconSize = $derived(size === "md" ? "w-3.5 h-3.5" : "w-3 h-3");
</script>

<span
  class="inline-flex items-center rounded-full font-semibold border {sizeClasses}"
  style="color: {statusVar}; background-color: color-mix(in srgb, {statusVar} 12%, transparent); border-color: color-mix(in srgb, {statusVar} 30%, transparent);"
>
  {#if status === "delivered"}
    <svg class={iconSize} fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
  {:else if status === "in_transit"}
    <svg class={iconSize} fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10m10 0H3m10 0a2 2 0 104 0m-4 0a2 2 0 114 0m6-6v6a1 1 0 01-1 1h-1m-6-1a2 2 0 104 0" /></svg>
  {:else if status === "out_for_delivery"}
    <svg class={iconSize} fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
  {:else if status === "failed"}
    <svg class={iconSize} fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
  {:else if status === "expired"}
    <svg class={iconSize} fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
  {:else if status === "info_received"}
    <svg class={iconSize} fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
  {:else}
    <span class="relative inline-flex {dotSize}">
      <span class="absolute inset-0 rounded-full animate-ping opacity-60" style="background-color: {statusVar};"></span>
      <span class="relative rounded-full {dotSize}" style="background-color: {statusVar};"></span>
    </span>
  {/if}
  {label}
</span>
