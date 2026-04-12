import { writable } from "svelte/store";
import type { Forecast } from "../types/models";
import { predictForecast } from "../lib/api";

export const forecasts = writable<Forecast[]>([]);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export async function generateForecast(
  fieldId: number,
  pestId: number,
): Promise<Forecast> {
  loading.set(true);
  error.set(null);
  try {
    const data = await predictForecast(fieldId, pestId);
    forecasts.update((current) => [
      ...current,
      { ...data, fieldId, pestId, created_at: new Date().toISOString() },
    ]);
    return data;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  } finally {
    loading.set(false);
  }
}
