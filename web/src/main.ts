import "./app.css";
import App from "./App.svelte";
import { mount } from "svelte";
import { getLocale } from "./lib/i18n.svelte";
import { initPush } from "./lib/push.svelte";

// Initialize theme from localStorage or system preference
const savedTheme = localStorage.getItem("theme");
if (
  savedTheme === "dark" ||
  (!savedTheme && window.matchMedia("(prefers-color-scheme: dark)").matches)
) {
  document.documentElement.setAttribute("data-theme", "dark");
}

// Initialize locale on HTML element
document.documentElement.setAttribute("lang", getLocale());

const app = mount(App, { target: document.getElementById("app")! });

initPush();

export default app;
