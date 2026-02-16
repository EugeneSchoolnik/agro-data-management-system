<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import Button from "./Button.svelte";

  type CloseEvent = CustomEvent<void>;

  export let open: boolean = false;
  export let title: string = "";

  const dispatch = createEventDispatcher<{ close: void; open: boolean }>();

  function close(): void {
    dispatch("close");
    dispatch("open", false);
  }

  function handleBackdropClick(event: MouseEvent): void {
    if (event.target === event.currentTarget) {
      close();
    }
  }

  function handleKeydown(e: KeyboardEvent): void {
    if (e.key === "Escape") {
      close();
    }
  }
</script>

{#if open}
  <div
    class="modal-backdrop"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
    role="button"
    tabindex="0"
  >
    <div
      class="modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <div class="modal-header">
        <h2 id="modal-title">{title}</h2>
        <button class="close-btn" on:click={close} aria-label="Закрити"
          >×</button
        >
      </div>
      <div class="modal-body">
        <slot />
      </div>
      {#if $$slots.footer}
        <div class="modal-footer">
          <slot name="footer" />
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(20, 35, 15, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    padding: 1.5rem;
  }

  .modal {
    background: #fbf5e8;
    border-radius: 1rem;
    width: min(92vw, 560px);
    max-height: 90vh;
    overflow-y: auto;
    border: 1px solid rgba(125, 115, 82, 0.16);
    box-shadow: 0 34px 70px rgba(31, 48, 16, 0.2);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.2rem 1.3rem;
    border-bottom: 1px solid rgba(133, 121, 87, 0.18);
  }

  .modal-header h2 {
    margin: 0;
    font-size: 1.25rem;
    color: #32481f;
  }

  .close-btn {
    background: #e7dcc0;
    border: 1px solid rgba(132, 114, 77, 0.2);
    border-radius: 50%;
    font-size: 1.35rem;
    cursor: pointer;
    width: 2.2rem;
    height: 2.2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #5e512f;
  }

  .modal-body {
    padding: 1.3rem;
  }

  .modal-footer {
    padding: 1rem 1.3rem 1.3rem;
    border-top: 1px solid rgba(133, 121, 87, 0.18);
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }
</style>
