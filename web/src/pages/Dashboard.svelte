<script lang="ts">
  import { push } from "svelte-spa-router";
  import { listParcels, type Parcel } from "../lib/api";
  import ParcelCard from "../components/ParcelCard.svelte";
  import Navbar from "../components/Navbar.svelte";

  let parcels = $state<Parcel[]>([]);
  let loading = $state(true);
  let search = $state("");
  let statusFilter = $state("");
  let showArchived = $state(false);

  const filtered = $derived(
    parcels.filter((p) => {
      if (!showArchived && p.archived) return false;
      if (showArchived && !p.archived) return false;
      if (statusFilter && p.status !== statusFilter) return false;
      if (search) {
        const q = search.toLowerCase();
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
    try {
      parcels = await listParcels();
    } catch {
      // handled by api redirect
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadParcels();
  });
</script>

<Navbar />

<main class="max-w-4xl mx-auto px-4 py-6">
  <!-- Header -->
  <div class="flex items-center justify-between mb-6">
    <div>
      <h1 class="text-2xl font-bold text-[var(--color-text-primary)]">
        {showArchived ? "Archive" : "My Parcels"}
      </h1>
      <p class="text-sm text-[var(--color-text-secondary)] mt-0.5">
        {filtered.length} parcel{filtered.length !== 1 ? "s" : ""}
      </p>
    </div>
    <button
      onclick={() => push("/parcels/new")}
      class="inline-flex items-center gap-2 px-4 py-2.5 bg-[var(--color-accent)] hover:bg-[var(--color-accent-hover)] text-white font-medium rounded-xl transition-colors cursor-pointer"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      <span class="hidden sm:inline">Add Parcel</span>
    </button>
  </div>

  <!-- Filters -->
  <div class="flex flex-col sm:flex-row gap-3 mb-6">
    <div class="flex-1 relative">
      <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-muted)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        type="text"
        bind:value={search}
        placeholder="Search parcels..."
        class="w-full pl-10 pr-4 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-xl text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] transition-all"
      />
    </div>

    <select
      bind:value={statusFilter}
      class="px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-xl text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] transition-all cursor-pointer"
    >
      <option value="">All statuses</option>
      <option value="unknown">Unknown</option>
      <option value="info_received">Info Received</option>
      <option value="in_transit">In Transit</option>
      <option value="out_for_delivery">Out for Delivery</option>
      <option value="delivered">Delivered</option>
      <option value="failed">Failed</option>
    </select>

    <button
      onclick={() => { showArchived = !showArchived; }}
      class="px-3 py-2.5 rounded-xl border transition-colors cursor-pointer {showArchived
        ? 'bg-[var(--color-accent-light)] border-[var(--color-accent)] text-[var(--color-accent)]'
        : 'bg-[var(--color-surface-alt)] border-[var(--color-border)] text-[var(--color-text-secondary)]'}"
    >
      <svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
      </svg>
    </button>
  </div>

  <!-- Parcel List -->
  {#if loading}
    <div class="space-y-3">
      {#each [1, 2, 3] as _}
        <div class="bg-[var(--color-surface-alt)] rounded-xl p-4 animate-pulse">
          <div class="flex items-start justify-between">
            <div class="space-y-2 flex-1">
              <div class="h-5 bg-[var(--color-border)] rounded w-1/3"></div>
              <div class="h-4 bg-[var(--color-border)] rounded w-1/2"></div>
              <div class="h-3 bg-[var(--color-border)] rounded w-1/4"></div>
            </div>
            <div class="h-6 w-20 bg-[var(--color-border)] rounded-full"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if filtered.length === 0}
    <div class="text-center py-16">
      <svg class="w-16 h-16 mx-auto mb-4 text-[var(--color-text-muted)] opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
      </svg>
      {#if search || statusFilter}
        <p class="text-[var(--color-text-secondary)] text-lg">No parcels match your filters</p>
        <button
          onclick={() => { search = ""; statusFilter = ""; }}
          class="mt-3 text-[var(--color-accent)] hover:underline cursor-pointer"
        >
          Clear filters
        </button>
      {:else}
        <p class="text-[var(--color-text-secondary)] text-lg">No parcels yet</p>
        <p class="text-[var(--color-text-muted)] mt-1">Add your first parcel to start tracking</p>
        <button
          onclick={() => push("/parcels/new")}
          class="mt-4 inline-flex items-center gap-2 px-5 py-2.5 bg-[var(--color-accent)] hover:bg-[var(--color-accent-hover)] text-white font-medium rounded-xl transition-colors cursor-pointer"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Parcel
        </button>
      {/if}
    </div>
  {:else}
    <div class="space-y-3">
      {#each filtered as parcel (parcel.id)}
        <ParcelCard {parcel} />
      {/each}
    </div>
  {/if}
</main>
