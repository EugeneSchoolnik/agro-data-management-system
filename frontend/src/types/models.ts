// API Types
export interface Crop {
  id: number;
  name: string;
  variety: string;
  description?: string;
}

export interface Field {
  id: number;
  name: string;
  area: number;
  location?: string;
  crop_id: number;
  crop_name?: string;
  pest_id?: number;
  sensors?: Sensor[];
  created_at?: string;
}

export interface Sensor {
  id: number;
  field_id: number;
  sensor_type: "temperature" | "humidity" | "vision_node";
  status: "active" | "inactive" | "error" | "testing";
  last_sync?: string;
}

export interface Metric {
  id: number;
  sensor_id: number;
  value: number;
  recorded_at: string;
}

export interface Pest {
  id: number;
  name: string;
  scientific_name?: string;
  description?: string;
}

export interface Forecast {
  id: number;
  field_id: number;
  pest_id: number;
  probability: number;
  recommendation: string;
  created_at: string;
}

export interface WeatherParameter {
  id: number;
  param_id: number;
  name: string;
  unit: string;
  description: string;
}

export interface WeatherObservation {
  id: number;
  station_id: number;
  weather_parameter_id: number;
  station_param: number;
  value: number;
  recorded_at: string;
  created_at?: string;
  weather_parameter: WeatherParameter;
}

export interface WeatherParameterSummary {
  parameter: WeatherParameter;
  value: number;
  station_param: number;
  recorded_at: string;
}

export interface WeatherParameterAggregate {
  parameter: WeatherParameter;
  average: number;
  min: number;
  max: number;
  count: number;
}

export interface HourlyTrendPoint {
  hour: string;
  value: number;
}

export interface WeatherParameterTrend {
  parameter: WeatherParameter;
  points: HourlyTrendPoint[];
}

export interface WeatherStationSummary {
  station: WeatherStation;
  latest: WeatherParameterSummary[];
  daily: WeatherParameterAggregate[];
  hourlyTrend: WeatherParameterTrend[];
  updated_at: string;
}

export interface WeatherForecastResult {
  temperature: number;
  hours_ahead: number;
  recommendation: string;
}

export interface WeatherStation {
  id: number;
  external_id: number;
  name: string;
  region: string;
  active: boolean;
  last_seen?: string;
  created_at?: string;
}

export interface MetricSummary {
  avg: number;
  min: number;
  max: number;
}

export interface FieldReport {
  field_name: string;
  temperature: MetricSummary;
  air_humidity: MetricSummary;
  forecast_average_probability: number;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export interface FieldForm {
  name: string;
  area: string | number;
  location: string;
  crop_id: number | null;
}

export interface RiskLevel {
  level: "high" | "medium" | "low" | "unknown";
  color: string;
  percentage: number;
}

export interface User {
  id: number;
  email: string;
  role: string;
  password_hash?: string;
  created_at?: string;
  updated_at?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}
