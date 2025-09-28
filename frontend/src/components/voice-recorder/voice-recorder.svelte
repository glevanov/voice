<script lang="ts">
  import { onDestroy } from "svelte";
  import { websocketStore } from "../../store/websocket";

  let text: "Record" | "Stop" = "Record";
  let mediaRecorder: MediaRecorder | null = null;
  let chunks: Blob[] = [];

  async function startRecording() {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    mediaRecorder = new MediaRecorder(stream);

    mediaRecorder.ondataavailable = (e) => {
      if (e.data.size > 0) chunks.push(e.data);
    };

    mediaRecorder.onstop = async () => {
      const blob = new Blob(chunks, { type: "audio/webm" });
      chunks = [];

      websocketStore.sendVoice(blob);
      mediaRecorder = null;
    };
    mediaRecorder.start();
  }

  function stopRecording() {
    mediaRecorder?.stop();
  }

  async function handleClick() {
    const isRecording = text === "Stop";
    if (isRecording) {
      stopRecording();
      text = "Record";
    } else {
      await startRecording();
      text = "Stop";
    }
  }

  onDestroy(() => {
    if (mediaRecorder) {
      mediaRecorder = null;
    }
  });
</script>

<button on:click={handleClick}>{text}</button>
