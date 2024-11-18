const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ", "../../../components/**/*.templ"],
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
      transitionProperty: {
        'backdrop': 'backdrop-filter',
      }
    },
  },
  prefix: "o-",
  safelist: ["quail-image-wrapper"],
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

