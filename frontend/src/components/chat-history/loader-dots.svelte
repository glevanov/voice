<script>
</script>

<span class="first">●</span>
<span class="second">●</span>
<span class="third">●</span>

<style>
  .first {
    animation: blink 2s infinite;
  }

  .second {
    animation: blink 2s infinite 0.2s;
  }

  .third {
    animation: blink 2s infinite 0.4s;
  }

  @keyframes blink {
    0%,
    100% {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }
</style>
