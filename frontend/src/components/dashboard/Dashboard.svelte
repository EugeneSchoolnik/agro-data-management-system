<script lang="ts">
  import { onMount } from "svelte";
  import type { Field, FieldForm as FieldFormType } from "../../types/models";
  import {
    fields,
    loading,
    error,
    loadFields,
    addField,
  } from "../../stores/fields";
  import { sensors, loadSensors } from "../../stores/sensors";
  import FieldCard from "../fields/FieldCard.svelte";
  import FieldFormComponent from "../fields/FieldForm.svelte";
  import FieldsManagement from "../fields/FieldsManagement.svelte";
  import CropsManagement from "../crops/CropsManagement.svelte";
  import PestsManagement from "../pests/PestsManagement.svelte";
  import Modal from "../common/Modal.svelte";
  import Button from "../common/Button.svelte";

  type ViewType = "dashboard" | "fields" | "crops" | "pests";

  let selectedView: ViewType = "dashboard";
  let showCreateFieldModal: boolean = false;
  let createFieldLoading: boolean = false;

  onMount(async () => {
    await loadFields();
    await loadSensors();
  });

  $: fieldSensors = ($fields ?? []).map((field) => ({
    ...field,
    sensors: ($sensors ?? []).filter(
      (s): s is NonNullable<typeof s> => s != null && s.field_id === field.id,
    ),
  }));

  async function handleCreateField(
    event: CustomEvent<FieldFormType>,
  ): Promise<void> {
    createFieldLoading = true;
    try {
      await addField(event.detail);
      showCreateFieldModal = false;
    } catch (err) {
      console.error("Failed to create field:", err);
    } finally {
      createFieldLoading = false;
    }
  }
</script>

<div class="dashboard">
  <header class="dashboard-header">
    <h1>Агро-Дані: Управління Полями</h1>
    <nav class="nav">
      <Button
        variant={selectedView === "dashboard" ? "primary" : "secondary"}
        on:click={() => (selectedView = "dashboard")}
      >
        Дашборд
      </Button>
      <Button
        variant={selectedView === "fields" ? "primary" : "secondary"}
        on:click={() => (selectedView = "fields")}
      >
        Поля
      </Button>
      <Button
        variant={selectedView === "crops" ? "primary" : "secondary"}
        on:click={() => (selectedView = "crops")}
      >
        Культури
      </Button>
      <Button
        variant={selectedView === "pests" ? "primary" : "secondary"}
        on:click={() => (selectedView = "pests")}
      >
        Шкідники
      </Button>
    </nav>
  </header>

  <main class="dashboard-content">
    {#if selectedView === "dashboard"}
      <div class="overview">
        <div class="overview-header">
          <h2>Огляд</h2>
          <Button
            variant="primary"
            on:click={() => (showCreateFieldModal = true)}
          >
            + Додати поле
          </Button>
        </div>
        {#if $loading}
          <p>Завантаження...</p>
        {:else if $error}
          <p class="error">Помилка: {$error}</p>
        {:else}
          <div class="stats">
            <div class="stat-card">
              <h3>Поля</h3>
              <p class="stat-number">{$fields.length}</p>
            </div>
            <div class="stat-card">
              <h3>Активні сенсори</h3>
              <p class="stat-number">
                {$sensors.filter((s) => s.status === "active").length}
              </p>
            </div>
          </div>

          <div class="fields-grid">
            {#each fieldSensors as field}
              <FieldCard {field} />
            {/each}
          </div>
        {/if}
      </div>
    {:else if selectedView === "fields"}
      <FieldsManagement />
    {:else if selectedView === "crops"}
      <CropsManagement />
    {:else if selectedView === "pests"}
      <PestsManagement />
    {/if}
  </main>
</div>

<Modal
  bind:open={showCreateFieldModal}
  on:close={() => (showCreateFieldModal = false)}
  title="Створити нове поле"
>
  <FieldFormComponent
    on:submit={handleCreateField}
    on:cancel={() => (showCreateFieldModal = false)}
    loading={createFieldLoading}
  />
</Modal>

<style>
  .dashboard {
    min-height: 100vh;
    background: transparent;
  }

  .dashboard-header {
    background: rgba(255, 255, 255, 0.94);
    padding: 1.4rem 2rem;
    box-shadow: 0 24px 50px rgba(53, 74, 28, 0.14);
    display: flex;
    justify-content: space-between;
    align-items: center;
    border: 1px solid rgba(134, 126, 87, 0.18);
    border-radius: 1.25rem;
    backdrop-filter: blur(6px);
  }

  .dashboard-header h1 {
    margin: 0;
    color: #314b1f;
    font-size: clamp(1.9rem, 2.6vw, 2.8rem);
  }

  .nav {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .dashboard-content {
    padding: 1.5rem 0 0;
  }

  .overview h2 {
    margin-bottom: 1rem;
    color: #2f4d1f;
  }

  .overview-header {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
    padding: 1.2rem 1.2rem;
    border-radius: 1rem;
    background: linear-gradient(
      135deg,
      rgba(243, 241, 226, 0.96),
      rgba(255, 255, 255, 0.95)
    );
    border: 1px solid rgba(136, 147, 96, 0.16);
  }

  .overview-header h2 {
    margin: 0;
  }

  .stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    background: linear-gradient(180deg, #ffffff 0%, #f8f6ea 100%);
    padding: 1.5rem;
    border-radius: 1rem;
    box-shadow: 0 16px 36px rgba(62, 82, 30, 0.1);
    text-align: center;
    border: 1px solid rgba(128, 127, 89, 0.16);
  }

  .stat-card h3 {
    margin: 0 0 0.55rem 0;
    color: #5d6c38;
    font-size: 0.95rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .stat-number {
    font-size: 2.2rem;
    font-weight: 700;
    color: #30481f;
    margin: 0;
  }

  .fields-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
  }

  .error {
    color: #842029;
    background: #f8d7da;
    padding: 1rem;
    border-radius: 0.85rem;
    border: 1px solid #f5c2c7;
  }
</style>
