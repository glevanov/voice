<script lang="ts">
  import StatusIndicator from "../status-indicator/status-indicator.svelte";
  import ClearChat from "../clear-chat/clear-chat.svelte";
</script>

<div class="chat-header">
  <StatusIndicator />
  <ClearChat />
</div>

<style>
  .chat-header {
    display: flex;
    justify-content: space-between;
  }
</style>
