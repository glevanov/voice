import { writable } from "svelte/store";

export const messages = writable([]);

export function addUserMessage(content) {
  messages.update((msgs) => [
    ...msgs,
    {
      role: "user",
      content: content,
      timestamp: new Date().toISOString(),
    },
  ]);
}

export function addAssistantMessage(content) {
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
