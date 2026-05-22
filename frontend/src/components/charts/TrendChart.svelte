<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import Chart from "chart.js/auto";
  import type { HourlyTrendPoint } from "../../types/models";

  export let points: HourlyTrendPoint[] = [];
  export let label = "";
  export let unit = "";

  let canvasEl: HTMLCanvasElement | null = null;
  let chart: Chart | null = null;

  function formatLabel(hour: string) {
    try {
      return new Date(hour).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    } catch {
      return String(hour);
    }
  }

  onMount(() => {
    const labels = points.map((p) => formatLabel(p.hour));
    const data = points.map((p) => p.value);

    if (!canvasEl) return;

    chart = new Chart(canvasEl, {
      type: "line",
      data: {
        labels,
        datasets: [
          {
            label: `${label} ${unit ? `(${unit})` : ""}`,
            data,
            borderColor: "#2f8f2f",
            backgroundColor: "rgba(47,143,47,0.08)",
            tension: 0.25,
            pointRadius: 3,
            borderWidth: 2,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
          tooltip: { mode: "index", intersect: false },
        },
        scales: {
          x: { display: true },
          y: { display: true, beginAtZero: false },
        },
      },
    });
  });

  onDestroy(() => {
    chart?.destroy();
    chart = null;
  });
</script>

<div class="trend-chart">
  <div class="chart-header">
    <strong>{label}</strong>
    {#if unit}
      <small>{unit}</small>
    {/if}
  </div>
  <div class="chart-canvas">
    <canvas bind:this={canvasEl}></canvas>
  </div>
</div>

<style>
  .trend-chart {
    width: 100%;
    height: 200px;
  }

  .chart-canvas {
    width: 100%;
    height: calc(100% - 22px);
  }

  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
    color: #32511e;
  }
</style>
