/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./node_modules/apexcharts/**/*.js",
    "./node_modules/flyonui/dist/js/helper-apexcharts.js",
    "./views/**/*.{templ,html,js}"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("flyonui")
  ],
  flyonui: {
    themes: ["dark"],
    darkTheme: "dark",
    vendors: true
  }
}
