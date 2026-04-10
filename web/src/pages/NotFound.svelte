<script lang="ts">
  import { push } from "svelte-spa-router";
  import Navbar from "../components/Navbar.svelte";
  import { t } from "../lib/i18n.svelte";
</script>

<Navbar />

<main class="max-w-2xl mx-auto px-4 py-6">
  <div class="text-center py-16">
    <svg class="w-16 h-16 mx-auto mb-4 text-[var(--color-text-muted)] opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mb-2">{t("notFound.title")}</h1>
    <p class="text-[var(--color-text-secondary)] mb-6">{t("notFound.message")}</p>
    <button
      onclick={() => push("/")}
      class="px-5 py-2.5 bg-[var(--color-accent)] hover:bg-[var(--color-accent-hover)] text-white font-medium rounded-xl transition-colors cursor-pointer"
    >
      {t("notFound.goHome")}
    </button>
  </div>
</main>
