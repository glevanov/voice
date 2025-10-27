<script lang="ts">
  import { connectionStatus } from "../../store/connection-status";
  import { i18n } from "../../service/i18n/i18n";

  let status = $derived(i18n(`connection.${$connectionStatus}`));
</script>

<div class="status">
  <span
    class="dot pulse"
    class:danger={["disconnected", "error"].includes($connectionStatus)}
    class:warning={["connecting", "reconnecting"].includes($connectionStatus)}
    class:success={$connectionStatus === "connected"}
    class:pulse={["connecting", "reconnecting"].includes($connectionStatus)}
  ></span>
  {status}
</div>

<style>
  .status {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .dot {
    display: inline-block;
    width: 16px;
    height: 16px;
    border-radius: 50%;

    background-color: transparent;
  }

  .danger {
    background-color: var(--status-danger);
  }

  .warning {
    background-color: var(--status-warning);
  }

  .success {
    background-color: var(--status-success);
  }

  .pulse {
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 1;
    }
    50% {
      opacity: 0.4;
    }
    100% {
      opacity: 1;
    }
  }
</style>
