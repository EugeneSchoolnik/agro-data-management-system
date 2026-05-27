<script lang="ts">
  import Dashboard from "./components/dashboard/Dashboard.svelte";
  import Login from "./components/auth/Login.svelte";
  import { authStore } from "./stores/auth";
</script>

{#if $authStore.isAuthenticated}
  <Dashboard />
{:else}
  <Login />
{/if}
