<script lang="ts">
  import { onMount } from "svelte";
  import type { Pest } from "../../types/models";
  import {
    pests,
    loading,
    error,
    loadPests,
    addPest,
    editPest,
    removePest,
  } from "../../stores/pests";
  import Modal from "../common/Modal.svelte";
  import Button from "../common/Button.svelte";
  import PestForm from "./PestForm.svelte";

  let showCreateModal = false;
  let showEditModal = false;
  let createLoading = false;
  let editLoading = false;
  let deleteConfirmId: number | null = null;
  let showDeleteConfirm = false;
  let deleteLoading = false;
  let searchQuery = "";
  let selectedPest: Pest | null = null;

  onMount(async () => {
    await loadPests();
  });

  $: filteredPests = $pests.filter((pest) => {
    const query = searchQuery.trim().toLowerCase();
    if (!query) return true;
    return (
      pest.name.toLowerCase().includes(query) ||
      pest.scientific_name?.toLowerCase().includes(query) ||
      pest.description?.toLowerCase().includes(query)
    );
  });

  async function handleCreatePest(event: CustomEvent<Partial<Pest>>) {
    createLoading = true;
    try {
      await addPest(event.detail as Omit<Pest, "id">);
      showCreateModal = false;
    } catch (err) {
      console.error("Не вдалося додати шкідника:", err);
    } finally {
      createLoading = false;
    }
  }

  function openEditModal(pest: Pest) {
    selectedPest = pest;
    showEditModal = true;
  }

  async function handleEditPest(event: CustomEvent<Partial<Pest>>) {
    if (!selectedPest) return;
    editLoading = true;
    try {
      await editPest(selectedPest.id, event.detail);
      showEditModal = false;
      selectedPest = null;
    } catch (err) {
      console.error("Не вдалося оновити шкідника:", err);
    } finally {
      editLoading = false;
    }
  }

  function confirmDelete(pestId: number) {
    deleteConfirmId = pestId;
    showDeleteConfirm = true;
  }

  async function handleDeletePest() {
    if (deleteConfirmId === null) return;
    deleteLoading = true;
    try {
      await removePest(deleteConfirmId);
      deleteConfirmId = null;
    } catch (err) {
      console.error("Не вдалося видалити шкідника:", err);
    } finally {
      deleteLoading = false;
    }
  }
</script>

<div class="pests-management">
  <div class="management-header">
    <h2>Управління шкідниками</h2>
    <Button variant="primary" on:click={() => (showCreateModal = true)}>
      + Додати шкідника
    </Button>
  </div>

  <div class="search-bar">
    <input
      type="text"
      placeholder="Пошук за назвою або описом..."
      bind:value={searchQuery}
    />
  </div>

  {#if $loading}
    <p class="loading">Завантаження шкідників...</p>
  {:else if $error}
    <p class="error">Помилка: {$error}</p>
  {:else if filteredPests.length === 0}
    <div class="empty-state">
      <p>{searchQuery ? "Шкідники не знайдені" : "Поки що немає шкідників"}</p>
    </div>
  {:else}
    <div class="pests-table">
      <table>
        <thead>
          <tr>
            <th>Назва</th>
            <th>Наукова назва</th>
            <th>Опис</th>
            <th>Дії</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredPests as pest (pest.id)}
            <tr>
              <td>{pest.name}</td>
              <td>{pest.scientific_name || "—"}</td>
              <td>{pest.description || "—"}</td>
              <td>
                <div class="actions">
                  <Button
                    variant="secondary"
                    on:click={() => openEditModal(pest)}
                  >
                    Редагувати
                  </Button>
                  <Button
                    variant="danger"
                    on:click={() => confirmDelete(pest.id)}
                    disabled={deleteLoading}
                  >
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

  <Modal
    bind:open={showCreateModal}
    on:close={() => (showCreateModal = false)}
    title="Додати шкідника"
  >
    <PestForm
      on:submit={handleCreatePest}
      on:cancel={() => (showCreateModal = false)}
      loading={createLoading}
    />
  </Modal>

  <Modal
    bind:open={showEditModal}
    on:close={() => {
      showEditModal = false;
      selectedPest = null;
    }}
    title="Редагувати шкідника"
  >
    {#if selectedPest}
      <PestForm
        pest={selectedPest}
        isEditMode={true}
        on:submit={handleEditPest}
        on:cancel={() => {
          showEditModal = false;
          selectedPest = null;
        }}
        loading={editLoading}
      />
    {/if}
  </Modal>

  <Modal
    bind:open={showDeleteConfirm}
    on:close={() => {
      deleteConfirmId = null;
      showDeleteConfirm = false;
    }}
    title="Підтвердження видалення"
  >
    <div class="delete-confirmation">
      <p>Ви впевнені, що хочете видалити цього шкідника?</p>
      <p class="warning">Цю дію неможливо скасувати.</p>
    </div>
    <svelte:fragment slot="footer">
      <Button
        variant="secondary"
        on:click={() => {
          deleteConfirmId = null;
          showDeleteConfirm = false;
        }}
        disabled={deleteLoading}
      >
        Скасувати
      </Button>
      <Button
        variant="danger"
        on:click={handleDeletePest}
        loading={deleteLoading}
      >
        Видалити
      </Button>
    </svelte:fragment>
  </Modal>
</div>

<style>
  .pests-management {
    padding: 1.5rem;
  }

  .management-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .management-header h2 {
    margin: 0;
    color: #333;
    font-size: 1.5rem;
  }

  .search-bar {
    margin-bottom: 1.25rem;
  }

  .search-bar input {
    width: 100%;
    max-width: 420px;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 0.25rem;
    font-size: 1rem;
  }

  .pests-table {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    background: white;
    border-radius: 0.5rem;
    overflow: hidden;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04);
  }

  th,
  td {
    padding: 0.9rem 1rem;
    border-bottom: 1px solid #f1f1f1;
    text-align: left;
  }

  th {
    background: #f7f9fb;
    font-weight: 600;
    color: #444;
  }

  tr:last-child td {
    border-bottom: none;
  }

  .actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
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

  .delete-confirmation p {
    margin: 0.75rem 0;
  }

  .warning {
    color: #c00;
    font-weight: 600;
  }
</style>
