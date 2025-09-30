<script lang="ts">
  import { clearMessages } from "../../store/messages";
  import StatusIndicator from "../status-indicator/status-indicator.svelte";
</script>

<div class="chat-header">
  <StatusIndicator />
  <button on:click={clearMessages} class="clear">Clear History</button>
</div>

<style>
  .chat-header {
    display: flex;
    justify-content: space-between;
  }

  .clear {
    padding: 5px 10px;
    background-color: #dc3545;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9em;
  }

  .clear:hover {
    background-color: #c82333;
  }
</style>
