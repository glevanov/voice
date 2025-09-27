<script lang="ts">
  import { messages, clearMessages } from "../../store/messages.js";
</script>

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
</style>
