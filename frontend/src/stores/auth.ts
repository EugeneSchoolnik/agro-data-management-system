import { writable } from "svelte/store";

export interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  userEmail: string | null;
  error: string | null;
}

const initialState: AuthState = {
  isAuthenticated: false,
  token: localStorage.getItem("jwt_token"),
  userEmail: localStorage.getItem("user_email"),
  error: null,
};

// Check if token is still valid on initialization
if (initialState.token) {
  initialState.isAuthenticated = true;
}

export const authStore = writable<AuthState>(initialState);

export function setAuth(token: string, email: string) {
  localStorage.setItem("jwt_token", token);
  localStorage.setItem("user_email", email);
  authStore.set({
    isAuthenticated: true,
    token,
    userEmail: email,
    error: null,
  });
}

export function clearAuth() {
  localStorage.removeItem("jwt_token");
  localStorage.removeItem("user_email");
  authStore.set({
    isAuthenticated: false,
    token: null,
    userEmail: null,
    error: null,
  });
}

export function setAuthError(error: string) {
  authStore.update((state) => ({
    ...state,
    error,
  }));
}
