<script lang="ts">
  import { onMount } from "svelte";
  import type { Crop } from "../../types/models";
  import {
    crops,
    loading,
    error,
    loadCrops,
    addCrop,
    editCrop,
    removeCrop,
  } from "../../stores/crops";
  import CropForm from "./CropForm.svelte";
  import Modal from "../common/Modal.svelte";
  import Button from "../common/Button.svelte";

  let showCreateModal: boolean = false;
  let showEditModal: boolean = false;
  let showDeleteModal: boolean = false;
  let editingCrop: Partial<Crop> | null = null;
  let editingCropId: number | null = null;
  let createLoading: boolean = false;
  let editLoading: boolean = false;
  let deleteLoading: boolean = false;
  let searchQuery: string = "";
  let deleteConfirmId: number | null = null;

  onMount(() => {
    loadCrops();
  });

  $: filteredCrops = $crops.filter(
    (crop) =>
      (crop.name || "").toLowerCase().includes(searchQuery.toLowerCase()) ||
      (crop.variety || "").toLowerCase().includes(searchQuery.toLowerCase()),
  );

  function openCreateModal(): void {
    showCreateModal = true;
  }

  function closeCreateModal(): void {
    showCreateModal = false;
  }

  function openEditModal(crop: Crop): void {
    editingCrop = { ...crop };
    editingCropId = crop.id;
    showEditModal = true;
  }

  function closeEditModal(): void {
    showEditModal = false;
    editingCrop = null;
    editingCropId = null;
  }

  function handleCancelCreate(): void {
    showCreateModal = false;
  }

  function handleCancelEdit(): void {
    showEditModal = false;
    editingCrop = null;
    editingCropId = null;
  }

  async function handleCreateCrop(
    event: CustomEvent<Partial<Crop>>,
  ): Promise<void> {
    const cropData = event.detail;
    if (!cropData.name || !cropData.variety) return;
    createLoading = true;
    try {
      await addCrop({
        name: cropData.name,
        variety: cropData.variety,
        description: cropData.description,
      });
      closeCreateModal();
    } catch (err) {
      console.error("Failed to create crop:", err);
    } finally {
      createLoading = false;
    }
  }

  async function handleEditCrop(
    event: CustomEvent<Partial<Crop>>,
  ): Promise<void> {
    if (!editingCropId) return;
    editLoading = true;
    try {
      await editCrop(editingCropId, event.detail);
      closeEditModal();
    } catch (err) {
      console.error("Failed to edit crop:", err);
    } finally {
      editLoading = false;
    }
  }

  function handleDeleteCrop(id: number): void {
    deleteConfirmId = id;
    showDeleteModal = true;
  }

  async function confirmDeleteCrop(): Promise<void> {
    if (deleteConfirmId === null) return;
    deleteLoading = true;
    try {
      await removeCrop(deleteConfirmId);
      deleteConfirmId = null;
      showDeleteModal = false;
    } catch (err) {
      console.error("Failed to delete crop:", err);
    } finally {
      deleteLoading = false;
    }
  }
</script>

<div class="crops-management">
  <div class="management-header">
    <h2>Управління культурами</h2>
    <Button variant="primary" on:click={openCreateModal}>
      + Додати культуру
    </Button>
  </div>

  <div class="search-box">
    <input
      type="text"
      placeholder="Пошук по назві або сорту..."
      bind:value={searchQuery}
    />
  </div>

  {#if $loading}
    <p class="loading">Завантаження...</p>
  {:else if $error}
    <p class="error">Помилка: {$error}</p>
  {:else if filteredCrops.length === 0}
    <p class="empty">
      {searchQuery ? "Культури не знайдені" : "Немає культур. Додайте першу!"}
    </p>
  {:else}
    <div class="crops-table">
      <table>
        <thead>
          <tr>
            <th>Назва</th>
            <th>Сорт</th>
            <th>Опис</th>
            <th>Дії</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredCrops as crop (crop.id)}
            <tr>
              <td class="crop-name">{crop.name}</td>
              <td class="crop-variety">{crop.variety}</td>
              <td class="crop-description">
                {crop.description || "—"}
              </td>
              <td class="crop-actions">
                <Button
                  variant="secondary"
                  on:click={() => openEditModal(crop)}
                >
                  Редагувати
                </Button>
                <Button
                  variant="danger"
                  on:click={() => handleDeleteCrop(crop.id)}
                >
                  Видалити
                </Button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  {#if showCreateModal}
    <Modal
      bind:open={showCreateModal}
      on:close={closeCreateModal}
      title="Додати культуру"
    >
      <CropForm
        crop={{
          name: "",
          variety: "",
          description: "",
        }}
        loading={createLoading}
        isEditMode={false}
        on:submit={handleCreateCrop}
        on:cancel={handleCancelCreate}
      />
    </Modal>
  {/if}

  {#if showEditModal && editingCrop}
    <Modal
      bind:open={showEditModal}
      on:close={closeEditModal}
      title="Редагувати культуру"
    >
      <CropForm
        crop={editingCrop}
        loading={editLoading}
        isEditMode={true}
        on:submit={handleEditCrop}
        on:cancel={handleCancelEdit}
      />
    </Modal>
  {/if}

  <Modal
    bind:open={showDeleteModal}
    on:close={() => {
      deleteConfirmId = null;
      showDeleteModal = false;
    }}
    title="Підтвердження видалення"
  >
    <div class="delete-confirmation">
      <p>
        Ви впевнені, що хочете видалити культуру <strong
          >{$crops.find((c) => c.id === deleteConfirmId)?.name}</strong
        >?
      </p>
      <p class="warning">Цю дію неможливо отримати назад.</p>
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
        on:click={confirmDeleteCrop}
        loading={deleteLoading}
      >
        Видалити
      </Button>
    </svelte:fragment>
  </Modal>
</div>

<style>
  .crops-management {
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
    font-size: 1.5rem;
    color: #333;
  }

  .search-box {
    margin-bottom: 1.5rem;
  }

  .search-box input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }

  .search-box input:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
  }

  .loading,
  .error,
  .empty {
    text-align: center;
    padding: 2rem;
    font-size: 1.1rem;
  }

  .loading {
    color: #666;
  }

  .error {
    color: #dc3545;
  }

  .empty {
    color: #999;
  }

  .crops-table {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    background: white;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border-radius: 4px;
    overflow: hidden;
  }

  thead {
    background-color: #f5f5f5;
  }

  th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: #333;
    border-bottom: 2px solid #ddd;
  }

  td {
    padding: 1rem;
    border-bottom: 1px solid #ddd;
  }

  tr:hover {
    background-color: #f9f9f9;
  }

  .crop-name {
    font-weight: 500;
    color: #333;
  }

  .crop-variety {
    color: #666;
  }

  .crop-description {
    color: #999;
    font-size: 0.9rem;
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .crop-actions {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
</style>
