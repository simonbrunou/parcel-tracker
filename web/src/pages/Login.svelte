<script lang="ts">
  import { push } from "svelte-spa-router";
  import { login, setup, checkAuth } from "../lib/api";
  import { t } from "../lib/i18n.svelte";

  let password = $state("");
  let confirmPassword = $state("");
  let error = $state("");
  let loading = $state(false);
  let configured = $state(true);
  let checking = $state(true);

  $effect(() => {
    checkAuth().then((res) => {
      if (res.authenticated) {
        push("/");
        return;
      }
      configured = res.configured;
      checking = false;
    });
  });

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      if (!configured) {
        if (password !== confirmPassword) {
          error = t("login.passwordsMismatch");
          loading = false;
          return;
        }
        if (password.length < 8) {
          error = t("login.passwordTooShort");
          loading = false;
          return;
        }
        await setup(password);
      } else {
        await login(password);
      }
      push("/");
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : t("login.failed");
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center p-4">
  {#if checking}
    <div class="animate-pulse text-[var(--color-text-muted)]">{t("common.loading")}</div>
  {:else}
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <div class="w-16 h-16 mx-auto mb-4 bg-[var(--color-accent-light)] rounded-2xl flex items-center justify-center">
          <svg class="w-8 h-8 text-[var(--color-accent)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-[var(--color-text-primary)]">{t("app.name")}</h1>
        <p class="text-[var(--color-text-secondary)] mt-1">
          {configured ? t("login.signIn") : t("login.createPassword")}
        </p>
      </div>

      <form onsubmit={handleSubmit} class="space-y-4">
        {#if error}
          <div class="bg-red-50 dark:bg-red-900/20 text-[var(--color-danger)] text-sm p-3 rounded-lg border border-red-200 dark:border-red-800">
            {error}
          </div>
        {/if}

        <div>
          <label for="password" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
            {t("login.password")}
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            required
            class="w-full px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
            placeholder={t("login.passwordPlaceholder")}
          />
        </div>

        {#if !configured}
          <div>
            <label for="confirm" class="block text-sm font-medium text-[var(--color-text-secondary)] mb-1.5">
              {t("login.confirmPassword")}
            </label>
            <input
              id="confirm"
              type="password"
              bind:value={confirmPassword}
              required
              class="w-full px-3 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
              placeholder={t("login.confirmPasswordPlaceholder")}
            />
          </div>
        {/if}

        <button
          type="submit"
          disabled={loading}
          class="w-full py-2.5 bg-[var(--color-accent)] hover:bg-[var(--color-accent-hover)] text-white font-medium rounded-lg transition-colors disabled:opacity-50 cursor-pointer"
        >
          {loading ? "..." : configured ? t("login.signInButton") : t("login.getStarted")}
        </button>
      </form>
    </div>
  {/if}
</div>
