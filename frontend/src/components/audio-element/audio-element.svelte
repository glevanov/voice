<script lang="ts">
  import { onDestroy } from "svelte";

  import { play, cleanup } from "../../service/audio";
  import { addHandler } from "../../service/websocket";

  let audio: string | null = $state(null);
  let audioElement: HTMLAudioElement | null = $state(null);

  addHandler("assistant-message", (_: string) => {
    play(audioElement, "answer.wav")
      .then((audioUrl) => {
        audio = audioUrl;
      })
      .catch((error) => {
        console.error("Error playing audio:", error);
      });
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
