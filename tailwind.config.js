/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{templ,html,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require("flyonui")
  ],
  flyonui: {
    themes: ["dark"],
    darkTheme: "dark"
  }
}
