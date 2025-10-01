<script lang="ts">
  import { get } from "svelte/store";

  import { messages, addUserMessage } from "../../store/messages";
  import { websocketStore, connectionStatus } from "../../store/websocket";
  import VoiceRecorder from "../voice-recorder/voice-recorder.svelte";
  import SubmitIcon from "./submit-icon.svelte";

  let value: string = "";

  $: disabled = $connectionStatus !== "Connected";

  function handleInput(event: Event) {
    const target = event.target as HTMLTextAreaElement;
    value = target.value;
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      if (value.trim()) {
        handleMessageSend(value);
      }
    }
  }

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    if (value.trim()) {
      handleMessageSend(value);
    }
  }

  function handleMessageSend(messageText: string): void {
    addUserMessage(messageText);

    const currentMessages = get(messages);
    const payload = { messages: currentMessages };

    const sent = websocketStore.sendMessages(payload);
    if (!sent) {
      alert("WebSocket‑anslutningen är inte öppen.");
    }
    value = "";
  }
</script>

<form class="input-section" on:submit={handleSubmit}>
  <textarea
    bind:value
    placeholder="Skriv ditt meddelande här."
    {disabled}
    rows={3}
    on:input={handleInput}
    on:keydown={handleKeydown}
    class="textarea"
  ></textarea>

  <div class="controls">
    <VoiceRecorder />

    <button
      type="submit"
      disabled={disabled || !value.trim()}
      class="round-outline-button"
    >
      <SubmitIcon />
    </button>
  </div>
</form>

<style>
  .input-section {
    --border-color: var(--neutral-mid);

    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 16px;

    border-radius: 12px;
    border: 1px solid var(--border-color);
  }

  .textarea {
    height: 60px;

    border: none;
    background-color: transparent;

    font-family: inherit;
    font-size: inherit;
    line-height: inherit;
    color: inherit;

    outline: none;
    resize: none;
  }

  .controls {
    display: flex;
    gap: 8px;
    align-self: end;
  }

  :root {
    --color: var(--purple-300);
    @media (prefers-color-scheme: dark) {
      --color: var(--purple-200);
    }
  }

  .round-outline-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 50px;
    height: 50px;
    box-sizing: border-box;

    background-color: transparent;
    border: 1px solid var(--color);
    border-radius: 50%;

    color: var(--color);

    transition:
      border-color 0.2s ease-in,
      color 0.2s ease-in;
    cursor: pointer;
  }

  .round-outline-button:hover,
  .round-outline-button:focus-visible {
    --color: var(--purple-400);
    @media (prefers-color-scheme: dark) {
      --color: var(--purple-100);
    }
  }
</style>
