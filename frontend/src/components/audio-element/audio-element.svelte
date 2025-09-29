<script lang="ts">
  export let audioElement: HTMLAudioElement | null = null;
  export let audio: string | null = null;
</script>

<audio class="visually-hidden" bind:this={audioElement} src={audio} />

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
