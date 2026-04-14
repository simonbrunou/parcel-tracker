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

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      onCancel();
    }
  }
</script>

{#if show}
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    role="dialog"
    aria-modal="true"
    aria-label={title}
    onkeydown={handleKeydown}
  >
    <div class="bg-[var(--color-surface)] rounded-xl p-6 max-w-sm w-full shadow-xl">
      <h3 class="text-lg font-semibold text-[var(--color-text-primary)] mb-2">{title}</h3>
      <p class="text-[var(--color-text-secondary)] text-sm mb-5">{message}</p>
      <div class="flex gap-3">
        <button
          onclick={onCancel}
          class="flex-1 py-2 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer"
        >
          {t("common.cancel")}
        </button>
        <button
          onclick={onConfirm}
          class="flex-1 py-2 rounded-lg transition-colors cursor-pointer {confirmClass}"
        >
          {confirmLabel}
        </button>
      </div>
    </div>
  </div>
{/if}
