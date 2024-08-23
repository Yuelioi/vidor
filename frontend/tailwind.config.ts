import type { Config } from 'tailwindcss'
import daisyui from 'daisyui'

import { addDynamicIconSelectors } from '@iconify/tailwind'

export default {
  darkMode: ['selector', '[data-theme="dark"]'],
  content: ['./index.html', './src/**/*.{js,ts,vue,css}'],
  theme: {
    extend: {}
  },
  plugins: [daisyui, addDynamicIconSelectors()],
  safelist: ['hover:text-primary', 'hover:text-secondary', 'hover:text-accent'],

  daisyui: {
    themes: [
      {
        light: {
          primary: '#ff6699',
          secondary: '#00aeec',
          accent: '#a36ffd',
          neutral: '#a3a7af',
          'neutral-content': '#f8fafc',

          'base-100': '#ffffff',
          'base-200': '#f6f8fa',
          'base-300': '#a6adbb',
          'base-content': '#0f172a',

          info: '#3056d3',
          success: '#00bd8d',
          warning: '#ffa200',
          error: '#f53135'
        }
      },
      {
        dark: {
          primary: '#ff6699',
          secondary: '#4ac7ff',
          accent: '#a36ffd',
          neutral: '#22212c',
          'neutral-content': '#f8fafc',

          'base-100': '#16181d',
          'base-200': '#2a2e37',
          'base-300': '#2f323c',
          'base-content': '#f8fafc',

          info: '#3056d3',
          success: '#00bd8d',
          warning: '#ffa200',
          error: '#f53135'
        }
      }
    ]
  }
} satisfies Config
