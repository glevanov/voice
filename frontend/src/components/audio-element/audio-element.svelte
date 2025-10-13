<script lang="ts">
  import { onDestroy } from "svelte";

  import { websocketStore } from "../../store/websocket";
  import { addUserMessage, addAssistantMessage } from "../../store/messages.js";
  import { play, cleanup } from "../../service/audio.js";

  let audio: string | null = $state(null);
  let audioElement: HTMLAudioElement | null = $state(null);

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
    cleanup();
  });
</script>

<audio class="visually-hidden" bind:this={audioElement} src={audio}></audio>

<style>
  .visually-hidden {
    border: 0;
    clip: rect(0 0 0 0);
    clip-path: inset(100%);
    font-size: 1px;
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }
</style>
