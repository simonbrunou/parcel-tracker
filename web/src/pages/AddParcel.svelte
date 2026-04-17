<script lang="ts">
  import { push } from "svelte-spa-router";
  import { createParcel, getHealth, type CarrierInfo } from "../lib/api";
  import Navbar from "../components/Navbar.svelte";
  import { t } from "../lib/i18n.svelte";
  import { addToast } from "../lib/toast.svelte";

  let trackingNumber = $state("");
  let carrier = $state("manual");
  let name = $state("");
  let notes = $state("");
  let error = $state("");
  let loading = $state(false);
  let carriers = $state<CarrierInfo[]>([]);

  $effect(() => {
    getHealth().then((h) => {
      carriers = h.carriers || [];
    });
  });

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      const parcel = await createParcel({
        tracking_number: trackingNumber,
        carrier: carrier,
        name: name,
        notes: notes,
      });
      addToast(t("toast.parcelCreated"));
      push(`/parcels/${parcel.id}`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : t("addParcel.failed");
      error = msg;
    } finally {
      loading = false;
    }
  }
</script>

<Navbar />

<main class="max-w-xl mx-auto px-4 py-6">
  <div class="mb-6 animate-fade-in">
    <button
      onclick={() => push("/")}
      class="inline-flex items-center gap-1 text-sm font-medium text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)] transition-colors cursor-pointer -ml-1 px-2 py-1 rounded-md hover:bg-[var(--color-surface-hover)]"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      {t("common.back")}
    </button>
    <h1 class="text-2xl sm:text-3xl font-bold tracking-tight text-[var(--color-text-primary)] mt-2">{t("addParcel.title")}</h1>
    <p class="text-sm text-[var(--color-text-secondary)] mt-1">{t("addParcel.subtitle")}</p>
  </div>

  <form onsubmit={handleSubmit} class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-5 sm:p-6 shadow-[var(--shadow-sm)] space-y-5 animate-slide-up">
    {#if error}
      <div class="flex items-start gap-2 bg-[var(--color-danger-light)] text-[var(--color-danger)] text-sm p-3 rounded-lg border border-[var(--color-danger)]/20 animate-scale-in">
        <svg class="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{error}</span>
      </div>
    {/if}

    <div>
      <label for="tracking" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        {t("addParcel.trackingNumber")} <span class="text-[var(--color-danger)]">*</span>
      </label>
      <input
        id="tracking"
        type="text"
        bind:value={trackingNumber}
        required
        autocomplete="off"
        spellcheck="false"
        class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] font-mono placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
        placeholder={t("addParcel.trackingPlaceholder")}
      />
    </div>

    <div>
      <label for="carrier" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        {t("addParcel.carrier")}
      </label>
      <select
        id="carrier"
        bind:value={carrier}
        class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all cursor-pointer"
      >
        {#each carriers as c}
          <option value={c.code}>{c.name}</option>
        {/each}
        {#if carriers.length === 0}
          <option value="manual">{t("common.manual")}</option>
        {/if}
      </select>
      <p class="text-xs text-[var(--color-text-muted)] mt-1.5">{t("addParcel.carrierHint")}</p>
    </div>

    <div>
      <label for="name" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        {t("addParcel.customName")}
      </label>
      <input
        id="name"
        type="text"
        bind:value={name}
        class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
        placeholder={t("addParcel.namePlaceholder")}
      />
    </div>

    <div>
      <label for="notes" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        {t("addParcel.notes")}
      </label>
      <textarea
        id="notes"
        bind:value={notes}
        rows="3"
        class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all resize-none"
        placeholder={t("addParcel.notesPlaceholder")}
      ></textarea>
    </div>

    <div class="flex gap-3 pt-1">
      <button
        type="button"
        onclick={() => push("/")}
        class="flex-1 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] text-[var(--color-text-primary)] font-medium rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer"
      >
        {t("common.cancel")}
      </button>
      <button
        type="submit"
        disabled={loading}
        class="flex-1 py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-semibold rounded-lg shadow-md shadow-indigo-500/20 hover:shadow-lg hover:shadow-indigo-500/30 transition-all disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer flex items-center justify-center gap-2"
      >
        {#if loading}
          <span class="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin"></span>
          {t("addParcel.adding")}
        {:else}
          {t("addParcel.submit")}
        {/if}
      </button>
    </div>
  </form>
</main>
