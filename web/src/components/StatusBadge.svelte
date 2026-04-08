<script lang="ts">
  import { STATUS_COLORS } from "../lib/utils";
  import { getStatusLabel } from "../lib/i18n.svelte";

  let { status }: { status: string } = $props();

  const label = $derived(getStatusLabel(status));
  const colorClass = $derived(STATUS_COLORS[status] || STATUS_COLORS["unknown"]);
</script>

<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold text-white {colorClass}">
  {label}
</span>
