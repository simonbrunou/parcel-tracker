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

<div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
  <!-- Decorative background orbs -->
  <div aria-hidden="true" class="pointer-events-none absolute inset-0">
    <div class="absolute -top-24 -left-16 w-72 h-72 rounded-full bg-indigo-400/25 blur-3xl"></div>
    <div class="absolute bottom-0 -right-16 w-80 h-80 rounded-full bg-violet-400/20 blur-3xl"></div>
  </div>

  {#if checking}
    <div class="flex flex-col items-center gap-3 text-[var(--color-text-muted)]">
      <div class="w-8 h-8 border-2 border-[var(--color-border)] border-t-[var(--color-accent)] rounded-full animate-spin"></div>
      <span>{t("common.loading")}</span>
    </div>
  {:else}
    <div class="w-full max-w-sm relative animate-slide-up">
      <div class="text-center mb-8">
        <div class="w-16 h-16 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-indigo-500 to-violet-600 flex items-center justify-center shadow-lg shadow-indigo-500/30 animate-pop">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <h1 class="text-3xl font-bold tracking-tight bg-gradient-to-br from-[var(--color-text-primary)] via-indigo-500 to-violet-600 bg-clip-text text-transparent">{t("app.name")}</h1>
        <p class="text-[var(--color-text-secondary)] mt-2">
          {configured ? t("login.signIn") : t("login.createPassword")}
        </p>
      </div>

      <div class="bg-[var(--color-surface)]/80 backdrop-blur-xl border border-[var(--color-border)] rounded-2xl p-6 shadow-[var(--shadow-lg)]">
        <form onsubmit={handleSubmit} class="space-y-4">
          {#if error}
            <div class="flex items-start gap-2 bg-[var(--color-danger-light)] text-[var(--color-danger)] text-sm p-3 rounded-lg border border-[var(--color-danger)]/20 animate-scale-in">
              <svg class="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{error}</span>
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
              autocomplete={configured ? "current-password" : "new-password"}
              class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
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
                autocomplete="new-password"
                class="w-full px-3.5 py-2.5 bg-[var(--color-surface-alt)] border border-[var(--color-border)] rounded-lg text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)] focus:border-transparent transition-all"
                placeholder={t("login.confirmPasswordPlaceholder")}
              />
            </div>
          {/if}

          <button
            type="submit"
            disabled={loading}
            class="w-full py-2.5 bg-gradient-to-br from-indigo-500 to-violet-600 hover:from-indigo-600 hover:to-violet-700 text-white font-semibold rounded-lg shadow-md shadow-indigo-500/20 hover:shadow-lg hover:shadow-indigo-500/30 transition-all disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer flex items-center justify-center gap-2"
          >
            {#if loading}
              <span class="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin"></span>
            {:else}
              {configured ? t("login.signInButton") : t("login.getStarted")}
            {/if}
          </button>
        </form>
      </div>
    </div>
  {/if}
</div>
