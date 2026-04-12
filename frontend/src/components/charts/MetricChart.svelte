<script lang="ts">
  import { onMount } from "svelte";
  import { Chart } from "chart.js/auto";
  import type { Metric, Sensor } from "../../types/models";
  import { getMetricsHistory } from "../../stores/sensors";

  export let sensorId: number;
  export let sensorType: Sensor["sensor_type"] = "temperature";

  let canvas: HTMLCanvasElement;
  let chart: Chart | null = null;
  let data: Metric[] = [];
  let loading: boolean = false;
  let mounted = false;

  async function loadData(): Promise<void> {
    loading = true;
    try {
      const to = new Date().toISOString();
      const from = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString();
      data = await getMetricsHistory(sensorId, from, to);
      updateChart();
    } catch (err) {
      console.error("Failed to load metrics:", err);
    } finally {
      loading = false;
    }
  }

  function updateChart(): void {
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    if (chart) {
      chart.destroy();
    }

    const labels = data.map((d) =>
      // new Date(d.recorded_at).toLocaleDateString(),
      new Date(d.recorded_at).toLocaleString(),
    );
    const values = data.map((d) => d.value);

    const metricLabel =
      sensorType === "temperature"
        ? "Температура (°C)"
        : sensorType === "humidity"
          ? "Вологість (%)"
          : "Метрика";

    const chartTitle =
      sensorType === "temperature"
        ? "Історія температури"
        : sensorType === "humidity"
          ? "Історія вологості"
          : sensorType === "vision_node"
            ? "Історія візуального вузла"
            : "Історія метрик";

    const palette =
      sensorType === "temperature"
        ? { border: "#dc3545", background: "rgba(220, 53, 69, 0.1)" }
        : sensorType === "humidity"
          ? { border: "#007bff", background: "rgba(0, 123, 255, 0.1)" }
          : { border: "#6c757d", background: "rgba(108, 117, 125, 0.1)" };

    chart = new Chart(ctx, {
      type: "line",
      data: {
        labels,
        datasets: [
          {
            label: metricLabel,
            data: values,
            borderColor: palette.border,
            backgroundColor: palette.background,
            tension: 0.1,
          },
        ],
      },
      options: {
        responsive: true,
        plugins: {
          legend: {
            position: "top",
          },
          title: {
            display: true,
            text: chartTitle,
          },
        },
        scales: {
          y: {
            beginAtZero: false,
          },
        },
      },
    });
  }

  onMount(() => {
    mounted = true;
  });

  $: if (mounted && sensorId) {
    loadData();
  }
</script>

<div class="chart-container">
  <canvas bind:this={canvas} width="400" height="200"></canvas>
  {#if loading}
    <div class="loading-overlay">
      <p>Завантаження даних...</p>
    </div>
  {/if}
</div>

<style>
  .chart-container {
    position: relative;
    width: 100%;
    height: 300px;
    margin: 1rem 0;
  }

  canvas {
    max-width: 100%;
    height: auto;
  }

  .loading-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.85);
    color: #333;
    font-weight: 600;
  }
</style>
