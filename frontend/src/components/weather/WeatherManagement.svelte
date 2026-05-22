<script lang="ts">
  import { onMount } from "svelte";
  import type {
    WeatherStation,
    WeatherStationSummary,
  } from "../../types/models";
  import { getWeatherStations, getWeatherStationSummary } from "../../lib/api";
  import Button from "../common/Button.svelte";
  import TrendChart from "../charts/TrendChart.svelte";

  let stations: WeatherStation[] = [];
  let selectedStation: WeatherStation | null = null;
  let stationSummary: WeatherStationSummary | null = null;
  let loadingStations = false;
  let loadingSummary = false;
  let stationsError = "";
  let summaryError = "";

  async function loadStations(): Promise<void> {
    loadingStations = true;
    stationsError = "";
    stationSummary = null;
    selectedStation = null;

    try {
      stations = await getWeatherStations();
      if (stations.length > 0) {
        selectedStation = stations[0];
        await loadStationSummary();
      }
    } catch (err) {
      stationsError =
        err instanceof Error
          ? err.message
          : "Не вдалося завантажити список станцій";
      console.error("Failed to load weather stations:", err);
    } finally {
      loadingStations = false;
    }
  }

  async function loadStationSummary(): Promise<void> {
    summaryError = "";
    stationSummary = null;
    if (!selectedStation) {
      return;
    }

    loadingSummary = true;
    try {
      stationSummary = await getWeatherStationSummary(
        selectedStation.external_id,
      );
    } catch (err) {
      summaryError =
        err instanceof Error
          ? err.message
          : "Не вдалося завантажити зведення погоди";
      console.error("Failed to load station weather summary:", err);
    } finally {
      loadingSummary = false;
    }
  }

  function selectStation(event: Event): void {
    const value = Number((event.target as HTMLSelectElement).value);
    selectedStation = stations.find((s) => s.external_id === value) ?? null;
    if (selectedStation) {
      loadStationSummary();
    }
  }

  onMount(loadStations);

  $: hourlyTrends = stationSummary?.hourlyTrend ?? [];
  $: trendsWithPoints = hourlyTrends.filter((t) => (t.points ?? []).length > 0);
</script>

<div class="weather-management">
  <div class="weather-header">
    <div>
      <h2>Погодні станції</h2>
      <p>Оберіть станцію з бази для перегляду останніх спостережень.</p>
    </div>
    <Button variant="primary" on:click={loadStations} loading={loadingStations}>
      Оновити список станцій
    </Button>
  </div>

  {#if loadingStations}
    <p>Завантаження станцій...</p>
  {:else if stationsError}
    <p class="error">{stationsError}</p>
  {:else if stations.length === 0}
    <p>Станції не знайдені.</p>
  {:else}
    <div class="weather-controls">
      <div class="station-select">
        <label for="station-select">Станція</label>
        <select id="station-select" on:change={selectStation}>
          {#each stations as station}
            <option value={station.external_id}>
              {station.name} ({station.region})
            </option>
          {/each}
        </select>
      </div>
      <div class="control-actions">
        <Button
          variant="secondary"
          on:click={loadStationSummary}
          loading={loadingSummary}
        >
          Оновити зведення
        </Button>
      </div>
    </div>

    {#if selectedStation}
      <div class="station-card">
        <h3>{selectedStation.name}</h3>
        <p><strong>Регіон:</strong> {selectedStation.region}</p>
        <p><strong>External ID:</strong> {selectedStation.external_id}</p>
        <p><strong>Активна:</strong> {selectedStation.active ? "Так" : "Ні"}</p>
      </div>

      <div class="weather-summary">
        <h3>Останнє зведення</h3>
        {#if loadingSummary}
          <p>Завантаження зведення...</p>
        {:else if summaryError}
          <p class="error">{summaryError}</p>
        {:else if !stationSummary}
          <p>Зведення не знайдено.</p>
        {:else}
          <p class="summary-updated">
            Оновлено: {new Date(stationSummary.updated_at).toLocaleString()}
          </p>
          <div class="latest-grid">
            {#each stationSummary.latest as entry}
              <div class="latest-card" title={entry.parameter.name}>
                <div class="obs-line compact">
                  <span>{entry.parameter.description}</span>
                  <span>{entry.value.toFixed(1)} {entry.parameter.unit}</span>
                </div>
              </div>
            {/each}
          </div>

          <div class="weather-aggregates">
            <h3>Добові показники</h3>
            {#if (stationSummary.daily ?? []).length === 0}
              <p>Дані агрегатів відсутні.</p>
            {:else}
              <table>
                <thead>
                  <tr>
                    <th>Параметр</th>
                    <th>Середнє</th>
                    <th>Мін</th>
                    <th>Макс</th>
                    <th>К-сть</th>
                  </tr>
                </thead>
                <tbody>
                  {#each stationSummary.daily ?? [] as agg}
                    <tr>
                      <td>{agg.parameter.name} ({agg.parameter.unit})</td>
                      <td>{agg.average.toFixed(1)}</td>
                      <td>{agg.min.toFixed(1)}</td>
                      <td>{agg.max.toFixed(1)}</td>
                      <td>{agg.count}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            {/if}
          </div>

          <div class="weather-trend">
            <h3>Годинний тренд</h3>
            {#if trendsWithPoints.length === 0}
              <p>Дані тренду відсутні.</p>
            {:else}
              <div class="trend-grid">
                {#each trendsWithPoints as trend}
                  <div class="trend-block">
                    <TrendChart
                      points={trend.points}
                      label={trend.parameter.name}
                      unit={trend.parameter.unit}
                    />
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {/if}
  {/if}
</div>

<style>
  .weather-management {
    display: grid;
    gap: 1.5rem;
  }

  .weather-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .weather-header h2 {
    margin: 0;
    color: #27421d;
  }

  .weather-controls {
    display: grid;
    grid-template-columns: 1fr 260px;
    gap: 1rem;
    align-items: center;
  }

  .station-select {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .control-actions {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  label {
    font-weight: 600;
    color: #2f4d1f;
  }

  select {
    width: 100%;
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    border: 1px solid rgba(76, 99, 57, 0.15);
    background: white;
    color: #2d3f21;
  }

  .station-card,
  .weather-summary {
    padding: 1.25rem;
    border-radius: 1rem;
    background: rgba(244, 250, 236, 0.95);
    border: 1px solid rgba(110, 138, 73, 0.16);
  }

  .station-card h3,
  .weather-summary h3 {
    margin-top: 0;
    color: #2f561c;
  }

  .latest-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 0.75rem;
  }

  .latest-card {
    padding: 0.8rem 1rem;
    border-radius: 0.8rem;
    background: white;
    border: 1px solid rgba(101, 132, 76, 0.12);
  }

  .obs-line.compact {
    display: flex;
    justify-content: space-between;
    gap: 0.5rem;
    flex-wrap: wrap;
    font-weight: 600;
    margin-bottom: 0.3rem;
  }

  .obs-description {
    margin: 0;
    color: #5e6f4d;
    font-size: 0.9rem;
  }

  .weather-aggregates,
  .weather-trend {
    margin-top: 1.5rem;
  }

  .trend-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 1rem;
  }

  .weather-aggregates table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 0.85rem;
  }

  .weather-aggregates th,
  .weather-aggregates td {
    text-align: left;
    padding: 0.75rem;
    border: 1px solid rgba(101, 132, 76, 0.12);
  }

  .weather-aggregates th {
    background: rgba(241, 249, 235, 0.95);
  }

  .trend-block {
    margin-top: 1rem;
    padding: 1rem;
    background: white;
    border-radius: 0.85rem;
    border: 1px solid rgba(101, 132, 76, 0.12);
  }

  .trend-block h4 {
    margin: 0 0 0.75rem 0;
    color: #32511e;
  }

  .trend-block ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    gap: 0.55rem;
  }

  .trend-block li {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    color: #3b4b32;
  }

  .obs-line {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
    color: #3b4b32;
  }

  .summary-updated {
    margin-bottom: 1rem;
    color: #586f44;
    font-size: 0.95rem;
  }

  .error {
    color: #842029;
    background: #f8d7da;
    padding: 1rem;
    border-radius: 0.85rem;
    border: 1px solid #f5c2c7;
  }
</style>
