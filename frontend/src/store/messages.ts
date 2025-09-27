import { writable } from "svelte/store";

export type Message = {
  role: "user" | "assistant";
  content: string;
  timestamp: string;
};

export const messages = writable<Message[]>([]);

export function addUserMessage(content: string) {
  messages.update((msgs) => [
    ...msgs,
    {
      role: "user",
      content: content,
      timestamp: new Date().toISOString(),
    },
  ]);
}

export function addAssistantMessage(content: string) {
  messages.update((msgs) => [
    ...msgs,
    {
      role: "assistant",
      content: content,
      timestamp: new Date().toISOString(),
    },
  ]);
}

export function clearMessages() {
  messages.set([]);
}
