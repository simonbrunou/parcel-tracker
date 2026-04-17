<script lang="ts">
  import { getToasts, dismissToast } from "../lib/toast.svelte";

  const toasts = $derived(getToasts());

  function typeClasses(type: string): string {
    switch (type) {
      case "success":
        return "bg-emerald-500 text-white border-emerald-400/50 shadow-emerald-500/20";
      case "error":
        return "bg-[var(--color-danger)] text-white border-red-400/50 shadow-red-500/20";
      case "info":
        return "bg-[var(--color-accent)] text-white border-indigo-400/50 shadow-indigo-500/20";
      default:
        return "bg-[var(--color-surface)] text-[var(--color-text-primary)] border-[var(--color-border)]";
    }
  }
</script>

{#if toasts.length > 0}
  <div class="fixed bottom-4 right-4 left-4 sm:left-auto z-50 flex flex-col gap-2 sm:max-w-sm pointer-events-none" aria-live="polite" role="status">
    {#each toasts as toast (toast.id)}
      <div
        class="flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg border backdrop-blur-xl text-sm font-medium animate-slide-in-right pointer-events-auto {typeClasses(toast.type)}"
        role="alert"
      >
        {#if toast.type === "success"}
          <div class="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        {:else if toast.type === "error"}
          <div class="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
        {:else}
          <div class="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        {/if}
        <span class="flex-1">{toast.message}</span>
        <button
          onclick={() => dismissToast(toast.id)}
          class="shrink-0 w-6 h-6 flex items-center justify-center rounded-md opacity-70 hover:opacity-100 hover:bg-white/10 transition-all cursor-pointer"
          aria-label="Dismiss"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    {/each}
  </div>
{/if}
