import { writable, get } from "svelte/store";
import type { Field, FieldForm } from "../types/models";
import { crops } from "./crops";
import {
  getFields,
  getField,
  createField,
  updateField,
  deleteField,
} from "../lib/api";

export const fields = writable<Field[]>([]);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export async function loadFields(): Promise<void> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getFields();
    const cropList = get(crops);
    fields.set(
      data.map((field) => ({
        ...field,
        crop_name:
          field.crop_name ??
          cropList.find((crop) => crop.id === field.crop_id)?.name,
      })),
    );
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
  } finally {
    loading.set(false);
  }
}

export async function loadField(id: number): Promise<Field | null> {
  loading.set(true);
  error.set(null);
  try {
    const data = await getField(id);
    return data;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    return null;
  } finally {
    loading.set(false);
  }
}

export async function addField(field: FieldForm): Promise<Field> {
  try {
    const newField = await createField(field);
    const cropName = get(crops).find(
      (crop) => crop.id === newField.crop_id,
    )?.name;
    const fieldWithCropName = { ...newField, crop_name: cropName };
    fields.update((current) => [...current, fieldWithCropName]);
    return fieldWithCropName;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function editField(
  id: number,
  field: Partial<FieldForm>,
): Promise<Field> {
  try {
    const response = await updateField(id, field);
    const currentFields = get(fields);
    const existing = currentFields.find((f) => f.id === id);
    const mergedField: Field = existing
      ? ({ ...existing, ...field } as Field)
      : ({ id, ...field } as Field);
    const updatedField = typeof response === "string" ? mergedField : response;

    fields.update((current) =>
      current.map((f) => (f.id === id ? updatedField : f)),
    );

    return updatedField;
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}

export async function removeField(id: number): Promise<void> {
  try {
    await deleteField(id);
    fields.update((current) => current.filter((f) => f.id !== id));
  } catch (err) {
    const message = err instanceof Error ? err.message : "Unknown error";
    error.set(message);
    throw err;
  }
}
