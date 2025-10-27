import { writable } from "svelte/store";

export type Status =
  | "connecting"
  | "connected"
  | "disconnected"
  | "reconnecting"
  | "error";
export const connectionStatus = writable<Status>("connecting");
