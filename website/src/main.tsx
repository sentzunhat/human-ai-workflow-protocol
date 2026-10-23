import { createRoot, hydrateRoot } from "react-dom/client";
import { App } from "./app";
import "./app.css";

const root = document.getElementById("root");
if (!root) throw new Error("Missing React root");
if (import.meta.env.DEV) createRoot(root).render(<App />);
else hydrateRoot(root, <App />);
