<script lang="ts">
  import { getToasts, dismissToast } from "../lib/toast.svelte";

  const toasts = $derived(getToasts());

  function typeClasses(type: string): string {
    switch (type) {
      case "success":
        return "bg-green-600 text-white";
      case "error":
        return "bg-[var(--color-danger)] text-white";
      case "info":
        return "bg-[var(--color-accent)] text-white";
      default:
        return "bg-[var(--color-surface-alt)] text-[var(--color-text-primary)]";
    }
  }
</script>

{#if toasts.length > 0}
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm">
    {#each toasts as toast (toast.id)}
      <div
        class="flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg text-sm font-medium {typeClasses(toast.type)}"
        role="alert"
      >
        {#if toast.type === "success"}
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        {:else if toast.type === "error"}
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        {:else}
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        {/if}
        <span class="flex-1">{toast.message}</span>
        <button
          onclick={() => dismissToast(toast.id)}
          class="shrink-0 opacity-70 hover:opacity-100 transition-opacity cursor-pointer"
          aria-label="Dismiss"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    {/each}
  </div>
{/if}
