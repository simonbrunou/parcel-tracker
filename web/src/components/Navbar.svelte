<script lang="ts">
  import { push } from "svelte-spa-router";
  import { logout } from "../lib/api";
  import { toggleTheme, isDarkTheme } from "../lib/utils";
  import { t, getLocale, setLocale, supportedLocales } from "../lib/i18n.svelte";
  import { isPushSupported, isPushEnabled, isPushLoading, togglePush } from "../lib/push.svelte";

  let dark = $state(isDarkTheme());

  function handleToggleTheme() {
    toggleTheme();
    dark = isDarkTheme();
  }

  async function handleLogout() {
    try {
      await logout();
    } catch {
      // Even if the server request fails, redirect to login.
      // The cookie will expire naturally.
    }
    push("/login");
  }

  function cycleLocale() {
    const codes = supportedLocales.map((l) => l.code);
    const idx = codes.indexOf(getLocale());
    setLocale(codes[(idx + 1) % codes.length]);
  }
</script>

<header class="sticky top-0 z-40 bg-[var(--color-surface)]/70 backdrop-blur-xl border-b border-[var(--color-border)]/60">
  <div class="max-w-4xl mx-auto px-4 h-14 flex items-center justify-between">
    <a
      href="#/"
      class="flex items-center gap-2.5 font-bold text-lg no-underline group"
    >
      <span class="relative flex items-center justify-center w-9 h-9 rounded-xl bg-gradient-to-br from-indigo-500 to-violet-600 shadow-md shadow-indigo-500/20 group-hover:shadow-indigo-500/40 group-hover:scale-105 transition-all duration-200">
        <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
      </span>
      <span class="bg-gradient-to-r from-[var(--color-text-primary)] to-[var(--color-text-primary)] bg-clip-text group-hover:from-indigo-500 group-hover:to-violet-600 group-hover:text-transparent transition-all">
        {t("app.name")}
      </span>
    </a>

    <div class="flex items-center gap-0.5">
      <button
        onclick={cycleLocale}
        class="w-9 h-9 flex items-center justify-center rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer text-xs font-bold uppercase tracking-wide"
        aria-label={t("nav.language")}
        title={t("nav.language")}
      >
        {getLocale()}
      </button>

      {#if isPushSupported()}
        <button
          onclick={togglePush}
          disabled={isPushLoading()}
          class="w-9 h-9 flex items-center justify-center rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors cursor-pointer disabled:opacity-50 {isPushEnabled() ? 'text-[var(--color-accent)]' : 'text-[var(--color-text-secondary)]'}"
          aria-label={t("nav.notifications")}
          title={isPushEnabled() ? t("nav.notificationsOn") : t("nav.notificationsOff")}
        >
          {#if isPushEnabled()}
            <svg class="w-[18px] h-[18px]" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/>
            </svg>
          {:else}
            <svg class="w-[18px] h-[18px]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
            </svg>
          {/if}
        </button>
      {/if}

      <button
        onclick={handleToggleTheme}
        class="w-9 h-9 flex items-center justify-center rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer"
        aria-label={t("nav.toggleTheme")}
      >
        {#if dark}
          <svg class="w-[18px] h-[18px]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        {:else}
          <svg class="w-[18px] h-[18px]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        {/if}
      </button>

      <div class="w-px h-5 bg-[var(--color-border)] mx-1"></div>

      <button
        onclick={handleLogout}
        class="w-9 h-9 flex items-center justify-center rounded-lg hover:bg-[var(--color-danger-light)] hover:text-[var(--color-danger)] transition-colors text-[var(--color-text-secondary)] cursor-pointer"
        aria-label={t("nav.logout")}
      >
        <svg class="w-[18px] h-[18px]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
      </button>
    </div>
  </div>
</header>
