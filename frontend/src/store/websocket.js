import { writable, derived } from "svelte/store";

export const connectionStatus = writable("Connecting...");
export const isConnected = derived(
  connectionStatus,
  ($status) => $status === "Connected",
);

let ws = null;
let reconnectTimeout = null;
let reconnectAttempts = 0;
const maxReconnectDelay = 30000;
const baseReconnectDelay = 1000;

export function createWebSocketStore() {
  const { subscribe, set, update } = writable(null);

  const connect = () => {
    try {
      ws = new WebSocket("ws://localhost:3002/ws");
      setupWebSocketHandlers();
      set(ws);
    } catch (error) {
      console.error("Failed to create WebSocket:", error);
      handleReconnect();
    }
  };

  const setupWebSocketHandlers = () => {
    ws.onopen = () => {
      connectionStatus.set("Connected");
      reconnectAttempts = 0;
      console.log("WebSocket connected");
    };

    ws.onclose = (event) => {
      connectionStatus.set("Disconnected");
      console.log("WebSocket closed:", event.code, event.reason);
      if (event.code !== 1000) {
        handleReconnect();
      }
    };

    ws.onerror = (error) => {
      connectionStatus.set("Error");
      console.error("WebSocket error:", error);
    };
  };

  const handleReconnect = () => {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
    }

    reconnectAttempts++;
    const delay = Math.min(
      baseReconnectDelay * Math.pow(2, reconnectAttempts - 1),
      maxReconnectDelay,
    );
    connectionStatus.set("Reconnecting...");

    console.log(
      `Attempting to reconnect in ${delay}ms (attempt ${reconnectAttempts})`,
    );

    reconnectTimeout = setTimeout(() => {
      connect();
    }, delay);
  };

  const disconnect = () => {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
    }
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.close(1000, "Disconnected by user");
    }
  };

  const send = (data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data));
      return true;
    }
    return false;
  };

  const onMessage = (callback) => {
    if (ws) {
      ws.onmessage = callback;
    }
  };

  return {
    subscribe,
    connect,
    disconnect,
    send,
    onMessage,
  };
}

export const websocketStore = createWebSocketStore();
