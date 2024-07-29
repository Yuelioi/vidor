import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'

import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        AutoImport({
            dts: 'auto-imports.d.ts',
            include: [
                /\.[t]sx?$/, // .ts, .tsx
                /\.vue$/,
                /\.vue\?vue/, // .vue
                /\.md$/ // .md
            ],
            eslintrc: {
                enabled: true
            },
            imports: [
                'vue',
                'vue-router',
                {
                    pinia: ['defineStore', 'storeToRefs']
                },
                {
                    from: './src/models/ui',
                    imports: ['Tab'],
                    type: true
                },
                {
                    from: './src/models/go',
                    imports: ['Config', 'MessageData'],
                    type: true
                }
            ],
            resolvers: [IconsResolver({})],

            dirs: [
                './src/hooks/**',
                './src/stores/',
                './src/plugins/*',
                './wailsjs/go/app/',
                './wailsjs/runtime/'
            ]
        }),
        Components({
            dts: 'components.d.ts',
            dirs: ['src/components/**', 'src/assets/icons', 'src/views'],
            extensions: ['vue'],
            include: [/\.vue$/, /\.vue\?vue/, /\.md$/],
            exclude: [/[\\/]node_modules[\\/]/, /[\\/]\.git[\\/]/],
            resolvers: [IconsResolver({})]
        }),
        Icons({
            compiler: 'vue3',
            autoInstall: true
        })
    ],
    resolve: {
        alias: {
            vue: 'vue/dist/vue.esm-bundler.js',
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            '@wailsjs': fileURLToPath(new URL('./wailsjs', import.meta.url))
        }
    },
    server: {
        hmr: {
            overlay: false
        }
    }
})
