/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				'blue': {
          50: '#B2CFEB',
          100: '#A2C5E6',
          200: '#82B1DE',
          300: '#629DD5',
          400: '#4189CD',
          500: '#3074B6',
          600: '#245889',
          700: '#193C5D',
          800: '#0D1F30',
          900: '#010304',
          950: '#000000'
				}
			}
		}
	},
	plugins: []
};
