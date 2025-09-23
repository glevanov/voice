<script>
  import { onDestroy } from "svelte";
  import "./app.css";
  import { get } from "svelte/store";
  import {
    messages,
    addUserMessage,
    addAssistantMessage,
    clearMessages,
  } from "./store/messages.js";
  import { websocketStore, connectionStatus } from "./store/websocket";
  import { play, cleanup } from "./service/audio.js";

  let text = "";
  let response = "";
  let audio = null;
  let audioElement = null;

  websocketStore.connect();

  websocketStore.onMessage((event) => {
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

  function handleMessageSend() {
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
  <h1>Voice Assistant</h1>
  <div class="status">Status: {$connectionStatus}</div>

  <div class="chat-history">
    <div class="chat-header">
      <h2>Chat History</h2>
      <button on:click={clearMessages} class="clear-btn">Clear History</button>
    </div>
    <div class="messages">
      {#each $messages as message}
        <div class="message {message.role}">
          <div class="role">
            {message.role === "user" ? "You" : "Assistant"}:
          </div>
          <div class="content">{message.content}</div>
        </div>
      {/each}
      {#if $messages.length === 0}
        <div class="no-messages">No messages yet...</div>
      {/if}
    </div>
  </div>

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
  .chat-history {
    margin: 20px 0;
    border: 1px solid #ccc;
    border-radius: 8px;
    max-height: 400px;
    overflow-y: auto;
  }

  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 15px;
    border-bottom: 1px solid #eee;
    background-color: #f8f9fa;
  }

  .chat-header h2 {
    margin: 0;
    font-size: 1.2em;
  }

  .clear-btn {
    padding: 5px 10px;
    background-color: #dc3545;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9em;
  }

  .clear-btn:hover {
    background-color: #c82333;
  }

  .messages {
    padding: 15px;
  }

  .message {
    margin-bottom: 15px;
    padding: 10px;
    border-radius: 8px;
  }

  .message.user {
    background-color: #e3f2fd;
    margin-left: 20px;
  }

  .message.assistant {
    background-color: #f1f8e9;
    margin-right: 20px;
  }

  .role {
    font-weight: bold;
    margin-bottom: 5px;
    color: #555;
  }

  .content {
    white-space: pre-wrap;
  }

  .no-messages {
    text-align: center;
    color: #666;
    font-style: italic;
    padding: 20px;
  }

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
