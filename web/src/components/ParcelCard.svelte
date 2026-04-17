<script lang="ts">
  import type { Parcel } from "../lib/api";
  import { CARRIER_LABELS, formatRelativeTime, formatDateTime } from "../lib/utils";
  import { t } from "../lib/i18n.svelte";
  import StatusBadge from "./StatusBadge.svelte";
  import { push } from "svelte-spa-router";

  let { parcel }: { parcel: Parcel } = $props();

  const accentVar = $derived(`var(--color-status-${parcel.status.replace(/_/g, "-")}, var(--color-status-unknown))`);

  function navigate() {
    push(`/parcels/${parcel.id}`);
  }
</script>

<button
  onclick={navigate}
  style="--accent: {accentVar}"
  class="relative w-full text-left bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-4 pl-5 hover:bg-[var(--color-surface-hover)] hover:border-[var(--color-border-strong)] hover:-translate-y-0.5 hover:shadow-[var(--shadow-md)] transition-all duration-200 cursor-pointer group overflow-hidden"
>
  <!-- Status color accent bar -->
  <span
    class="absolute left-0 top-3 bottom-3 w-1 rounded-r-full"
    style="background-color: var(--accent);"
    aria-hidden="true"
  ></span>

  <div class="flex items-start justify-between gap-3">
    <div class="flex-1 min-w-0">
      <h3 class="font-semibold text-[var(--color-text-primary)] truncate group-hover:text-[var(--color-accent)] transition-colors">
        {parcel.name || parcel.tracking_number}
      </h3>
      {#if parcel.name}
        <p class="text-sm text-[var(--color-text-secondary)] truncate mt-0.5 font-mono">
          {parcel.tracking_number}
        </p>
      {/if}
      <div class="flex items-center flex-wrap gap-x-2 gap-y-1 text-xs text-[var(--color-text-muted)] mt-2">
        <span class="inline-flex items-center gap-1 font-medium">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 17a2 2 0 11-4 0 2 2 0 014 0zM19 17a2 2 0 11-4 0 2 2 0 014 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10h10zm0 0h6m-6-6h6l3 3v3h-3"/>
          </svg>
          {CARRIER_LABELS[parcel.carrier] || parcel.carrier}
        </span>
        <span class="text-[var(--color-border-strong)]">·</span>
        <span>{formatRelativeTime(parcel.updated_at)}</span>
      </div>
      {#if parcel.estimated_delivery && parcel.status !== "delivered"}
        <p class="inline-flex items-center gap-1.5 text-xs text-[var(--color-text-secondary)] mt-2 bg-[var(--color-surface-alt)] px-2 py-1 rounded-md">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
          {t("detail.estimatedDelivery")}: {formatDateTime(parcel.estimated_delivery)}
        </p>
      {/if}
    </div>
    <div class="flex flex-col items-end gap-2 shrink-0">
      <StatusBadge status={parcel.status} />
      <svg class="w-4 h-4 text-[var(--color-text-muted)] opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
      </svg>
    </div>
  </div>
</button>
