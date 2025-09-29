<script lang="ts">
  import { onDestroy } from "svelte";
  import { websocketStore } from "../../store/websocket";
  import { messages } from "../../store/messages";

  let isRecording = false;
  let mediaRecorder: MediaRecorder | null = null;
  let chunks: Blob[] = [];
  let pulseAnimation = false;

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
        pulseAnimation = false;
      };

      mediaRecorder.start();
      isRecording = true;
      pulseAnimation = true;
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
  class:recording={isRecording}
  class:pulse={pulseAnimation}
>
  {#if isRecording}<span class="icon">⬤</span>{/if}
  <span class="text">{isRecording ? "Stop Recording" : "Start Recording"}</span>
</button>

<style>
  button {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 140px;
    justify-content: center;
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

  button:hover:not(:disabled) {
    background-color: #0056b3;
  }

  button:disabled {
    background-color: #6c757d;
    cursor: not-allowed;
    opacity: 0.6;
  }

  button:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }

  .icon {
    font-size: 12px;
    transition: color 0.2s ease;
  }

  .recording {
    background-color: #dc3545 !important;
  }

  .recording:hover {
    background-color: #bd2130 !important;
  }

  .pulse .icon {
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
