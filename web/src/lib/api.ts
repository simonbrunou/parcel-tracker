const BASE = "/api";

class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    ...options,
    credentials: "same-origin",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (res.status === 401) {
    window.location.hash = "#/login";
    throw new ApiError(401, "Unauthorized");
  }

  if (res.status === 204) {
    return undefined as T;
  }

  const data = await res.json();

  if (!res.ok) {
    throw new ApiError(res.status, data.error || "Unknown error");
  }

  return data as T;
}

// Auth
export const checkAuth = () =>
  request<{ authenticated: boolean; configured: boolean }>("/auth/check");

export const login = (password: string) =>
  request<{ status: string }>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ password }),
  });

export const logout = () =>
  request<{ status: string }>("/auth/logout", { method: "POST" });

export const setup = (password: string) =>
  request<{ status: string }>("/auth/setup", {
    method: "POST",
    body: JSON.stringify({ password }),
  });

// Parcels
export interface Parcel {
  id: string;
  tracking_number: string;
  carrier: string;
  name: string;
  notes: string;
  status: string;
  archived: boolean;
  last_check: string | null;
  created_at: string;
  updated_at: string;
}

export interface TrackingEvent {
  id: string;
  parcel_id: string;
  status: string;
  message: string;
  location: string;
  timestamp: string;
  created_at: string;
}

export interface CarrierInfo {
  code: string;
  name: string;
}

export const listParcels = (params?: Record<string, string>) => {
  const query = params
    ? "?" + new URLSearchParams(params).toString()
    : "";
  return request<Parcel[]>(`/parcels${query}`);
};

export const getParcel = (id: string) => request<Parcel>(`/parcels/${id}`);

export const createParcel = (data: Partial<Parcel>) =>
  request<Parcel>("/parcels", {
    method: "POST",
    body: JSON.stringify(data),
  });

export const updateParcel = (id: string, data: Partial<Parcel>) =>
  request<Parcel>(`/parcels/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });

export const deleteParcel = (id: string) =>
  request<void>(`/parcels/${id}`, { method: "DELETE" });

export const refreshParcel = (id: string) =>
  request<Parcel>(`/parcels/${id}/refresh`, { method: "POST" });

// Events
export const listEvents = (parcelId: string) =>
  request<TrackingEvent[]>(`/parcels/${parcelId}/events`);

export const createEvent = (
  parcelId: string,
  data: Partial<TrackingEvent>
) =>
  request<TrackingEvent>(`/parcels/${parcelId}/events`, {
    method: "POST",
    body: JSON.stringify(data),
  });

export const deleteEvent = (parcelId: string, eventId: string) =>
  request<void>(`/parcels/${parcelId}/events/${eventId}`, {
    method: "DELETE",
  });

// Health
export const getHealth = () =>
  request<{ status: string; carriers: CarrierInfo[] }>("/health");
