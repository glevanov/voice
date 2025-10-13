<script lang="ts">
  import { connectionStatus, type Status } from "../../store/websocket";

  const connectionStatusMap: Record<Status, string> = {
    ["Connecting"]: "Ansluter",
    ["Connected"]: "Ansluten",
    ["Disconnected"]: "Frånkopplad",
    ["Reconnecting"]: "Återansluter",
    ["Error"]: "Fel",
  } as const;
</script>

<div class="status">
  <span
    class="dot pulse"
    class:danger={["Disconnected", "Error"].includes($connectionStatus)}
    class:warning={["Connecting", "Reconnecting"].includes($connectionStatus)}
    class:success={$connectionStatus === "Connected"}
    class:pulse={["Connecting", "Reconnecting"].includes($connectionStatus)}
  ></span>
  {connectionStatusMap[$connectionStatus]}
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
