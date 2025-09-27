import App from "./app.svelte";

const target = document.getElementById("app");
if (!target) {
  throw new Error("No element with id 'app' found");
}
const app = new App({
  target,
});

export default app;
