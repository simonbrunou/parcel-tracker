<script lang="ts">
  import { push } from "svelte-spa-router";
  import { logout } from "../lib/api";
  import { toggleTheme, isDarkTheme } from "../lib/utils";

  let dark = $state(isDarkTheme());

  function handleToggleTheme() {
    toggleTheme();
    dark = isDarkTheme();
  }

  async function handleLogout() {
    await logout();
    push("/login");
  }
</script>

<header class="sticky top-0 z-50 bg-[var(--color-surface)]/80 backdrop-blur-lg border-b border-[var(--color-border)]">
  <div class="max-w-4xl mx-auto px-4 h-14 flex items-center justify-between">
    <a href="#/" class="flex items-center gap-2 font-bold text-lg text-[var(--color-text-primary)] hover:text-[var(--color-accent)] transition-colors no-underline">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
      </svg>
      Parcel Tracker
    </a>

    <div class="flex items-center gap-1">
      <button
        onclick={handleToggleTheme}
        class="p-2 rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer"
        aria-label="Toggle theme"
      >
        {#if dark}
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        {:else}
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        {/if}
      </button>

      <button
        onclick={handleLogout}
        class="p-2 rounded-lg hover:bg-[var(--color-surface-hover)] transition-colors text-[var(--color-text-secondary)] cursor-pointer"
        aria-label="Logout"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
      </button>
    </div>
  </div>
</header>
