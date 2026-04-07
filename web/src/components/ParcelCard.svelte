<script lang="ts">
  import type { Parcel } from "../lib/api";
  import { CARRIER_LABELS, formatRelativeTime } from "../lib/utils";
  import StatusBadge from "./StatusBadge.svelte";
  import { push } from "svelte-spa-router";

  let { parcel }: { parcel: Parcel } = $props();

  function navigate() {
    push(`/parcels/${parcel.id}`);
  }
</script>

<button
  onclick={navigate}
  class="w-full text-left bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-4 hover:bg-[var(--color-surface-hover)] transition-all duration-200 cursor-pointer hover:shadow-md group"
>
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
      <p class="text-xs text-[var(--color-text-muted)] mt-1.5">
        {CARRIER_LABELS[parcel.carrier] || parcel.carrier} &middot; {formatRelativeTime(parcel.updated_at)}
      </p>
    </div>
    <StatusBadge status={parcel.status} />
  </div>
</button>
