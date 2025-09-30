<script lang="ts">
  import { onDestroy } from "svelte";

  import "./app.css";
  import { addUserMessage, addAssistantMessage } from "./store/messages.js";
  import { websocketStore } from "./store/websocket";
  import { play, cleanup } from "./service/audio.js";
  import ChatHistory from "./components/chat-history/chat-history.svelte";
  import InputSection from "./components/input-section/input-section.svelte";
  import AudioElement from "./components/audio-element/audio-element.svelte";
  import ChatHeader from "./components/chat-header/chat-header.svelte";

  let audio: string | null = null;
  let audioElement: HTMLAudioElement | null = null;

  websocketStore.connect();

  websocketStore.onMessage((event: MessageEvent): void => {
    try {
      const data = JSON.parse(event.data);

      if (data.type === "user") {
        addUserMessage(data.text);
      } else if (data.type === "assistant") {
        addAssistantMessage(data.text);

        play(audioElement, "answer.wav")
          .then((audioUrl) => {
            audio = audioUrl;
          })
          .catch((error) => {
            console.error("Error playing audio:", error);
          });
      }
    } catch (error) {
      console.error("Error parsing response:", error);
      addAssistantMessage(event.data);
    }
  });

  onDestroy(() => {
    websocketStore.disconnect();
    cleanup();
  });
</script>

<div class="app">
  <div class="chat">
    <ChatHeader />
    <ChatHistory />
    <InputSection />
    <AudioElement bind:audioElement {audio} />
  </div>
</div>

<style>
  .app {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100svh;
    width: 100vw;
    padding: 20px;
    box-sizing: border-box;
  }

  .chat {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: 20px;
    width: 100%;
    max-width: 1200px;
    padding: 32px;

    background: var(--card-bg-color);
    border-radius: 36px;
  }
</style>
