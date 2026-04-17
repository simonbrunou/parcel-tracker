<script lang="ts">
  import { t } from "../lib/i18n.svelte";

  let {
    show,
    title,
    message,
    confirmLabel,
    confirmClass = "bg-[var(--color-accent)] hover:bg-[var(--color-accent-hover)] text-white",
    onConfirm,
    onCancel,
  }: {
    show: boolean;
    title: string;
    message: string;
    confirmLabel: string;
    confirmClass?: string;
    onConfirm: () => void;
    onCancel: () => void;
  } = $props();

  const isDanger = $derived(confirmClass.includes("danger"));

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      onCancel();
    }
  }

  function handleBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) onCancel();
  }
</script>

{#if show}
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-fade-in"
    role="dialog"
    aria-modal="true"
    aria-label={title}
    tabindex="-1"
    onkeydown={handleKeydown}
    onclick={handleBackdrop}
  >
    <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-6 max-w-sm w-full shadow-[var(--shadow-lg)] animate-scale-in">
      <div class="flex items-start gap-3 mb-2">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 {isDanger ? 'bg-[var(--color-danger-light)] text-[var(--color-danger)]' : 'bg-[var(--color-accent-light)] text-[var(--color-accent)]'}">
          {#if isDanger}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          {:else}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          {/if}
        </div>
        <div class="flex-1 pt-1">
          <h3 class="text-base font-semibold text-[var(--color-text-primary)]">{title}</h3>
          <p class="text-[var(--color-text-secondary)] text-sm mt-1">{message}</p>
        </div>
      </div>
      <div class="flex gap-2 mt-5">
        <button
          onclick={onCancel}
          class="flex-1 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg font-medium hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer"
        >
          {t("common.cancel")}
        </button>
        <button
          onclick={onConfirm}
          class="flex-1 py-2.5 rounded-lg font-medium transition-all cursor-pointer shadow-md {confirmClass}"
        >
          {confirmLabel}
        </button>
      </div>
    </div>
  </div>
{/if}
