export type Handler = (payload: string) => void;

type MessageBase<Name extends string, Payload> = {
  name: Name;
  payload: Payload;
};

type UserMessage = MessageBase<"user-message", string>;
type AssistantMessage = MessageBase<"assistant-message", string>;

export type ServerMessage = UserMessage | AssistantMessage;

export const isValidMessage = (data: unknown): data is ServerMessage => {
  if (typeof data !== "object" || data === null) return false;
  if (!("name" in data) || !("payload" in data)) return false;
  if (data.name === "assistant-message" && typeof data.payload === "string") {
    return true;
  }
  if (data.name === "user-message" && typeof data.payload === "string") {
    return true;
  }
  return false;
};
