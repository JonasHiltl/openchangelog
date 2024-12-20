const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./admin/views/**/*.templ", "./icons/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  prefix: "o-",
  plugins: [
    require('tailwind-scrollbar-hide'),
    require('daisyui'),
  ]
}

