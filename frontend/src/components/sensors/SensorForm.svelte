<script lang="ts">
  import { onMount, createEventDispatcher } from "svelte";
  import type { Sensor, Field } from "../../types/models";
  import { fields, loadFields } from "../../stores/fields";
  import Button from "../common/Button.svelte";

  export let sensor: Partial<Sensor> = {
    sensor_type: "temperature",
    status: "active",
  };
  export let fieldId: number | null = null;
  export let loading: boolean = false;

  let selectedFieldId: number | null = sensor.field_id ?? fieldId;
  $: if (fieldId !== null && selectedFieldId !== fieldId) {
    selectedFieldId = fieldId;
  }

  let sensorType: Sensor["sensor_type"] = sensor.sensor_type ?? "temperature";
  let status: Sensor["status"] = sensor.status ?? "active";
  let showFieldError = false;
  let showTypeError = false;

  const dispatch = createEventDispatcher<{
    submit: Partial<Sensor>;
    cancel: void;
  }>();

  onMount(() => {
    loadFields();
  });

  function submit(): void {
    showFieldError = false;
    showTypeError = false;

    if (!selectedFieldId) {
      showFieldError = true;
      return;
    }

    if (!sensorType) {
      showTypeError = true;
      return;
    }

    dispatch("submit", {
      field_id: selectedFieldId,
      sensor_type: sensorType,
      status,
    });
  }

  function cancel(): void {
    dispatch("cancel");
  }
</script>

<form on:submit|preventDefault={submit}>
  <div class="form-group">
    <label for="field">Поле:</label>
    <select
      id="field"
      bind:value={selectedFieldId}
      required
      disabled={fieldId !== null}
      class:selected-error={showFieldError}
    >
      <option value={null}>Оберіть поле</option>
      {#each $fields as field}
        <option value={field.id}>{field.name}</option>
      {/each}
    </select>
    {#if showFieldError}
      <p class="validation-error">Оберіть поле для сенсора</p>
    {/if}
  </div>

  <div class="form-group">
    <label for="type">Тип сенсора:</label>
    <select
      id="type"
      bind:value={sensorType}
      required
      class:selected-error={showTypeError}
    >
      <option value="temperature">Температурний</option>
      <option value="humidity">Вологість</option>
      <option value="vision_node">Візійний</option>
    </select>
    {#if showTypeError}
      <p class="validation-error">Оберіть тип сенсора</p>
    {/if}
  </div>

  <div class="form-group">
    <label for="status">Статус:</label>
    <select id="status" bind:value={status}>
      <option value="active">Active</option>
      <option value="inactive">Inactive</option>
      <option value="error">Error</option>
      <option value="testing">Testing</option>
    </select>
  </div>

  <div class="form-actions">
    <Button type="button" variant="secondary" on:click={cancel}>
      Скасувати
    </Button>
    <Button type="submit" {loading}>
      {loading ? "Збереження..." : "Додати сенсор"}
    </Button>
  </div>
</form>

<style>
  form {
    max-width: 480px;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #333;
  }

  select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 0.25rem;
    font-size: 1rem;
  }

  select:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.18);
  }

  .validation-error {
    margin: 0.5rem 0 0;
    color: #d9534f;
    font-size: 0.9rem;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1.5rem;
  }
</style>
