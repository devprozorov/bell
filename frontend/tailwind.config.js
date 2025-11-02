/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./components/**/*.{vue,js,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./app.vue",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#2563EB", // основной синий (blue-600)
          light: "#3B82F6",   // чуть ярче
          dark: "#1E40AF",    // темнее
        },
        background: {
          DEFAULT: "#FFFFFF",
          light: "#F9FAFB",
          dark: "#E5E7EB",
        },
        text: {
          DEFAULT: "#1E293B", // почти чёрный с синим оттенком
          light: "#334155",
        },
      },
    },
  },
  plugins: [],
}
