<script lang="ts">
  import { push } from "svelte-spa-router";
  import {
    getParcel,
    listEvents,
    deleteParcel,
    refreshParcel,
    createEvent,
    updateParcel,
    type Parcel,
    type TrackingEvent,
  } from "../lib/api";
  import { CARRIER_LABELS, STATUS_LABELS, formatRelativeTime } from "../lib/utils";
  import StatusBadge from "../components/StatusBadge.svelte";
  import ParcelTimeline from "../components/ParcelTimeline.svelte";
  import Navbar from "../components/Navbar.svelte";

  let { params }: { params: { id: string } } = $props();

  let parcel = $state<Parcel | null>(null);
  let events = $state<TrackingEvent[]>([]);
  let loading = $state(true);
  let refreshing = $state(false);
  let showDelete = $state(false);
  let showAddEvent = $state(false);
  let editing = $state(false);

  // Add event form
  let eventMessage = $state("");
  let eventStatus = $state("in_transit");
  let eventLocation = $state("");

  // Edit form
  let editName = $state("");
  let editNotes = $state("");
  let editTrackingNumber = $state("");
  let editCarrier = $state("");

  async function load() {
    loading = true;
    try {
      const [p, e] = await Promise.all([
        getParcel(params.id),
        listEvents(params.id),
      ]);
      parcel = p;
      events = e;
    } catch {
      push("/");
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    load();
  });

  async function handleRefresh() {
    refreshing = true;
    try {
      parcel = await refreshParcel(params.id);
      events = await listEvents(params.id);
    } catch {}
    refreshing = false;
  }

  async function handleDelete() {
    await deleteParcel(params.id);
    push("/");
  }

  async function handleArchive() {
    if (!parcel) return;
    parcel = await updateParcel(params.id, { ...parcel, archived: !parcel.archived });
  }

  async function handleAddEvent(e: Event) {
    e.preventDefault();
    await createEvent(params.id, {
      status: eventStatus,
      message: eventMessage,
      location: eventLocation,
    });
    eventMessage = "";
    eventLocation = "";
    showAddEvent = false;
    await load();
  }

  function startEdit() {
    if (!parcel) return;
    editName = parcel.name;
    editNotes = parcel.notes;
    editTrackingNumber = parcel.tracking_number;
    editCarrier = parcel.carrier;
    editing = true;
  }

  async function handleEdit(e: Event) {
    e.preventDefault();
    if (!parcel) return;
    parcel = await updateParcel(params.id, {
      ...parcel,
      name: editName,
      notes: editNotes,
      tracking_number: editTrackingNumber,
      carrier: editCarrier,
    });
    editing = false;
  }
</script>

<Navbar />

<main class="max-w-2xl mx-auto px-4 py-6">
  {#if loading}
    <div class="animate-pulse space-y-4">
      <div class="h-8 bg-[var(--color-border)] rounded w-1/3"></div>
      <div class="h-4 bg-[var(--color-border)] rounded w-1/2"></div>
      <div class="h-48 bg-[var(--color-border)] rounded-xl"></div>
    </div>
  {:else if parcel}
    <!-- Back + Actions -->
    <div class="flex items-center justify-between mb-6">
      <button
        onclick={() => push("/")}
        class="inline-flex items-center gap-1 text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)] transition-colors cursor-pointer"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Back
      </button>

      <div class="flex gap-2">
        <button
          onclick={handleRefresh}
          disabled={refreshing}
          class="p-2 rounded-lg border border-[var(--color-border)] hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer disabled:opacity-50"
          title="Refresh tracking"
        >
          <svg class="w-4 h-4 {refreshing ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
        <button
          onclick={startEdit}
          class="p-2 rounded-lg border border-[var(--color-border)] hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer"
          title="Edit"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button
          onclick={handleArchive}
          class="p-2 rounded-lg border border-[var(--color-border)] hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer"
          title={parcel.archived ? "Unarchive" : "Archive"}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
          </svg>
        </button>
        <button
          onclick={() => { showDelete = true; }}
          class="p-2 rounded-lg border border-[var(--color-danger)]/30 hover:bg-[var(--color-danger)]/10 transition-colors text-[var(--color-danger)] cursor-pointer"
          title="Delete"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Parcel Info -->
    {#if editing}
      <form onsubmit={handleEdit} class="bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-xl p-5 mb-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1">Name</label>
          <input bind:value={editName} class="w-full px-3 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]" />
        </div>
        <div>
          <label class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1">Tracking Number</label>
          <input bind:value={editTrackingNumber} class="w-full px-3 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] font-mono focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]" />
        </div>
        <div>
          <label class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1">Notes</label>
          <textarea bind:value={editNotes} rows="2" class="w-full px-3 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] resize-none"></textarea>
        </div>
        <div class="flex gap-2">
          <button type="button" onclick={() => { editing = false; }} class="flex-1 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer">Cancel</button>
          <button type="submit" class="flex-1 py-2 bg-[var(--color-accent)] text-white rounded-lg hover:bg-[var(--color-accent-hover)] transition-colors cursor-pointer">Save</button>
        </div>
      </form>
    {:else}
      <div class="bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-xl p-5 mb-6">
        <div class="flex items-start justify-between gap-3 mb-3">
          <h1 class="text-xl font-bold text-[var(--color-text-primary)]">
            {parcel.name || parcel.tracking_number}
          </h1>
          <StatusBadge status={parcel.status} />
        </div>

        <div class="space-y-1.5 text-sm">
          {#if parcel.name}
            <div class="flex gap-2">
              <span class="text-[var(--color-text-muted)] w-24 shrink-0">Tracking</span>
              <span class="font-mono text-[var(--color-text-primary)]">{parcel.tracking_number}</span>
            </div>
          {/if}
          <div class="flex gap-2">
            <span class="text-[var(--color-text-muted)] w-24 shrink-0">Carrier</span>
            <span class="text-[var(--color-text-primary)]">{CARRIER_LABELS[parcel.carrier] || parcel.carrier}</span>
          </div>
          <div class="flex gap-2">
            <span class="text-[var(--color-text-muted)] w-24 shrink-0">Status</span>
            <span class="text-[var(--color-text-primary)]">{STATUS_LABELS[parcel.status] || parcel.status}</span>
          </div>
          <div class="flex gap-2">
            <span class="text-[var(--color-text-muted)] w-24 shrink-0">Added</span>
            <span class="text-[var(--color-text-primary)]">{formatRelativeTime(parcel.created_at)}</span>
          </div>
          {#if parcel.notes}
            <div class="flex gap-2">
              <span class="text-[var(--color-text-muted)] w-24 shrink-0">Notes</span>
              <span class="text-[var(--color-text-primary)]">{parcel.notes}</span>
            </div>
          {/if}
          {#if parcel.archived}
            <div class="mt-2 text-xs text-[var(--color-text-muted)] bg-[var(--color-surface-hover)] inline-block px-2 py-1 rounded">Archived</div>
          {/if}
        </div>
      </div>
    {/if}

    <!-- Tracking Events -->
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-[var(--color-text-primary)]">Tracking History</h2>
      <button
        onclick={() => { showAddEvent = !showAddEvent; }}
        class="text-sm text-[var(--color-accent)] hover:text-[var(--color-accent-hover)] transition-colors cursor-pointer"
      >
        {showAddEvent ? "Cancel" : "+ Add Event"}
      </button>
    </div>

    {#if showAddEvent}
      <form onsubmit={handleAddEvent} class="bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-xl p-4 mb-6 space-y-3">
        <div>
          <label class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1">Message *</label>
          <input
            bind:value={eventMessage}
            required
            class="w-full px-3 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]"
            placeholder="e.g. Package arrived at sorting facility"
          />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1">Status</label>
            <select bind:value={eventStatus} class="w-full px-3 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] cursor-pointer">
              <option value="unknown">Unknown</option>
              <option value="info_received">Info Received</option>
              <option value="in_transit">In Transit</option>
              <option value="out_for_delivery">Out for Delivery</option>
              <option value="delivered">Delivered</option>
              <option value="failed">Failed</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1">Location</label>
            <input
              bind:value={eventLocation}
              class="w-full px-3 py-2 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]"
              placeholder="Optional"
            />
          </div>
        </div>
        <button type="submit" class="w-full py-2 bg-[var(--color-accent)] text-white rounded-lg hover:bg-[var(--color-accent-hover)] transition-colors cursor-pointer">
          Add Event
        </button>
      </form>
    {/if}

    <ParcelTimeline {events} />

    <!-- Delete confirmation -->
    {#if showDelete}
      <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-[var(--color-surface)] rounded-xl p-6 max-w-sm w-full shadow-xl">
          <h3 class="text-lg font-semibold text-[var(--color-text-primary)] mb-2">Delete parcel?</h3>
          <p class="text-[var(--color-text-secondary)] text-sm mb-5">
            This will permanently delete this parcel and all its tracking events.
          </p>
          <div class="flex gap-3">
            <button
              onclick={() => { showDelete = false; }}
              class="flex-1 py-2 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              onclick={handleDelete}
              class="flex-1 py-2 bg-[var(--color-danger)] hover:bg-[var(--color-danger-hover)] text-white rounded-lg transition-colors cursor-pointer"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</main>
