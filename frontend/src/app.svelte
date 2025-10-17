<script lang="ts">
  import { onDestroy } from "svelte";

  import "./app.css";
  import { addHandler, connect, disconnect } from "./service/websocket";
  import ChatHistory from "./components/chat-history/chat-history.svelte";
  import InputSection from "./components/input-section/input-section.svelte";
  import AudioElement from "./components/audio-element/audio-element.svelte";
  import ChatHeader from "./components/chat-header/chat-header.svelte";
  import { addAssistantMessage, addUserMessage } from "./store/messages";

  connect();
  addHandler("user-message", (payload: string) => addUserMessage(payload));
  addHandler("assistant-message", (payload: string) =>
    addAssistantMessage(payload),
  );

  onDestroy(() => {
    disconnect();
  });
</script>

<div class="app">
  <div class="chat">
    <ChatHeader />
    <ChatHistory />
    <InputSection />
    <AudioElement />
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
