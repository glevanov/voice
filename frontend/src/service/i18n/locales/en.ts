import type { Locale } from "../i18n";

export const en: Locale = {
  clearChat: "Clear chat",
  chatPlaceholder: "What do you want to talk about today?",
  connection: {
    connecting: "Connecting",
    connected: "Connected",
    disconnected: "Disconnected",
    reconnecting: "Reconnecting",
    error: "Error",
  },
  recording: {
    start: "Start recording",
    stop: "Stop recording",
    microphonePermissionError:
      "Could not access microphone. Make sure you have granted the necessary permissions.",
  },
  input: {
    connectionError: "WebSocket connection is not open.",
    placeholder: "Type your message here",
    placeholderDisconnected: "Disconnected",
    submit: "Submit",
    textAreaLabel: "Write your message here",
  },
};
