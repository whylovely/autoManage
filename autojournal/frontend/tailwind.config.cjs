/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{vue,js,ts}'],
  theme: {
    extend: {
      colors: {
        ink: '#18212f',
        canvas: '#f4f6f8',
        brand: {
          50: '#effaf7',
          100: '#d8f3ea',
          500: '#169b78',
          600: '#0f7d62',
          700: '#0d6551',
        },
      },
      boxShadow: {
        panel: '0 12px 35px rgba(30, 42, 55, 0.08)',
      },
    },
  },
  plugins: [],
}
