<script lang="ts">
  import { onMount } from "svelte";
  import type { Metric, Sensor } from "../../types/models";
  import { fields, loadFields } from "../../stores/fields";
  import {
    sensors,
    loading,
    error,
    loadSensors,
    addSensor,
    changeSensorStatus,
    removeSensor,
    getLatestMetrics,
  } from "../../stores/sensors";

  export let fieldId: number | null = null;
  export let fieldName: string | null = null;
  import Modal from "../common/Modal.svelte";

  $: activeField =
    fieldId !== null ? $fields.find((field) => field.id === fieldId) : null;
  import Button from "../common/Button.svelte";
  import SensorForm from "./SensorForm.svelte";
  import MetricChart from "../charts/MetricChart.svelte";

  let showCreateModal = false;
  let createLoading = false;
  let statusLoadingId: number | null = null;
  let deleteConfirmId: number | null = null;
  let showDeleteModal = false;
  let deleteLoading = false;
  let searchQuery = "";
  let selectedFieldId: number | null = null;
  let showLatestModal = false;
  let latestLoading = false;
  let latestMetric: Metric | null = null;
  let latestError: string | null = null;
  let showHistoryModal = false;
  let selectedHistorySensor: Sensor | null = null;

  let sensorsInitialized = false;

  onMount(async () => {
    await loadFields();
    if (fieldId !== null) {
      await loadSensors([fieldId]);
    } else {
      await loadSensors();
    }
    sensorsInitialized = true;
  });

  $: if (sensorsInitialized && fieldId !== null) {
    (async () => {
      await loadSensors([fieldId]);
    })();
  }

  $: fieldSensors = ($fields ?? []).map((field) => ({
    ...field,
    sensors: ($sensors ?? []).filter((sensor) => sensor.field_id === field.id),
  }));

  function matchesSensor(sensor: Sensor, fieldName: string): boolean {
    const query = searchQuery.toLowerCase();
    return (
      fieldName.includes(query) ||
      sensor.sensor_type.toLowerCase().includes(query) ||
      sensor.status.toLowerCase().includes(query) ||
      sensor.id.toString().includes(query)
    );
  }

  $: filteredFieldSensors = (
    fieldId !== null
      ? fieldSensors.filter((field) => field.id === fieldId)
      : fieldSensors
  )
    .map((field) => ({
      ...field,
      sensors: field.sensors.filter((sensor) =>
        matchesSensor(sensor, field.name.toLowerCase()),
      ),
    }))
    .filter((field) => field.sensors.length > 0 || searchQuery.trim() === "");

  function openCreateModal(fieldId: number) {
    selectedFieldId = fieldId;
    showCreateModal = true;
  }

  function formatStatus(status: Sensor["status"]): string {
    switch (status) {
      case "active":
        return "Активний";
      case "inactive":
        return "Неактивний";
      case "error":
        return "Помилка";
      case "testing":
        return "Тестування";
      default:
        return status;
    }
  }

  async function handleCreateSensor(
    event: CustomEvent<Partial<Sensor>>,
  ): Promise<void> {
    createLoading = true;
    try {
      await addSensor(event.detail);
      showCreateModal = false;
    } catch (err) {
      console.error("Не вдалося додати сенсор:", err);
    } finally {
      createLoading = false;
    }
  }

  async function handleStatusChange(
    sensor: Sensor,
    event: Event,
  ): Promise<void> {
    const select = event.target as HTMLSelectElement;
    const newStatus = select.value as Sensor["status"];
    statusLoadingId = sensor.id;
    try {
      await changeSensorStatus(sensor.id, newStatus);
    } catch (err) {
      console.error("Не вдалося оновити статус сенсора:", err);
    } finally {
      statusLoadingId = null;
    }
  }

  async function handleShowLatest(sensor: Sensor): Promise<void> {
    latestLoading = true;
    latestError = null;
    latestMetric = null;
    try {
      latestMetric = await getLatestMetrics(sensor.id);
      showLatestModal = true;
    } catch (err) {
      latestError = err instanceof Error ? err.message : "Unknown error";
      console.error("Не вдалося отримати останню метрику:", err);
    } finally {
      latestLoading = false;
    }
  }

  function closeLatestModal(): void {
    showLatestModal = false;
    latestMetric = null;
    latestError = null;
  }

  function handleShowHistory(sensor: Sensor): void {
    selectedHistorySensor = sensor;
    showHistoryModal = true;
  }

  async function confirmDeleteSensor(): Promise<void> {
    if (deleteConfirmId === null) return;
    deleteLoading = true;
    try {
      await removeSensor(deleteConfirmId);
      deleteConfirmId = null;
      showDeleteModal = false;
    } catch (err) {
      console.error("Не вдалося видалити сенсор:", err);
    } finally {
      deleteLoading = false;
    }
  }
</script>

<div class="sensors-management">
  <div class="management-header">
    <h2>
      {#if fieldId !== null}
        Сенсори для поля «{fieldName ?? activeField?.name ?? fieldId}»
      {:else}
        Управління сенсорами
      {/if}
    </h2>
  </div>

  <div class="search-box">
    <input
      type="text"
      placeholder="Пошук за полем, типом або статусом..."
      bind:value={searchQuery}
    />
  </div>

  {#if $loading}
    <p class="loading">Завантаження сенсорів...</p>
  {:else if $error}
    <p class="error">Помилка: {$error}</p>
  {:else if $fields.length === 0}
    <div class="empty-state">
      <p>Створіть поле, щоб додавати сенсори до нього.</p>
    </div>
  {:else if filteredFieldSensors.length === 0}
    <div class="empty-state">
      <p>{searchQuery ? "Сенсори не знайдені" : "Поки що немає сенсорів"}</p>
    </div>
  {:else}
    {#each filteredFieldSensors as field (field.id)}
      <section class="field-group">
        <div class="field-group-header">
          <div>
            <h3>{field.name}</h3>
            <p>{field.sensors.length} сенсорів</p>
          </div>
          <Button variant="primary" on:click={() => openCreateModal(field.id)}>
            + Додати сенсор до поля
          </Button>
        </div>

        <div class="sensors-table">
          {#if field.sensors.length > 0}
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Тип</th>
                  <th>Статус</th>
                  <th>Останнє синхр.</th>
                  <th>Дії</th>
                </tr>
              </thead>
              <tbody>
                {#each field.sensors as sensor (sensor.id)}
                  <tr>
                    <td>{sensor.id}</td>
                    <td>{sensor.sensor_type}</td>
                    <td>
                      <select
                        value={sensor.status}
                        on:change={(event) => handleStatusChange(sensor, event)}
                        disabled={statusLoadingId === sensor.id}
                      >
                        <option value="active">Активний</option>
                        <option value="inactive">Неактивний</option>
                        <option value="error">Помилка</option>
                        <option value="testing">Тестування</option>
                      </select>
                    </td>
                    <td>{sensor.last_sync ?? "—"}</td>
                    <td>
                      <Button
                        variant="secondary"
                        on:click={() => handleShowLatest(sensor)}
                        disabled={latestLoading}
                      >
                        Остання метрика
                      </Button>
                      <Button
                        variant="secondary"
                        on:click={() => handleShowHistory(sensor)}
                      >
                        Гісторія
                      </Button>
                      <Button
                        variant="danger"
                        on:click={() => {
                          deleteConfirmId = sensor.id;
                          showDeleteModal = true;
                        }}
                        disabled={deleteLoading}
                      >
                        Видалити
                      </Button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {:else}
            <p class="empty-field-message">
              У цьому полі поки що немає сенсорів.
            </p>
          {/if}
        </div>
      </section>
    {/each}
  {/if}

  <Modal
    bind:open={showCreateModal}
    on:close={() => {
      showCreateModal = false;
      selectedFieldId = null;
    }}
    title="Додати сенсор"
  >
    <SensorForm
      fieldId={selectedFieldId}
      on:submit={handleCreateSensor}
      on:cancel={() => {
        showCreateModal = false;
        selectedFieldId = null;
      }}
      loading={createLoading}
    />
  </Modal>

  <Modal
    bind:open={showLatestModal}
    on:close={closeLatestModal}
    title="Остання метрика"
  >
    {#if latestLoading}
      <p>Завантаження...</p>
    {:else if latestError}
      <p class="error">Помилка: {latestError}</p>
    {:else if latestMetric}
      <div class="metric-details">
        <p><strong>Значення:</strong> {latestMetric.value}</p>
        <p>
          <strong>Час:</strong>
          {new Date(latestMetric.recorded_at).toLocaleString()}
        </p>
      </div>
    {:else}
      <p>Дані недоступні</p>
    {/if}
    <svelte:fragment slot="footer">
      <Button variant="secondary" on:click={closeLatestModal}>Закрити</Button>
    </svelte:fragment>
  </Modal>

  <Modal
    bind:open={showHistoryModal}
    on:close={() => {
      showHistoryModal = false;
      selectedHistorySensor = null;
    }}
    title="Історія метрик"
  >
    {#if selectedHistorySensor}
      <MetricChart
        sensorId={selectedHistorySensor.id}
        sensorType={selectedHistorySensor.sensor_type}
      />
    {:else}
      <p>Виберіть сенсор для перегляду історії.</p>
    {/if}
    <svelte:fragment slot="footer">
      <Button
        variant="secondary"
        on:click={() => {
          showHistoryModal = false;
          selectedHistorySensor = null;
        }}
      >
        Закрити
      </Button>
    </svelte:fragment>
  </Modal>

  <Modal
    bind:open={showDeleteModal}
    on:close={() => {
      deleteConfirmId = null;
      showDeleteModal = false;
    }}
    title="Підтвердження видалення"
  >
    <div class="delete-confirmation">
      <p>Ви впевнені, що хочете видалити сенсор #{deleteConfirmId}?</p>
      <p class="warning">Цю дію неможливо скасувати.</p>
    </div>
    <svelte:fragment slot="footer">
      <Button
        variant="secondary"
        on:click={() => {
          deleteConfirmId = null;
          showDeleteModal = false;
        }}
        disabled={deleteLoading}
      >
        Скасувати
      </Button>
      <Button
        variant="danger"
        on:click={confirmDeleteSensor}
        loading={deleteLoading}
      >
        Видалити
      </Button>
    </svelte:fragment>
  </Modal>
</div>

<style>
  .sensors-management {
    padding: 1.5rem;
  }

  .management-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .management-header h2 {
    margin: 0;
    color: #333;
    font-size: 1.5rem;
  }

  .search-box {
    margin-bottom: 1.25rem;
  }

  .search-box input {
    width: 100%;
    max-width: 420px;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 0.25rem;
    font-size: 1rem;
  }

  .loading,
  .error,
  .empty-state {
    padding: 1.5rem;
    background: white;
    border-radius: 0.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  }

  .error {
    color: #c00;
  }

  .field-group {
    margin-bottom: 1.5rem;
    background: white;
    border-radius: 0.75rem;
    padding: 1rem;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  }

  .field-group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .field-group-header h3 {
    margin: 0;
    font-size: 1.1rem;
    color: #2d5016;
  }

  .field-group-header p {
    margin: 0.25rem 0 0;
    color: #666;
    font-size: 0.95rem;
  }

  .sensors-table table {
    width: 100%;
    border-collapse: collapse;
    background: white;
    border-radius: 0.5rem;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
  }

  th,
  td {
    padding: 0.95rem 1rem;
    text-align: left;
    border-bottom: 1px solid #f1f1f1;
  }

  th {
    background: #f7f9fb;
    font-weight: 600;
    color: #444;
  }

  tr:last-child td {
    border-bottom: none;
  }

  select {
    padding: 0.5rem 0.6rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
    font-size: 0.95rem;
    min-width: 150px;
  }

  .empty-field-message {
    padding: 1rem;
    margin: 0;
    color: #555;
  }

  .delete-confirmation {
    padding: 1rem 0;
  }

  .warning {
    margin-top: 0.75rem;
    color: #8a1f11;
  }
</style>
