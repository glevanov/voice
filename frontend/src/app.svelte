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
  import InputSection from "./components/input-section/input-section.svelte";

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

  onDestroy(() => {
    websocketStore.disconnect();
    cleanup();
  });
</script>

<div class="app">
  <div class="status">Status: {$connectionStatus}</div>

  <ChatHistory />

  <InputSection />

  {#if response}
    <h2>Latest Response:</h2>
    <div class="response">{response}</div>
  {/if}

  <audio bind:this={audioElement} controls src={audio} />
</div>

<style>
  .status {
    margin: 10px 0;
    padding: 10px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 10px;
  }
</style>
