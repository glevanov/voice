import { writable } from "svelte/store";

export type Message = {
  role: "user" | "assistant";
  content: string;
  timestamp: string;
};

export const messages = writable<Message[]>([]);

export const isGeneratingMessage = writable(false);

export function addUserMessage(content: string) {
  isGeneratingMessage.set(true);

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
  isGeneratingMessage.set(false);

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
  isGeneratingMessage.set(false);
}
