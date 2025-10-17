<script lang="ts">
  import { messages, isGeneratingMessage } from "../../store/messages";
  import LoaderDots from "./loader-dots.svelte";
</script>

<div class="chat-history">
  {#each $messages as message}
    <div class="message content {message.role}">
      {message.content}
    </div>
  {/each}
  {#if $isGeneratingMessage === true}
    <div class="message content assistant">
      <LoaderDots></LoaderDots>
    </div>
  {/if}
  {#if $messages.length === 0 && $isGeneratingMessage === false}
    <div class="no-messages">Vad vill du prata om idag?</div>
  {/if}
</div>

<style>
  :root {
    --user-card-color: var(--purple-300);
    --user-text-color: var(--neutral-light);
    --assistant-card-color: var(--neutral-mid);
    --assistant-text-color: var(--neutral-darkest);

    @media (prefers-color-scheme: dark) {
      --user-card-user-color: var(--purple-200);
      --assistant-card-color: var(--neutral-darkest);
      --assistant-text-color: var(--neutral-light);
    }
  }

  .chat-history {
    display: flex;
    flex-direction: column;
    gap: 20px;
    margin: -10px;
    padding: 10px;
    overflow-y: auto;
    height: calc(100svh - 450px);
  }

  .message {
    max-width: fit-content;
    padding: 12px 18px;

    border-radius: 20px;
  }

  .user {
    align-self: end;
    margin-left: 60px;

    background-color: var(--user-card-color);
    border-bottom-right-radius: 2px;

    color: var(--neutral-light);
  }

  .assistant {
    align-self: start;
    margin-right: 60px;

    background-color: var(--assistant-card-color);
    border-bottom-left-radius: 2px;

    color: var(--default-text-color);
  }

  .content {
    white-space: pre-wrap;
  }

  .no-messages {
    display: flex;
    flex-grow: 1;
    justify-content: center;
    align-items: center;
  }
</style>
