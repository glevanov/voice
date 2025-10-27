import { connectionStatus } from "../store/connection-status";
import type { Message } from "../store/messages";
import {
  type Handler,
  isValidMessage,
  type ServerMessage,
} from "./websocket.types";

let ws: WebSocket | null = null;
let reconnectTimeout: number | null = null;
let reconnectAttempts = 0;
const maxReconnectDelay = 30000;
const baseReconnectDelay = 1000;

const listeners = new Map<string, Set<Handler>>();

const handleMessage = ({ name, payload }: ServerMessage) => {
  const callbacks = listeners.get(name);
  if (callbacks) {
    callbacks.forEach((callback) => callback(payload));
  }
};

export const addHandler = (name: string, handler: Handler) => {
  if (!listeners.has(name)) {
    listeners.set(name, new Set());
  }
  listeners.get(name)!.add(handler);

  const unsubscribe = () => {
    listeners.get(name)?.delete(handler);
  };

  return unsubscribe;
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
  connectionStatus.set("reconnecting");

  console.log(
    `Attempting to reconnect in ${delay}ms (attempt ${reconnectAttempts})`,
  );

  reconnectTimeout = setTimeout(() => {
    connect();
  }, delay);
};

export const sendMessages = (data: { messages: Message[] }) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(data));
    return true;
  }
  return false;
};

export const sendVoice = async (data: Blob, messages: Message[]) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    try {
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

export const connect = () => {
  try {
    ws = new WebSocket("ws://localhost:3002/ws");

    ws.onopen = () => {
      connectionStatus.set("connected");
      reconnectAttempts = 0;
    };

    ws.onclose = (event) => {
      connectionStatus.set("disconnected");
      if (event.code !== 1000) {
        handleReconnect();
      }
    };

    ws.onerror = (error) => {
      connectionStatus.set("error");
      console.error("WebSocket error:", error);
    };

    ws.onmessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (!isValidMessage(data)) {
          throw new Error("Invalid message format");
        }
        handleMessage(data);
      } catch (e) {
        console.error("Failed to process message", e);
      }
    };
  } catch (error) {
    console.error("Failed to create WebSocket:", error);
    handleReconnect();
  }
};

export const disconnect = () => {
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
  }
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.close(1000, "Disconnected by user");
    ws = null;
  }
};
