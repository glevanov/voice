<script lang="ts">
  import type { HTMLButtonAttributes } from "svelte/elements";

  type Fill = "filled" | "outlined";
  type Color = "primary" | "danger";

  interface Props extends HTMLButtonAttributes {
    fill?: Fill;
    isRound?: boolean;
    color?: Color;
  }

  let {
    fill = "filled",
    isRound = false,
    color = "primary",
    children,
    ...rest
  }: Props = $props();

  let merged = $derived(`base ${rest.class ?? ""}`);
</script>

<button
  class={merged}
  class:fill-filled={fill === "filled"}
  class:fill-outlined={fill === "outlined"}
  class:shape-rectangular={isRound === false}
  class:shape-circular={isRound === true}
  class:color-primary={color === "primary"}
  class:color-danger={color === "danger"}
  {...rest}
>
  {@render children?.()}
</button>

<style>
  :root {
    --primary: var(--purple-300);
    --primary-hover: var(--purple-400);
    --danger: var(--red-300);
    --danger-hover: var(--red-400);
    --disabled: var(--neutral-mid);

    @media (prefers-color-scheme: dark) {
      --primary: var(--purple-200);
      --primary-hover: var(--purple-100);
      --danger: var(--red-200);
      --danger-hover: var(--red-100);
    }
  }

  .base {
    font-size: inherit;
    line-height: inherit;

    cursor: pointer;
    outline: 3px solid transparent;
    outline-offset: 2px;
    transition: outline-color 0.2s ease;
    &:focus-visible:not(:disabled) {
      outline-color: var(--outline-color);
    }
  }

  :disabled {
    cursor: not-allowed;
  }

  .shape-rectangular {
    padding: 5px 10px;

    border-radius: 4px;
  }

  .shape-circular {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 50px;
    height: 50px;
    box-sizing: border-box;

    border-radius: 50%;
  }

  .fill-filled {
    border: none;

    color: var(--neutral-light);

    transition: background-color 0.2s ease;

    &.color-primary:not(:disabled) {
      background-color: var(--primary);

      &:hover {
        background-color: var(--primary-hover);
      }
    }

    &.color-danger:not(:disabled) {
      background-color: var(--danger);

      &:hover {
        background-color: var(--danger-hover);
      }
    }

    &:disabled {
      background-color: var(--disabled);
    }
  }

  .fill-outlined {
    background-color: transparent;
    border-style: solid;
    border-width: 1px;

    transition:
      border-color 0.2s ease,
      color 0.2s ease;

    &.color-primary:not(:disabled) {
      color: var(--primary);
      border-color: var(--primary);

      &:hover {
        color: var(--primary);
        border-color: var(--primary);
      }
    }

    &.color-danger:not(:disabled) {
      color: var(--danger);
      border-color: var(--danger);

      &:hover {
        color: var(--danger-hover);
        border-color: var(--danger-hover);
      }
    }

    &:disabled {
      color: var(--disabled);
      border-color: var(--disabled);
    }
  }
</style>
