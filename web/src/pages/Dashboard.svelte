<script lang="ts">
  import { push } from "svelte-spa-router";
  import { listParcels, type Parcel } from "../lib/api";
  import ParcelCard from "../components/ParcelCard.svelte";
  import Navbar from "../components/Navbar.svelte";
  import { t, getStatusLabel } from "../lib/i18n.svelte";

  let parcels = $state<Parcel[]>([]);
  let loading = $state(true);
  let error = $state("");
  let search = $state("");
  let debouncedSearch = $state("");
  let statusFilter = $state("");
  let showArchived = $state(false);

  // Debounce the search input to avoid filtering on every keystroke.
  let debounceTimer: ReturnType<typeof setTimeout>;
  $effect(() => {
    const val = search;
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      debouncedSearch = val;
    }, 300);
    return () => clearTimeout(debounceTimer);
  });

  const activeParcels = $derived(parcels.filter((p) => !p.archived));

  const stats = $derived({
    total: activeParcels.length,
    inTransit: activeParcels.filter((p) => p.status === "in_transit" || p.status === "out_for_delivery").length,
    delivered: activeParcels.filter((p) => p.status === "delivered").length,
    attention: activeParcels.filter((p) => p.status === "failed" || p.status === "expired").length,
  });

  const filtered = $derived(
    parcels.filter((p) => {
      if (!showArchived && p.archived) return false;
      if (showArchived && !p.archived) return false;
      if (statusFilter && p.status !== statusFilter) return false;
      if (debouncedSearch) {
        const q = debouncedSearch.toLowerCase();
        return (
          p.name.toLowerCase().includes(q) ||
          p.tracking_number.toLowerCase().includes(q)
        );
      }
      return true;
    })
  );

  async function loadParcels() {
    loading = true;
    error = "";
    try {
      parcels = (await listParcels()).data;
    } catch (err: unknown) {
      if (err instanceof Error && "status" in err && (err as any).status === 401) return;
      const msg = err instanceof Error ? err.message : t("dashboard.loadFailed");
      error = t("dashboard.loadFailed") + ": " + msg;
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadParcels();
  });
</script>

<Navbar />

<main class="max-w-4xl mx-auto px-4 py-6 sm:py-8">
  <!-- Header -->
  <div class="flex items-center justify-between mb-6 animate-slide-up">
    <div>
      <h1 class="text-2xl sm:text-3xl font-bold tracking-tight text-[var(--color-text-primary)]">
        {showArchived ? t("dashboard.archive") : t("dashboard.title")}
      </h1>
      <p class="text-sm text-[var(--color-text-secondary)] mt-1">
        {filtered.length !== 1 ? t("dashboard.parcelCountPlural", { count: filtered.length }) : t("dashboard.parcelCount", { count: filtered.length })}
      </p>
    </div>
    <button
      onclick={() => push("/parcels/new")}
      class="inline-flex items-center gap-2 px-4 py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-medium rounded-xl shadow-md shadow-indigo-500/20 hover:shadow-lg hover:shadow-indigo-500/30 hover:-translate-y-0.5 active:translate-y-0 transition-all cursor-pointer"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
      </svg>
      <span class="hidden sm:inline">{t("dashboard.addParcel")}</span>
    </button>
  </div>

  <!-- Stats -->
  {#if !showArchived && !loading && !error && parcels.length > 0}
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5 mb-6 animate-slide-up">
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-3">
        <div class="text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wide">{t("dashboard.statTotal")}</div>
        <div class="text-2xl font-bold text-[var(--color-text-primary)] mt-0.5">{stats.total}</div>
      </div>
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-3">
        <div class="text-xs font-medium uppercase tracking-wide" style="color: var(--color-status-in-transit);">{t("dashboard.statInTransit")}</div>
        <div class="text-2xl font-bold text-[var(--color-text-primary)] mt-0.5">{stats.inTransit}</div>
      </div>
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-3">
        <div class="text-xs font-medium uppercase tracking-wide" style="color: var(--color-status-delivered);">{t("dashboard.statDelivered")}</div>
        <div class="text-2xl font-bold text-[var(--color-text-primary)] mt-0.5">{stats.delivered}</div>
      </div>
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-3">
        <div class="text-xs font-medium uppercase tracking-wide" style="color: var(--color-status-failed);">{t("dashboard.statAttention")}</div>
        <div class="text-2xl font-bold text-[var(--color-text-primary)] mt-0.5">{stats.attention}</div>
      </div>
    </div>
  {/if}

  <!-- Filters -->
  <div class="flex flex-col sm:flex-row gap-2.5 mb-5">
    <div class="flex-1 relative">
      <svg class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-muted)] pointer-events-none" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        type="text"
        bind:value={search}
        placeholder={t("dashboard.searchPlaceholder")}
        class="w-full pl-10 pr-10 py-2.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
      />
      {#if search}
        <button
          onclick={() => { search = ""; }}
          class="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded-md text-[var(--color-text-muted)] hover:text-[var(--color-text-primary)] hover:bg-[var(--color-surface-hover)] cursor-pointer"
          aria-label={t("dashboard.clearFilters")}
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      {/if}
    </div>

    <select
      bind:value={statusFilter}
      class="px-3.5 py-2.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all cursor-pointer"
    >
      <option value="">{t("dashboard.allStatuses")}</option>
      <option value="unknown">{getStatusLabel("unknown")}</option>
      <option value="info_received">{getStatusLabel("info_received")}</option>
      <option value="in_transit">{getStatusLabel("in_transit")}</option>
      <option value="out_for_delivery">{getStatusLabel("out_for_delivery")}</option>
      <option value="delivered">{getStatusLabel("delivered")}</option>
      <option value="failed">{getStatusLabel("failed")}</option>
    </select>

    <button
      onclick={() => { showArchived = !showArchived; statusFilter = ""; }}
      title={showArchived ? t("dashboard.title") : t("dashboard.archive")}
      aria-label={showArchived ? t("dashboard.title") : t("dashboard.archive")}
      class="px-3.5 py-2.5 rounded-xl border transition-all cursor-pointer inline-flex items-center justify-center gap-2 {showArchived
        ? 'bg-[var(--color-accent-light)] border-[var(--color-accent)] text-[var(--color-accent)]'
        : 'bg-[var(--color-surface)] border-[var(--color-border)] text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-hover)]'}"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
      </svg>
      <span class="hidden sm:inline text-sm font-medium">{showArchived ? t("dashboard.title") : t("dashboard.archive")}</span>
    </button>
  </div>

  <!-- Error State -->
  {#if error}
    <div class="text-center py-16">
      <div class="w-16 h-16 mx-auto mb-4 rounded-2xl bg-[var(--color-danger-light)] flex items-center justify-center">
        <svg class="w-8 h-8 text-[var(--color-danger)]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <p class="text-[var(--color-danger)] text-lg mb-4 font-medium">{error}</p>
      <button
        onclick={() => loadParcels()}
        class="px-4 py-2 bg-[var(--color-accent)] text-white rounded-lg hover:bg-[var(--color-accent-hover)] transition-colors cursor-pointer"
      >
        {t("common.retry")}
      </button>
    </div>
  <!-- Parcel List -->
  {:else if loading}
    <div class="space-y-3">
      {#each [1, 2, 3] as _}
        <div class="skeleton rounded-2xl p-4 border border-[var(--color-border)]">
          <div class="flex items-start justify-between">
            <div class="space-y-2 flex-1">
              <div class="h-5 bg-[var(--color-border)]/60 rounded w-1/3"></div>
              <div class="h-4 bg-[var(--color-border)]/60 rounded w-1/2"></div>
              <div class="h-3 bg-[var(--color-border)]/60 rounded w-1/4"></div>
            </div>
            <div class="h-6 w-20 bg-[var(--color-border)]/60 rounded-full"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if filtered.length === 0}
    <div class="text-center py-16 animate-fade-in">
      <div class="w-20 h-20 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-indigo-100 to-violet-100 dark:from-indigo-950/40 dark:to-violet-950/40 flex items-center justify-center">
        <svg class="w-10 h-10 text-[var(--color-accent)]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
      </div>
      {#if search || statusFilter}
        <p class="text-[var(--color-text-primary)] text-lg font-semibold">{t("dashboard.noMatch")}</p>
        <button
          onclick={() => { search = ""; statusFilter = ""; }}
          class="mt-3 text-[var(--color-accent)] hover:underline cursor-pointer font-medium"
        >
          {t("dashboard.clearFilters")}
        </button>
      {:else}
        <p class="text-[var(--color-text-primary)] text-lg font-semibold">{t("dashboard.noParcels")}</p>
        <p class="text-[var(--color-text-muted)] mt-1">{t("dashboard.addFirst")}</p>
        <button
          onclick={() => push("/parcels/new")}
          class="mt-5 inline-flex items-center gap-2 px-5 py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-medium rounded-xl shadow-md shadow-indigo-500/20 hover:shadow-lg hover:shadow-indigo-500/30 hover:-translate-y-0.5 transition-all cursor-pointer"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          {t("dashboard.addParcel")}
        </button>
      {/if}
    </div>
  {:else}
    <div class="space-y-3 stagger">
      {#each filtered as parcel (parcel.id)}
        <ParcelCard {parcel} />
      {/each}
    </div>
  {/if}
</main>
