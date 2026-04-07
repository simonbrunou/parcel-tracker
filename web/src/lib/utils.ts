export function formatRelativeTime(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffSec < 60) return "just now";
  if (diffMin < 60) return `${diffMin}m ago`;
  if (diffHour < 24) return `${diffHour}h ago`;
  if (diffDay < 7) return `${diffDay}d ago`;

  return date.toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: date.getFullYear() !== now.getFullYear() ? "numeric" : undefined,
  });
}

export function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export const STATUS_LABELS: Record<string, string> = {
  unknown: "Unknown",
  info_received: "Info Received",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  failed: "Failed",
  expired: "Expired",
};

export const STATUS_COLORS: Record<string, string> = {
  unknown: "bg-[var(--color-status-unknown)]",
  info_received: "bg-[var(--color-status-info-received)]",
  in_transit: "bg-[var(--color-status-in-transit)]",
  out_for_delivery: "bg-[var(--color-status-out-for-delivery)]",
  delivered: "bg-[var(--color-status-delivered)]",
  failed: "bg-[var(--color-status-failed)]",
  expired: "bg-[var(--color-status-expired)]",
};

export const CARRIER_LABELS: Record<string, string> = {
  manual: "Manual",
  mock: "Mock (Demo)",
  usps: "USPS",
  fedex: "FedEx",
  ups: "UPS",
  dhl: "DHL",
  postnl: "PostNL",
  colissimo: "Colissimo",
  laposte: "La Poste",
};

export function toggleTheme(): void {
  const current = document.documentElement.getAttribute("data-theme");
  const next = current === "dark" ? "light" : "dark";
  if (next === "dark") {
    document.documentElement.setAttribute("data-theme", "dark");
  } else {
    document.documentElement.removeAttribute("data-theme");
  }
  localStorage.setItem("theme", next);
}

export function isDarkTheme(): boolean {
  return document.documentElement.getAttribute("data-theme") === "dark";
}
