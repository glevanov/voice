<script lang="ts">
  import { clearMessages } from "../../store/messages.js";
</script>

<button on:click={clearMessages} class="button danger-outline">
  Clear chat
</button>

<style>
  .button {
    padding: 5px 10px;

    border-radius: 4px;

    font-size: inherit;
    line-height: inherit;

    cursor: pointer;
  }

  .danger-outline {
    --color: var(--red-300);
    @media (prefers-color-scheme: dark) {
      --color: var(--red-200);
    }

    background-color: transparent;
    border: 1px solid var(--color);

    color: var(--color);

    transition:
      border-color 0.2s ease-in,
      color 0.2s ease-in;
  }

  .danger-outline:hover,
  .danger-outline:focus {
    --color: var(--red-400);
    @media (prefers-color-scheme: dark) {
      --color: var(--red-100);
    }
  }
</style>
