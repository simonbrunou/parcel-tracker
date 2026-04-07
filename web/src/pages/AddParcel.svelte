<script lang="ts">
  import { push } from "svelte-spa-router";
  import { createParcel, getHealth, type CarrierInfo } from "../lib/api";
  import Navbar from "../components/Navbar.svelte";

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
      push(`/parcels/${parcel.id}`);
    } catch (err: any) {
      error = err.message || "Failed to create parcel";
    } finally {
      loading = false;
    }
  }
</script>

<Navbar />

<main class="max-w-xl mx-auto px-4 py-6">
  <div class="mb-6">
    <button
      onclick={() => push("/")}
      class="inline-flex items-center gap-1 text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)] transition-colors cursor-pointer"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Back
    </button>
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mt-2">Add Parcel</h1>
  </div>

  <form onsubmit={handleSubmit} class="space-y-5">
    {#if error}
      <div class="bg-red-50 dark:bg-red-900/20 text-[var(--color-danger)] text-sm p-3 rounded-lg border border-red-200 dark:border-red-800">
        {error}
      </div>
    {/if}

    <div>
      <label for="tracking" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        Tracking Number <span class="text-[var(--color-danger)]">*</span>
      </label>
      <input
        id="tracking"
        type="text"
        bind:value={trackingNumber}
        required
        class="w-full px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] font-mono placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] transition-all"
        placeholder="e.g. 1Z999AA10123456784"
      />
    </div>

    <div>
      <label for="carrier" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        Carrier
      </label>
      <select
        id="carrier"
        bind:value={carrier}
        class="w-full px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] transition-all cursor-pointer"
      >
        {#each carriers as c}
          <option value={c.code}>{c.name}</option>
        {/each}
        {#if carriers.length === 0}
          <option value="manual">Manual</option>
        {/if}
      </select>
    </div>

    <div>
      <label for="name" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        Custom Name
      </label>
      <input
        id="name"
        type="text"
        bind:value={name}
        class="w-full px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] transition-all"
        placeholder="e.g. New Headphones"
      />
    </div>

    <div>
      <label for="notes" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
        Notes
      </label>
      <textarea
        id="notes"
        bind:value={notes}
        rows="3"
        class="w-full px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] transition-all resize-none"
        placeholder="Optional notes..."
      ></textarea>
    </div>

    <div class="flex gap-3 pt-2">
      <button
        type="button"
        onclick={() => push("/")}
        class="flex-1 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] text-[var(--color-text-primary)] font-medium rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer"
      >
        Cancel
      </button>
      <button
        type="submit"
        disabled={loading}
        class="flex-1 py-2.5 bg-[var(--color-accent)] hover:bg-[var(--color-accent-hover)] text-white font-medium rounded-lg transition-colors disabled:opacity-50 cursor-pointer"
      >
        {loading ? "Adding..." : "Add Parcel"}
      </button>
    </div>
  </form>
</main>
