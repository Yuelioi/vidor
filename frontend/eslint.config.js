import globals from 'globals'
import tseslint from '@typescript-eslint/eslint-plugin'
import pluginVue from 'eslint-plugin-vue'

export default [
  {
    files: ['**/*.{js,mjs,cjs,ts,vue}'],
    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        parser: tseslint.parser
      }
    },
    rules: {
      '@typescript-eslint/no-namespace': 'off'
    }
  },
  ...tseslint.configs.recommended,
  ...pluginVue.configs['flat/essential'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser
      }
    }
  }
]
