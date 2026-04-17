<script lang="ts">
  import type { TrackingEvent } from "../lib/api";
  import { formatDateTime } from "../lib/utils";
  import { getStatusLabel, t } from "../lib/i18n.svelte";

  let { events }: { events: TrackingEvent[] } = $props();

  function statusVar(status: string): string {
    return `var(--color-status-${status.replace(/_/g, "-")}, var(--color-status-unknown))`;
  }
</script>

{#if events.length === 0}
  <div class="text-center py-10 text-[var(--color-text-muted)]">
    <div class="w-14 h-14 mx-auto mb-3 rounded-2xl bg-[var(--color-surface-alt)] flex items-center justify-center">
      <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </div>
    <p class="text-sm">{t("timeline.empty")}</p>
  </div>
{:else}
  <ol class="relative">
    {#each events as event, i}
      <li class="relative pl-9 {i < events.length - 1 ? 'pb-5' : ''}">
        <!-- Connecting line -->
        {#if i < events.length - 1}
          <div class="absolute left-[13px] top-7 bottom-0 w-px bg-gradient-to-b from-[var(--color-border-strong)] to-[var(--color-border)]"></div>
        {/if}

        <!-- Dot -->
        <div
          class="absolute left-0 top-1 w-[26px] h-[26px] rounded-full flex items-center justify-center {i === 0 ? '' : ''}"
          style="background-color: {i === 0 ? statusVar(event.status) : 'var(--color-surface-alt)'}; border: 2px solid {i === 0 ? 'color-mix(in srgb, ' + statusVar(event.status) + ' 30%, transparent)' : 'var(--color-border)'};"
        >
          {#if i === 0}
            <span class="w-1.5 h-1.5 rounded-full bg-white"></span>
          {:else}
            <span class="w-1.5 h-1.5 rounded-full" style="background-color: {statusVar(event.status)};"></span>
          {/if}
        </div>

        <!-- Content -->
        <div class="{i > 0 ? 'opacity-75' : ''}">
          <p class="font-medium text-[var(--color-text-primary)] leading-snug">
            {event.message}
          </p>
          <div class="flex flex-wrap items-center gap-x-2 gap-y-1 mt-1 text-xs">
            <span
              class="inline-flex items-center gap-1 font-medium px-1.5 py-0.5 rounded"
              style="color: {statusVar(event.status)}; background-color: color-mix(in srgb, {statusVar(event.status)} 10%, transparent);"
            >
              {getStatusLabel(event.status)}
            </span>
            {#if event.location}
              <span class="text-[var(--color-text-muted)] inline-flex items-center gap-1">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
                {event.location}
              </span>
            {/if}
            <span class="text-[var(--color-text-muted)]">· {formatDateTime(event.timestamp)}</span>
          </div>
        </div>
      </li>
    {/each}
  </ol>
{/if}
