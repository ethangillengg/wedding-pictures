/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.templ", "./**/*.go"],
  theme: {
    extend: {
      gridTemplateColumns: {
        // Complex site-specific column configuration
        gallery: "repeat( auto-fill, 24rem)",
      },
    },
  },
  plugins: [],
};
