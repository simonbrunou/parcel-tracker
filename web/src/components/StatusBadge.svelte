<script lang="ts">
  import { STATUS_LABELS, STATUS_COLORS } from "../lib/utils";

  let { status }: { status: string } = $props();

  const label = $derived(STATUS_LABELS[status] || status);
  const colorClass = $derived(STATUS_COLORS[status] || STATUS_COLORS["unknown"]);
</script>

<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold text-white {colorClass}">
  {label}
</span>
