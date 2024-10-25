import { defineConfig } from 'tsup';

export default defineConfig({
    target: 'esnext',
    clean: true,
    dts: true,
    entry: ['src/index.ts'],
    keepNames: true,
    minify: true,
    sourcemap: true,
    format: ['cjs'],
})
