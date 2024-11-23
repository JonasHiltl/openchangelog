const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "./icons/**/*.templ", "../../../components/**/*.templ"],
  darkMode: ['selector', '[color-scheme="dark"]'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans],
      },
      colors: (theme) => ({
        "primary": theme.colors.blue[500],
        "caption": theme.colors.gray[400]
      }),
      keyframes: {
        "slide-bottom": {
          "0%": {
            transform: "translateY(100%)"
          },
          "100%": {
            transform: "translateY(0)"
          }
        },
        "slide-right": {
          "0%": {
            transform: "translateX(0)"
          },
          "100%": {
            transform: "translateX(100%)"
          }
        }
      },
      animation: {
        "slide-bottom": "slide-bottom 0.2s ease-out",
        "slide-right": "slide-right 0.2s ease-in",
      },
    },
  },
  prefix: "o-",
  safelist: ["quail-image-wrapper"],
  plugins: [
    require('tailwind-scrollbar-hide'),
    require('@tailwindcss/typography'),
    plugin(function ({ addVariant }) {
      addVariant('htmx-settling', ['&[class~="htmx-settling"]'])
      addVariant('htmx-request', ['&[class~="htmx-request"]'])
      addVariant('htmx-swapping', ['&[class~="htmx-swapping"]'])
      addVariant('htmx-added', ['&[class~="htmx-added"]'])
    }),
  ],
}

