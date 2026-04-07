import "./app.css";
import App from "./App.svelte";
import { mount } from "svelte";

// Initialize theme from localStorage or system preference
const savedTheme = localStorage.getItem("theme");
if (
  savedTheme === "dark" ||
  (!savedTheme && window.matchMedia("(prefers-color-scheme: dark)").matches)
) {
  document.documentElement.setAttribute("data-theme", "dark");
}

const app = mount(App, { target: document.getElementById("app")! });

export default app;
