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
  import { CARRIER_LABELS, formatRelativeTime, formatDateTime } from "../lib/utils";
  import { t, getStatusLabel } from "../lib/i18n.svelte";
  import { addToast } from "../lib/toast.svelte";
  import StatusBadge from "../components/StatusBadge.svelte";
  import ParcelTimeline from "../components/ParcelTimeline.svelte";
  import ConfirmDialog from "../components/ConfirmDialog.svelte";
  import Navbar from "../components/Navbar.svelte";

  let { params }: { params: { id: string } } = $props();

  let parcel = $state<Parcel | null>(null);
  let events = $state<TrackingEvent[]>([]);
  let loading = $state(true);
  let error = $state("");
  let refreshing = $state(false);
  let showDelete = $state(false);
  let showArchiveConfirm = $state(false);
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

  const progressSteps = ["info_received", "in_transit", "out_for_delivery", "delivered"];

  const currentStepIndex = $derived.by(() => {
    if (!parcel) return -1;
    if (parcel.status === "failed" || parcel.status === "expired") return -1;
    return progressSteps.indexOf(parcel.status);
  });

  const accentVar = $derived(
    parcel
      ? `var(--color-status-${parcel.status.replace(/_/g, "-")}, var(--color-status-unknown))`
      : "var(--color-status-unknown)"
  );

  async function load() {
    loading = true;
    error = "";
    try {
      const [p, e] = await Promise.all([
        getParcel(params.id),
        listEvents(params.id),
      ]);
      parcel = p;
      events = e;
    } catch (err: unknown) {
      if (err instanceof Error && "status" in err && (err as any).status === 404) {
        push("/");
        return;
      }
      const msg = err instanceof Error ? err.message : t("detail.loadFailed");
      error = msg;
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
      addToast(t("toast.trackingRefreshed"));
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : t("toast.error");
      addToast(t("detail.refreshFailed") + ": " + msg, "error");
    } finally {
      refreshing = false;
    }
  }

  async function handleDelete() {
    try {
      await deleteParcel(params.id);
      addToast(t("toast.parcelDeleted"));
      push("/");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : t("toast.error");
      addToast(msg, "error");
      showDelete = false;
    }
  }

  async function handleArchive() {
    if (!parcel) return;
    const wasArchived = parcel.archived;
    try {
      parcel = await updateParcel(params.id, { ...parcel, archived: !parcel.archived });
      addToast(t(wasArchived ? "toast.parcelUnarchived" : "toast.parcelArchived"));
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : t("toast.error");
      addToast(msg, "error");
    } finally {
      showArchiveConfirm = false;
    }
  }

  async function handleAddEvent(e: Event) {
    e.preventDefault();
    try {
      await createEvent(params.id, {
        status: eventStatus,
        message: eventMessage,
        location: eventLocation,
      });
      addToast(t("toast.eventAdded"));
      eventMessage = "";
      eventLocation = "";
      showAddEvent = false;
      await load();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : t("toast.error");
      addToast(msg, "error");
    }
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
    try {
      parcel = await updateParcel(params.id, {
        ...parcel,
        name: editName,
        notes: editNotes,
        tracking_number: editTrackingNumber,
        carrier: editCarrier,
      });
      addToast(t("toast.parcelUpdated"));
      editing = false;
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : t("toast.error");
      addToast(msg, "error");
    }
  }
</script>

<Navbar />

<main class="max-w-2xl mx-auto px-4 py-6">
  {#if loading}
    <div class="space-y-4">
      <div class="skeleton h-8 rounded-lg w-1/3"></div>
      <div class="skeleton h-4 rounded-lg w-1/2"></div>
      <div class="skeleton h-48 rounded-2xl"></div>
      <div class="skeleton h-64 rounded-2xl"></div>
    </div>
  {:else if error}
    <div class="text-center py-16">
      <div class="w-16 h-16 mx-auto mb-4 rounded-2xl bg-[var(--color-danger-light)] flex items-center justify-center">
        <svg class="w-8 h-8 text-[var(--color-danger)]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <p class="text-[var(--color-danger)] text-lg mb-4 font-medium">{error}</p>
      <div class="flex gap-3 justify-center">
        <button
          onclick={() => load()}
          class="px-4 py-2 bg-[var(--color-accent)] text-white rounded-lg hover:bg-[var(--color-accent-hover)] transition-colors cursor-pointer"
        >
          {t("common.retry")}
        </button>
        <button
          onclick={() => push("/")}
          class="px-4 py-2 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer"
        >
          {t("common.back")}
        </button>
      </div>
    </div>
  {:else if parcel}
    <!-- Back + Actions -->
    <div class="flex items-center justify-between mb-5 animate-fade-in">
      <button
        onclick={() => push("/")}
        class="inline-flex items-center gap-1 text-sm font-medium text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)] transition-colors cursor-pointer -ml-1 px-2 py-1 rounded-md hover:bg-[var(--color-surface-hover)]"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        {t("common.back")}
      </button>

      <div class="flex gap-1.5">
        <button
          onclick={handleRefresh}
          disabled={refreshing}
          class="w-9 h-9 flex items-center justify-center rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] hover:text-[var(--color-accent)] cursor-pointer disabled:opacity-50"
          title={t("detail.refreshTracking")}
          aria-label={t("detail.refreshTracking")}
        >
          <svg class="w-4 h-4 {refreshing ? 'animate-spin' : ''}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
        <button
          onclick={startEdit}
          class="w-9 h-9 flex items-center justify-center rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] hover:text-[var(--color-accent)] cursor-pointer"
          title={t("common.edit")}
          aria-label={t("common.edit")}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button
          onclick={() => { showArchiveConfirm = true; }}
          class="w-9 h-9 flex items-center justify-center rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] hover:text-[var(--color-accent)] cursor-pointer"
          title={parcel.archived ? t("detail.unarchive") : t("detail.archive")}
          aria-label={parcel.archived ? t("detail.unarchive") : t("detail.archive")}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
          </svg>
        </button>
        <button
          onclick={() => { showDelete = true; }}
          class="w-9 h-9 flex items-center justify-center rounded-lg border border-[var(--color-danger)]/30 bg-[var(--color-surface)] hover:bg-[var(--color-danger-light)] transition-colors text-[var(--color-danger)] cursor-pointer"
          title={t("common.delete")}
          aria-label={t("common.delete")}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Parcel Info -->
    {#if editing}
      <form onsubmit={handleEdit} class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-5 mb-6 space-y-4 shadow-[var(--shadow-sm)] animate-scale-in">
        <div>
          <label for="edit-name" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">{t("detail.name")}</label>
          <input id="edit-name" bind:value={editName} class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all" />
        </div>
        <div>
          <label for="edit-tracking" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">{t("detail.trackingNumber")}</label>
          <input id="edit-tracking" bind:value={editTrackingNumber} class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] font-mono focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all" />
        </div>
        <div>
          <label for="edit-notes" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">{t("detail.notes")}</label>
          <textarea id="edit-notes" bind:value={editNotes} rows="2" class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all resize-none"></textarea>
        </div>
        <div class="flex gap-2 pt-1">
          <button type="button" onclick={() => { editing = false; }} class="flex-1 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg font-medium hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer">{t("common.cancel")}</button>
          <button type="submit" class="flex-1 py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-medium rounded-lg shadow-md shadow-indigo-500/20 transition-all cursor-pointer">{t("common.save")}</button>
        </div>
      </form>
    {:else}
      <!-- Hero card -->
      <div
        class="relative bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-6 mb-6 overflow-hidden shadow-[var(--shadow-sm)] animate-slide-up"
      >
        <!-- Status-tinted background accent -->
        <div
          aria-hidden="true"
          class="absolute inset-x-0 top-0 h-24 opacity-10 pointer-events-none"
          style="background: linear-gradient(180deg, {accentVar}, transparent);"
        ></div>

        <div class="relative">
          <div class="flex items-start justify-between gap-3 mb-1">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2 text-xs font-semibold text-[var(--color-text-muted)] uppercase tracking-wide mb-1">
                <span>{CARRIER_LABELS[parcel.carrier] || parcel.carrier}</span>
                {#if parcel.archived}
                  <span class="px-1.5 py-0.5 rounded bg-[var(--color-surface-hover)] text-[var(--color-text-secondary)] text-[10px]">{t("detail.archived")}</span>
                {/if}
              </div>
              <h1 class="text-2xl font-bold text-[var(--color-text-primary)] leading-tight">
                {parcel.name || parcel.tracking_number}
              </h1>
              {#if parcel.name}
                <p class="text-sm font-mono text-[var(--color-text-secondary)] mt-1 select-all">
                  {parcel.tracking_number}
                </p>
              {/if}
            </div>
            <StatusBadge status={parcel.status} size="md" />
          </div>

          {#if parcel.estimated_delivery && parcel.status !== "delivered"}
            <div class="mt-4 inline-flex items-center gap-2 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg px-3 py-2 text-sm">
              <svg class="w-4 h-4 text-[var(--color-accent)]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
              <span class="text-[var(--color-text-secondary)]">{t("detail.estimatedDelivery")}</span>
              <span class="font-semibold text-[var(--color-text-primary)]">{formatDateTime(parcel.estimated_delivery)}</span>
            </div>
          {/if}

          <!-- Progress steps -->
          {#if currentStepIndex >= 0}
            <div class="mt-6">
              <div class="flex items-center justify-between gap-1">
                {#each progressSteps as step, i}
                  <div class="flex-1 flex flex-col items-center gap-1.5 relative">
                    {#if i > 0}
                      <div class="absolute right-1/2 top-[11px] h-0.5 w-full {i <= currentStepIndex ? '' : 'opacity-100'}"
                        style="background: {i <= currentStepIndex ? accentVar : 'var(--color-border)'};"></div>
                    {/if}
                    <div
                      class="relative z-10 w-6 h-6 rounded-full flex items-center justify-center transition-all {i <= currentStepIndex ? 'scale-100' : 'scale-90'}"
                      style="background: {i <= currentStepIndex ? accentVar : 'var(--color-surface-hover)'}; box-shadow: {i === currentStepIndex ? `0 0 0 4px color-mix(in srgb, ${accentVar} 20%, transparent)` : 'none'};"
                    >
                      {#if i < currentStepIndex}
                        <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                        </svg>
                      {:else if i === currentStepIndex}
                        <span class="w-2 h-2 rounded-full bg-white"></span>
                      {:else}
                        <span class="w-1.5 h-1.5 rounded-full bg-[var(--color-text-muted)]/40"></span>
                      {/if}
                    </div>
                    <span class="text-[10px] sm:text-xs text-center font-medium {i <= currentStepIndex ? 'text-[var(--color-text-primary)]' : 'text-[var(--color-text-muted)]'}">
                      {getStatusLabel(step)}
                    </span>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Meta row -->
          <div class="mt-5 pt-5 border-t border-[var(--color-border)] grid grid-cols-2 gap-4 text-sm">
            <div>
              <div class="text-xs text-[var(--color-text-muted)] uppercase tracking-wide font-medium">{t("detail.added")}</div>
              <div class="text-[var(--color-text-primary)] mt-0.5">{formatRelativeTime(parcel.created_at)}</div>
            </div>
            <div>
              <div class="text-xs text-[var(--color-text-muted)] uppercase tracking-wide font-medium">{t("detail.status")}</div>
              <div class="text-[var(--color-text-primary)] mt-0.5">{getStatusLabel(parcel.status)}</div>
            </div>
          </div>

          {#if parcel.notes}
            <div class="mt-4 bg-[var(--color-surface-alt)] rounded-lg p-3 text-sm">
              <div class="text-xs text-[var(--color-text-muted)] uppercase tracking-wide font-medium mb-1">{t("detail.notes")}</div>
              <p class="text-[var(--color-text-primary)] whitespace-pre-wrap">{parcel.notes}</p>
            </div>
          {/if}
        </div>
      </div>
    {/if}

    <!-- Tracking Events -->
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-[var(--color-text-primary)]">{t("detail.trackingHistory")}</h2>
      <button
        onclick={() => { showAddEvent = !showAddEvent; }}
        class="text-sm font-medium text-[var(--color-accent)] hover:text-[var(--color-accent-hover)] transition-colors cursor-pointer px-2 py-1 rounded-md hover:bg-[var(--color-accent-light)]"
      >
        {showAddEvent ? t("common.cancel") : t("detail.addEvent")}
      </button>
    </div>

    {#if showAddEvent}
      <form onsubmit={handleAddEvent} class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-4 mb-6 space-y-3 shadow-[var(--shadow-sm)] animate-scale-in">
        <div>
          <label for="event-message" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">{t("detail.message")}</label>
          <input
            id="event-message"
            bind:value={eventMessage}
            required
            class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
            placeholder={t("detail.messagePlaceholder")}
          />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="event-status" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">{t("detail.eventStatus")}</label>
            <select id="event-status" bind:value={eventStatus} class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all cursor-pointer">
              <option value="unknown">{getStatusLabel("unknown")}</option>
              <option value="info_received">{getStatusLabel("info_received")}</option>
              <option value="in_transit">{getStatusLabel("in_transit")}</option>
              <option value="out_for_delivery">{getStatusLabel("out_for_delivery")}</option>
              <option value="delivered">{getStatusLabel("delivered")}</option>
              <option value="failed">{getStatusLabel("failed")}</option>
            </select>
          </div>
          <div>
            <label for="event-location" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">{t("detail.location")}</label>
            <input
              id="event-location"
              bind:value={eventLocation}
              class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
              placeholder={t("detail.locationPlaceholder")}
            />
          </div>
        </div>
        <button type="submit" class="w-full py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-medium rounded-lg shadow-md shadow-indigo-500/20 transition-all cursor-pointer">
          {t("detail.submitEvent")}
        </button>
      </form>
    {/if}

    <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-5 shadow-[var(--shadow-sm)]">
      <ParcelTimeline {events} />
    </div>

    <ConfirmDialog
      show={showDelete}
      title={t("detail.deleteTitle")}
      message={t("detail.deleteMessage")}
      confirmLabel={t("common.delete")}
      confirmClass="bg-[var(--color-danger)] hover:bg-[var(--color-danger-hover)] text-white"
      onConfirm={handleDelete}
      onCancel={() => { showDelete = false; }}
    />

    <ConfirmDialog
      show={showArchiveConfirm}
      title={parcel.archived ? t("detail.unarchiveTitle") : t("detail.archiveTitle")}
      message={parcel.archived ? t("detail.unarchiveMessage") : t("detail.archiveMessage")}
      confirmLabel={parcel.archived ? t("detail.unarchive") : t("detail.archive")}
      onConfirm={handleArchive}
      onCancel={() => { showArchiveConfirm = false; }}
    />
  {/if}
</main>
