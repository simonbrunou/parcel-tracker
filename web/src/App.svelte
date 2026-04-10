<script lang="ts">
  import Router from "svelte-spa-router";
  import Dashboard from "./pages/Dashboard.svelte";
  import Login from "./pages/Login.svelte";
  import AddParcel from "./pages/AddParcel.svelte";
  import ParcelDetail from "./pages/ParcelDetail.svelte";
  import NotFound from "./pages/NotFound.svelte";
  import ToastContainer from "./components/ToastContainer.svelte";

  const routes = {
    "/": Dashboard,
    "/login": Login,
    "/parcels/new": AddParcel,
    "/parcels/:id": ParcelDetail,
    "*": NotFound,
  };
</script>

<Router {routes} />
<ToastContainer />
