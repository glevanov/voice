import App from "./app.svelte";
import { mount } from "svelte";

const target = document.getElementById("app");
if (!target) {
  throw new Error("No element with id 'app' found");
}
const app = mount(App, {
  target,
});

export default app;
