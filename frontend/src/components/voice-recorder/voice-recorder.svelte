<script lang="ts">
  import { onDestroy } from "svelte";

  import { websocketStore } from "../../store/websocket";
  import { messages } from "../../store/messages";
  import RecordIcon from "./record-icon.svelte";
  import StopIcon from "./stop-icon.svelte";

  let isRecording = false;
  let mediaRecorder: MediaRecorder | null = null;
  let chunks: Blob[] = [];

  async function startRecording() {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      mediaRecorder = new MediaRecorder(stream);

      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) chunks.push(e.data);
      };

      mediaRecorder.onstop = async () => {
        const blob = new Blob(chunks, { type: "audio/webm" });
        chunks = [];

        // Close audio tracks to release microphone
        for (const track of stream.getTracks()) {
          track.stop();
        }

        await websocketStore.sendVoice(blob, $messages);
        mediaRecorder = null;
      };

      mediaRecorder.start();
      isRecording = true;
    } catch (error) {
      console.error("Error accessing microphone:", error);
      alert(
        "Could not access microphone. Please ensure you've granted the necessary permissions.",
      );
    }
  }

  function stopRecording() {
    mediaRecorder?.stop();
    isRecording = false;
  }

  async function handleClick() {
    if (isRecording) {
      stopRecording();
    } else {
      await startRecording();
    }
  }

  onDestroy(() => {
    if (mediaRecorder) {
      stopRecording();
    }
  });
</script>

<button
  on:click={handleClick}
  class="round-button"
  aria-label={isRecording ? "Stop recording" : "Start recording"}
>
  {#if isRecording}
    <div class="pulse">
      <StopIcon />
    </div>
  {:else}
    <RecordIcon />
  {/if}
</button>

<style>
  :root {
    --bg-color: var(--purple-300);
    @media (prefers-color-scheme: dark) {
      --bg-color: var(--purple-200);
    }
  }

  .round-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 50px;
    height: 50px;
    box-sizing: border-box;

    background-color: var(--bg-color);
    border: none;
    border-radius: 50%;

    color: var(--neutral-light);

    transition: background-color 0.2s ease-in;
    cursor: pointer;
  }

  .round-button:hover,
  .round-button:focus-visible {
    --bg-color: var(--purple-400);
    @media (prefers-color-scheme: dark) {
      --bg-color: var(--purple-100);
    }
  }

  .pulse {
    width: 24px;
    height: 24px;

    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 1;
    }
    50% {
      opacity: 0.4;
    }
    100% {
      opacity: 1;
    }
  }
</style>
