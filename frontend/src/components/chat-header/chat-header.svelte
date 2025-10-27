<script lang="ts">
  import StatusIndicator from "../status-indicator/status-indicator.svelte";
  import Button from "../button/button.svelte";
  import { clearMessages } from "../../store/messages.js";
  import { connectionStatus } from "../../store/connection-status";
  import { i18n } from "../../service/i18n/i18n";

  let disabled = $derived($connectionStatus !== "connected");
</script>

<div class="chat-header">
  <StatusIndicator />

  <Button onclick={clearMessages} fill="outlined" color="danger" {disabled}>
    {i18n("clearChat")}
  </Button>
</div>

<style>
  .chat-header {
    display: flex;
    justify-content: space-between;
  }
</style>
