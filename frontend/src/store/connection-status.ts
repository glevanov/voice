import { writable } from "svelte/store";

export type Status =
  | "Connecting"
  | "Connected"
  | "Disconnected"
  | "Reconnecting"
  | "Error";
export const connectionStatus = writable<Status>("Connecting");
