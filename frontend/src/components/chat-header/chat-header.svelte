<script lang="ts">
  import StatusIndicator from "../status-indicator/status-indicator.svelte";
  import Button from "../button/button.svelte";
  import { clearMessages } from "../../store/messages.js";
</script>

<div class="chat-header">
  <StatusIndicator />

  <Button on:click={clearMessages} fill="outlined" color="danger"
    >Rensa chatt</Button
  >
</div>

<style>
  .chat-header {
    display: flex;
    justify-content: space-between;
  }
</style>
