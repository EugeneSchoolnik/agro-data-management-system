<script lang="ts">
  type ButtonVariant = "primary" | "secondary" | "danger";
  type ButtonType = "button" | "submit" | "reset";

  export let type: ButtonType = "button";
  export let variant: ButtonVariant = "primary";
  export let disabled: boolean = false;
  export let loading: boolean = false;
</script>

<button
  {type}
  class="btn btn-{variant}"
  {disabled}
  class:disabled
  class:loading
  on:click
>
  {#if loading}
    <span class="spinner"></span>
  {/if}
  <slot />
</button>

<style>
  .btn {
    padding: 0.75rem 1.3rem;
    border: 1px solid transparent;
    border-radius: 999px;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 600;
    transition:
      transform 0.2s ease,
      filter 0.2s ease,
      box-shadow 0.2s ease;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 3rem;
    box-shadow: 0 12px 24px rgba(74, 91, 36, 0.12);
  }

  .btn:hover:not(.disabled) {
    transform: translateY(-1px);
  }

  .btn-primary {
    background: linear-gradient(135deg, #4b7331 0%, #7d9c4a 100%);
    color: #fff;
    border-color: #41682d;
  }

  .btn-primary:hover:not(.disabled) {
    filter: brightness(1.07);
  }

  .btn-secondary {
    background: linear-gradient(135deg, #f5e5c0 0%, #e2ce97 100%);
    color: #3f3a26;
    border-color: #c9b17e;
  }

  .btn-secondary:hover:not(.disabled) {
    filter: brightness(1.03);
  }

  .btn-danger {
    background: linear-gradient(135deg, #c23f25 0%, #e0704a 100%);
    color: #fff;
    border-color: #99331f;
  }

  .btn-danger:hover:not(.disabled) {
    filter: brightness(1.06);
  }

  .btn.disabled {
    opacity: 0.65;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }

  .btn.loading {
    position: relative;
  }

  .spinner {
    display: inline-block;
    width: 1rem;
    height: 1rem;
    border: 2px solid transparent;
    border-top: 2px solid currentColor;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-right: 0.5rem;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>
