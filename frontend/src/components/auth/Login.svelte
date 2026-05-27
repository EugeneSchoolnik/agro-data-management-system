<script lang="ts">
  import { login } from "../../lib/api";
  import { setAuth, setAuthError, authStore } from "../../stores/auth";
  import type { LoginRequest } from "../../types/models";

  let email = "";
  let password = "";
  let isLoading = false;
  let error = "";

  $: error = $authStore.error || "";

  async function handleLogin() {
    if (!email || !password) {
      error = "Please enter both email and password";
      return;
    }

    isLoading = true;
    error = "";

    try {
      const credentials: LoginRequest = { email, password };
      const response = await login(credentials);

      if (response.token) {
        setAuth(response.token, email);
        // The component will be replaced by Dashboard after auth state updates
      }
    } catch (err) {
      error =
        err instanceof Error ? err.message : "Login failed. Please try again.";
      setAuthError(error);
    } finally {
      isLoading = false;
    }
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      handleLogin();
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <h1>Agro Data Management System</h1>
    <p class="subtitle">Sign in to your account</p>

    <form on:submit|preventDefault={handleLogin}>
      <div class="form-group">
        <label for="email">Email</label>
        <input
          id="email"
          type="email"
          bind:value={email}
          placeholder="Enter your email"
          disabled={isLoading}
          on:keydown={handleKeyDown}
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          placeholder="Enter your password"
          disabled={isLoading}
          on:keydown={handleKeyDown}
        />
      </div>

      {#if error}
        <div class="error-message">
          {error}
        </div>
      {/if}

      <button
        type="submit"
        disabled={isLoading || !email || !password}
        class="login-button"
      >
        {isLoading ? "Signing in..." : "Sign In"}
      </button>
    </form>
  </div>
</div>

<style>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    font-family:
      system-ui,
      -apple-system,
      sans-serif;
  }

  .login-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    padding: 40px;
    width: 100%;
    max-width: 400px;
  }

  h1 {
    margin: 0 0 8px 0;
    color: #333;
    font-size: 24px;
    text-align: center;
  }

  .subtitle {
    margin: 0 0 30px 0;
    color: #666;
    text-align: center;
    font-size: 14px;
  }

  form {
    display: flex;
    flex-direction: column;
  }

  .form-group {
    margin-bottom: 20px;
    display: flex;
    flex-direction: column;
  }

  label {
    margin-bottom: 8px;
    font-weight: 500;
    color: #333;
    font-size: 14px;
  }

  input {
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
    transition: border-color 0.2s;
  }

  input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  input:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
  }

  .error-message {
    background-color: #fee;
    border: 1px solid #fcc;
    border-radius: 4px;
    color: #c33;
    padding: 10px 12px;
    margin-bottom: 20px;
    font-size: 14px;
  }

  .login-button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 12px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .login-button:hover:not(:disabled) {
    opacity: 0.9;
  }

  .login-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
