<script lang="ts">
  import { push } from "svelte-spa-router";
  import Navbar from "../components/Navbar.svelte";
  import { t } from "../lib/i18n.svelte";
</script>

<Navbar />

<main class="max-w-2xl mx-auto px-4 py-6">
  <div class="text-center py-20 animate-fade-in">
    <div class="w-20 h-20 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-indigo-100 to-violet-100 dark:from-indigo-950/40 dark:to-violet-950/40 flex items-center justify-center">
      <span class="text-4xl">📦</span>
    </div>
    <h1 class="text-3xl font-bold tracking-tight text-[var(--color-text-primary)] mb-2">{t("notFound.title")}</h1>
    <p class="text-[var(--color-text-secondary)] mb-6 max-w-sm mx-auto">{t("notFound.message")}</p>
    <button
      onclick={() => push("/")}
      class="inline-flex items-center gap-2 px-5 py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-medium rounded-xl shadow-md shadow-indigo-500/20 hover:shadow-lg hover:shadow-indigo-500/30 hover:-translate-y-0.5 transition-all cursor-pointer"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
      {t("notFound.goHome")}
    </button>
  </div>
</main>
