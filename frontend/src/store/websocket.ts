import { writable, derived } from "svelte/store";
import type { Message } from "./messages";

export const connectionStatus = writable("Connecting...");
export const isConnected = derived(
  connectionStatus,
  ($status) => $status === "Connected",
);

let ws: WebSocket | null = null;
let reconnectTimeout: number | null = null;
let reconnectAttempts = 0;
const maxReconnectDelay = 30000;
const baseReconnectDelay = 1000;

export function createWebSocketStore() {
  const { subscribe, set } = writable<WebSocket | null>(null);

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
    if (!ws) {
      console.warn("WebSocket unavailable");
      return;
    }
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

  const sendMessages = (data: { messages: Message[] }) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data));
      return true;
    }
    return false;
  };

  const sendVoice = async (data: Blob, messages: Message[]) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      try {
        // Convert blob to base64
        const reader = new FileReader();
        const base64Promise = new Promise<string>((resolve, reject) => {
          reader.onload = () => {
            const arrayBuffer = reader.result as ArrayBuffer;
            const bytes = new Uint8Array(arrayBuffer);
            let binary = "";
            for (let i = 0; i < bytes.byteLength; i++) {
              binary += String.fromCharCode(bytes[i]);
            }
            resolve(btoa(binary));
          };
          reader.onerror = reject;
        });
        reader.readAsArrayBuffer(data);

        const base64Audio = await base64Promise;

        const audioMessage = {
          type: "audio",
          audioData: base64Audio,
          messages: messages.map((msg) => ({
            role: msg.role,
            content: msg.content,
          })),
        };

        ws.send(JSON.stringify(audioMessage));
        return true;
      } catch (error) {
        console.error("Error encoding audio:", error);
        return false;
      }
    }
    return false;
  };

  const onMessage = (callback: (event: MessageEvent) => void) => {
    if (ws) {
      ws.onmessage = callback;
    }
  };

  return {
    subscribe,
    connect,
    disconnect,
    sendMessages,
    sendVoice,
    onMessage,
  };
}

export const websocketStore = createWebSocketStore();
