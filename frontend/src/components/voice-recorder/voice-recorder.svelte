<script lang="ts">
  import { onDestroy } from "svelte";

  import { websocketStore } from "../../store/websocket";
  import { messages } from "../../store/messages";
  import Button from "../button/button.svelte";
  import RecordIcon from "./record-icon.svelte";
  import StopIcon from "./stop-icon.svelte";

  interface Props {
    disabled?: boolean;
  }

  let { disabled = false }: Props = $props();

  let isRecording = $state(false);
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
        "Kunde inte komma åt mikrofonen. Se till att du har gett nödvändiga behörigheter.",
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

<Button
  isRound
  {disabled}
  onclick={handleClick}
  aria-label={isRecording ? "Sluta spela in." : "Starta inspelning."}
>
  {#if isRecording}
    <div class="pulse">
      <StopIcon />
    </div>
  {:else}
    <RecordIcon />
  {/if}
</Button>

<style>
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
