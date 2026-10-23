import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  base: './',
  plugins: [react(), tailwindcss()],
  publicDir: "static",
  build: { outDir: "build", emptyOutDir: true },
});
