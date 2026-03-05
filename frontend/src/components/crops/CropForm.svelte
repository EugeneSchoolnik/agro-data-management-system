<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { Crop } from "../../types/models";
  import Button from "../common/Button.svelte";

  export let crop: Partial<Crop> = {
    name: "",
    variety: "",
    description: "",
  };
  export let loading: boolean = false;
  export let isEditMode: boolean = false;

  const dispatch = createEventDispatcher<{
    submit: Partial<Crop>;
    cancel: void;
  }>();

  function submit(): void {
    if (!crop.name || !crop.name.trim()) {
      alert("Введіть назву культури");
      return;
    }
    if (!crop.variety || !crop.variety.trim()) {
      alert("Введіть сорт культури");
      return;
    }
    dispatch("submit", crop);
  }

  function cancel(): void {
    dispatch("cancel");
  }
</script>

<form on:submit|preventDefault={submit}>
  <div class="form-group">
    <label for="name">Назва культури:</label>
    <input
      id="name"
      type="text"
      bind:value={crop.name}
      required
      placeholder="Введіть назву (наприклад: Пшениця)"
    />
  </div>

  <div class="form-group">
    <label for="variety">Сорт:</label>
    <input
      id="variety"
      type="text"
      bind:value={crop.variety}
      required
      placeholder="Введіть сорт культури"
    />
  </div>

  <div class="form-group">
    <label for="description">Опис (опціонально):</label>
    <textarea
      id="description"
      bind:value={crop.description}
      placeholder="Введіть опис культури"
      rows="3"
    ></textarea>
  </div>

  <div class="form-actions">
    <Button variant="primary" type="submit" disabled={loading}>
      {loading ? "Збереження..." : isEditMode ? "Оновити" : "Додати"}
    </Button>
    <Button variant="secondary" on:click={cancel} disabled={loading}>
      Скасувати
    </Button>
  </div>
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    font-weight: 500;
    color: #333;
  }

  .form-group input,
  .form-group textarea {
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    font-family: inherit;
  }

  .form-group input:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }
</style>
