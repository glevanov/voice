<script lang="ts">
  import { onDestroy } from "svelte";
  import "./app.css";
  import { get } from "svelte/store";
  import {
    messages,
    addUserMessage,
    addAssistantMessage,
  } from "./store/messages.js";
  import { websocketStore, connectionStatus } from "./store/websocket";
  import { play, cleanup } from "./service/audio.js";
  import ChatHistory from "./components/chat-history/chat-history.svelte";

  let text = "";
  let response = "";
  let audio: string | null = null;
  let audioElement: HTMLAudioElement | null = null;

  websocketStore.connect();

  websocketStore.onMessage((event: MessageEvent): void => {
    try {
      const data = JSON.parse(event.data);
      response = data.text;
      addAssistantMessage(data.text);

      if (data.audio) {
        play(audioElement, data.audio)
          .then((audioUrl) => {
            audio = audioUrl;
          })
          .catch((error) => {
            console.error("Error playing audio:", error);
          });
      }
    } catch (error) {
      console.error("Error parsing response:", error);
      response = event.data;
      addAssistantMessage(event.data);
    }
  });

  function handleMessageSend(): void {
    const userMessage = text;
    addUserMessage(userMessage);

    const currentMessages = get(messages);
    const payload = { messages: currentMessages };

    const sent = websocketStore.send(payload);
    if (!sent) {
      alert("WebSocket connection is not open");
    }
    text = "";
  }

  onDestroy(() => {
    websocketStore.disconnect();
    cleanup();
  });
</script>

<div class="app">
  <div class="status">Status: {$connectionStatus}</div>

  <ChatHistory />

  <div class="input-section">
    <textarea bind:value={text} placeholder="Type your message here..."
    ></textarea>
    <button
      on:click={handleMessageSend}
      disabled={$connectionStatus !== "Connected"}>Send</button
    >
  </div>

  {#if response}
    <h2>Latest Response:</h2>
    <div class="response">{response}</div>
  {/if}

  <audio bind:this={audioElement} controls src={audio} />
</div>

<style>
  .input-section {
    margin: 20px 0;
  }

  .input-section textarea {
    width: 100%;
    margin-bottom: 10px;
  }

  .status {
    margin: 10px 0;
    padding: 10px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 10px;
  }
</style>
