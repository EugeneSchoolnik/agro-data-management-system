import { writable, get } from "svelte/store";
import type { Crop } from "../types/models";
import {
  getCrops,
  getCrop,
  createCrop,
  updateCrop,
  deleteCrop,
} from "../lib/api";

export const crops = writable<Crop[]>([]);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export async function loadCrops(): Promise<void> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getCrops();
    crops.set(data);
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
  } finally {
    loading.set(false);
  }
}

export async function loadCrop(id: number): Promise<Crop | null> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getCrop(id);
    return data;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return null;
  } finally {
    loading.set(false);
  }
}

export async function addCrop(crop: Omit<Crop, "id">): Promise<Crop> {
  try {
    const newCrop = await createCrop(crop as Crop);
    crops.update((current) => [...current, newCrop]);
    return newCrop;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function editCrop(id: number, crop: Partial<Crop>): Promise<Crop> {
  try {
    const response = await updateCrop(id, crop);
    const currentCrops = get(crops);
    const existing = currentCrops.find((c) => c.id === id);
    const mergedCrop: Crop = existing
      ? { ...existing, ...crop }
      : ({ id, ...crop } as Crop);
    const updatedCrop = typeof response === "string" ? mergedCrop : response;

    crops.update((current) =>
      current.map((c) => (c.id === id ? updatedCrop : c)),
    );
    return updatedCrop;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function removeCrop(id: number): Promise<void> {
  try {
    await deleteCrop(id);
    crops.update((current) => current.filter((c) => c.id !== id));
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}
