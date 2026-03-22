import { get } from "svelte/store";
import { writable } from "svelte/store";
import type { Sensor, Metric } from "../types/models";
import { fields } from "./fields";
import {
  getSensor,
  createSensor,
  updateSensorStatus,
  deleteSensor,
  getFieldSensors,
  getSensorLatestMetrics,
  getSensorMetricsHistory,
} from "../lib/api";

export const sensors = writable<Sensor[]>([]);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export async function loadSensor(id: number): Promise<Sensor | null> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getSensor(id);
    return data;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return null;
  } finally {
    loading.set(false);
  }
}

export async function loadSensors(fieldIds?: number[]): Promise<void> {
  loading.set(true);
  error.set(null);
  try {
    const ids = fieldIds?.length
      ? fieldIds
      : get(fields).map((field) => field.id);

    if (ids.length === 0) {
      sensors.set([]);
      return;
    }

    const results = await Promise.all(ids.map((id) => getFieldSensors(id)));
    sensors.set(
      results
        .flat()
        .filter((item): item is NonNullable<typeof item> => item != null),
    );
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
  } finally {
    loading.set(false);
  }
}

export async function loadFieldSensors(fieldId: number): Promise<Sensor[]> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getFieldSensors(fieldId);
    sensors.update((current) => {
      const otherSensors = current.filter((s) => s.field_id !== fieldId);
      return [...otherSensors, ...data];
    });
    return data;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return [];
  } finally {
    loading.set(false);
  }
}

export async function addSensor(sensor: Partial<Sensor>): Promise<Sensor> {
  try {
    const newSensor = await createSensor(sensor);
    sensors.update((current) => [...current, newSensor]);
    return newSensor;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function changeSensorStatus(
  id: number,
  status: Sensor["status"],
): Promise<Sensor> {
  try {
    const updatedSensor = await updateSensorStatus(id, status);
    sensors.update((current) =>
      current.map((s) => (s.id === id ? updatedSensor : s)),
    );
    return updatedSensor;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function removeSensor(id: number): Promise<void> {
  try {
    await deleteSensor(id);
    sensors.update((current) => current.filter((s) => s.id !== id));
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function getLatestMetrics(
  sensorId: number,
): Promise<Metric | null> {
  try {
    return await getSensorLatestMetrics(sensorId);
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return null;
  }
}

export async function getMetricsHistory(
  sensorId: number,
  from: string,
  to: string,
): Promise<Metric[]> {
  try {
    return await getSensorMetricsHistory(sensorId, from, to);
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return [];
  }
}
