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
    daisyui: {
        themes: [
            {
                light: {
                    primary: '#ff6699',
                    secondary: '#a36ffd',
                    accent: '#ffbe00',
                    neutral: '#a3a7af',

                    'base-100': '#ffffff',
                    'base-200': '#f3f4f6',
                    'base-300': '#d1d5db',
                    'base-content': '#0f172a',

                    info: '#3056d3',
                    success: '#00bd8d',
                    warning: '#ffa200',
                    error: '#ff6a3d'
                }
            },
            {
                dark: {
                    primary: '#ff6699',
                    secondary: '#a36ffd',
                    accent: '#ffd043',
                    neutral: '#22212c',

                    'base-100': '#16181d',
                    'base-200': '#2a2e37',
                    'base-300': '#2f323c',
                    'base-content': '#f8fafc',

                    info: '#3056d3',
                    success: '#00bd8d',
                    warning: '#ffa200',
                    error: '#ff6a3d'
                }
            }
        ]
    }
} satisfies Config
