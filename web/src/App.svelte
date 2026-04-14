<script lang="ts">
  import Router from "svelte-spa-router";
  import { wrap } from "svelte-spa-router/wrap";
  import { checkAuth } from "./lib/api";
  import Dashboard from "./pages/Dashboard.svelte";
  import Login from "./pages/Login.svelte";
  import AddParcel from "./pages/AddParcel.svelte";
  import ParcelDetail from "./pages/ParcelDetail.svelte";
  import NotFound from "./pages/NotFound.svelte";
  import ToastContainer from "./components/ToastContainer.svelte";

  let authChecked = false;
  let authenticated = false;

  async function requireAuth(): Promise<boolean> {
    if (!authChecked) {
      try {
        const res = await checkAuth();
        authenticated = res.authenticated;
      } catch {
        authenticated = false;
      }
      authChecked = true;
    }
    if (!authenticated) {
      window.location.hash = "#/login";
      return false;
    }
    return true;
  }

  // Reset auth cache on navigation to login (user may have logged out).
  function allowAll(): boolean {
    authChecked = false;
    authenticated = false;
    return true;
  }

  const routes = {
    "/": wrap({ component: Dashboard, conditions: [requireAuth] }),
    "/login": wrap({ component: Login, conditions: [allowAll] }),
    "/parcels/new": wrap({ component: AddParcel, conditions: [requireAuth] }),
    "/parcels/:id": wrap({ component: ParcelDetail, conditions: [requireAuth] }),
    "*": NotFound,
  };
</script>

<Router {routes} />
<ToastContainer />
