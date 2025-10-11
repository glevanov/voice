<script lang="ts">
  import { get } from "svelte/store";

  import { messages, addUserMessage } from "../../store/messages";
  import { websocketStore, connectionStatus } from "../../store/websocket";
  import Button from "../button/button.svelte";
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
    aria-label="Skriv ditt meddelande här."
    placeholder={disabled ? "Ej ansluten" : "Skriv ditt meddelande här."}
    {disabled}
    rows={3}
    on:input={handleInput}
    on:keydown={handleKeydown}
    class="textarea"
  ></textarea>

  <div class="controls">
    <VoiceRecorder {disabled} />

    <Button
      type="submit"
      isRound
      fill="outlined"
      disabled={disabled || !value.trim()}><SubmitIcon /></Button
    >
  </div>
</form>

<style>
  :root {
    --primary: var(--purple-300);

    @media (prefers-color-scheme: dark) {
      --primary: var(--purple-200);
    }
  }

  .input-section {
    --border-color: var(--neutral-mid);

    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 16px;

    border-radius: 12px;
    border: 1px solid var(--border-color);

    outline: 3px solid transparent;
    outline-offset: 2px;
    transition: outline-color 0.2s ease;

    &:has(.textarea:focus) {
      border-color: var(--primary);
    }
    &:has(.textarea:focus-visible) {
      outline-color: var(--outline-color);
    }
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
</style>
