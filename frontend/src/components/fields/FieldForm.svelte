<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import type { Crop, FieldForm } from "../../types/models";
  import { getCrops } from "../../lib/api";
  import Button from "../common/Button.svelte";

  export let field: FieldForm = {
    name: "",
    area: "",
    location: "",
    crop_id: null,
  };
  export let loading: boolean = false;

  let locationLat: string | number = "";
  let locationLon: string | number = "";
  let crops: Crop[] = [];
  let cropLoadError = false;
  let previousLocation = "";

  function syncLocationInputs(): void {
    if (!field.location) {
      locationLat = "";
      locationLon = "";
      previousLocation = field.location;
      return;
    }

    if (field.location === previousLocation) {
      return;
    }

    previousLocation = field.location;
    const [lat = "", lon = ""] = field.location.trim().split(/[\s,]+/);
    locationLat = lat;
    locationLon = lon;
  }

  $: if (field.location && field.location !== previousLocation) {
    syncLocationInputs();
  }

  onMount(() => {
    syncLocationInputs();
  });

  const dispatch = createEventDispatcher<{ submit: FieldForm; cancel: void }>();

  async function loadCrops(): Promise<void> {
    try {
      crops = await getCrops();
    } catch (err) {
      console.error("Failed to load crops:", err);
      cropLoadError = true;
    }
  }

  loadCrops();

  function submit(): void {
    if (!field.name.trim()) {
      alert("Введіть назву поля");
      return;
    }
    if (!field.area || Number(field.area) <= 0) {
      alert("Введіть коректну площу поля");
      return;
    }
    const latString = String(locationLat ?? "").trim();
    const lonString = String(locationLon ?? "").trim();

    if (!latString || Number.isNaN(Number(latString))) {
      alert("Введіть коректну широту поля");
      return;
    }
    if (!lonString || Number.isNaN(Number(lonString))) {
      alert("Введіть коректну довготу поля");
      return;
    }
    if (!field.crop_id) {
      alert("Оберіть культуру");
      return;
    }

    field.location = `${latString} ${lonString}`;
    dispatch("submit", field);
  }

  function cancel(): void {
    dispatch("cancel");
  }
</script>

<form on:submit|preventDefault={submit}>
  <div class="form-group">
    <label for="name">Назва поля:</label>
    <input
      id="name"
      type="text"
      bind:value={field.name}
      required
      placeholder="Введіть назву поля"
    />
  </div>

  <div class="form-group">
    <label for="area">Площа (га):</label>
    <input
      id="area"
      type="number"
      step="0.01"
      bind:value={field.area}
      required
      placeholder="Введіть площу в гектарах"
    />
  </div>

  <div class="form-group location-row">
    <div class="location-field">
      <label for="location-lat">Широта:</label>
      <input
        id="location-lat"
        type="number"
        step="any"
        bind:value={locationLat}
        placeholder="Широта"
      />
    </div>
    <div class="location-field">
      <label for="location-lon">Довгота:</label>
      <input
        id="location-lon"
        type="number"
        step="any"
        bind:value={locationLon}
        placeholder="Довгота"
      />
    </div>
  </div>

  <div class="form-group">
    <label for="crop">Культура:</label>
    {#if cropLoadError}
      <div class="error-msg">Помилка завантаження культур</div>
    {:else if crops.length === 0}
      <select id="crop" bind:value={field.crop_id} disabled>
        <option>Завантаження...</option>
      </select>
    {:else}
      <select id="crop" bind:value={field.crop_id} required>
        <option value={null}>Оберіть культуру</option>
        {#each crops as crop}
          <option value={crop.id}>{crop.name} - {crop.variety}</option>
        {/each}
      </select>
    {/if}
  </div>

  <div class="form-actions">
    <Button type="button" variant="secondary" on:click={cancel}>
      Скасувати
    </Button>
    <Button type="submit" {loading}>
      {loading ? "Збереження..." : "Зберегти"}
    </Button>
  </div>
</form>

<style>
  form {
    max-width: 520px;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  label {
    display: block;
    margin-bottom: 0.55rem;
    font-weight: 600;
    color: #3d5127;
  }

  input,
  select {
    width: 100%;
    padding: 0.85rem 1rem;
    border: 1px solid rgba(126, 127, 94, 0.25);
    border-radius: 0.95rem;
    font-size: 1rem;
    background: #fff;
    transition:
      border-color 0.2s,
      box-shadow 0.2s;
  }

  input:focus,
  select:focus {
    outline: none;
    border-color: #7b8f44;
    box-shadow: 0 0 0 4px rgba(123, 143, 68, 0.16);
  }

  input::placeholder,
  select::placeholder {
    color: #8a876b;
  }

  input:disabled,
  select:disabled {
    background: #f3f0e4;
    cursor: not-allowed;
  }

  .error-msg {
    padding: 0.85rem 1rem;
    background: #f8d7da;
    color: #721c24;
    border: 1px solid #f5c2c7;
    border-radius: 0.85rem;
    font-size: 0.95rem;
  }

  .form-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    margin-top: 1.5rem;
  }

  .location-row {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
  }

  @media (max-width: 640px) {
    .location-row {
      grid-template-columns: 1fr;
    }
  }
</style>
