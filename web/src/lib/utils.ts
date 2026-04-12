import { t, getLocaleCode } from "./i18n.svelte";

export function formatRelativeTime(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffSec < 60) return t("time.justNow");
  if (diffMin < 60) return t("time.minutesAgo", { n: diffMin });
  if (diffHour < 24) return t("time.hoursAgo", { n: diffHour });
  if (diffDay < 7) return t("time.daysAgo", { n: diffDay });

  return date.toLocaleDateString(getLocaleCode(), {
    month: "short",
    day: "numeric",
    year: date.getFullYear() !== now.getFullYear() ? "numeric" : undefined,
  });
}

export function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString(getLocaleCode(), {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

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
  chronopost: "Chronopost",
  laposte: "La Poste",
  vintedgo: "Vinted Go",
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
