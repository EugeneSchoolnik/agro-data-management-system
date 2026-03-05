import type {
  Crop,
  Field,
  Sensor,
  Metric,
  Pest,
  Forecast,
  ApiResponse,
  FieldForm,
  FieldReport,
} from "../types/models";

const API_BASE = "http://localhost:8080/api/v1";

interface RequestOptions extends RequestInit {
  headers?: Record<string, string>;
}

export async function apiRequest<T>(
  endpoint: string,
  options: RequestOptions = {},
): Promise<T> {
  const url = `${API_BASE}${endpoint}`;
  const config: RequestInit = {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    const data: ApiResponse<T> = await response.json();

    if (!response.ok) {
      throw new Error(data.error || `HTTP error! status: ${response.status}`);
    }

    return data.data as T;
  } catch (error) {
    console.error("API request failed:", error);
    throw error;
  }
}

// Crops
export const getCrops = () => apiRequest<Crop[]>("/crops");
export const getCrop = (id: number) => apiRequest<Crop>(`/crops/${id}`);
export const createCrop = (crop: Crop) =>
  apiRequest<Crop>("/crops", { method: "POST", body: JSON.stringify(crop) });
export const updateCrop = (id: number, crop: Partial<Crop>) =>
  apiRequest<Crop | string>(`/crops/${id}`, {
    method: "PUT",
    body: JSON.stringify(crop),
  });
export const deleteCrop = (id: number) =>
  apiRequest<void>(`/crops/${id}`, { method: "DELETE" });

// Fields
export const getFields = () => apiRequest<Field[]>("/fields");
export const getField = (id: number) => apiRequest<Field>(`/fields/${id}`);
export const createField = (field: FieldForm) =>
  apiRequest<Field>("/fields", { method: "POST", body: JSON.stringify(field) });
export const updateField = (id: number, field: Partial<FieldForm>) =>
  apiRequest<Field | string>(`/fields/${id}`, {
    method: "PUT",
    body: JSON.stringify(field),
  });
export const deleteField = (id: number) =>
  apiRequest<void>(`/fields/${id}`, { method: "DELETE" });
export const getFieldSensors = (id: number) =>
  apiRequest<Sensor[]>(`/fields/${id}/sensors`);
export const getFieldLatestForecast = (id: number) =>
  apiRequest<Forecast>(`/fields/${id}/forecast/latest`);
export const getFieldReport = (id: number, from: string, to: string) =>
  apiRequest<FieldReport>(
    `/reports/fields/${id}?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`,
  );

// Sensors
export const getSensor = (id: number) => apiRequest<Sensor>(`/sensors/${id}`);
export const createSensor = (sensor: Partial<Sensor>) =>
  apiRequest<Sensor>("/sensors", {
    method: "POST",
    body: JSON.stringify(sensor),
  });
export const updateSensorStatus = (id: number, status: Sensor["status"]) =>
  apiRequest<Sensor>(`/sensors/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
export const deleteSensor = (id: number) =>
  apiRequest<void>(`/sensors/${id}`, { method: "DELETE" });
export const getSensorLatestMetrics = (id: number) =>
  apiRequest<Metric>(`/sensors/${id}/metrics/latest`);
export const getSensorMetricsHistory = (id: number, from: string, to: string) =>
  apiRequest<Metric[]>(`/sensors/${id}/metrics/history?from=${from}&to=${to}`);

// Metrics
export const createMetric = (metric: Partial<Metric>) =>
  apiRequest<Metric>("/metrics", {
    method: "POST",
    body: JSON.stringify(metric),
  });

// Pests
export const getPests = () => apiRequest<Pest[]>("/pests");
export const getPest = (id: number) => apiRequest<Pest>(`/pests/${id}`);
export const createPest = (pest: Partial<Pest>) =>
  apiRequest<Pest>("/pests", { method: "POST", body: JSON.stringify(pest) });
export const updatePest = (id: number, pest: Partial<Pest>) =>
  apiRequest<Pest | string>(`/pests/${id}`, {
    method: "PUT",
    body: JSON.stringify(pest),
  });
export const deletePest = (id: number) =>
  apiRequest<void>(`/pests/${id}`, { method: "DELETE" });

// Forecasts
export const predictForecast = (fieldId: number, pestId: number) =>
  apiRequest<Forecast>("/forecasts/predict", {
    method: "POST",
    body: JSON.stringify({ field_id: fieldId, pest_id: pestId }),
  });
