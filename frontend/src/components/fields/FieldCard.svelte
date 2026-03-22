<script lang="ts">
  import { onMount } from "svelte";
  import type {
    Field,
    Forecast,
    RiskLevel,
    FieldReport,
  } from "../../types/models";
  import {
    getFieldLatestForecast,
    predictForecast,
    getFieldReport,
  } from "../../lib/api";
  import { downloadFieldReportPDF } from "../../lib/pdf";
  import Button from "../common/Button.svelte";
  import Modal from "../common/Modal.svelte";

  export let field: Field;

  let forecast: Forecast | null = null;
  let forecastLoading: boolean = false;
  let showDetailsModal = false;
  let showReportModal = false;
  let report: FieldReport | null = null;
  let reportLoading: boolean = false;
  let reportError: string = "";
  const today = new Date().toISOString().slice(0, 10);
  let fromDate: string = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
    .toISOString()
    .slice(0, 10);
  let toDate: string = today;

  async function loadLatestForecast(): Promise<void> {
    forecastLoading = true;
    try {
      forecast = await getFieldLatestForecast(field.id);
    } catch (err) {
      console.error("Failed to load latest forecast:", err);
    } finally {
      forecastLoading = false;
    }
  }

  async function refreshForecast(): Promise<void> {
    forecastLoading = true;
    try {
      const pestId = field.pest_id ?? 1;
      if (!field.pest_id) {
        console.warn(
          "Field has no pest_id; using default pest_id=1 for forecast refresh.",
        );
      }
      forecast = await predictForecast(field.id, pestId);
    } catch (err) {
      console.error("Failed to refresh forecast:", err);
    } finally {
      forecastLoading = false;
    }
  }

  function openDetails(): void {
    showDetailsModal = true;
  }

  function closeDetails(): void {
    showDetailsModal = false;
  }

  function openReport(): void {
    report = null;
    reportError = "";
    showReportModal = true;
  }

  function closeReport(): void {
    showReportModal = false;
  }

  async function loadReport(): Promise<void> {
    reportError = "";
    report = null;

    if (!fromDate || !toDate) {
      reportError = "Вкажіть період звіту";
      return;
    }

    if (fromDate > toDate) {
      reportError = "Дата початку не може бути пізніше дати кінця";
      return;
    }

    reportLoading = true;
    try {
      report = await getFieldReport(
        field.id,
        `${fromDate}T00:00:00Z`,
        `${toDate}T23:59:59Z`,
      );
    } catch (err) {
      reportError =
        err instanceof Error ? err.message : "Не вдалося завантажити звіт";
    } finally {
      reportLoading = false;
    }
  }

  async function downloadReport(): Promise<void> {
    if (!report) {
      return;
    }

    try {
      await downloadFieldReportPDF(report, field.name, fromDate, toDate);
    } catch (err) {
      console.error("Не вдалося завантажити PDF:", err);
      reportError =
        err instanceof Error ? err.message : "Не вдалося згенерувати PDF";
    }
  }

  onMount(() => {
    loadLatestForecast();
  });

  function getRiskLevel(): Record<string, string> {
    if (!forecast) return { level: "unknown", color: "#6c757d" };

    const probability = forecast.probability;
    if (probability > 0.75) {
      return { level: "high", color: "#dc3545" };
    } else if (probability > 0.5) {
      return { level: "medium", color: "#ff9800" };
    } else {
      return { level: "low", color: "#28a745" };
    }
  }

  $: riskLevel = getRiskLevel();
  $: riskColor = riskLevel.color;
</script>

<div class="field-card">
  <div class="field-header">
    <h3>{field.name}</h3>
    <span class="crop-badge">{field.crop_name || "Невідома культура"}</span>
  </div>

  <div class="field-info">
    <p><strong>Площа:</strong> {field.area} га</p>
    <p><strong>Локація:</strong> {field.location || "Невідомо"}</p>
    <p><strong>Сенсори:</strong> {field.sensors?.length || 0}</p>
  </div>

  <div class="field-risk">
    <h4>Ризик шкідників</h4>
    {#if forecastLoading}
      <p>Завантаження прогнозу...</p>
    {:else if forecast}
      <div class="risk-indicator" style="background-color: {riskColor}">
        {Math.round(forecast.probability * 100)}% ймовірність
      </div>
      <p class="recommendation">{forecast.recommendation}</p>
    {:else}
      <p>Прогноз недоступний</p>
    {/if}
  </div>

  <div class="field-actions">
    <Button variant="secondary" on:click={openDetails}>Деталі</Button>
    <Button variant="secondary" on:click={openReport}>Звіт</Button>
    <Button
      variant="primary"
      on:click={refreshForecast}
      loading={forecastLoading}
    >
      Оновити прогноз
    </Button>
  </div>
</div>

<Modal bind:open={showReportModal} title="Звіт по полю" on:close={closeReport}>
  <div class="report-form">
    <div class="report-row">
      <label for="report-from">Початок</label>
      <input id="report-from" type="date" bind:value={fromDate} max={today} />
    </div>
    <div class="report-row">
      <label for="report-to">Кінець</label>
      <input id="report-to" type="date" bind:value={toDate} max={today} />
    </div>

    {#if reportError}
      <div class="report-error">{reportError}</div>
    {/if}

    <div class="report-actions">
      <Button variant="primary" on:click={loadReport} loading={reportLoading}>
        Створити звіт
      </Button>
      {#if report}
        <Button variant="secondary" on:click={downloadReport}>
          Завантажити PDF
        </Button>
      {/if}
    </div>

    {#if report}
      <div class="report-result">
        <h4>{report.field_name}</h4>
        <div class="report-grid">
          <div>
            <strong>Температура</strong>
            <p>Середнє: {report.temperature.avg.toFixed(1)}</p>
            <p>Мінімум: {report.temperature.min.toFixed(1)}</p>
            <p>Максимум: {report.temperature.max.toFixed(1)}</p>
          </div>
          <div>
            <strong>Вологість</strong>
            <p>Середнє: {report.air_humidity.avg.toFixed(1)}</p>
            <p>Мінімум: {report.air_humidity.min.toFixed(1)}</p>
            <p>Максимум: {report.air_humidity.max.toFixed(1)}</p>
          </div>
          <div>
            <strong>Прогноз</strong>
            <p>
              {Math.round(report.forecast_average_probability * 100)}% середня
              ймовірність
            </p>
          </div>
        </div>
      </div>
    {/if}
  </div>
</Modal>

<Modal bind:open={showDetailsModal} title="Деталі поля" on:close={closeDetails}>
  <div class="field-details">
    <p><strong>Назва:</strong> {field.name}</p>
    <p><strong>Культура:</strong> {field.crop_name || "Невідома"}</p>
    <p><strong>Площа:</strong> {field.area} га</p>
    <p><strong>Локація:</strong> {field.location || "Невідомо"}</p>
    <p><strong>Поточний прогноз:</strong></p>
    {#if forecast}
      <ul>
        <li>Ймовірність: {Math.round(forecast.probability * 100)}%</li>
        <li>Рекомендація: {forecast.recommendation}</li>
        <li>Створено: {new Date(forecast.created_at).toLocaleString()}</li>
      </ul>
    {:else}
      <p>Прогноз недоступний</p>
    {/if}

    <div class="sensor-details">
      <h4>Сенсори поля</h4>
      {#if field.sensors?.length}
        <ul>
          {#each field.sensors as sensor}
            <li>
              <strong>{sensor.sensor_type}</strong> — {sensor.status}
            </li>
          {/each}
        </ul>
      {:else}
        <p>Сенсори не знайдені</p>
      {/if}
    </div>
  </div>
</Modal>

<style>
  .field-card {
    background: white;
    border-radius: 0.5rem;
    padding: 1.5rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s;
  }

  .field-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }

  .field-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .field-header h3 {
    margin: 0;
    color: #2d5016;
  }

  .crop-badge {
    background: #e9ecef;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    font-size: 0.8rem;
    color: #495057;
  }

  .field-info p {
    margin: 0.5rem 0;
    color: #666;
  }

  .field-risk {
    margin: 1rem 0;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 0.25rem;
  }

  .field-risk h4 {
    margin: 0 0 0.5rem 0;
    color: #333;
  }

  .risk-indicator {
    display: inline-block;
    padding: 0.5rem 1rem;
    border-radius: 0.25rem;
    color: white;
    font-weight: bold;
    margin-bottom: 0.5rem;
  }

  .recommendation {
    font-size: 0.9rem;
    color: #666;
    margin: 0;
  }

  .field-actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  .report-form {
    display: grid;
    gap: 1rem;
  }

  .report-row {
    display: grid;
    gap: 0.5rem;
  }

  .report-actions {
    display: flex;
    justify-content: flex-end;
  }

  .report-error {
    color: #842029;
    background: #f8d7da;
    padding: 0.85rem 1rem;
    border-radius: 0.75rem;
  }

  .report-result {
    margin-top: 1.25rem;
    padding: 1rem;
    background: #f7f8f3;
    border-radius: 0.85rem;
    border: 1px solid rgba(126, 127, 94, 0.18);
  }

  .report-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
  }
</style>
