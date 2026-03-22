<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { Pest } from "../../types/models";
  import Button from "../common/Button.svelte";

  export let pest: Partial<Pest> = {
    name: "",
    scientific_name: "",
    description: "",
  };
  export let loading: boolean = false;
  export let isEditMode: boolean = false;

  const dispatch = createEventDispatcher<{
    submit: Partial<Pest>;
    cancel: void;
  }>();

  function submit(): void {
    if (!pest.name || !pest.name.trim()) {
      alert("Введіть назву шкідника");
      return;
    }
    dispatch("submit", pest);
  }

  function cancel(): void {
    dispatch("cancel");
  }
</script>

<form on:submit|preventDefault={submit}>
  <div class="form-group">
    <label for="name">Назва шкідника:</label>
    <input
      id="name"
      type="text"
      bind:value={pest.name}
      required
      placeholder="Наприклад: Тля"
    />
  </div>

  <div class="form-group">
    <label for="scientific_name">Наукова назва:</label>
    <input
      id="scientific_name"
      type="text"
      bind:value={pest.scientific_name}
      placeholder="Наприклад: Aphis"
    />
  </div>

  <div class="form-group">
    <label for="description">Опис (опціонально):</label>
    <textarea
      id="description"
      bind:value={pest.description}
      placeholder="Короткий опис шкідника"
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
