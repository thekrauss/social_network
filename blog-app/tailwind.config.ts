import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background)", // Utilise les variables CSS
        foreground: "var(--foreground)", // Utilise les variables CSS
      },
    },
  },
  darkMode: "class", // Utilise la classe `dark` pour basculer le mode
  plugins: [],
};

export default config;
