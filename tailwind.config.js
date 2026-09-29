import colors from 'tailwindcss/colors'
import forms from '@tailwindcss/forms'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: colors.indigo,
        surface: colors.slate,
      }
    },
  },
  plugins: [
    forms,
  ],
}


