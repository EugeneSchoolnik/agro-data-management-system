<script lang="ts">
  import type { Field, FieldForm as FieldFormType } from "../../types/models";
  import {
    fields,
    loading,
    error,
    editField,
    removeField,
    addField,
  } from "../../stores/fields";
  import Modal from "../common/Modal.svelte";
  import Button from "../common/Button.svelte";
  import FieldFormComponent from "./FieldForm.svelte";
  import SensorsManagement from "../sensors/SensorsManagement.svelte";

  let showCreateModal = false;
  let showEditModal = false;
  let showDeleteConfirm = false;
  let selectedField: Field | null = null;
  let selectedFieldForSensors: Field | null = null;
  let formLoading = false;
  let deleteLoading = false;
  let searchQuery = "";

  $: filteredFields = $fields.filter(
    (field) =>
      field.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (field.location &&
        field.location.toLowerCase().includes(searchQuery.toLowerCase())),
  );

  function handleEdit(field: Field) {
    selectedField = field;
    showEditModal = true;
  }

  function handleDelete(field: Field) {
    selectedField = field;
    showDeleteConfirm = true;
  }

  async function confirmDelete() {
    if (!selectedField) return;
    deleteLoading = true;
    try {
      await removeField(selectedField.id);
      showDeleteConfirm = false;
      selectedField = null;
    } catch (err) {
      console.error("Failed to delete field:", err);
    } finally {
      deleteLoading = false;
    }
  }

  async function handleCreateField(event: CustomEvent<FieldFormType>) {
    formLoading = true;
    try {
      await addField(event.detail);
      showCreateModal = false;
    } catch (err) {
      console.error("Failed to create field:", err);
    } finally {
      formLoading = false;
    }
  }

  async function handleUpdateField(event: CustomEvent<FieldFormType>) {
    if (!selectedField) return;
    formLoading = true;
    try {
      await editField(selectedField.id, event.detail);
      showEditModal = false;
      selectedField = null;
    } catch (err) {
      console.error("Failed to update field:", err);
    } finally {
      formLoading = false;
    }
  }

  function openSensorManagement(field: Field) {
    selectedFieldForSensors = field;
  }

  function closeSensorManagement() {
    selectedFieldForSensors = null;
  }

  function closeModals() {
    showCreateModal = false;
    showEditModal = false;
    showDeleteConfirm = false;
    selectedField = null;
  }
</script>

<div class="fields-management">
  <div class="management-header">
    <h2>Управління полями</h2>
    <Button variant="primary" on:click={() => (showCreateModal = true)}>
      + Додати нове поле
    </Button>
  </div>

  <div class="search-bar">
    <input
      type="text"
      placeholder="Пошук за назвою або локацією..."
      bind:value={searchQuery}
    />
  </div>

  {#if $loading}
    <div class="loading">Завантаження полів...</div>
  {:else if $error}
    <div class="error-message">Помилка: {$error}</div>
  {:else if filteredFields.length === 0}
    <div class="empty-state">
      <p>Поля не знайдені</p>
      <Button variant="primary" on:click={() => (showCreateModal = true)}>
        Створити перше поле
      </Button>
    </div>
  {:else}
    <div class="fields-table">
      <table>
        <thead>
          <tr>
            <th>Назва</th>
            <th>Культура</th>
            <th>Площа (га)</th>
            <th>Локація</th>
            <th>Дії</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredFields as field (field.id)}
            <tr>
              <td class="field-name">
                <strong>{field.name}</strong>
              </td>
              <td>{field.crop_name || "—"}</td>
              <td class="area">{field.area}</td>
              <td>{field.location || "—"}</td>
              <td class="actions">
                <div class="action-buttons">
                  <Button
                    variant="secondary"
                    on:click={() => handleEdit(field)}
                  >
                    Редагувати
                  </Button>
                  <Button
                    variant="secondary"
                    on:click={() => openSensorManagement(field)}
                  >
                    Сенсори
                  </Button>
                  <Button variant="danger" on:click={() => handleDelete(field)}>
                    Видалити
                  </Button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  {#if selectedFieldForSensors}
    <div class="sensor-management-panel">
      <div class="panel-header">
        <h2>Сенсори для поля «{selectedFieldForSensors.name}»</h2>
        <Button variant="secondary" on:click={closeSensorManagement}>
          Закрити
        </Button>
      </div>
      <SensorsManagement
        fieldId={selectedFieldForSensors.id}
        fieldName={selectedFieldForSensors.name}
      />
    </div>
  {/if}
</div>

<Modal
  bind:open={showCreateModal}
  title="Створити нове поле"
  on:close={closeModals}
>
  <FieldFormComponent
    on:submit={handleCreateField}
    on:cancel={() => (showCreateModal = false)}
    loading={formLoading}
  />
</Modal>

<Modal bind:open={showEditModal} title="Редагувати поле" on:close={closeModals}>
  {#if selectedField}
    <FieldFormComponent
      field={{
        name: selectedField.name,
        area: selectedField.area,
        location: selectedField.location || "",
        crop_id: selectedField.crop_id,
      }}
      on:submit={handleUpdateField}
      on:cancel={() => (showEditModal = false)}
      loading={formLoading}
    />
  {/if}
</Modal>

<Modal
  bind:open={showDeleteConfirm}
  title="Підтвердження видалення"
  on:close={closeModals}
>
  <div class="delete-confirmation">
    <p>
      Ви впевнені, що хочете видалити поле <strong>{selectedField?.name}</strong
      >?
    </p>
    <p class="warning">Цю дію неможливо отримати назад.</p>
  </div>
  <svelte:fragment slot="footer">
    <Button
      variant="secondary"
      on:click={() => (showDeleteConfirm = false)}
      disabled={deleteLoading}
    >
      Скасувати
    </Button>
    <Button variant="danger" on:click={confirmDelete} loading={deleteLoading}>
      Видалити
    </Button>
  </svelte:fragment>
</Modal>

<style>
  .fields-management {
    padding: 0 2rem 2rem 2rem;
  }

  .management-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .management-header h2 {
    margin: 0;
    color: #333;
    font-size: 1.5rem;
  }

  .search-bar {
    margin-bottom: 1.5rem;
  }

  .search-bar input {
    width: 100%;
    max-width: 400px;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 0.25rem;
    font-size: 1rem;
  }

  .search-bar input:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }

  .loading {
    padding: 2rem;
    text-align: center;
    color: #666;
  }

  .error-message {
    padding: 1rem;
    background: #f8d7da;
    color: #721c24;
    border: 1px solid #f5c6cb;
    border-radius: 0.25rem;
    margin-bottom: 1rem;
  }

  .empty-state {
    text-align: center;
    padding: 3rem 2rem;
    background: white;
    border-radius: 0.5rem;
    border: 1px solid #dee2e6;
  }

  .sensor-management-panel {
    margin-top: 2rem;
    padding: 1.5rem;
    border-radius: 0.75rem;
    background: #fff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .panel-header h2 {
    margin: 0;
    font-size: 1.25rem;
    color: #2d5016;
  }

  .empty-state p {
    color: #666;
    margin: 0 0 1rem 0;
  }

  .fields-table {
    background: white;
    border-radius: 0.5rem;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: #f8f9fa;
    border-bottom: 2px solid #dee2e6;
  }

  th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #495057;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  td {
    padding: 1rem;
    border-bottom: 1px solid #dee2e6;
  }

  tbody tr:hover {
    background: #f8f9fa;
  }

  .field-name {
    color: #2d5016;
  }

  .area {
    text-align: right;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
    justify-content: flex-end;
  }

  .action-buttons {
    display: flex;
    gap: 0.5rem;
  }

  .action-buttons :global(.btn) {
    padding: 0.35rem 0.75rem;
    font-size: 0.85rem;
  }

  .delete-confirmation {
    padding: 1rem 0;
  }

  .delete-confirmation p {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .delete-confirmation .warning {
    color: #dc3545;
    font-size: 0.9rem;
  }
</style>
