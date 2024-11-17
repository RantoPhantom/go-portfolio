/** @type {import('tailwindcss').Config} */
	module.exports = {
		content: ["./views/*.{html,js}"],
		variants: {
			extend: {
				display: ["group-hover"],
			},
		},
		theme: {
			extend: {},
		},
		plugins: [],
	}
