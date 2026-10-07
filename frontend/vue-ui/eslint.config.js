import js from '@eslint/js'
import tseslint from 'typescript-eslint'
import vue from 'eslint-plugin-vue'
import vueParser from 'vue-eslint-parser'

export default tseslint.config(
  {
    ignores: [
      '**/node_modules/**',
      '**/dist/**',
      '**/dist-ssr/**',
      '**/coverage/**',
      '**/generated/**',
      '**/*.d.ts',
    ],
  },
  js.configs.recommended,
  {
    files: ['**/*.{ts,tsx,mts,cts,vue}'],
    extends: [...tseslint.configs.recommended],
  },
  ...vue.configs['flat/essential'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: { parser: tseslint.parser },
    },
    // ESLint's assignment analysis cannot see reads from Vue templates.
    // TypeScript and the Vue-aware unused-variable rule still check bindings.
    rules: { 'no-useless-assignment': 'off' },
  },
  {
    files: ['**/*.{js,mjs,cjs,ts,tsx,mts,cts,vue}'],
    rules: { complexity: ['error', 15] },
  },
)
