<script lang="ts">
  import type { TrackingEvent } from "../lib/api";
  import { STATUS_COLORS, formatDateTime } from "../lib/utils";
  import { getStatusLabel, t } from "../lib/i18n.svelte";

  let { events }: { events: TrackingEvent[] } = $props();
</script>

{#if events.length === 0}
  <div class="text-center py-8 text-[var(--color-text-muted)]">
    <svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <p>{t("timeline.empty")}</p>
  </div>
{:else}
  <ol class="relative">
    {#each events as event, i}
      <li class="relative pl-8 {i < events.length - 1 ? 'pb-6' : ''}">
        <!-- Connecting line -->
        {#if i < events.length - 1}
          <div class="absolute left-[11px] top-6 bottom-0 w-0.5 bg-[var(--color-border)]"></div>
        {/if}

        <!-- Dot -->
        <div class="absolute left-0 top-1 w-6 h-6 rounded-full flex items-center justify-center {i === 0 ? STATUS_COLORS[event.status] || 'bg-[var(--color-status-unknown)]' : 'bg-[var(--color-border)]'}">
          {#if i === 0}
            <div class="w-2 h-2 rounded-full bg-white"></div>
          {:else}
            <div class="w-2 h-2 rounded-full bg-[var(--color-text-muted)]"></div>
          {/if}
        </div>

        <!-- Content -->
        <div>
          <p class="font-medium text-[var(--color-text-primary)] {i > 0 ? 'opacity-70' : ''}">
            {event.message}
          </p>
          <div class="flex flex-wrap gap-2 mt-1 text-xs text-[var(--color-text-muted)]">
            <span class="inline-flex items-center gap-1">
              <span class="w-1.5 h-1.5 rounded-full {STATUS_COLORS[event.status] || 'bg-[var(--color-status-unknown)]'}"></span>
              {getStatusLabel(event.status)}
            </span>
            {#if event.location}
              <span>&middot; {event.location}</span>
            {/if}
            <span>&middot; {formatDateTime(event.timestamp)}</span>
          </div>
        </div>
      </li>
    {/each}
  </ol>
{/if}
