<script lang="ts">
  import { get } from "svelte/store";
  import { messages, addUserMessage } from "../../store/messages";
  import { websocketStore, connectionStatus } from "../../store/websocket";

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

  function handleSubmit() {
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
      alert("WebSocket connection is not open");
    }
    value = "";
  }
</script>

<div class="input-section">
  <div class="input-container">
    <textarea
      bind:value
      placeholder="Type your message here..."
      {disabled}
      rows={3}
      on:input={handleInput}
      on:keydown={handleKeydown}
      class="textarea"
    ></textarea>
    <button
      on:click={handleSubmit}
      disabled={disabled || !value.trim()}
      class="submit-button"
    >
      Send
    </button>
  </div>
</div>

<style>
  .input-section {
    width: 100%;
    margin: 20px 0;
  }

  .input-container {
    display: flex;
    gap: 8px;
    flex-direction: column;
  }

  .textarea {
    flex: 1;
    min-height: 60px;
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-family: inherit;
    font-size: 14px;
    line-height: 1.4;
    resize: vertical;
    box-sizing: border-box;
    transition: border-color 0.2s ease;
  }

  .textarea:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }

  .textarea:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
    opacity: 0.6;
  }

  .textarea::placeholder {
    color: #999;
  }

  .submit-button {
    padding: 12px 20px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
    white-space: nowrap;
    height: fit-content;
    width: fit-content;
  }

  .submit-button:hover:not(:disabled) {
    background-color: #0056b3;
  }

  .submit-button:disabled {
    background-color: #6c757d;
    cursor: not-allowed;
    opacity: 0.6;
  }

  .submit-button:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }
</style>
