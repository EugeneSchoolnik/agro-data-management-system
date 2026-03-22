import { writable, get } from "svelte/store";
import type { Pest } from "../types/models";
import {
  getPests,
  getPest,
  createPest,
  updatePest,
  deletePest,
} from "../lib/api";

export const pests = writable<Pest[]>([]);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export async function loadPests(): Promise<void> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getPests();
    pests.set(data);
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
  } finally {
    loading.set(false);
  }
}

export async function loadPest(id: number): Promise<Pest | null> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getPest(id);
    return data;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return null;
  } finally {
    loading.set(false);
  }
}

export async function addPest(pest: Omit<Pest, "id">): Promise<Pest> {
  try {
    const newPest = await createPest(pest as Partial<Pest>);
    pests.update((current) => [...current, newPest]);
    return newPest;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function editPest(id: number, pest: Partial<Pest>): Promise<Pest> {
  try {
    const response = await updatePest(id, pest);
    const currentPests = get(pests);
    const existing = currentPests.find((item) => item.id === id);
    const mergedPest: Pest = existing
      ? { ...existing, ...pest }
      : ({ id, ...pest } as Pest);
    const updatedPest = typeof response === "string" ? mergedPest : response;

    pests.update((current) =>
      current.map((item) => (item.id === id ? updatedPest : item)),
    );
    return updatedPest;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function removePest(id: number): Promise<void> {
  try {
    await deletePest(id);
    pests.update((current) => current.filter((item) => item.id !== id));
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}
